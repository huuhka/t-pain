package speechtotext

import (
	"fmt"
	"os"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func sessionStartedHandler(event speech.SessionEventArgs) {
	defer event.Close()
	fmt.Println("Session Started (ID=", event.SessionID, ")")
}

func sessionStoppedHandler(event speech.SessionEventArgs) {
	defer event.Close()
	fmt.Println("Session Stopped (ID=", event.SessionID, ")")
}

func recognizingHandler(event speech.SpeechRecognitionEventArgs) {
	defer event.Close()
	fmt.Println("Recognizing:", event.Result.Text)
}

func recognizedHandler(event speech.SpeechRecognitionEventArgs) {
	defer event.Close()
	fmt.Println("Recognized:", event.Result.Text)
}

func cancelledHandler(event speech.SpeechRecognitionCanceledEventArgs) {
	defer event.Close()
	fmt.Println("Received a cancellation: ", event.ErrorDetails)
	fmt.Println("Did you set the speech resource key and region values?")
}

type Recognizer struct {
	SpeechConfig *speech.SpeechConfig
}

func NewRecognizer() (*Recognizer, error) {
	speechKey := os.Getenv("SPEECH_KEY")
	speechRegion := os.Getenv("SPEECH_REGION")
	speechConfig, err := speech.NewSpeechConfigFromSubscription(speechKey, speechRegion)
	if err != nil {
		return nil, err
	}
	return &Recognizer{
		SpeechConfig: speechConfig,
	}, err
}

func (r *Recognizer) HandleAudioLink(url string) (string, error) {
	// Download and convert
	wavFile, err := handleAudioFileSetup(url)
	defer deleteFromDisk(wavFile)

	audioConfig, err := audio.NewAudioConfigFromWavFileInput(wavFile)
	if err != nil {
		return "", fmt.Errorf("handleAudioLink: error creating newAudioConfig, %w", err)
	}
	defer audioConfig.Close()
	speechRecognizer, err := speech.NewSpeechRecognizerFromConfig(r.SpeechConfig, audioConfig)
	if err != nil {
		return "", fmt.Errorf("handleAudioLink: error creating speechRecognizer, %w", err)
	}
	defer speechRecognizer.Close()

	speechRecognizer.SessionStarted(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Started (ID=", event.SessionID, ")")
	})
	speechRecognizer.SessionStopped(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Stopped (ID=", event.SessionID, ")")
	})

	task := speechRecognizer.RecognizeOnceAsync()
	var outcome speech.SpeechRecognitionOutcome
	select {
	case outcome = <-task:
	case <-time.After(40 * time.Second):
		return "", fmt.Errorf("handleAudioLink: timed out waiting for response")
	}

	defer outcome.Close()
	if outcome.Error != nil {
		return "", fmt.Errorf("handleAudioLink: error recognizing speech, %w", outcome.Error)
	}

	return outcome.Result.Text, nil
}