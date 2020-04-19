package bot

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	rpg "github.com/pconcepcion/dice"
	storage "github.com/pconcepcion/telegram_dice_bot/storage"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	//valid "github.com/asaskevich/govalidator"
)

var log = logrus.New()

func init() {
	log.Formatter = &logrus.TextFormatter{DisableLevelTruncation: true, FullTimestamp: true, PadLevelText: true}
	log.Level = logrus.DebugLevel

}

var apiToken string

func getAPIToken() {
	apiToken = viper.GetString("api_token")
}

// roll the dice expession contained on the message
func roll(message string) (rpg.ExpressionResult, error) {
	toRoll := rpg.NewSimpleExpression(message)
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
	// TODO: Handle extra argument for the roll identifier: example "/d100 hide"
	var response = "Unknown command"
	log.Infof("Bot command: %v, arguments %v", m.Command(), m.CommandArguments())
	switch m.Command() {
	// Dice expression
	case "de":
		diceExpression := m.CommandArguments()
		dicesResult, err := roll(diceExpression)
		if err != nil {
			// TODO: handle the error gracefully
			log.Error(err)
			response = "Dice expression Error"
		} else {
			response = fmt.Sprintf("%s: %v -> %s", diceExpression, dicesResult.GetResults(), dicesResult)
		}
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
	// Session Handling
	case "startSession":
		sessionName , err := storage.StartSession(m.CommandArguments()) 
		if err != nil {
			response = fmt.Sprintf("Failed to create Session, invalid session name")
			log.Errorf("Failed to create Session, invalid session arguments: %s", m.CommandArguments())
			return response	
		}
		response = fmt.Sprintf("Starting Session: %s", sessionName)
		log.Info(response)
		// TODO: Store session info
		// TODO: Set session timeout
	case "endSession":
		// TODO: Close session and add info on which session is closed
		response = "Session End"
		log.Info(response)

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
		log.Debugf("---\n%+v\n---", update.Message)
		log.Debugf("---\n%+v\n---", update.Message.Chat)
		response := handleMessage(update.Message)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, sendErr := bot.Send(msg); sendErr != nil {
			log.Error(sendErr)
		}
	}
}
