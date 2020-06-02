package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	storage "github.com/pconcepcion/telegram_dice_bot/storage"
	"github.com/pkg/errors"
)

var (
	// ErrInvalidTelegramUser it's used when a Telegram User it's not valid
	ErrInvalidTelegramUser = errors.New("Invalid Telegram User")
	// ErrBotUser it's used when a Telegram User it's a Bot
	ErrBotUser = errors.New("Telegram User it's a bot and bot's are not allowed to interact with telegram_dice_bot")
)

const (
	// DefaultPlayerColor just a default color to use with players
	DefaultPlayerColor = "#333"
)

func (b *bot) getActivePlayer(player *tgbotapi.User) (*storage.Player, error) {
	if player == nil {
		return nil, errors.Wrap(ErrInvalidTelegramUser, "User it's nil")
	}
	if player.IsBot {
		return nil, ErrBotUser
	}
	log.Debugf("received player: %#v", player)
	activePlayer, ok := b.ActivePlayers[player.ID]
	if ok == false {
		log.Infof("Player ID %d not found among acitve players, adding it", player.ID)
		var err error
		activePlayer, err = b.storage.Player(player.FirstName, player.UserName, player.ID, DefaultPlayerColor)
		if err != nil {
			log.Error(err)
			// TODO Use an specific error
			return nil, ErrInvalidTelegramUser
		}
		log.Infof("Stored player %#v", activePlayer)
		b.ActivePlayers[player.ID] = activePlayer
		log.Infof("Stored player 2 %#v", activePlayer)
	} else {
		log.Debug("Player ID %d already among active players", player.ID)
	}
	log.Warn("Active Player to return:", activePlayer)
	return activePlayer, nil
}
