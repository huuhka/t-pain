package speechtotext

import (
	"fmt"
	"os"
)

func deleteFromDisk(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		fmt.Println("Error deleting file: ", err)
	}
}