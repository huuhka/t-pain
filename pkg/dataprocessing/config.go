package dataprocessing

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// CreateUrl creates a url for the request
func CreateUrl(endpoint, deploymentName string) string {
	apiVersion := "2023-03-15-preview"
	return fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s", endpoint, deploymentName, apiVersion)
}

// Config is the configuration for calling the Azure OpenAI API
type Config struct {
	Url             string
	ApiKey          string
	SystemContext   Message
	AzureCredential *azidentity.DefaultAzureCredential
}

func NewConfig(endpoint, deploymentName string, opts ...ConfigOpt) (*Config, error) {
	systemRole := "Assistant is an AI chatbot that helps users turn a natural language description of their pain levels into JSON format. " +
		"After users inputs a description of their pain levels, location of the pain, optional numbness description and further description of their feelings, it will provide a json representation in the format preceded by ####.\n" +
		"- If the user does not give a direct 0-10 number for their pain level, the assistant makes an estimate of the level on that range based on the given description.\n" +
		"- The location field should be an integer mapping to the following chart delmited by ```. If no body part is mentioned directly, try to map the pain from the description to the closest body part in the mapping.\n" +
		"- The assistant's response should always include \"####\" before the JSON representation\n" +
		"neck: 1\nshoulder: 2\nupper back: 3\nlower back: 4\nleft arm: 5\nright arm: 6\nleft leg: 7\nright leg: 8\n´´´\n" +
		"JSON Format:\n{\ntimeStamp: \nlevel: int\nlocation: int\ndescription: string\nnumbness: bool\nnumbnessDescription: string\n}"

	c := Config{
		Url:           CreateUrl(endpoint, deploymentName),
		SystemContext: NewSystemMessage(systemRole),
	}

	for _, opt := range opts {
		err := opt(&c)
		if err != nil {
			return nil, err
		}
	}

	if c.AzureCredential == nil && c.ApiKey == "" {
		return nil, fmt.Errorf("no authentication method provided, please provide an API key or Azure credential")
	}

	return &c, nil
}

type ConfigOpt func(*Config) error

func WithApiKey(apiKey string) ConfigOpt {
	return func(c *Config) error {
		c.ApiKey = apiKey
		return nil
	}
}

func WithAzureCredential() ConfigOpt {
	return func(c *Config) error {
		cred, err := LoginWithDefaultCredential()
		if err != nil {
			return err
		}

		c.AzureCredential = cred
		return nil
	}
}