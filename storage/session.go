package storage

import (
	"errors"
	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"time"
)

// Session stores information on the game session, including an UUID, a name and a Starting time for the Session
type Session struct {
	UUID     uuid.UUID
	Name     string `valid:"alphanum,required,runelength(2|32)"`
	SartTime time.Time
}

const (
	MinSessionNameLength = 2
	MaxSessionNameLength = 32
)

// StartSession creates a new Session object with a current sarting time and the given name
func StartSession(name string) (*Session, error) {
	trimedname := valid.Trim(name, "") // Remove starting an tailing whitespace
	if valid.IsAlphanumeric(trimedname) && len(trimedname) <= MaxSessionNameLength && len(trimedname) >= MinSessionNameLength {
		session := Session{UUID: uuid.New(), Name: trimedname, SartTime: time.Now()}
		log.Infof("Start Session: %v", session)
		return &session, nil
	} else {
		return nil, errors.New("Invalid Session Name")
	}
}
