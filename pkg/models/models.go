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
	for ID, bodyPart := range bp {
		result.WriteString(fmt.Sprintf("%d: %s\n", ID, bodyPart))
	}
	result.WriteString("```\n")
	return result.String()
}

func (bp BodyParts) StringNameFirst() string {
	var result strings.Builder
	result.WriteString("```\n")
	for ID, bodyPart := range bp {
		result.WriteString(fmt.Sprintf("%s: %d\n", bodyPart, ID))
	}
	result.WriteString("```\n")
	return result.String()
}

var BodyPartMapping = BodyParts{
	1:  "Head",
	2:  "Neck",
	3:  "Shoulder",
	4:  "Arm",
	5:  "Elbow",
	6:  "Wrist",
	7:  "Hand",
	8:  "Upper Back",
	9:  "Lower Back",
	10: "Hip",
	11: "Leg",
	12: "Knee",
	13: "Ankle",
	14: "Foot",
	15: "Chest",
	16: "Abdomen",
	17: "Pelvis",
	18: "Genitals",
	19: "Thigh",
	20: "Calf",
	21: "Toes",
}

type Sides map[int]string

func (sd Sides) String() string {
	var result strings.Builder
	result.WriteString("```\n")
	for ID, side := range sd {
		result.WriteString(fmt.Sprintf("%d: %s\n", ID, side))
	}
	result.WriteString("```\n")
	return result.String()
}

func (sd Sides) StringNameFirst() string {
	var result strings.Builder
	result.WriteString("```\n")
	for ID, side := range sd {
		result.WriteString(fmt.Sprintf("%s: %d\n", side, ID))
	}
	result.WriteString("```\n")
	return result.String()
}

var SideMap = Sides{
	1: "Both",
	2: "Left",
	3: "Right",
}

type PainDescription struct {
	Timestamp           time.Time `json:"timestamp,omitempty"`
	Level               int       `json:"level"`
	LocationId          int       `json:"locationId"`
	SideId              int       `json:"sideId"`
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

func (p *PainDescription) MapToLogEntry(userId int64) (PainDescriptionLogEntry, error) {
	locationName, locOk := BodyPartMapping[p.LocationId]
	sideName, sideOk := SideMap[p.SideId]
	userName, userOk := UserIDs[userId]

	if !locOk {
		return PainDescriptionLogEntry{}, fmt.Errorf("invalid LocationId: %d", p.LocationId)
	}
	if !sideOk {
		return PainDescriptionLogEntry{}, fmt.Errorf("invalid SideId: %d", p.SideId)
	}
	if !userOk {
		return PainDescriptionLogEntry{}, fmt.Errorf("invalid UserId: %d", userId)
	}

	p.Timestamp = p.Timestamp.UTC()

	pdLog := PainDescriptionLogEntry{
		PainDescription: *p,
		LocationName:    locationName,
		SideName:        sideName,
		UserName:        userName,
	}

	return pdLog, nil
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

		if i == objType.NumField()-1 {
			sb.WriteString(fmt.Sprintf("\t\"%s\": %s\n", field.Name, fieldValue.Type()))
		} else {
			sb.WriteString(fmt.Sprintf("\t\"%s\": %s,\n", field.Name, fieldValue.Type()))
		}
	}

	sb.WriteString("}\n")

	return sb.String()
}

func PrintSinglePainDescriptionJSONFormat(p PainDescription) string {
	sb := strings.Builder{}

	sb.WriteString("{\n")

	sb.WriteString(fmt.Sprintf("\t\"%s\": %d,\n", "level", p.Level))
	sb.WriteString(fmt.Sprintf("\t\"%s\": %d,\n", "locationId", p.LocationId))
	sb.WriteString(fmt.Sprintf("\t\"%s\": %d,\n", "sideId", p.SideId))
	sb.WriteString(fmt.Sprintf("\t\"%s\": %q,\n", "description", p.Description))
	sb.WriteString(fmt.Sprintf("\t\"%s\": %t,\n", "numbness", p.Numbness))
	sb.WriteString(fmt.Sprintf("\t\"%s\": %q\n", "numbnessDescription", p.NumbnessDescription))

	sb.WriteString("}")

	return sb.String()
}

type PainDescriptionLogEntry struct {
	PainDescription
	LocationName string `json:"locationName"`
	SideName     string `json:"sideName"`
	UserName     string `json:"userName"`
}