package speechtotext

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Wrapper interface {
	StartContinuous(handler func(event *SDKWrapperEvent)) error
	StopContinuous() error
	Writer
}

func HandleAudioLink(url string, wrapper Wrapper) (string, error) {
	// Download and convert
	wavFile, err := handleAudioFileSetup(url)
	if err != nil {
		return "", err
	}
	defer deleteFromDisk(wavFile)

	stop := make(chan int)
	ready := make(chan struct{})
	go PumpFileToStream(stop, wavFile, wrapper)

	var resultText []string

	log.Println("Starting continuous recognition")
	err = wrapper.StartContinuous(func(event *SDKWrapperEvent) {
		defer event.Close()
		switch event.EventType {
		case Recognized:
			log.Println("Got a recognized event")
			resultText = append(resultText, event.Recognized.Result.Text)
		case Recognizing:
		case Cancellation:
			log.Println("Got a cancellation event. Reason: ", event.Cancellation.Reason)
			close(ready)
			if event.Cancellation.Reason.String() == "Error" {
				log.Println("ErrorCode:" + event.Cancellation.ErrorCode.String() + " ErrorDetails: " + event.Cancellation.ErrorDetails)
			}
			// TODO: If we receive an error here, the writing to the stream should be stopped. Currently that does not seem to happen.
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
	case <-time.After(120 * time.Second):
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