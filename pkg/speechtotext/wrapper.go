package speechtotext

import "github.com/Microsoft/cognitive-services-speech-sdk-go/speech"

// See https://github.com/microsoft/cognitive-services-speech-sdk-go/blob/v1.29.0/samples/recognizer/wrapper.go

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
)

type SDKWrapperEventType int

const (
	Cancellation SDKWrapperEventType = iota
	Recognizing
	Recognized
)

type SDKWrapperEvent struct {
	EventType    SDKWrapperEventType
	Cancellation *speech.SpeechRecognitionCanceledEventArgs
	Recognized   *speech.SpeechRecognitionEventArgs
	Recognizing  *speech.SpeechRecognitionEventArgs
}

func (event *SDKWrapperEvent) Close() {
	if event.Cancellation != nil {
		event.Cancellation.Close()
	}
	if event.Recognizing != nil {
		event.Recognizing.Close()
	}
	if event.Recognized != nil {
		event.Recognized.Close()
	}
}

type SDKWrapper struct {
	stream     *audio.PushAudioInputStream
	recognizer *speech.SpeechRecognizer
	started    int32
}

func NewWrapper(subscription string, region string) (*SDKWrapper, error) {
	format, err := audio.GetDefaultInputFormat()
	if err != nil {
		return nil, err
	}
	defer format.Close()
	stream, err := audio.CreatePushAudioInputStreamFromFormat(format)
	if err != nil {
		return nil, err
	}
	audioConfig, err := audio.NewAudioConfigFromStreamInput(stream)
	if err != nil {
		stream.Close()
		return nil, err
	}
	defer audioConfig.Close()

	config, err := speech.NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		stream.Close()
		return nil, err
	}
	defer config.Close()

	autodetect, err := speech.NewAutoDetectSourceLanguageConfigFromLanguages([]string{
		"en-US",
		"en-GB",
		"fi-FI",
	})
	if err != nil {
		stream.Close()
		return nil, err
	}

	recognizer, err := speech.NewSpeechRecognizerFomAutoDetectSourceLangConfig(config, autodetect, audioConfig)
	if err != nil {
		stream.Close()
		return nil, err
	}
	wrapper := new(SDKWrapper)
	wrapper.recognizer = recognizer
	wrapper.stream = stream
	return wrapper, nil
}

func (wrapper *SDKWrapper) Close() {
	wrapper.stream.CloseStream()
	<-wrapper.recognizer.StopContinuousRecognitionAsync()
	wrapper.stream.Close()
	wrapper.recognizer.Close()
}

func (wrapper *SDKWrapper) Write(buffer []byte) error {
	if atomic.LoadInt32(&wrapper.started) != 1 {
		return fmt.Errorf("Trying to write when recognizer is stopped")
	}
	return wrapper.stream.Write(buffer)
}

func (wrapper *SDKWrapper) StartContinuous(callback func(*SDKWrapperEvent)) error {
	if atomic.SwapInt32(&wrapper.started, 1) == 1 {
		return nil
	}
	wrapper.recognizer.Recognized(func(event speech.SpeechRecognitionEventArgs) {
		wrapperEvent := new(SDKWrapperEvent)
		wrapperEvent.EventType = Recognized
		wrapperEvent.Recognized = &event
		callback(wrapperEvent)
	})
	wrapper.recognizer.Recognizing(func(event speech.SpeechRecognitionEventArgs) {
		wrapperEvent := new(SDKWrapperEvent)
		wrapperEvent.EventType = Recognizing
		wrapperEvent.Recognizing = &event
		callback(wrapperEvent)
	})
	wrapper.recognizer.Canceled(func(event speech.SpeechRecognitionCanceledEventArgs) {
		wrapperEvent := new(SDKWrapperEvent)
		wrapperEvent.EventType = Cancellation
		wrapperEvent.Cancellation = &event
		callback(wrapperEvent)
	})
	return <-wrapper.recognizer.StartContinuousRecognitionAsync()
}

func (wrapper *SDKWrapper) StopContinuous() error {
	if atomic.SwapInt32(&wrapper.started, 0) == 0 {
		return nil
	}
	var empty = []byte{}
	wrapper.stream.Write(empty)
	wrapper.recognizer.Recognized(nil)
	wrapper.recognizer.Recognizing(nil)
	wrapper.recognizer.Canceled(nil)
	return <-wrapper.recognizer.StopContinuousRecognitionAsync()
}

func PumpFileContinuously(stop chan int, ready chan<- struct{}, filename string, wrapper *SDKWrapper) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	buffer := make([]byte, 3200)
	for {
		select {
		case <-stop:
			fmt.Println("Stopping pump...")
			close(ready)
			return
		case <-time.After(1 * time.Millisecond):
		}
		n, err := reader.Read(buffer)
		if err == io.EOF {
			close(ready)
			err = wrapper.Write(buffer[0:n])
			if err != nil {
				fmt.Println("Error writing last data chunk to the stream")
			}
			return
		}
		if err != nil {
			fmt.Println("Error reading file: ", err)
			break
		}
		err = wrapper.Write(buffer[0:n])
		if err != nil {
			fmt.Println("Error writing to the stream")
		}
	}
}

func PumpFileContinuously2(stop chan int, filename string, wrapper *SDKWrapper) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	buffer := make([]byte, 3200)
	for {
		select {
		case <-stop:
			fmt.Println("Stopping pump...")
			return
		case <-time.After(100 * time.Millisecond):
		}
		n, err := reader.Read(buffer)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("Error reading file: ", err)
			break
		}
		err = wrapper.Write(buffer[0:n])
		if err != nil {
			fmt.Println("Error writing to the stream")
		}
	}
}