package main

import (
	"fmt"
	"log"
	"os"
	"t-pain/pkg/tgbot"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")

	speechKey := os.Getenv("SPEECH_KEY")
	speechRegion := os.Getenv("SPEECH_REGION")

	openAIKey := os.Getenv("OPENAI_KEY")
	openAIEndpoint := os.Getenv("OPENAI_ENDPOINT")
	openAiDeployment := os.Getenv("OPENAI_DEPLOYMENT")

	dcEndpoint := os.Getenv("DATA_COLLECTION_ENDPOINT")
	dcRuleId := os.Getenv("DATA_COLLECTION_RULE_ID")
	dcStreamName := os.Getenv("DATA_COLLECTION_STREAM_NAME")

	conf, err := tgbot.NewConfig(
		botToken,
		speechKey,
		speechRegion,
		openAIKey,
		openAIEndpoint,
		openAiDeployment,
		dcEndpoint,
		dcRuleId,
		dcStreamName,
	)
	if err != nil {
		log.Fatalln(fmt.Errorf("error creating config. Often relates to missing env variables in ALL_CAPS_SNAKE_CASE: %w", err))
	}

	err = tgbot.Run(conf)
	if err != nil {
		log.Fatalln(err)
	}
}