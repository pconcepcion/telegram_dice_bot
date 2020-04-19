package storage

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Session  database model that stores information on the game session,
// including an UUID, a name and a Starting time for the Session
type Session struct {
	gorm.Model
	UUID     uuid.UUID
	Name     string `valid:"alphanum,required,runelength(2|32)"`
	SartTime time.Time
	EndTime  time.Time
}

const (
	// MinSessionNameLength minimum length for the Session name
	MinSessionNameLength = 2
	// MaxSessionNameLength maximum length for the Session name
	MaxSessionNameLength = 32
)

// StartSession creates a new Session object with a current sarting time and the given name
func StartSession(name string) (*Session, error) {
	trimedname := valid.Trim(name, "") // Remove starting an tailing whitespace
	if valid.IsAlphanumeric(trimedname) && len(trimedname) <= MaxSessionNameLength && len(trimedname) >= MinSessionNameLength {
		session := Session{UUID: uuid.New(), Name: trimedname, SartTime: time.Now()}
		log.Infof("Start Session: %v", session)
		db.Create(session)
		// TODO: Set as active session
		return &session, nil
	}
	return nil, errors.New("Invalid Session Name")
}

// EndSession stores the end time for the session and unset it as the active session
func (s *Session) EndSession() {
	s.EndTime = time.Now()
	db.Save(s)
	// TODO: unset as active session
}
