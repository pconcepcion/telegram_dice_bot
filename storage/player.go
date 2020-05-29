package storage

import (
	"errors"

	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Character database model that reprents a game character
type Character struct {
	gorm.Model
	UUID          uuid.UUID
	CharacterName string `valid:"alphanum,required,runelength(2|32)"`
	Color         string `valid:"hexcolor,required"`
}

// Player database model that represents a player in the game that can manage several Characters
type Player struct {
	gorm.Model
	UUID       uuid.UUID
	Name       string      `valid:"alphanum,required,runelength(2|32)"`
	UserName   string      `valid:"alphanum,required,runelength(2|32)"`
	Color      string      `valid:"hexcolor,required"`
	Characters []Character ``
	Rolls      []Roll      `gorm:"foreignkey:PlayerUUID"`
}

// Errors
var (
	// ErrInvalidCharacterName
	ErrInvalidCharacterName = errors.New("Invalid Character Name")
	// ErrInvalidPlayerName
	ErrInvalidPlayerName = errors.New("Invalid Player Name")
)

const (
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

// NewPlayer creates a new player
func (sqliteStorage SQLiteStorage) NewPlayer(name string, username string, color string) (*Player, error) {
	trimedname := valid.Trim(name, "") // Remove starting an tailing whitespace
	// TODO: finish validations
	if valid.IsAlphanumeric(trimedname) && len(trimedname) <= MaxSessionNameLength && len(trimedname) >= MinSessionNameLength {
		player := Player{UUID: uuid.New(), Name: trimedname, UserName: username}
		sqliteStorage.db.Save(player)
		log.Infof("Registered Player: %v", player)
		return &player, nil
	}
	return nil, ErrInvalidPlayerName
}

// RegisterCharacter register a Character controled by a player
func (sqliteStorage SQLiteStorage) RegisterCharacter(p *Player, charactername string, color string) (*Character, error) {
	trimedCharactername := valid.Trim(charactername, "") // Remove starting an tailing whitespace
	trimedColor := valid.Trim(color, "")
	if valid.IsAlphanumeric(trimedCharactername) && len(trimedCharactername) <= MaxSessionNameLength &&
		len(trimedCharactername) >= MinSessionNameLength && valid.IsHexcolor(trimedColor) {
		character := Character{UUID: uuid.New(), CharacterName: trimedCharactername, Color: trimedColor}
		p.Characters = append(p.Characters, character)
		sqliteStorage.db.Save(character)
		sqliteStorage.db.Save(p)
		log.Infof("Player %s registered Character: %s", p.Name, character.CharacterName)
		return &character, nil
	}
	return nil, ErrInvalidCharacterName
}
