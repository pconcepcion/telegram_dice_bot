package storage

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pconcepcion/dice"
)

// Session  database model that stores information on the game session,
// including an UUID, a name and a Starting time for the Session
type Session struct {
	gorm.Model
	UUID           uuid.UUID
	Name           string `valid:"alphanum,required,runelength(2|32)"`
	SartTime       time.Time
	EndTime        time.Time
	Rolls          []Roll `gorm:"foreignkey:SessionUUID"`
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

const (
	// MinSessionNameLength minimum length for the Session name
	MinSessionNameLength = 2
	// MaxSessionNameLength maximum length for the Session name
	MaxSessionNameLength = 32
)

// StartSession creates a new Session object with a current sarting time and the given name
func (sqliteStorage *SQLiteStorage) StartSession(name string, chatID int64) (*Session, error) {
	trimedname := valid.Trim(name, "") // Remove starting an tailing whitespace
	if valid.IsAlphanumeric(trimedname) && len(trimedname) <= MaxSessionNameLength && len(trimedname) >= MinSessionNameLength {
		session := &Session{UUID: uuid.New(), Name: trimedname, SartTime: time.Now(), ChatTelegramID: chatID}
		log.Infof("Start Session: %v", session)
		sqliteStorage.db.Create(session)
		return session, nil
	}
	return nil, errors.New("Invalid Session Name")
}

// EndSession stores the end time for the session and unset it as the active session
func (sqliteStorage *SQLiteStorage) EndSession(ses *Session) {
	ses.EndTime = time.Now()
	sqliteStorage.db.Save(ses)
}

// RegisterRoll registers the result of rolling a dice expression on the storage backend
func (s Session) RegisterRoll(rollResult dice.ExpressionResult) {

}
