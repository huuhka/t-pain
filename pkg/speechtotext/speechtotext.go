package speechtotext

import (
	"fmt"
	"strings"
)

func (wrapper *SDKWrapper) HandleAudioLink(url string) (string, error) {
	// Download and convert
	wavFile, err := handleAudioFileSetup(url)
	defer deleteFromDisk(wavFile)

	stop := make(chan int)
	ready := make(chan struct{})
	go PumpFileContinuously(stop, ready, wavFile, wrapper)

	var resultText []string

	err = wrapper.StartContinuous(func(event *SDKWrapperEvent) {
		defer event.Close()
		switch event.EventType {
		case Recognized:
			resultText = append(resultText, event.Recognized.Result.Text)
		case Recognizing:
			fmt.Println("Got a recognizing event")
		case Cancellation:
			fmt.Println("Got a cancellation event")
		}
	})
	defer wrapper.StopContinuous()
	if err != nil {
		return "", err
	}
	select {
	case <-ready:
		wrapper.StopContinuous()
	}

	return strings.Join(resultText, " "), nil
}

//func RecognizeContinuousUsingWrapper(subscription string, region string, file string) {
//	/* If running this in a server, each worker thread should run something similar to this */
//	wrapper, err := NewWrapper(subscription, region)
//	if err != nil {
//		fmt.Println("Got an error: ", err)
//	}
//	defer wrapper.Close()
//	stop := make(chan int)
//	ready := make(chan struct{})
//	go PumpFileContinuously(stop, ready, file, wrapper)
//	fmt.Println("Starting Continuous...")
//	wrapper.StartContinuous(func(event *SDKWrapperEvent) {
//		defer event.Close()
//		switch event.EventType {
//		case Recognized:
//			fmt.Println("Got a recognized event")
//		case Recognizing:
//			fmt.Println("Got a recognizing event")
//		case Cancellation:
//			fmt.Println("Got a cancellation event")
//		}
//	})
//	<-time.After(10 * time.Second)
//	stop <- 1
//	fmt.Println("Stopping Continuous...")
//	wrapper.StopContinuous()
//	fmt.Println("Exiting...")
//}