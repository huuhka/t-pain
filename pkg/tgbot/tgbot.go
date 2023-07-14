package tgbot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"t-pain/pkg/speechtotext"
)

type Bot struct {
	bot        *tgbotapi.BotAPI
	recognizer *speechtotext.Recognizer
}

func Run() error {
	botToken, err := getBotToken()
	if err != nil {
		return err
	}

	b, err := NewBot(botToken)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			go b.processMessage(update)
		}
	}

	return err
}

func NewBot(botToken string) (*Bot, error) {

	botObj := &Bot{}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	//env := os.Getenv("ENV")
	//if env == "Development" {
	//	bot.Debug = true
	//}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	botObj.bot = bot

	recognizer, err := speechtotext.NewRecognizer()
	if err != nil {
		return nil, err
	}

	botObj.recognizer = recognizer

	return botObj, nil
}

func (b *Bot) processMessage(update tgbotapi.Update) {
	messageText := ""
	msgFormat := ""
	if update.Message.VideoNote != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.VideoNote.FileID)
		msgFormat = "video"
	} else if update.Message.Voice != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Voice.FileID)
		fileLink, err := b.bot.GetFileDirectURL(update.Message.Voice.FileID)
		msgFormat, err = b.recognizer.HandleAudioLink(fileLink)
		if err != nil {
			log.Printf("Error handling audio link: %v", err)
		}
	} else if update.Message.Text != "" {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msgFormat = "text"
	} else {
		messageText = fmt.Sprintf("This bot can only handle text, voice and videoNote messages.")
	}

	if messageText == "" {
		messageText = fmt.Sprintf("You sent me a %s", msgFormat)
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.Text = messageText

	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func getBotToken() (string, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return "", fmt.Errorf("unable to get bot token, BOT_TOKEN env variable is empty")
	}
	return botToken, nil
}