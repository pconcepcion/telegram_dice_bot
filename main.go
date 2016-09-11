package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/pconcepcion/dice"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const (
	apiToken = "YOUR API TOKEN HERE"
)

var log = logrus.New()

func init() {
	log.Formatter = new(prefixed.TextFormatter)
	log.Level = logrus.DebugLevel

}

// roll the dice expession contained on the message
func roll(message string) (rpg.DiceExpressionResult, error) {
	toRoll := rpg.NewSimpleDiceExpression(message)
	return toRoll.Roll()
}

// Do the authorization of the bot and set the bot mode
func authorizeBot(debug bool) *tgbotapi.BotAPI {
	log.Printf("Authorizing bog with token %s", apiToken)
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Panic(errors.Wrap(err, "Couldn't acces the API wiht token "+apiToken))
	}

	bot.Debug = debug

	log.Printf("Authorized on account %s", bot.Self.UserName)
	log.Printf("debug = %t", debug)

	return bot

}

func main() {
	// Athorize the bot with deboug mode
	bot := authorizeBot(true)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	//dice := rpg.NewDice(6)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		diceExpression := update.Message.Text
		//response := fmt.Sprintf("d6: %d", dice.Roll())
		dicesResult, err := roll(diceExpression)
		if err != nil {
			// TODO: handle the error gracefully
			//log.Panic(err)
			log.Error(err)
		}
		response := fmt.Sprintf("%s: %v -> %s", diceExpression, dicesResult.GetResults(), dicesResult)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
