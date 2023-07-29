package speechtotext_test

import (
	"os"
	"strings"
	"t-pain/pkg/speechtotext"
	"testing"
	"time"
)

type MockSDKWrapper struct {
	receivedData        string
	StartContinuousFunc func(eventHandler func(event *speechtotext.SDKWrapperEvent)) error
	StopContinuousFunc  func() error
}

func (m *MockSDKWrapper) StartContinuous(eventHandler func(event *speechtotext.SDKWrapperEvent)) error {
	return m.StartContinuousFunc(eventHandler)
}

func (m *MockSDKWrapper) StopContinuous() error {
	return m.StopContinuousFunc()
}

func (m *MockSDKWrapper) Write(data []byte) error {
	m.receivedData += string(data)
	return nil
}

func TestPumpFileToStream(t *testing.T) {
	// Create a temporary file with data
	tempFile, err := os.CreateTemp(os.TempDir(), "testfile-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer os.Remove(tempFile.Name())

	text := "This is a test text"
	if _, err = tempFile.Write([]byte(text)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close the temp file: %v", err)
	}

	mockWrapper := &MockSDKWrapper{}

	// Create stop channel and close it after a small delay to stop the pump
	stop := make(chan int)
	go func() {
		time.Sleep(500 * time.Millisecond)
		close(stop)
	}()

	// Pump the file to stream
	speechtotext.PumpFileToStream(stop, tempFile.Name(), mockWrapper)

	// Check if the pumped text matches the original text
	if strings.Compare(text, mockWrapper.receivedData) != 0 {
		t.Errorf("Pumped text doesn't match original text. Got %v, expected %v", mockWrapper.receivedData, text)
	}
}