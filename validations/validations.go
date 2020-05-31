package validations

import (
	"fmt"

	valid "github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

var (
	// ErrInvalidName it's trown when a name contains invalid characters
	ErrInvalidName = errors.New("Invalid name Error")
	// ErrSessionNameTooLong it's trown when a session name is longer tha MaxSessionNameLength
	ErrSessionNameTooLong = fmt.Errorf("Session Name Too Long. Maximum length %d", MaxSessionNameLength)
	// ErrSessionNameTooShort it's trown when a session name is shorter tha MinSessionNameLength
	ErrSessionNameTooShort = fmt.Errorf("Session Name Too Short, minimum session name is %d", MinSessionNameLength)
	// ErrInvalidCharacterName it's trhown when a character name contains invalid characters
	ErrInvalidCharacterName = errors.New("Invalid Character Name")
	// ErrSessionNameTooLong it's trown when a session name is longer tha MaxSessionNameLength
	ErrCharacterNameTooLong = fmt.Errorf("Character Name Too Long. Maximum length %d", MaxCharacterNameLength)
	// ErrCharacterNameTooShort it's trown when a session name is shorter tha MinCharacterNameLength
	ErrCharacterNameTooShort = fmt.Errorf("Character Name Too Short, minimum character name is %d", MinCharacterNameLength)
	// ErrInvalidPlayerName it's trhown when a character name contains invalid characters
	ErrInvalidPlayerName = errors.New("Invalid Player Name")
	// ErrPlayerNameTooLong it's trown when a session name is longer tha MaxPlayerNameLength
	ErrPlayerNameTooLong = fmt.Errorf("Player Name Too Long. Maximum length %d", MaxPlayerNameLength)
	// ErrPlayerNameTooShort it's trown when a session name is shorter tha MinPlayerNameLength
	ErrPlayerNameTooShort = fmt.Errorf("Player Name Too Short, minimum player name is %d", MinPlayerNameLength)
)

const (
	// MinSessionNameLength minimum length for the Session name
	MinSessionNameLength = 2
	// MaxSessionNameLength maximum length for the Session name
	MaxSessionNameLength = 32
	// MinCharacterNameLength minimum length for the Character name
	MinCharacterNameLength = 2
	// MaxCharacterNameLength maximum length for the Character name
	MaxCharacterNameLength = 32
	// MinPlayerNameLength minimum length for the Player name
	MinPlayerNameLength = 2
	// MaxPlayerNameLength maximum length for the Player name
	MaxPlayerNameLength = 32
	// MinUsernameLength minimum length for the Username
	MinUsernameLength = 2
	// MaxUsernameNameLength maximum length for the Username
	MaxUsernameNameLength = 32
)

// ValidateSessionName ensures a sesion name is valid
func ValidateSessionName(name string) (string, error) {
	trimedname := valid.Trim(name, "") // Remove starting an tailing whitespace
	if !valid.IsAlphanumeric(trimedname) {
		return name, ErrInvalidName
	}
	if len(trimedname) > MaxSessionNameLength {
		return name, ErrSessionNameTooLong
	}
	if len(trimedname) < MinSessionNameLength {
		return name, ErrSessionNameTooShort
	}
	return trimedname, nil
}

// ValidatePlayerName ensures a player name is valid
func ValidatePlayerName(name string) (string, error) {
	trimedname := valid.Trim(name, "") // Remove starting an tailing whitespace
	if !valid.IsAlphanumeric(trimedname) {
		return name, ErrInvalidPlayerName
	}
	if len(trimedname) > MaxPlayerNameLength {
		return name, ErrPlayerNameTooLong
	}
	if len(trimedname) < MinPlayerNameLength {
		return name, ErrPlayerNameTooShort
	}
	return trimedname, nil
}

// ValidateCharacterName ensures a character name is valid
func ValidateCharacterName(name string) (string, error) {
	trimedname := valid.Trim(name, "") // Remove starting an tailing whitespace
	if !valid.IsAlphanumeric(trimedname) {
		return name, ErrInvalidCharacterName
	}
	if len(trimedname) > MaxCharacterNameLength {
		return name, ErrCharacterNameTooLong
	}
	if len(trimedname) < MinCharacterNameLength {
		return name, ErrCharacterNameTooShort
	}
	return trimedname, nil
}
