package openai_test

import (
	"bytes"
	"io"
	"net/http"
	"t-pain/pkg/openai"
	"testing"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestNewClient(t *testing.T) {
	config := &openai.Config{
		AzureCredential: nil,
		ApiKey:          "test-api-key",
		Url:             "test-url",
		SystemContext:   *openai.NewConversation(openai.NewSystemMessage("test"), openai.NewUserMessage("test")),
	}

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{}, nil
		},
	}

	client, err := openai.NewClient(config, openai.WithDoer(mockClient))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if client == nil {
		t.Error("Expected client not to be nil")
	}
}

func TestClient_GetPainDescriptionObject_ShouldSucceedOnNormalResponse(t *testing.T) {
	t.Parallel()
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("{\"choices\":[{\"message\":{\"content\":\"{\n\"level\": 7,\n\"locationId\": 1,\n\"sideId\": 2,\n\"description\": \"Pain in Head\",\n\"numbness\": true,\n\"numbnessDescription\": \"Numbness in Head\"},\"}}]}")), // sample response
			}, nil
		},
	}

	config := &openai.Config{
		AzureCredential: nil,
		ApiKey:          "test-api-key",
		Url:             "test-url",
		SystemContext:   *openai.NewConversation(openai.NewSystemMessage("test"), openai.NewUserMessage("test")),
	}

	client, _ := openai.NewClient(config, openai.WithDoer(mockClient))

	_, err := client.GetPainDescriptionObject("test pain description")
	if err == nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClient_GetPainDescriptionObject_ShouldErrorOnFailedOutputFromModel(t *testing.T) {
	t.Parallel()
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("{\"choices\":[{\"message\":{\"content\":\"####\"}}]}")), // sample response
			}, nil
		},
	}

	config := &openai.Config{
		AzureCredential: nil,
		ApiKey:          "test-api-key",
		Url:             "test-url",
		SystemContext:   *openai.NewConversation(openai.NewSystemMessage("test"), openai.NewUserMessage("test")),
	}

	client, _ := openai.NewClient(config, openai.WithDoer(mockClient))

	_, err := client.GetPainDescriptionObject("test pain description")
	if err == nil {
		t.Errorf("Expected failing parse of painDescription error, got nil")
	}
}