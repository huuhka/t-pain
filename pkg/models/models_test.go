package models_test

import (
	"t-pain/pkg/models"
	"testing"
	"time"
)

func TestMapToLogEntryWithInvalidDataShouldReturnError(t *testing.T) {
	testCases := map[string]struct {
		painDesc models.PainDescription
		userId   int64
	}{
		"Non-existent location ID": {
			painDesc: models.PainDescription{
				Timestamp:           time.Now(),
				Level:               1,
				LocationId:          999, // Non-existent location ID
				SideId:              1,
				Description:         "Pain description",
				Numbness:            true,
				NumbnessDescription: "Numbness description",
			},
			userId: 1,
		},
		"Non-existent side ID": {
			painDesc: models.PainDescription{
				Timestamp:           time.Now(),
				Level:               1,
				LocationId:          1,
				SideId:              999, // Non-existent side ID
				Description:         "Pain description",
				Numbness:            true,
				NumbnessDescription: "Numbness description",
			},
			userId: 1,
		},
		"Non-existent user ID": {
			painDesc: models.PainDescription{
				Timestamp:           time.Now(),
				Level:               1,
				LocationId:          1,
				SideId:              1,
				Description:         "Pain description",
				Numbness:            true,
				NumbnessDescription: "Numbness description",
			},
			userId: 9999999999999, // Non-existent user ID
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := tc.painDesc.MapToLogEntry(tc.userId)

			if err == nil {
				t.Error("expected an error for non-existent IDs, got none")
			}
		})
	}
}