package bot

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	storage "github.com/pconcepcion/telegram_dice_bot/storage"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	//valid "github.com/asaskevich/govalidator"
)

var log = logrus.New()

const (
	apiTokenMinLength     = 5
	responseArrowTemplate = "\u27A1 %d"
	rigthArrow            = "\u27A1"
	label                 = "\U0001F3F7"
)

func init() {
	log.Formatter = &logrus.TextFormatter{DisableLevelTruncation: true, FullTimestamp: true, PadLevelText: true}
	log.Level = logrus.DebugLevel

}

// SessionsMap maps the Chat ID to the Session name
type SessionsMap map[int64]*storage.Session

// bot holds the configuration and the reference to the API for the bot
type bot struct {
	apiToken       string
	api            *tgbotapi.BotAPI
	timeout        int // Bot timeout in seconds
	updateConfig   tgbotapi.UpdateConfig
	storage        *storage.SQLiteStorage
	ActiveSessions SessionsMap
}

// getAPIToken gets the API token from the configuration and does some basic validation
func (b *bot) getAPIToken() {
	b.apiToken = viper.GetString("api_token")
	if b.apiToken == "" || len(b.apiToken) < apiTokenMinLength {
		log.Error("API Token not found")
		os.Exit(-1)
	}
	log.Debugf("Found API Token: %sXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\n", b.apiToken[0:3])
}

// authorize does the authorization of the bot and set the bot mode
func (b *bot) authorize(debug bool) {
	var err error
	b.getAPIToken()
	log.Printf("Authorizing bot with token '%sXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX'", b.apiToken[0:3])
	b.api, err = tgbotapi.NewBotAPI(b.apiToken)
	if err != nil {
		errorMsg := fmt.Sprintf("Couldn't access the API with token starting with '%s'", b.apiToken[0:3])
		log.Panic(errors.Wrap(err, errorMsg))
	}

	b.api.Debug = debug

	log.Printf("Authorized on account %s", b.api.Self.UserName)
	log.Printf("debug = %t", debug)
}

func botSetup(debug bool) *bot {
	// Authorize the bot with debug mode
	bot := bot{timeout: 30}
	bot.ActiveSessions = make(map[int64]*storage.Session, 10)
	bot.authorize(debug)
	bot.updateConfig = tgbotapi.NewUpdate(0)
	bot.updateConfig.Timeout = bot.timeout

	self, err := bot.api.GetMe()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("Bot info: %#v", self)
	}

	return &bot
}

// setupStorage setups the storage confiuration and keeps a referenc on the bot
func (b *bot) setupStorage(storageAccessString string) {
	log.Debugln("About to connect to Storage: ", storageAccessString)
	b.storage = storage.Connect(storageAccessString)
}

// Run launches the bot, does the authorization process and starts to listen for messages
func Run() {
	storageAccessString := viper.GetString("storage")
	// setup the bot with debug mode
	bot := botSetup(true)
	updates, err := bot.api.GetUpdatesChan(bot.updateConfig)
	if err != nil {
		log.Panic(err)
	}
	bot.setupStorage(storageAccessString)
	defer bot.storage.Close()
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		log.Debugf("---\n%+v\n---", update.Message)
		log.Debugf("---\n%+v\n---", update.Message.Chat)
		response := bot.handleMessage(update.Message)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		msg.ParseMode = "MarkdownV2"
		msg.ReplyToMessageID = update.Message.MessageID
		if _, sendErr := bot.api.Send(msg); sendErr != nil {
			log.Error(sendErr)
		}
		log.Info("Response:", msg)
	}
}
