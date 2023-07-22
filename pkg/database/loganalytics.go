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

type LogAnalyticsClient struct {
	client     *azingest.Client
	ruleId     string
	streamName string
}

func NewLogAnalyticsClient(endpoint, ruleId, streamName string) (*LogAnalyticsClient, error) {
	cred, err := getCredential()
	if err != nil {
		return nil, fmt.Errorf("unable to get credential: %w", err)
	}

	client, err := azingest.NewClient(endpoint, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %w", err)
	}

	return &LogAnalyticsClient{
		client:     client,
		ruleId:     ruleId,
		streamName: streamName,
	}, nil
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