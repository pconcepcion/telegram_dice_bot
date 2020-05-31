package storage

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pconcepcion/dice"
	"github.com/pconcepcion/telegram_dice_bot/validations"
	"github.com/pkg/errors"
)

// Session  database model that stores information on the game session,
// including an UUID, a name and a Starting time for the Session
type Session struct {
	gorm.Model
	UUID           uuid.UUID
	Name           string `valid:"alphanum,required,runelength(2|32)"`
	SartTime       time.Time
	EndTime        time.Time
	Rolls          []Roll   `gorm:"foreignkey:SessionUUID"`
	Players        []Player `gorm:"foreignkey:PlayerUUID"`
	ChatTelegramID int64
}

// Roll database model to sotre a single dice expression roll
type Roll struct {
	gorm.Model
	SesssionUUID     uuid.UUID
	PlayerUUID       uuid.UUID
	PlayerTelegramID int64
	Results          string
	Total            int
	Description      string
}

// StartSession creates a new Session object with a current sarting time and the given name
func (sqliteStorage *SQLiteStorage) StartSession(name string, chatID int64) (*Session, error) {
	validName, err := validations.ValidateSessionName(name) //
	if err != nil {
		return nil, errors.Wrap(err, "Invalid Session Name")
	}
	session := &Session{UUID: uuid.New(), Name: validName, SartTime: time.Now(), ChatTelegramID: chatID}
	log.Infof("Start Session: %v", session)
	sqliteStorage.db.Create(session)
	return session, nil
}

// RenameSession changes the name of the session and stores it on the SQLite DB
func (sqliteStorage *SQLiteStorage) RenameSession(ses *Session, name string) error {
	validName, err := validations.ValidateSessionName(name) //
	if err != nil {
		return errors.Wrap(err, "Invalid Session Name")
	}
	ses.Name = validName
	sqliteStorage.db.Save(ses)
	return nil
}

// EndSession stores the end time for the session and unset it as the active session
func (sqliteStorage *SQLiteStorage) EndSession(ses *Session) {
	ses.EndTime = time.Now()
	sqliteStorage.db.Save(ses)
}

// RegisterRoll registers the result of rolling a dice expression on the storage backend
func (s Session) RegisterRoll(rollResult dice.ExpressionResult) {

}
