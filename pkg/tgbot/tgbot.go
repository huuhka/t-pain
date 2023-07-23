package tgbot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"t-pain/pkg/database"
	"t-pain/pkg/models"
	"t-pain/pkg/openai"
	"t-pain/pkg/speechtotext"
	"time"
	_ "time/tzdata"
)

type Bot struct {
	bot                *tgbotapi.BotAPI
	speechConfig       *speechtotext.Config
	openAIClient       *openai.OpenAiClient
	logAnalyticsClient *database.LogAnalyticsClient
}

func Run(c *Config) error {
	b, err := NewBot(c)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if _, ok := models.UserIDs[update.Message.From.ID]; !ok {
				log.Printf("Unauthorized user tried to use the bot: %v", update.Message.From)
				b.reply(update, "You are not authorized to use this bot")
				continue
			}
			go b.processMessage(update)
		}
	}

	return err
}

func NewBot(c *Config) (*Bot, error) {

	botObj := &Bot{}

	bot, err := tgbotapi.NewBotAPI(c.BotToken)
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
	botObj.speechConfig = speechtotext.NewConfig(c.SpeechKey, c.SpeechRegion)

	// OPENAI
	oaiConf, err := openai.NewConfig(c.openAiEndpoint, c.openAiDeploymentName, openai.WithApiKey(c.openAiKey))
	if err != nil {
		return nil, err
	}
	openAIClient, err := openai.NewOpenAiClient(oaiConf)
	if err != nil {
		return nil, err
	}
	botObj.openAIClient = openAIClient

	// DATA SAVING
	dcClient, err := database.NewLogAnalyticsClient(c.dataCollectionEndpoint, c.dataCollectionRuleId, c.dataCollectionStreamName)
	if err != nil {
		return nil, err
	}
	botObj.logAnalyticsClient = dcClient

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

	err = b.saveDataToLogAnalytics(update.Message.From.ID, painDesc)
	if err != nil {
		log.Printf("Error saving data to log analytics: %v", err)
		b.reply(update, "Error saving data. Please contact Pasi and try again later.")
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	b.reply(update, fmtReply(painDesc))
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

func (b *Bot) saveDataToLogAnalytics(userId int64, pd []models.PainDescription) error {
	var data []models.PainDescriptionLogEntry
	for _, pain := range pd {
		data = append(data, pain.MapToLogEntry(userId))
	}
	err := b.logAnalyticsClient.SavePainDescriptionsToLogAnalytics(data)
	if err != nil {
		return fmt.Errorf("saveDataToLogAnalytics: %w", err)
	}
	return nil
}

func fmtReply(pd []models.PainDescription) string {
	var result strings.Builder
	if len(pd) == 0 {
		return ""
	}

	first := pd[0]

	loc, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		log.Printf("Error loading location: %v", err)
		return "Error generating reply from description. Your data has been saved. Please contact Pasi"
	}
	tstamp := first.Timestamp.Round(time.Minute).In(loc).Format("02-01-2006 15:04")

	result.WriteString(fmt.Sprintf("Timestamp: %s\n", tstamp))
	result.WriteString("Pains:\n")
	for _, pain := range pd {
		result.WriteString(fmt.Sprintf("\t- Location: %s, Side: %s, Level: %d\n", models.BodyPartMapping[pain.LocationId], models.SideMap[pain.SideId], pain.Level))
	}
	result.WriteString(fmt.Sprintf("Description: %s\n", first.Description))
	result.WriteString(fmt.Sprintf("Numbness: %t\n", first.Numbness))
	result.WriteString(fmt.Sprintf("Numbness Description: %s\n", first.NumbnessDescription))
	return result.String()
}