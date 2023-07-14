package speechtotext

import (
	"fmt"
	"os"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

type Recognizer struct {
	speechConfig *speech.SpeechConfig
	languages    []string
}

func NewRecognizer() (*Recognizer, error) {
	speechKey := os.Getenv("SPEECH_KEY")
	speechRegion := os.Getenv("SPEECH_REGION")

	speechConfig, err := speech.NewSpeechConfigFromSubscription(speechKey, speechRegion)
	if err != nil {
		return nil, err
	}
	// speechConfig.EnableAudioLogging()

	return &Recognizer{
		speechConfig: speechConfig,
		languages: []string{
			"en-US",
			"en-GB",
			"fi-FI",
		},
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

	autodetect, err := speech.NewAutoDetectSourceLanguageConfigFromLanguages(r.languages)
	if err != nil {
		return "", err
	}

	speechRecognizer, err := speech.NewSpeechRecognizerFomAutoDetectSourceLangConfig(r.speechConfig, autodetect, audioConfig)
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