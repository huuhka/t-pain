package speechtotext

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// handleAudioFileSetup downloads the audio file from the url, converts it to wav, and returns the wav file name
func handleAudioFileSetup(url string) (string, error) {
	// Generate new guid
	newGuid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	// Download file from url
	oggFileName := fmt.Sprintf("%s.ogg", newGuid.String())
	err = downloadFile(url, oggFileName)
	if err != nil {
		return "", err
	}
	log.Println("handleAudioFileSetup: Downloaded file to disk")
	defer deleteFromDisk(oggFileName)

	// Convert file to wav
	wavFileName := fmt.Sprintf("%s.wav", newGuid.String())
	err = convertOggToWav(fmt.Sprintf("%s.ogg", newGuid.String()), wavFileName)
	if err != nil {
		log.Println("handleAudioFileSetup: Error converting file to wav: ", err)
		return "", err
	}

	log.Println("handleAudioFileSetup: Converted file to wav at " + wavFileName)

	return wavFileName, nil
}

// convertOggToWav converts an Ogg audio file to a WAV file using FFmpeg.
func convertOggToWav(inputFile string, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-acodec", "pcm_s16le", "-ar", "16000", "-ac", "1", outputFile)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func downloadFile(url string, fileName string) error {
	// Open file
	file, err := os.Create(fileName)
	defer file.Close()

	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("downloadFile: error sending request, %w", err)
	}
	defer resp.Body.Close()

	// Create a new reader from the response body
	reader := io.Reader(resp.Body)
	buffer := make([]byte, 1000)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			_, err = file.Write(buffer[0:n])
			if err != nil {
				fmt.Println("Error writing last bytes to the file")
			}
			fmt.Println("Done reading file.")
			break
		}
		if err != nil {
			fmt.Println("Error reading file: ", err)
			break
		}
		_, err = file.Write(buffer[0:n])
		if err != nil {
			fmt.Println("Error writing to the file")
		}
	}

	return nil
}