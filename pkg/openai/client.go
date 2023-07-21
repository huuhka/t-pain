package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"io"
	"log"
	"net/http"
	"strings"
	"t-pain/pkg/models"
	"time"
)

type OpenAiClient struct {
	config     *Config
	httpClient *http.Client
}

func NewOpenAiClient(config *Config, opts ...OpenAiClientOption) (*OpenAiClient, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	// create http client
	httpClient := http.Client{}

	if config.AzureCredential != nil {
		accessToken, err := config.AzureCredential.GetToken(context.Background(), policy.TokenRequestOptions{
			Scopes: []string{"https://openai.azure.com/.default"},
		})
		if err != nil {
			return nil, fmt.Errorf("unable to get access token: %w", err)
		}
		transport := BearerTokenRoundTripper{
			Transport:   http.DefaultTransport,
			BearerToken: accessToken.Token,
		}
		httpClient.Transport = transport
	} else {
		transport := ApiKeyRoundTripper{
			Transport: http.DefaultTransport,
			ApiKey:    config.ApiKey,
		}
		httpClient.Transport = transport
	}

	client := OpenAiClient{
		config:     config,
		httpClient: &httpClient,
	}

	for _, opt := range opts {
		err := opt(&client)
		if err != nil {
			return nil, err
		}
	}

	return &client, nil
}

type OpenAiClientOption func(*OpenAiClient) error

func WithTimeout(timeout time.Duration) OpenAiClientOption {
	return func(c *OpenAiClient) error {
		c.httpClient.Timeout = timeout
		return nil
	}
}

func WithHttpClient(client *http.Client) OpenAiClientOption {
	return func(c *OpenAiClient) error {
		c.httpClient = client
		return nil
	}
}

type BearerTokenRoundTripper struct {
	Transport   http.RoundTripper
	BearerToken string
}

func (rt BearerTokenRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rt.BearerToken))
	return rt.Transport.RoundTrip(req)
}

type ApiKeyRoundTripper struct {
	Transport http.RoundTripper
	ApiKey    string
}

func (rt ApiKeyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("api-key", rt.ApiKey)
	return rt.Transport.RoundTrip(req)
}

func (c OpenAiClient) GetPainDescriptionObject(painDescription string) ([]models.PainDescription, error) {
	// create conversation with config.SystemContext and diff
	painDescMsg := NewUserMessage(painDescription)
	conversation := NewConversation(c.config.SystemContext)
	conversation.Messages = append(conversation.Messages, painDescMsg)

	var painDescObj []models.PainDescription

	// create request
	req, cancel, err := c.createRequest(conversation)
	defer cancel()
	if err != nil {
		return painDescObj, fmt.Errorf("unable to create request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return painDescObj, fmt.Errorf("unable to send request: %w", err)
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return painDescObj, fmt.Errorf("request failed with status code %d and body %s", resp.StatusCode, responseBody)
	}

	// parse response
	var parsedResp OpenAiCompletionResponse
	err = json.Unmarshal(responseBody, &parsedResp)
	if err != nil {
		return painDescObj, fmt.Errorf("unable to parse response: %w", err)
	}
	split := strings.Split(parsedResp.Choices[0].Message.Content, "####")
	oaiText := split[len(split)-1]
	b := []byte(oaiText)

	err = json.Unmarshal(b, &painDescObj)
	if err != nil {
		log.Println("unable to parse to PainDescObject:", err)
		return painDescObj, fmt.Errorf(oaiText)
	}
	for i := range painDescObj {
		now := time.Now()
		painDescObj[i].Timestamp = now
	}
	return painDescObj, nil
}

// createRequest creates a request for the OpenAI API
func (c OpenAiClient) createRequest(conversation *Conversation) (*http.Request, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

	body, err := c.generateRequestBody(conversation)
	if err != nil {
		return nil, cancel, fmt.Errorf("unable to generate request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.Url, bytes.NewBuffer(body))

	return req, cancel, err
}

func (c OpenAiClient) generateRequestBody(conversation *Conversation) ([]byte, error) {
	body := OpenAiCompletionRequest{
		Messages: conversation.Messages,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal request body: %w", err)
	}
	return bodyBytes, nil
}