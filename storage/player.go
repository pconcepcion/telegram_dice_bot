package storage

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pconcepcion/telegram_dice_bot/validations"
	"github.com/pkg/errors"
)

// Character database model that reprents a game character
type Character struct {
	gorm.Model
	UUID          uuid.UUID
	PlayerUUID    uuid.UUID
	CharacterName string `valid:"alphanum,required,runelength(2|32)"`
	Color         string `valid:"hexcolor,required"`
}

// Player database model that represents a player in the game that can manage several Characters
type Player struct {
	gorm.Model
	UUID           uuid.UUID
	UserTelegramID int64       `valid:"type(int64),required"`
	Name           string      `valid:"alphanum,required,runelength(2|32)"`
	UserName       string      `valid:"alphanum,required,runelength(2|32)"`
	Color          string      `valid:"hexcolor,required"`
	Characters     []Character `gorm:"foreingkey:PlayerUUID`
	Rolls          []Roll      `gorm:"foreignkey:PlayerUUID"`
}

// Player get's a player from the storage or stores a new player and returns it
func (sqliteStorage SQLiteStorage) Player(name string, username string, telegramUserID int64, color string) (*Player, error) {
	var player Player
	//trimedcolor := valid.Trim(color, "") // Remove starting an tailing whitespace
	validName, err := validations.ValidatePlayerName(name)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't store player, invalid name")
	}
	validUserName, err := validations.ValidatePlayerName(username)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't store player, invalid username")
	}
	// TODO: finish validations
	sqliteStorage.db.Where(Player{UserTelegramID: telegramUserID}).Attrs(Player{UUID: uuid.New(), Name: validName, UserName: validUserName}).FirstOrInit(&player)
	log.Infof("Registered Player: %v", player)
	return &player, nil
}

// RegisterCharacter register a Character controled by a player
func (sqliteStorage SQLiteStorage) RegisterCharacter(p *Player, charactername string, color string) (*Character, error) {
	trimedColor := valid.Trim(color, "")
	validCharacterName, err := validations.ValidateCharacterName(charactername)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't store character, invalid character name")
	}
	if !valid.IsHexcolor(trimedColor) {
		return nil, errors.Wrap(err, "Couldn't store character, invalid color")
	}
	character := Character{UUID: uuid.New(), CharacterName: validCharacterName, Color: trimedColor}
	///sqliteStorage.db.Where(Player{UserTelegramID: telegramUserID}).Attrs(Player{UUID: uuid.New(), Name: trimedname, UserName: trimedUsername}).FirstOrInit(&player)
	p.Characters = append(p.Characters, character)
	sqliteStorage.db.Save(character)
	sqliteStorage.db.Save(p)
	log.Infof("Player %s registered Character: %s", p.Name, character.CharacterName)
	return &character, nil
}
