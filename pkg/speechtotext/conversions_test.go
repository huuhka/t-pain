package speechtotext

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Unsure if these should really just be deleted as wrapper_test.go already tests some of the same things from
// outside the package. Leaving them in for now as they at least found a bug already.

func TestHandleAudioFileSetup_WithValidValues(t *testing.T) {
	t.Parallel()
	// Create a test server with a handler that serves up a test file
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./testdata/working.ogg")
	}))
	defer server.Close()

	// Run handleAudioFileSetup with the URL from our test server
	filename, err := handleAudioFileSetup(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the resulting file exists and then clean it up
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatal("Wav file was not created")
	}
	defer os.Remove(filename)
}

func Test_ConvertOggToWav_ShouldErrorOnFileNotFound(t *testing.T) {
	t.Parallel()
	inputFile := "./testdata/doesnotexist.ogg"
	outputFile := "./testdata/doesnotexist.wav"

	err := convertOggToWav(inputFile, outputFile)
	if err == nil {
		defer os.Remove(outputFile)
		t.Errorf("Expected error, got nil")
	}
}

func Test_DownloadFile(t *testing.T) {
	t.Parallel()
	// Create a test server with a handler that serves up a test string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "This is test data")
	}))
	defer server.Close()

	filename := "./testdata/downloaded.txt"

	err := downloadFile(server.URL, filename)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the resulting file exists
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the file contains the data we sent
	if strings.TrimSpace(string(content)) != "This is test data" {
		t.Fatalf("Expected file to contain '%s', but it contained '%s'", "This is test data", string(content))
	}

	defer os.Remove(filename)
}

func Test_DeleteFromDisk(t *testing.T) {
	filename := "./testdata/deletetest.txt"

	// Create a test file
	err := os.WriteFile(filename, []byte("test data"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the test file
	deleteFromDisk(filename)

	// Check that the file was deleted
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		t.Fatal("File was not deleted")
	}
}