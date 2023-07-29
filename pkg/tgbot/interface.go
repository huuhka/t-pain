package tgbotapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
)

type BotAPI interface {
	MakeRequest(endpoint string, params tgbotapi.Params) (*tgbotapi.APIResponse, error)
	DecodeAPIResponse(responseBody io.Reader, resp *tgbotapi.APIResponse) ([]byte, error)
	UploadFiles(endpoint string, params tgbotapi.Params, files []tgbotapi.RequestFile) (*tgbotapi.APIResponse, error)
	GetFileDirectURL(fileID string) (string, error)
	GetMe() (tgbotapi.User, error)
	IsMessageToMe(message tgbotapi.Message) bool
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	SendMediaGroup(config tgbotapi.MediaGroupConfig) ([]tgbotapi.Message, error)
	GetUserProfilePhotos(config tgbotapi.UserProfilePhotosConfig) (tgbotapi.UserProfilePhotos, error)
	GetFile(config tgbotapi.FileConfig) (tgbotapi.File, error)
	GetUpdates(config tgbotapi.UpdateConfig) ([]tgbotapi.Update, error)
	GetWebhookInfo() (tgbotapi.WebhookInfo, error)
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	StopReceivingUpdates()
	ListenForWebhook(pattern string) tgbotapi.UpdatesChannel
	ListenForWebhookRespReqFormat(w http.ResponseWriter, r *http.Request) tgbotapi.UpdatesChannel
	HandleUpdate(r *http.Request) (*tgbotapi.Update, error)
	GetChat(config tgbotapi.ChatInfoConfig) (tgbotapi.Chat, error)
	GetChatAdministrators(config tgbotapi.ChatAdministratorsConfig) ([]tgbotapi.ChatMember, error)
	GetChatMembersCount(config tgbotapi.ChatMemberCountConfig) (int, error)
	GetChatMember(config tgbotapi.GetChatMemberConfig) (tgbotapi.ChatMember, error)
	GetGameHighScores(config tgbotapi.GetGameHighScoresConfig) ([]tgbotapi.GameHighScore, error)
	GetInviteLink(config tgbotapi.ChatInviteLinkConfig) (string, error)
	GetStickerSet(config tgbotapi.GetStickerSetConfig) (tgbotapi.StickerSet, error)
	StopPoll(config tgbotapi.StopPollConfig) (tgbotapi.Poll, error)
	GetMyCommands() ([]tgbotapi.BotCommand, error)
	GetMyCommandsWithConfig(config tgbotapi.GetMyCommandsConfig) ([]tgbotapi.BotCommand, error)
	CopyMessage(config tgbotapi.CopyMessageConfig) (tgbotapi.MessageID, error)
}