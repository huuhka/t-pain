package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	b, err := tgbot.NewDefaultBot(conf)
	if err != nil {
		log.Fatalln(err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-signalChan
		log.Println("Shutting down bot...")
		b.Stop()
	}()

	b.Run()
}