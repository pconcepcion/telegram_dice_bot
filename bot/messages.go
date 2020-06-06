package bot

import (
	"fmt"
	"strings"

	storage "github.com/pconcepcion/telegram_dice_bot/storage"
	valid "github.com/pconcepcion/telegram_dice_bot/validations"
	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	rpg "github.com/pconcepcion/dice"
)

func (b *bot) handleDiceExpression(m *tgbotapi.Message, player *storage.Player) string {
	diceExpression, rollMessage, err := separateExressionAndRollMessage(m.CommandArguments())
	if err != nil {
		// TODO: handle the error gracefully
		return "Dice expression Error"
	}
	dicesResult, err := roll(diceExpression)
	// TODO: handle the result outside here
	if err != nil {
		// TODO: handle the error gracefully
		log.Error(err)
		return "Dice expression Error"
	}
	activeSession := b.ActiveSessions[m.Chat.ID]
	results := fmt.Sprintf("%v", dicesResult.GetResults())
	// TODO: get the actual raw ressults
	rawresults := fmt.Sprintf("%v", dicesResult.GetResults())
	// TODO: Use an object for this call
	// The roll should somehow be marked as sent or not, this should probably happen after the response is sent
	err = b.storage.RegisterRoll(diceExpression, results, rawresults, rollMessage, dicesResult.GetTotal(), activeSession, player)
	if err != nil {
		log.Warn(" Error registering dice expression", err)
	}

	return composeResponse(m.From, diceExpression, rollMessage, dicesResult)
}

// handleMessage handles a telegram bot API message and returns a response on a string with Markdown_v2 format
// see https://core.telegram.org/bots/api#markdownv2-style
func (b *bot) handleMessage(m *tgbotapi.Message) string {
	var response = "Unknown command"
	log.Debugf("Bot received command: %v", m)
	log.Infof("Bot command: %v, arguments %v", m.Command(), m.CommandArguments())
	activePlayer, err := b.getActivePlayer(m.From)
	if err != nil {
		log.Warn("Unable to get active player", err)
		// TODO handle this gracefully
		return "Invalid user"
	}
	log.Debug("Active Player: ", activePlayer)
	activeSession, ok := b.ActiveSessions[m.Chat.ID]
	if ok {
		err := b.storage.AddPlayerToSessoionIfMissing(activePlayer, activeSession)
		if err != nil {
			log.Error("Error addig player to session")
		}
	}
	switch m.Command() {
	// Dice expression
	case "de":
		response = b.handleDiceExpression(m, activePlayer)
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
		response = b.handleEndSession(m.Chat.ID)
		return response
	default:
		return "Unknown Command"
	}
	log.Debug("response", response)
	return response
}

// Helper functions

// separates the dice expression and the following roll message on the first space
func separateExressionAndRollMessage(arguments string) (expression, rollMessage string, err error) {
	// TODO: Add some validation
	splitted := strings.SplitAfterN(arguments, " ", 2)
	expression = splitted[0]
	validExpression, err := valid.CheckInvalidDiceExpression(expression)
	if err != nil {
		return expression, rollMessage, errors.Wrap(err, "Obtained an invalid dice Expression")
	}
	// If there was a message get it
	if len(splitted) > 1 {
		rollMessage = splitted[1]
	}
	return validExpression, rollMessage, nil
}

// composeResponse preapres the response string to be sent to Telegram
// See: https://core.telegram.org/bots/api#markdownv2-style for MarkdownV2 style
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

// roll the dice expession contained on the message
func roll(message string) (rpg.ExpressionResult, error) {
	toRoll := rpg.NewSimpleExpression(message)
	return toRoll.Roll()
}
