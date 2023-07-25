package tgbot_test

import (
	"strings"
	"t-pain/pkg/tgbot"
	"testing"
)

func TestNewConfigShouldFailWithEmptyValues(t *testing.T) {
	_, err := tgbot.NewConfig(
		"notImportant",
		"notImportant",
		"notImportant",
		"notImportant",
		"notImportant",
		"notImportant",
		"",
		"",
		"notImportant",
	)

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if strings.Contains(err.Error(), "dataCollectionEndpoint") == false {
		t.Errorf("expected error to contain dataCollectionEndpoint, got %s", err.Error())
	}
	if len(strings.Split(err.Error(), ",")) != 3 {
		t.Errorf("expected error to contain exactly two field names, got %s", err.Error())
	}
}