package openai

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"t-pain/pkg/models"
)

// CreateUrl creates an url for the request
func CreateUrl(endpoint, deploymentName string) string {
	//if last character is /, remove it
	if endpoint[len(endpoint)-1:] == "/" {
		endpoint = endpoint[:len(endpoint)-1]
	}

	apiVersion := "2023-03-15-preview"
	return fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s", endpoint, deploymentName, apiVersion)
}

// Config is the configuration for calling the Azure OpenAI API
type Config struct {
	Url             string
	ApiKey          string
	SystemContext   Conversation
	AzureCredential *azidentity.DefaultAzureCredential
}

func NewConfig(endpoint, deploymentName string, opts ...ConfigOpt) (*Config, error) {
	systemRole := "Assistant is an AI chatbot that helps users turn a natural language description of their pain levels into valid JSON format. " +
		"After users inputs a description of their pain levels, location of the pain, optional numbness description and further description of their feelings, it will provide a json array representation in the format preceded by ####.\n" +
		"- Ignore any references to previous messages by the user. The pain description you return should only contain items from the latest message from the user. \n" +
		"- If the user does not give a direct 0-10 number for their pain level, the assistant makes an estimate of the level on that range based on the given description. The numbers should always be full integers rounded up.\n" +
		"- The location and side fields should be an integer mapping to the following chart= delmited by ```. If no body part is mentioned directly, try to map the pain from the description to the closest body part in the mapping. If no side is mentioned, set the value to both.\n" +
		"- The assistant's response should always include \"####\" before the JSON array representation and SHOULD NOT include any other text\n" +
		"- The description might have multiple pain areas, they should be considered in the json by adding a new object in the result array. However, if both pain areas would map to the same location, only add that location once." +
		"- Do not mention anything about the JSON format to the user" +
		"- If the user writes in Finnish, respond to them in Finnish. Do not modify the names of the properties in the JSON object in any situation" +
		"- If the user mentions pain radiating to other locations, add those locations to the response along with respective pain levels. Include full description on all entries.\n" +
		"Body Parts:\n" +
		fmt.Sprintf("%v", models.BodyPartMapping.StringNameFirst()) +
		"Sides:\n" +
		fmt.Sprintf("%v", models.SideMap.StringNameFirst()) +
		"JSON Format:\n" +
		"[\n" +
		fmt.Sprintf(models.PrintPainDescriptionJSONFormat()) +
		"]\n"

	examples := []Message{
		{
			Content: "My left arm is quite painful today. About level 5. The pain is radiating to my left shoulder also",
			Role:    "user",
		},
		{
			Content: "####\n" +
				"[\n" +
				models.PrintSinglePainDescriptionJSONFormat(models.PainDescription{
					Level:               5,
					LocationId:          4,
					SideId:              2,
					Description:         "My left arm is quite painful today. About level 5. The pain is radiating to my left shoulder also",
					Numbness:            false,
					NumbnessDescription: "",
				}) + ",\n" +
				models.PrintSinglePainDescriptionJSONFormat(models.PainDescription{
					Level:               5,
					LocationId:          3,
					SideId:              2,
					Description:         "My left arm is quite painful today. About level 5. The pain is radiating to my left shoulder also",
					Numbness:            false,
					NumbnessDescription: "",
				}) + "\n" +
				"]\n",
			Role: "assistant",
		},
		{
			Content: "Pitkät selkälihakset vähän krampissa. Alaselkä aika perustasoa lääkkeiden oton jälkeen. Tippu ehkä seiska puolikkaasta kutoseen. Istuessa.",
			Role:    "user",
		},
		{
			Content: "####\n" +
				"[\n" +
				models.PrintSinglePainDescriptionJSONFormat(models.PainDescription{
					Level:               7,
					LocationId:          8,
					SideId:              1,
					Description:         "Pitkät selkälihakset vähän krampissa. Alaselkä aika perustasoa lääkkeiden oton jälkeen. Tippu ehkä seiska puolikkaasta kutoseen. Istuessa.",
					Numbness:            false,
					NumbnessDescription: "",
				}) + ",\n" +
				models.PrintSinglePainDescriptionJSONFormat(models.PainDescription{
					Level:               7,
					LocationId:          9,
					SideId:              1,
					Description:         "Pitkät selkälihakset vähän krampissa. Alaselkä aika perustasoa lääkkeiden oton jälkeen. Tippu ehkä seiska puolikkaasta kutoseen. Istuessa.",
					Numbness:            false,
					NumbnessDescription: "",
				}) + "\n" +
				"]\n",
			Role: "assistant",
		},
		{
			Content: "Lisäyksenä edelliseen, myös oikea käsi on kipeä. Taso 3.",
			Role:    "user",
		},
		{
			Content: "####\n" +
				"[\n" +
				models.PrintSinglePainDescriptionJSONFormat(models.PainDescription{
					Level:               3,
					LocationId:          4,
					SideId:              3,
					Description:         "Lisäyksenä edelliseen, myös oikea käsi on kipeä ja turtunut. Taso 3. ",
					Numbness:            true,
					NumbnessDescription: "oikea käsi kipeä ja turtunut",
				}) + "\n" +
				"]\n",
			Role: "assistant",
		},
	}

	sc := NewConversation(NewSystemMessage(systemRole), examples...)

	c := Config{
		Url:           CreateUrl(endpoint, deploymentName),
		SystemContext: *sc,
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
	// TODO: This should be tested with a real environment
	return func(c *Config) error {
		cred, err := LoginWithDefaultCredential()
		if err != nil {
			return err
		}

		c.AzureCredential = cred
		return nil
	}
}