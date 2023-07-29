package database_test

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azingest"
	"t-pain/pkg/database"
	"t-pain/pkg/models"
	"testing"
)

type MockAzureClient struct {
	UploadFunc func(ctx context.Context, ruleId string, streamName string, logs []byte, opts *azingest.UploadOptions) (azingest.UploadResponse, error)
}

func (mac *MockAzureClient) Upload(ctx context.Context, ruleId string, streamName string, logs []byte, opts *azingest.UploadOptions) (azingest.UploadResponse, error) {
	return mac.UploadFunc(ctx, ruleId, streamName, logs, opts)
}

func TestSavePainDescriptionsToLogAnalytics(t *testing.T) {
	testCases := map[string]struct {
		painDesc  models.PainDescription
		userId    int64
		uploadErr error
		wantErr   bool
	}{
		"Upload succeeds": {
			painDesc: models.PainDescription{Level: 3, LocationId: 1, SideId: 2, Description: "Test"},
			userId:   1111111111111111111,
			wantErr:  false,
		},
		"Upload fails": {
			painDesc:  models.PainDescription{Level: 3, LocationId: 1, SideId: 2, Description: "Test"},
			userId:    1111111111111111111,
			uploadErr: errors.New("upload error"),
			wantErr:   true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockClient := &MockAzureClient{}
			logAnalyticsClient, err := database.NewLogAnalyticsClient("endpoint", "testRule", "testStream", database.WithCustomClient(mockClient))
			if err != nil {
				t.Fatalf("error creating client, got %v", err)
			}
			mockClient.UploadFunc = func(ctx context.Context, ruleId string, streamName string, logs []byte, opts *azingest.UploadOptions) (azingest.UploadResponse, error) {
				return azingest.UploadResponse{}, tc.uploadErr
			}

			pdLog, err := tc.painDesc.MapToLogEntry(tc.userId)
			if err != nil {
				t.Errorf("error mapping to log entry, got %v", err)
			}

			err = logAnalyticsClient.SavePainDescriptionsToLogAnalytics([]models.PainDescriptionLogEntry{pdLog})

			if (err != nil) != tc.wantErr {
				t.Errorf("SavePainDescriptionsToLogAnalytics() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}