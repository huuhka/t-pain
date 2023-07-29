package tgbot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotAPI is an interface for tgbotapi.BotAPI to make testing easier
type BotAPI interface {
	GetFileDirectURL(fileID string) (string, error)
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
}