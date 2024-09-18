package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDoxxier_Default(t *testing.T) {
	doxxier := NewDoxxier()

	assert.NotNil(t, doxxier)
	assert.NotEmpty(t, doxxier.Id)
	assert.Empty(t, doxxier.Parts)
	assert.Empty(t, doxxier.Recipient)
	assert.Zero(t, doxxier.TransmissionStart)
}

func TestNewDoxxier_WithId(t *testing.T) {
	id := uuid.New().String()
	doxxier := NewDoxxier(DoxxierWithId(id))

	assert.NotNil(t, doxxier)
	assert.Equal(t, id, doxxier.Id)
	assert.Empty(t, doxxier.Parts)
	assert.Empty(t, doxxier.Recipient)
	assert.Zero(t, doxxier.TransmissionStart)
}

func TestNewDoxxier_WithParts(t *testing.T) {
	parts := []DoxxierPart{
		{Id: "part1"},
		{Id: "part2"},
	}
	doxxier := NewDoxxier(DoxxierWithParts(parts))

	assert.NotNil(t, doxxier)
	assert.NotEmpty(t, doxxier.Id)
	assert.Equal(t, parts, doxxier.Parts)
	assert.Empty(t, doxxier.Recipient)
	assert.Zero(t, doxxier.TransmissionStart)
}

func TestNewDoxxier_WithMultipleOptions(t *testing.T) {
	id := uuid.New().String()
	parts := []DoxxierPart{
		{Id: "part1"},
		{Id: "part2"},
	}
	doxxier := NewDoxxier(DoxxierWithId(id), DoxxierWithParts(parts))

	assert.NotNil(t, doxxier)
	assert.Equal(t, id, doxxier.Id)
	assert.Equal(t, parts, doxxier.Parts)
	assert.Empty(t, doxxier.Recipient)
	assert.Zero(t, doxxier.TransmissionStart)
}
func TestDoxxier_ToJson(t *testing.T) {
	doxxier := NewDoxxier()
	jsonStr := doxxier.ToJson()

	assert.NotEmpty(t, jsonStr)
	assert.Contains(t, jsonStr, `"Id":"`)
	assert.Contains(t, jsonStr, `"Parts":`)
	assert.Contains(t, jsonStr, `"Recipient":"`)
	assert.Contains(t, jsonStr, `"TransmissionStart":"`)
	assert.Contains(t, jsonStr, `"TransmissionEnd":"`)
}

func TestDoxxier_ToJson_WithValues(t *testing.T) {
	parsedDate, _ := time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")
	id := uuid.New().String()
	parts := []DoxxierPart{
		{Id: "part1"},
		{Id: "part2"},
	}
	recipient := "test@example.com"
	transmissionStart := parsedDate
	transmissionEnd := transmissionStart.Add(1 * time.Hour)

	doxxier := &Doxxier{
		Id:                id,
		Parts:             parts,
		Recipient:         recipient,
		TransmissionStart: transmissionStart,
		TransmissionEnd:   transmissionEnd,
	}
	jsonStr := doxxier.ToJson()

	assert.NotEmpty(t, jsonStr)
	assert.Contains(t, jsonStr, `"id":"`+id+`"`)
	assert.Contains(t, jsonStr, `"recipient":"`+recipient+`"`)
	assert.Contains(t, jsonStr, `"transmission_start":"`+transmissionStart.Format(time.RFC3339)+`"`)
	assert.Contains(t, jsonStr, `"transmission_end":"`+transmissionEnd.Format(time.RFC3339)+`"`)

	err := json.Unmarshal([]byte(jsonStr), &doxxier)
	assert.NoError(t, err)
	assert.Equal(t, parts[0].Id, doxxier.Parts[0].Id)
}
func TestDoxxier_AddPart(t *testing.T) {
	doxxier := NewDoxxier()
	part := DoxxierPart{Id: "part1"}

	doxxier.AddPart(part)

	assert.Len(t, doxxier.Parts, 1)
	assert.Equal(t, part, doxxier.Parts[0])
}

func TestDoxxier_AddPart_MultipleParts(t *testing.T) {
	doxxier := NewDoxxier()
	part1 := DoxxierPart{Id: "part1"}
	part2 := DoxxierPart{Id: "part2"}

	doxxier.AddPart(part1)
	doxxier.AddPart(part2)

	assert.Len(t, doxxier.Parts, 2)
	assert.Equal(t, part1, doxxier.Parts[0])
	assert.Equal(t, part2, doxxier.Parts[1])
}
