package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"t-pain/pkg/database"
	"t-pain/pkg/models"
	"time"
)

func main() {
	var data []models.PainDescriptionLogEntry

	// Generate random data going back two years from today containing 1-2 pain descriptions per day
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := 0; i < 20; i++ {
		// Generate a date going back from today
		date := time.Now().AddDate(0, 0, -i)

		// Generate 1-2 pain descriptions per day
		numDescriptions := r.Intn(2) + 1

		for j := 0; j < numDescriptions; j++ {
			// Generate random level of pain between 1 and 10
			level := r.Intn(10) + 1

			// Generate random location of pain
			location := r.Intn(len(models.BodyPartMapping)) + 1

			// Generate random description
			description := "Pain in " + models.BodyPartMapping[location]

			// Generate random numbness
			numbness := r.Intn(2) == 1

			// Generate random side
			side := r.Intn(len(models.SideMap)) + 1

			// Generate numbness description
			numbnessDescription := ""
			if numbness {
				numbnessDescription = "Numbness in " + models.BodyPartMapping[location]
			}

			// Create PainDescription and append to data
			painDescription := models.PainDescription{
				Timestamp:           date,
				Level:               level,
				LocationId:          location,
				SideId:              side,
				Description:         description,
				Numbness:            numbness,
				NumbnessDescription: numbnessDescription,
			}

			pdLog := painDescription.MapToLogEntry(175255021)

			data = append(data, pdLog)
		}
	}

	client, err := database.NewLogAnalyticsClient(
		os.Getenv("DATA_COLLECTION_ENDPOINT"),
		os.Getenv("DATA_COLLECTION_RULE_ID"),
		os.Getenv("DATA_COLLECTION_STREAM_NAME"))
	if err != nil {
		panic(err)
	}

	err = client.SavePainDescriptionsToLogAnalytics(data)
	if err != nil {
		panic(err)
	}

	//err := SavePainDescriptionsToFile(data, "paindescriptions.json")
	//if err != nil {
	//	panic(err)
	//}
}

func SavePainDescriptionsToFile(data []models.PainDescriptionLogEntry, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}