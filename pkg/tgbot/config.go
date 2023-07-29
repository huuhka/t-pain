package tgbot

import (
	"fmt"
	"reflect"
)

type Config struct {
	botToken                 string
	speechRegion             string
	speechKey                string
	openAiKey                string
	openAiEndpoint           string
	openAiDeploymentName     string
	dataCollectionEndpoint   string
	dataCollectionRuleId     string
	dataCollectionStreamName string
}

// NewConfig creates a new Config struct that contains all the configurations required for the bot to run
func NewConfig(botToken, speechKey, speechRegion, openAiKey, openAiEndpoint, openAiDeploymentName, dataCollectionEndpoint, dataCollectionRuleId, dataCollectionStreamName string) (*Config, error) {
	c := &Config{
		botToken:                 botToken,
		speechKey:                speechKey,
		speechRegion:             speechRegion,
		openAiKey:                openAiKey,
		openAiEndpoint:           openAiEndpoint,
		openAiDeploymentName:     openAiDeploymentName,
		dataCollectionEndpoint:   dataCollectionEndpoint,
		dataCollectionRuleId:     dataCollectionRuleId,
		dataCollectionStreamName: dataCollectionStreamName,
	}

	err := checkEmptyFields(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func checkEmptyFields(c *Config) error {
	var emptyValues string

	// This won't work if the struct has a field that is not a string, but good enough for now

	v := reflect.ValueOf(*c)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() == "" {
			if emptyValues != "" && i != 0 {
				emptyValues += ", "
			}
			emptyValues += v.Type().Field(i).Name
		}
	}

	if emptyValues != "" {
		return fmt.Errorf("empty values for required config fields: %s", emptyValues)
	}

	return nil
}