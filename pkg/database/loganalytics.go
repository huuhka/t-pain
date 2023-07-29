package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azingest"
	"os"
	"t-pain/pkg/models"
)

type AzureClient interface {
	Upload(ctx context.Context, ruleID string, streamName string, logs []byte, options *azingest.UploadOptions) (azingest.UploadResponse, error)
}

type LogAnalyticsClient struct {
	client     AzureClient
	ruleId     string
	streamName string
}

func NewLogAnalyticsClient(endpoint, ruleId, streamName string, opts ...LogAnalyticsClientOption) (*LogAnalyticsClient, error) {
	options := &LogAnalyticsClientOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var cred azcore.TokenCredential
	var err error
	if options.CustomCredential != nil {
		cred = options.CustomCredential
	} else {
		cred, err = getCredential()
		if err != nil {
			return nil, fmt.Errorf("unable to get credential: %w", err)
		}
	}

	var client AzureClient
	if options.CustomClient != nil {
		client = options.CustomClient
	} else {
		azClient, err := azingest.NewClient(endpoint, cred, nil)
		if err != nil {
			return nil, fmt.Errorf("unable to create client: %w", err)
		}
		client = azClient
	}

	return &LogAnalyticsClient{
		client:     client,
		ruleId:     ruleId,
		streamName: streamName,
	}, nil
}

type LogAnalyticsClientOptions struct {
	CustomCredential azcore.TokenCredential
	CustomClient     AzureClient
}

type LogAnalyticsClientOption func(*LogAnalyticsClientOptions)

func WithCustomCredential(cred azcore.TokenCredential) LogAnalyticsClientOption {
	return func(options *LogAnalyticsClientOptions) {
		options.CustomCredential = cred
	}
}

func WithCustomClient(client AzureClient) LogAnalyticsClientOption {
	return func(options *LogAnalyticsClientOptions) {
		options.CustomClient = client
	}
}

func getCredential() (azcore.TokenCredential, error) {
	var cred azcore.TokenCredential
	var err error

	userAssignedId := os.Getenv("AZURE_CLIENT_ID")
	if userAssignedId != "" {
		cred, err = azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
			ID: azidentity.ClientID(userAssignedId),
		})
		if err != nil {
			return nil, fmt.Errorf("unable to get managed identity credential: %w", err)
		}
	} else {
		cred, err = azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil, fmt.Errorf("unable to get default credential: %w", err)
		}
	}
	return cred, nil

}

func (lac *LogAnalyticsClient) SavePainDescriptionsToLogAnalytics(pd []models.PainDescriptionLogEntry) error {

	logs, err := json.Marshal(pd)
	if err != nil {
		return fmt.Errorf("unable to marshal pain descriptions: %w", err)
	}

	_, err = lac.client.Upload(context.Background(), lac.ruleId, lac.streamName, logs, nil)
	if err != nil {
		return fmt.Errorf("unable to upload logs: %w", err)
	}

	return nil
}