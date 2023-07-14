package speechtotext

import (
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"io"
	"log"
	"net/http"
	"os"
)

func pumpURLIntoStream(url string, stream *audio.PushAudioInputStream) {
	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		return
	}
	defer resp.Body.Close()

	// Create a new reader from the response body
	reader := io.Reader(resp.Body)
	buffer := make([]byte, 1000)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			fmt.Println("Done reading file.")
			break
		}
		if err != nil {
			fmt.Println("Error reading file: ", err)
			break
		}
		err = stream.Write(buffer[0:n])
		if err != nil {
			fmt.Println("Error writing to the stream")
		}
	}
	stream.CloseStream()
}

func deleteFromDisk(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		log.Println("Error deleting file: ", err)
	}
}