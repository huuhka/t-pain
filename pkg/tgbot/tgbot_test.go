package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"t-pain/pkg/models"
	"t-pain/pkg/speechtotext"
	"testing"
	"time"
)

type MockBotAPI struct {
	mock.Mock
}

func (m *MockBotAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	args := m.Called(c)
	return args.Get(0).(tgbotapi.Message), args.Error(1)
}

func (m *MockBotAPI) GetFileDirectURL(fileID string) (string, error) {
	args := m.Called(fileID)
	return args.String(0), args.Error(1)
}

func (m *MockBotAPI) GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	args := m.Called(config)
	return args.Get(0).(tgbotapi.UpdatesChannel)
}

type MockAI struct {
	mock.Mock
}

func (m *MockAI) GetPainDescriptionObject(text string) ([]models.PainDescription, error) {
	args := m.Called(text)
	return args.Get(0).([]models.PainDescription), args.Error(1)
}

type MockLogAnalytics struct {
	mock.Mock
}

func (m *MockLogAnalytics) SavePainDescriptionsToLogAnalytics(data []models.PainDescriptionLogEntry) error {
	args := m.Called(data)
	return args.Error(0)
}

func generateTestUpdate() tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 1234,
			},
			From: &tgbotapi.User{
				UserName: "tester",
			},
		},
	}
}

func TestBotProcessMessage(t *testing.T) {
	t.Parallel()
	mockBotAPI := new(MockBotAPI)
	mockAI := new(MockAI)
	mockLogAnalytics := new(MockLogAnalytics)

	b := &Bot{
		Bot:                mockBotAPI,
		openAIClient:       mockAI,
		logAnalyticsClient: mockLogAnalytics,
		speechConfig:       speechtotext.NewConfig("key", "region"),
	}

	update := generateTestUpdate()
	update.Message.Text = "Test Message"

	mockAI.On("GetPainDescriptionObject", "Test Message").Return([]models.PainDescription{}, nil)
	mockLogAnalytics.On("SavePainDescriptionsToLogAnalytics", mock.Anything).Return(nil)
	mockBotAPI.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil)

	b.processMessage(update)

	mockAI.AssertExpectations(t)
	mockLogAnalytics.AssertExpectations(t)
	mockBotAPI.AssertExpectations(t)
}

func TestBotReply(t *testing.T) {
	t.Parallel()
	mockBotAPI := new(MockBotAPI)
	b := &Bot{Bot: mockBotAPI}

	update := generateTestUpdate()
	update.Message.Text = "Hello"

	mockBotAPI.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil)

	b.reply(update, "Test Reply")

	mockBotAPI.AssertExpectations(t)
}

func TestBotProcessToText(t *testing.T) {
	t.Parallel()
	mockBotAPI := new(MockBotAPI)
	b := &Bot{Bot: mockBotAPI, speechConfig: speechtotext.NewConfig("key", "region")}

	// Case 1: Text message
	update := generateTestUpdate()
	update.Message.Text = "Hello"
	text, err := b.processToText(update)
	assert.Nil(t, err)
	assert.Equal(t, "Hello", text)

	// Case 2: No text or voice message
	update = generateTestUpdate()
	update.Message.Video = &tgbotapi.Video{}

	_, err = b.processToText(update)
	assert.NotNil(t, err)
}

func TestBotSaveDataToLogAnalytics(t *testing.T) {
	t.Parallel()
	mockLogAnalytics := new(MockLogAnalytics)
	b := &Bot{logAnalyticsClient: mockLogAnalytics}

	painDesc := []models.PainDescription{{
		Timestamp:           time.Now(),
		LocationId:          1,
		SideId:              1,
		Level:               1,
		Description:         "Test pain",
		Numbness:            false,
		NumbnessDescription: "No numbness",
	}}
	userId := int64(1111111111111111111)

	mockLogAnalytics.On("SavePainDescriptionsToLogAnalytics", mock.Anything).Return(nil)

	err := b.saveDataToLogAnalytics(userId, painDesc)

	assert.Nil(t, err)
	mockLogAnalytics.AssertExpectations(t)
}

func TestFmtReplyShouldNotBeEmpty(t *testing.T) {
	t.Parallel()
	painDesc := []models.PainDescription{
		{
			Timestamp:           time.Now(),
			LocationId:          1,
			SideId:              1,
			Level:               1,
			Description:         "Test pain",
			Numbness:            false,
			NumbnessDescription: "No numbness",
		},
	}

	reply := fmtReply(painDesc)
	assert.NotEmpty(t, reply)
}