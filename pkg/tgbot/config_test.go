package tgbot_test

import (
	"strings"
	"t-pain/pkg/tgbot"
	"testing"
)

func TestNewConfigShouldFailWithEmptyValues2(t *testing.T) {
	tests := []struct {
		name                                                                                                                                                       string
		botToken, speechKey, speechRegion, openAiKey, openAiEndpoint, openAiDeploymentName, dataCollectionEndpoint, dataCollectionRuleId, dataCollectionStreamName string
		expectedErrorContains                                                                                                                                      string
		expectedErrorFieldNumber                                                                                                                                   int
	}{
		// Well, this got a bit out of hand. It's hideous!
		{"missing botToken", "", "x", "x", "x", "x", "x", "x", "x", "x", "botToken", 1},
		{"missing speechKey", "x", "", "x", "x", "x", "x", "x", "x", "x", "speechKey", 1},
		{"missing speechRegion", "x", "x", "", "x", "x", "x", "x", "x", "x", "speechRegion", 1},
		{"missing openAiKey", "x", "x", "x", "", "x", "x", "x", "x", "x", "openAiKey", 1},
		{"missing openAiEndpoint", "x", "x", "x", "x", "", "x", "x", "x", "x", "openAiEndpoint", 1},
		{"missing openAiDeploymentName", "x", "x", "x", "x", "x", "", "x", "x", "x", "openAiDeploymentName", 1},
		{"missing dataCollectionEndpoint", "x", "x", "x", "x", "x", "x", "", "x", "x", "dataCollectionEndpoint", 1},
		{"missing dataCollectionRuleId", "x", "x", "x", "x", "x", "x", "x", "", "x", "dataCollectionRuleId", 1},
		{"missing dataCollectionStreamName", "x", "x", "x", "x", "x", "x", "x", "x", "", "dataCollectionStreamName", 1},
		{"missing all", "", "", "", "", "", "", "", "", "", "speechKey", 9},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tgbot.NewConfig(tt.botToken, tt.speechKey, tt.speechRegion, tt.openAiKey, tt.openAiEndpoint, tt.openAiDeploymentName, tt.dataCollectionEndpoint, tt.dataCollectionRuleId, tt.dataCollectionStreamName)

			if err == nil {
				t.Errorf("expected error, got nil")
			}

			if strings.Contains(err.Error(), tt.expectedErrorContains) == false {
				t.Errorf("expected error to contain %s, got %s", tt.expectedErrorContains, err.Error())
			}
			if len(strings.Split(err.Error(), ",")) != tt.expectedErrorFieldNumber {
				t.Errorf("expected error to contain exactly %d field names, got %s", tt.expectedErrorFieldNumber, err.Error())
			}
		})
	}
}