package bot

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/pconcepcion/dice"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var log = logrus.New()

func init() {
	log.Formatter = new(prefixed.TextFormatter)
	log.Level = logrus.DebugLevel

}

var apiToken string

func getAPIToken() {
	apiToken = viper.GetString("api_token")
}

// roll the dice expession contained on the message
func roll(message string) (rpg.DiceExpressionResult, error) {
	toRoll := rpg.NewSimpleDiceExpression(message)
	return toRoll.Roll()
}

// Do the authorization of the bot and set the bot mode
func authorizeBot(debug bool) *tgbotapi.BotAPI {
	getAPIToken()
	if apiToken == "" {
		log.Error("API Token not found")
		os.Exit(-1)
	}
	log.Printf("Authorizing bot with token %s", apiToken)
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Panic(errors.Wrap(err, "Couldn't access the API with token "+apiToken))
	}

	bot.Debug = debug

	log.Printf("Authorized on account %s", bot.Self.UserName)
	log.Printf("debug = %t", debug)

	return bot

}

// Run launches the bot, does the authorization process and starts to listen for messages
func Run() {
	// Athorize the bot with deboug mode
	bot := authorizeBot(true)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		diceExpression := update.Message.Text
		dicesResult, err := roll(diceExpression)
		if err != nil {
			// TODO: handle the error gracefully
			log.Error(err)
		}
		response := fmt.Sprintf("%s: %v -> %s", diceExpression, dicesResult.GetResults(), dicesResult)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, sendErr := bot.Send(msg); sendErr != nil {
			log.Error(sendErr)
		}
	}
}
