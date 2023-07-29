package speechtotext_test

import (
	"errors"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"net/http"
	"net/http/httptest"
	"os"
	"t-pain/pkg/speechtotext"
	"testing"
	"time"
)

func TestHandleAudioLink_Success(t *testing.T) {
	t.Parallel()
	mockWrapper := &MockSDKWrapper{
		StartContinuousFunc: func(eventHandler func(event *speechtotext.SDKWrapperEvent)) error {
			go func() {
				// Simulate a Recognized event with some text
				eventHandler(&speechtotext.SDKWrapperEvent{
					EventType: speechtotext.Recognized,
					Recognized: &speech.SpeechRecognitionEventArgs{
						Result: speech.SpeechRecognitionResult{
							Text:       "test",
							Properties: &common.PropertyCollection{}, // Needs to exist for close to work
						},
					},
				})
				time.Sleep(time.Second) // delay to simulate processing time
				// Simulate a Cancellation event
				eventHandler(&speechtotext.SDKWrapperEvent{
					EventType: speechtotext.Cancellation,
					Cancellation: &speech.SpeechRecognitionCanceledEventArgs{
						Reason: common.CancellationReason(2), //EndOfStream
						SpeechRecognitionEventArgs: speech.SpeechRecognitionEventArgs{
							Result: speech.SpeechRecognitionResult{
								Properties: &common.PropertyCollection{}, // Needs to exist for close to work
							},
							RecognitionEventArgs: speech.RecognitionEventArgs{},
						},
					},
				})
			}()
			return nil
		},
		StopContinuousFunc: func() error {
			return nil
		},
	}

	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("testdata/working.ogg")
		if err != nil {
			t.Fatalf("unexpected error opening file: %v", err)
		}
		defer file.Close()

		http.ServeContent(w, r, "working.ogg", time.Now(), file)
	}))
	defer server.Close()

	// Use the URL of the test server as the audio link
	audioLink := server.URL

	result, err := speechtotext.HandleAudioLink(audioLink, mockWrapper)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "test" {
		t.Errorf("unexpected result, got: %v, want: %v", result, "test")
	}
}

func TestHandleAudioLink_ErrorFromDownload(t *testing.T) {
	t.Parallel()
	mockWrapper := &MockSDKWrapper{
		StartContinuousFunc: func(eventHandler func(event *speechtotext.SDKWrapperEvent)) error {
			return errors.New("mock error")
		},
		StopContinuousFunc: func() error {
			return nil
		},
	}

	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer server.Close()

	// Use the URL of the test server as the audio link
	audioLink := server.URL

	_, err := speechtotext.HandleAudioLink(audioLink, mockWrapper)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

// TODO: Add cases