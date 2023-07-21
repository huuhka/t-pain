package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type BodyParts map[int]string

func (bp BodyParts) String() string {
	var result strings.Builder
	result.WriteString("```\n")
	for bodyPart, ID := range bp {
		result.WriteString(fmt.Sprintf("%d: %s\n", ID, bodyPart))
	}
	result.WriteString("```\n")
	return result.String()
}

func (bp BodyParts) StringNameFirst() string {
	var result strings.Builder
	result.WriteString("```\n")
	for bodyPart, ID := range bp {
		result.WriteString(fmt.Sprintf("%s: %d\n", bodyPart, ID))
	}
	result.WriteString("```\n")
	return result.String()
}

var BodyPartMapping = BodyParts{
	1:  "Head",
	2:  "Neck",
	3:  "Left Shoulder",
	4:  "Right Shoulder",
	5:  "Left Arm",
	6:  "Right Arm",
	7:  "Left Elbow",
	8:  "Right Elbow",
	9:  "Left Wrist",
	10: "Right Wrist",
	11: "Left Hand",
	12: "Right Hand",
	13: "Upper Back",
	14: "Lower Back",
	15: "Hip",
	16: "Left Leg",
	17: "Right Leg",
	18: "Left Knee",
	19: "Right Knee",
	20: "Left Ankle",
	21: "Right Ankle",
	22: "Left Foot",
	23: "Right Foot",
	24: "Chest",
	25: "Abdomen",
	26: "Pelvis",
	27: "Genitals",
	28: "Left Thigh",
	29: "Right Thigh",
	30: "Left Calf",
	31: "Right Calf",
	32: "Left Toes",
	33: "Right Toes",
}

type PainDescription struct {
	Timestamp           time.Time `json:"timestamp,omitempty"`
	Level               int       `json:"level"`
	LocationId          int       `json:"location"`
	Description         string    `json:"description"`
	Numbness            bool      `json:"numbness"`
	NumbnessDescription string    `json:"numbnessDescription,omitempty"`
}

func NewPainDescription() PainDescription {
	return PainDescription{
		Timestamp: time.Now(),
	}
}

func (p *PainDescription) UnmarshalJSON(data []byte) error {
	type Alias PainDescription
	aux := &struct {
		Timestamp json.RawMessage `json:"timestamp,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux.Timestamp) > 0 {
		return json.Unmarshal(aux.Timestamp, &p.Timestamp)
	}

	return nil
}

func PrintPainDescriptionJSONFormat() string {
	sb := strings.Builder{}

	p := NewPainDescription()
	objType := reflect.TypeOf(p)

	sb.WriteString("{\n")

	objValue := reflect.ValueOf(p)
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		if field.Type == reflect.TypeOf(time.Time{}) {
			continue
		}
		fieldValue := objValue.Field(i)

		sb.WriteString(fmt.Sprintf("\t\"%s\": %s\n", field.Name, fieldValue.Type()))
	}

	sb.WriteString("}\n")

	return sb.String()
}

type PainDescriptionLogEntry struct {
	PainDescription
	LocationName string `json:"locationName"`
}

func (p PainDescriptionLogEntry) MarshalJSON() ([]byte, error) {
	type Alias PainDescriptionLogEntry
	return json.Marshal(&struct {
		TimeGenerated time.Time `json:"TimeGenerated"`
		*Alias
	}{
		TimeGenerated: p.Timestamp,
		Alias:         (*Alias)(&p),
	})
}