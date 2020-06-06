package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pconcepcion/telegram_dice_bot/validations"
	"github.com/pkg/errors"
)

var (
	// ErrMissingPlayer is used when the an operation that needs a player and the player it's missing
	ErrMissingPlayer = errors.New("Missing Player")
)

// Session  database model that stores information on the game session,
// including an UUID, a name and a Starting time for the Session
type Session struct {
	BaseModel
	Name           string `valid:"alphanum,required,runelength(2|32)"`
	StartTime      time.Time
	EndTime        time.Time
	Rolls          []*Roll   `gorm:"foreignkey:SessionUUID"`
	Players        []*Player `gorm:"many2many:PlayerSession"`
	ChatTelegramID int64
}

// Roll database model to sotre a single dice expression roll
type Roll struct {
	BaseModel
	SessionUUID      uuid.UUID `valid:"uuidv4"`
	PlayerUUID       uuid.UUID `valid:"uuidv4"`
	PlayerTelegramID int       `valid:"int,required"`
	Expression       string    `valid:diceexpression,required`
	RawResults       string    `valid:required,runelength(1|4096)`
	Results          string    `valid:required,runelength(1|4096)`
	Total            int       `valid:"int,required"`
	Description      string
}

// StartSession creates a new Session object with a current sarting time and the given name
func (sqliteStorage *SQLiteStorage) StartSession(name string, chatID int64) (*Session, error) {
	validName, err := validations.ValidateSessionName(name) //
	if err != nil {
		return nil, errors.Wrap(err, "Invalid Session Name")
	}
	session := &Session{Name: validName, StartTime: time.Now(), ChatTelegramID: chatID}
	log.Infof("Start Session: %v", session)
	if err := sqliteStorage.db.Create(session).Error; err != nil {
		log.Errorf("Error creating session on DB '%v' : %s", session, err)
		return nil, err
	}
	return session, nil
}

// RenameSession changes the name of the session and stores it on the SQLite DB
func (sqliteStorage *SQLiteStorage) RenameSession(ses *Session, name string) error {
	validName, err := validations.ValidateSessionName(name) //
	if err != nil {
		return errors.Wrap(err, "Invalid Session Name")
	}
	ses.Name = validName
	if err := sqliteStorage.db.Save(ses).Error; err != nil {
		log.Errorf("Error updating session on DB '%s' : %s", validName, err)
		return err
	}

	return nil
}

// EndSession stores the end time for the session and unset it as the active session
func (sqliteStorage *SQLiteStorage) EndSession(ses *Session) {
	ses.EndTime = time.Now()
	// TODO: Add error handling
	sqliteStorage.db.Save(ses)
}

// RegisterRoll registers the result of rolling a dice expression on the storage backend and links it to the session received
func (sqliteStorage *SQLiteStorage) RegisterRoll(expression, results, rawResults, rollMessage string, total int, session *Session, player *Player) error {
	// TODO: validate UUID
	if player == nil {
		log.Errorln("Missing player while storing a dice roll")
		return ErrMissingPlayer
	}
	roll := Roll{PlayerTelegramID: player.UserTelegramID, Expression: normalizeExpression(expression), Results: results, RawResults: rawResults, Total: total}
	validRollMessage, err := validations.ValidateDescription(rollMessage)
	if err != nil {
		log.Warn("Invalid Description: ", err)
	} else {
		roll.Description = validRollMessage
	}
	if session != nil {
		// TODO: Check if this can be done by the ORM
		roll.SessionUUID = session.ID
	}
	log.Debug("Registering Roll: ", roll)
	// Link the roll wiht the player
	if err := sqliteStorage.db.Model(player).Association("Rolls").Append(roll).Error; err != nil {
		log.Error("Error registering roll: ", err)
		return err
	}
	log.Debug("Stored Roll: ", roll, "Session: ", session)
	return nil
}

// Normalize Expression normalizes the expresion to store it on the DB
// So far it only adds a 1 when the number of dices was omited.
func normalizeExpression(expression string) string {
	if strings.HasPrefix(expression, "d") {
		return "1" + expression
	}
	return expression
}

// AddPlayerToSessoionIfMissing adds a player to a session if it was not alredy included
func (sqliteStorage *SQLiteStorage) AddPlayerToSessoionIfMissing(player *Player, session *Session) error {
	for _, p := range session.Players {
		if p.ID == player.ID {
			return nil
		}
	}
	if err := sqliteStorage.db.Model(session).Association("Players").Append(player).Error; err != nil {
		log.Error("Error adding player to session: ", err)
		return err
	}
	// TODO: Error handling
	log.Debugf("AddPlayerToSessoionIfMissing\nPlayer:\n %#v \nSession: %#v\n", player, session)
	return nil
}

// SummaryEntry holds the data from an entry of the summary
type SummaryEntry struct {
	Username   string
	Expression string
	NumRolls   int `gorm:"column:numrolls"`
	Minimum    int
	Maximum    int
	Average    float64
}

// String is the SummaryEntry Stringer
func (s *SummaryEntry) String() string {
	res := fmt.Sprintf("%15s  \\| %10s \\| %6d \\| %4d \\| %4d \\| %6f\n", s.Username, s.Expression, s.NumRolls, s.Minimum, s.Maximum, s.Average)
	return strings.Replace(res, ".", "\\.", -1)

}

// SummaryQuery is the raw query to obtain the summary data from the DB, use the username (or the name if the username is null)
var SummaryQuery = `SELECT  (p.user_name, p.name) as username, r.expression, COUNT(r.expression) as numrolls,
	MIN(r.total) as minimum, MAX(r.total) as maximum, AVG(r.total) as average
	FROM rolls r join players p ON r.player_uuid = p.id AND r.session_uuid = ?
	GROUP BY r.expression
	ORDER BY username, r.expression`

// Summary get's the summary for a session
func (sqliteStorage *SQLiteStorage) Summary(sesion *Session) (summary []SummaryEntry, err error) {
	var summaryEntry SummaryEntry
	var rows *sql.Rows
	sqliteStorage.db.LogMode(true)
	if rows, err = sqliteStorage.db.Raw(SummaryQuery, sesion.ID).Rows(); err != nil {
		return nil, errors.Wrap(err, "Error obtaining summary. ")
	}
	log.Debugf("ROWS!! %#v", rows)
	for rows.Next() {
		sqliteStorage.db.ScanRows(rows, &summaryEntry)
		log.Debugf("summaryEntry: %#v", summaryEntry)
		summary = append(summary, summaryEntry)
	}
	sqliteStorage.db.LogMode(false)
	return summary, nil
}
