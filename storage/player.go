package storage

import (
	"errors"
	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Character reprents a game character
type Character struct {
	gorm.Model
	UUID          uuid.UUID
	CharacterName string `valid:"alphanum,required,runelength(2|32)"`
	Color         string `valid:"hexcolor,required"`
}

// Player represents a player in the game that can manage several Characters
type Player struct {
	gorm.Model
	UUID       uuid.UUID
	Name       string      `valid:"alphanum,required,runelength(2|32)"`
	UserName   string      `valid:"alphanum,required,runelength(2|32)"`
	Color      string      `valid:"hexcolor,required"`
	Characters []Character ``
}

const (
	MinCharacterNameLength = 2
	MaxCharacterNameLength = 32
	MinPlayerNameLength    = 2
	MaxPlayerNameLength    = 32
	MinUsernameLength      = 2
	MaxUsernameNameLength  = 32
)

// RegisterPlayer register the data about a player
func RegisterPlayer(name string, username string, color string) (*Player, error) {
	trimedname := valid.Trim(name, "") // Remove starting an tailing whitespace
	// TODO: finish validations
	if valid.IsAlphanumeric(trimedname) && len(trimedname) <= MaxSessionNameLength && len(trimedname) >= MinSessionNameLength {
		player := Player{UUID: uuid.New(), Name: trimedname, UserName: username}
		log.Infof("Registered Player: %v", player)
		return &player, nil
	} else {
		return nil, errors.New("Invalid Player Name")
	}
}

// RegisterCharacter register a Character controled by a player
func (*Player) RegisterCharacter(charactername string, color string) (*Character, error) {
	trimedname := valid.Trim(charactername, "") // Remove starting an tailing whitespace
	if valid.IsAlphanumeric(trimedname) && len(trimedname) <= MaxSessionNameLength && len(trimedname) >= MinSessionNameLength {
		character := Character{UUID: uuid.New(), CharacterName: trimedname}
		log.Infof("Registered Character: %v", character)
		return &character, nil
	} else {
		return nil, errors.New("Invalid Character Name")
	}
}
