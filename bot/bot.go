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

func handleMessage(m *tgbotapi.Message) string {
	var response string = "Unknown command"
	log.Infof("Bot command: %v, arguments %v", m.Command(), m.CommandArguments())
	switch m.Command() {
	// Dice expression
	case "de":
		diceExpression := m.CommandArguments()
		dicesResult, err := roll(diceExpression)
		if err != nil {
			// TODO: handle the error gracefully
			log.Error(err)
		}
		response = fmt.Sprintf("%s: %v -> %s", diceExpression, dicesResult.GetResults(), dicesResult)
		// Basic dices
	case "d2":
		response = fmt.Sprintf("d2 -> %d", rpg.D2())
	case "d4":
		response = fmt.Sprintf("d4 -> %d", rpg.D4())
	case "d6":
		response = fmt.Sprintf("d6 -> %d", rpg.D6())
	case "d8":
		response = fmt.Sprintf("d8 -> %d", rpg.D8())
	case "d10":
		response = fmt.Sprintf("d10 -> %d", rpg.D10())
	case "d12":
		response = fmt.Sprintf("d12 -> %d", rpg.D12())
	case "d20":
		response = fmt.Sprintf("d20 -> %d", rpg.D20())
	case "d100":
		response = fmt.Sprintf("d100 -> %d", rpg.D100())
	}
	return response
}

// Run launches the bot, does the authorization process and starts to listen for messages
func Run() {
	// Athorize the bot with deboug mode
	bot := authorizeBot(true)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	self, err := bot.GetMe()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("Bot info: %v", self)
	}

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		log.Printf("---\n%+v\n---", update.Message)
		log.Printf("---\n%+v\n---", update.Message.Chat)
		response := handleMessage(update.Message)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, sendErr := bot.Send(msg); sendErr != nil {
			log.Error(sendErr)
		}
	}
}