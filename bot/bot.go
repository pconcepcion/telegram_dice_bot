package bot

import (
	"fmt"
	"os"
  "strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	rpg "github.com/pconcepcion/dice"
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

// roll the dice expession contained on the message
func roll(message string) (rpg.ExpressionResult, error) {
	toRoll := rpg.NewSimpleExpression(message)
	return toRoll.Roll()
  }

// handleMessage handles a telegram bot API message and returns a response on a string with Markdown_v2 format
// see https://core.telegram.org/bots/api#markdownv2-style
func (b *bot) handleMessage(m *tgbotapi.Message) string {
	// TODO: Handle extra argument for the roll identifier: example "/d100 hide"
	var response = "Unknown command"
  log.Debugf("Bot received command: %v", m)
	log.Infof("Bot command: %v, arguments %v", m.Command(), m.CommandArguments())
	switch m.Command() {
	// Dice expression
	case "de":
		diceExpression, rollMessage, err := separateExressionAndRollMessage(m.CommandArguments())
    if err != nil {
			// TODO: handle the error gracefully
			return "Dice expression Error"
    }
		dicesResult, err := roll(diceExpression)
		if err != nil {
			// TODO: handle the error gracefully
			log.Error(err)
			 return "Dice expression Error"
		}
		response = composeResponse(m.From, diceExpression, rollMessage, dicesResult)
	// Basic dices
	case "d2":
		response = fmt.Sprintf("d2 "+responseArrowTemplate, rpg.D2())
	case "d4":
		response = fmt.Sprintf("d4 "+responseArrowTemplate, rpg.D4())
	case "d6":
		response = fmt.Sprintf("d6 "+responseArrowTemplate, rpg.D6())
	case "d8":
		response = fmt.Sprintf("d8 "+responseArrowTemplate, rpg.D8())
	case "d10":
		response = fmt.Sprintf("d10 "+responseArrowTemplate, rpg.D10())
	case "d12":
		response = fmt.Sprintf("d12 "+responseArrowTemplate, rpg.D12())
	case "d20":
		response = fmt.Sprintf("d20 "+responseArrowTemplate, rpg.D20())
	case "d100":
		response = fmt.Sprintf("d100 "+responseArrowTemplate, rpg.D100())
	// Session Handling
	case "start_session":
		// Store Session info
		sessionName := m.CommandArguments()
		response = b.handleStartSession(m.Chat.ID, sessionName)
		// TODO: Set session timeout
	case "rename_session":
		// Updte Session info
		sessionName := m.CommandArguments()
		response = b.handleRenameSession(m.Chat.ID, sessionName)
	case "end_session":
		// TODO: Generate session Summary
		response = b.handleEndSession(m.Chat.ID)
		return response
	default:
		return "Unknown Command"
	}
	log.Debug("response", response)
	return response
}

func (b *bot) handleStartSession(chatID int64, sessionName string) string {
	var response string
	session, err := b.storage.StartSession(sessionName, chatID)
		if err != nil {
			response = fmt.Sprintf("Failed to create Session, invalid session name")
		log.Errorf("Failed to create Session, invalid session arguments: %s", sessionName)
		return response
	}
	b.ActiveSessions[chatID] = session
	log.Infof("Starting Session %v\n", session)
	response = fmt.Sprintf("Starting Session: \n \U0001F3F7  *_\\#%s_*", sessionName)
	return response
}

func (b *bot) handleRenameSession(chatID int64, name string) string {
	activeSession := b.ActiveSessions[chatID]
	if activeSession == nil {
		response := fmt.Sprintf("Failed to rename Session, no active session found")
		log.Errorf("Failed to rename Session session not found: %d", chatID)
		return response
	}

	oldName := activeSession.Name
	err := b.storage.RenameSession(activeSession, name)
	if err != nil {
		response := fmt.Sprintf("Failed to rename Session _\\#%s_, invalid session name", oldName)
		log.Errorf("Failed to rename Session _\\#%s_, invalid session arguments: %s", oldName, name)
		return response
	}
	return fmt.Sprintf("Session renamed: \n _\\#%s_ \U0001F3F7  *_\\#%s_*", oldName, name)
}

func (b *bot) handleEndSession(chatID int64) string {
	activeSession := b.ActiveSessions[chatID]
	if activeSession == nil {
		response := fmt.Sprintf("Failed to rename Session, no active session found")
		log.Errorf("Failed to rename Session session not found: %d", chatID)
			return response
		}
    log.Info("Closing Session: ", activeSession)
		b.storage.EndSession(activeSession)
	b.ActiveSessions[chatID] = nil
	response := fmt.Sprintf("\U0001F51A Session _\\#%s_ Finished", activeSession.Name)
		log.Info(response)
    return response
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

// Helper functions

// separates the dice expression and the following roll message on the first space
func separateExressionAndRollMessage(arguments string) (expression, rollMessage string, err error) {
	// TODO: Add some validation
	splitted := strings.SplitAfterN(arguments, " ", 2)
	expression = splitted[0]
	// If there was a message get it
	if len(splitted) > 1 {
		rollMessage = splitted[1]
	}
	return expression, rollMessage, nil
}

func composeResponse(user *tgbotapi.User, diceExpression, rollMessage string, result rpg.ExpressionResult) string {
	var message string
	if rollMessage != "" {
		message = fmt.Sprintf("*[@%s](tg://user?id=%d)* rolled *%s* and got _%v_ \n *_%s_* "+rigthArrow+" _%s_",
		user.UserName, user.ID, diceExpression, result.GetResults(), result, rollMessage)
	} else {
		message = fmt.Sprintf("*[@%s](tg://user?id=%d)* rolled *%s* and got _%v_ \n *_%s_* "+rigthArrow+" _%s_ ",
			user.UserName, user.ID, diceExpression, result.GetResults(), result, "unspecified roll")
	}
	log.Info("MESSAGE", message)
	return message

}
