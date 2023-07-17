package tgbot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"t-pain/pkg/models"
	"t-pain/pkg/openai"
	"t-pain/pkg/speechtotext"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	speechConfig *speechtotext.Config
	openAIClient *openai.OpenAiClient
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
			if _, ok := models.UserIDs[update.Message.From.ID]; !ok {
				b.reply(update, "You are not authorized to use this bot")
			}
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

	// SPEECH TO TEXT
	speechKey := os.Getenv("SPEECH_KEY")
	speechRegion := os.Getenv("SPEECH_REGION")
	botObj.speechConfig = speechtotext.NewConfig(speechKey, speechRegion)

	// OPENAI, TODO make this a bit clearer
	openAIKey := os.Getenv("OPENAI_KEY")
	openAIEndpoint := os.Getenv("OPENAI_ENDPOINT")
	openAiDeployment := os.Getenv("OPENAI_DEPLOYMENT")
	config, err := openai.NewConfig(openAIEndpoint, openAiDeployment, openai.WithApiKey(openAIKey))
	if err != nil {
		return nil, err
	}
	openAIClient, err := openai.NewOpenAiClient(config)
	if err != nil {
		return nil, err
	}
	botObj.openAIClient = openAIClient

	return botObj, nil
}

func (b *Bot) processMessage(update tgbotapi.Update) {
	receivedText, err := b.processToText(update)
	if err != nil {
		log.Printf("Error processing message: %v", err)
		if err.Error() == "This bot can only handle text and voice messages" {
			b.reply(update, "This bot can only handle text and voice messages")
		} else {
			b.reply(update, "Error processing message. Please contact Pasi")
		}
		return
	}

	painDesc, err := b.openAIClient.GetPainDescriptionObject(receivedText)
	if err != nil {
		log.Printf("Error processing message: %v", err)
		b.reply(update, err.Error())
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	b.reply(update, painDesc.StringFriendly())
}

func (b *Bot) reply(update tgbotapi.Update, replyText string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	//msg.ReplyToMessageID = update.Message.MessageID
	msg.Text = replyText

	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (b *Bot) processToText(update tgbotapi.Update) (string, error) {
	var text string
	if update.Message.Voice != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Voice.FileID)
		fileLink, err := b.bot.GetFileDirectURL(update.Message.Voice.FileID)
		recognizer, err := speechtotext.NewWrapper(b.speechConfig.Key, b.speechConfig.Region)
		if err != nil {
			return "", fmt.Errorf("processToText: recognizer creation: %w", err)
		}
		text, err = speechtotext.HandleAudioLink(fileLink, recognizer)
		if err != nil {
			log.Printf("processToText: Error handling audio: %v", err)
			return "", fmt.Errorf("processToText: Error handling audio: %w", err)
		}
	} else if update.Message.Text != "" {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		text = update.Message.Text
	} else {
		return "This bot can only handle text and voice messages", fmt.Errorf("this bot can only handle text and voice messages")
	}
	return text, nil
}

func getBotToken() (string, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return "", fmt.Errorf("unable to get bot token, BOT_TOKEN env variable is empty")
	}
	return botToken, nil
}