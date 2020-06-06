package storage

import (
	"github.com/google/uuid"
	"github.com/pconcepcion/telegram_dice_bot/validations"
	"github.com/pkg/errors"
)

// Character database model that reprents a game character
type Character struct {
	//gorm.Model
	BaseModel
	PlayerUUID    uuid.UUID `valid:"uuidv4"`
	CharacterName string    `valid:"alphanum,required,runelength(2|32)"`
	Color         string    `valid:"hexcolor,required"`
}

// Player database model that represents a player in the game that can manage several Characters
type Player struct {
	BaseModel
	UserTelegramID int          `valid:"type(int),required"`
	Name           string       `valid:"alphanum,required,runelength(2|32)"`
	UserName       string       `valid:"alphanum,required,runelength(2|32)"`
	Color          string       `valid:"hexcolor,required"`
	Characters     []*Character `gorm:"foreingkey:PlayerUUID`
	Rolls          []*Roll      `gorm:"foreignkey:PlayerUUID"`
	Sessions       []*Session   `gorm:"many2many:PlayerSession"`

	//`gorm:"many2many:player_sessions;"`
}

var (
	//ErrNameOrUsernameMissing error marks when a user has neither name nor username
	ErrNameOrUsernameMissing = errors.New("Name or Username Missing")
)

// Player get's a player from the storage or stores a new player and returns it
func (sqliteStorage SQLiteStorage) Player(name string, username string, telegramUserID int, color string) (*Player, error) {
	var player Player
	//trimedcolor := valid.Trim(color, "") // Remove starting an tailing whitespace
	validName, err := validations.ValidatePlayerName(name)
	if err != nil && validName != "" {
		return nil, errors.Wrap(err, "Couldn't store player, invalid name")
	}
	validUserName, err := validations.ValidatePlayerName(username)
	if err != nil && validUserName != "" {
		return nil, errors.Wrap(err, "Couldn't store player, invalid username")
	}
	if validName == "" && validUserName == "" {
		return nil, ErrNameOrUsernameMissing
	}
	validColor, err := validations.ValidateColor(color)
	if err != nil {
		// TODO: this can be handled gracefully by setting a default color
		return nil, errors.Wrap(err, "Couldn't store player, invalid color")
	}
	// TODO: finish validations
	if err := sqliteStorage.db.Where(Player{UserTelegramID: telegramUserID}).Attrs(Player{
		Name: validName, UserName: validUserName, Color: validColor}).FirstOrCreate(&player).Error; err != nil {
		log.Errorf("Error accessing or registering player: ", err)
		return nil, err
	}
	log.Debugf("Player: %#v", player)
	return &player, nil
}

// RegisterCharacter register a Character controled by a player
func (sqliteStorage SQLiteStorage) RegisterCharacter(p *Player, charactername string, color string) (*Character, error) {
	validCharacterName, err := validations.ValidateCharacterName(charactername)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't store character, invalid character name")
	}
	validColor, err := validations.ValidateColor(color)
	if err != nil {
		// TODO: this can be handled gracefully by setting a default color
		return nil, errors.Wrap(err, "Couldn't store character, invalid color")
	}
	character := Character{CharacterName: validCharacterName, Color: validColor}
	///sqliteStorage.db.Where(Player{UserTelegramID: telegramUserID}).Attrs(Player{UUID: uuid.New(), Name: trimedname, UserName: trimedUsername}).FirstOrInit(&player)
	p.Characters = append(p.Characters, &character)
	if err := sqliteStorage.db.Save(character).Error; err != nil {
		log.Error("Errror storing character", err)
		return nil, err
	}
	if err := sqliteStorage.db.Save(p).Error; err != nil {
		log.Error("Error storing player while registering the character. ", err)
	}
	log.Infof("Player %s registered Character: %s", p.Name, character.CharacterName)
	return &character, nil
}
