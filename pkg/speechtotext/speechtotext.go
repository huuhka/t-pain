package speechtotext

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func HandleAudioLink(url string, wrapper *SDKWrapper) (string, error) {
	// Download and convert
	wavFile, _ := handleAudioFileSetup(url)
	defer deleteFromDisk(wavFile)

	stop := make(chan int)
	ready := make(chan struct{})
	go PumpFileContinuously(stop, wavFile, wrapper)

	var resultText []string

	log.Println("Starting continuous recognition")
	err := wrapper.StartContinuous(func(event *SDKWrapperEvent) {
		defer event.Close()
		switch event.EventType {
		case Recognized:
			log.Println("Got a recognized event")
			resultText = append(resultText, event.Recognized.Result.Text)
		case Recognizing:
		case Cancellation:
			log.Println("Got a cancellation event. Reason: ", event.Cancellation.Reason)
			close(ready)
		}
	})

	if err != nil {
		return "", err
	}

	select {
	case <-ready:
		err := wrapper.StopContinuous()
		if err != nil {
			log.Println("Error stopping continuous: ", err)
		}
	case <-time.After(30 * time.Second):
		close(stop)
		_ = wrapper.StopContinuous()
		return "", fmt.Errorf("timeout")
	}
	defer wrapper.StopContinuous()

	if len(resultText) == 0 {
		return "", fmt.Errorf("only got empty results")
	}

	return strings.Join(resultText, " "), nil
}