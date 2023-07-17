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
	13: "Back",
	14: "Hip",
	15: "Left Leg",
	16: "Right Leg",
	17: "Left Knee",
	18: "Right Knee",
	19: "Left Ankle",
	20: "Right Ankle",
	21: "Left Foot",
	22: "Right Foot",
	23: "Chest",
	24: "Abdomen",
	25: "Pelvis",
	26: "Genitals",
	27: "Left Thigh",
	28: "Right Thigh",
	29: "Left Calf",
	30: "Right Calf",
	31: "Left Toes",
	32: "Right Toes",
}

type PainDescription struct {
	TimeStamp           time.Time `json:"timestamp,omitempty"`
	Level               []int     `json:"level"`
	Location            []int     `json:"location"`
	Description         string    `json:"description"`
	Numbness            bool      `json:"numbness"`
	NumbnessDescription string    `json:"numbnessDescription,omitempty"`
}

func NewPainDescription() PainDescription {
	return PainDescription{
		TimeStamp: time.Now(),
	}
}

func (p *PainDescription) UnmarshalJSON(data []byte) error {
	type Alias PainDescription
	aux := &struct {
		TimeStamp json.RawMessage `json:"timestamp,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux.TimeStamp) > 0 {
		return json.Unmarshal(aux.TimeStamp, &p.TimeStamp)
	}

	return nil
}

func (p *PainDescription) StringFriendly() string {
	var result strings.Builder

	loc, _ := time.LoadLocation("Europe/Helsinki")
	tstamp := p.TimeStamp.Round(time.Minute).In(loc).Format("02-01-2006 15:04")

	result.WriteString(fmt.Sprintf("Timestamp: %s\n", tstamp))
	result.WriteString("Pains:\n")
	for i := range p.Level {
		result.WriteString(fmt.Sprintf("\t- Location: %s, Level: %d\n", BodyPartMapping[p.Location[i]], p.Level[i]))
	}
	result.WriteString(fmt.Sprintf("Description: %s\n", p.Description))
	result.WriteString(fmt.Sprintf("Numbness: %t\n", p.Numbness))
	result.WriteString(fmt.Sprintf("Numbness Description: %s\n", p.NumbnessDescription))
	return result.String()
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