package models

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
)

type Doxxier struct {
	Id                string        `json:"id"`
	Description       string        `json:"description"`
	Parts             []DoxxierPart `json:"parts"`
	Recipient         string        `json:"recipient"`
	CreatedAt         time.Time     `json:"created_at"`
	TransmissionStart time.Time     `json:"transmission_start"`
	TransmissionEnd   time.Time     `json:"transmission_end"`
}

type DoxxierOption func(*Doxxier)

func DoxxierWithId(id string) DoxxierOption {
	return func(d *Doxxier) {
		d.Id = id
	}
}

func DoxxierWithParts(parts []DoxxierPart) DoxxierOption {
	return func(d *Doxxier) {
		d.Parts = parts
	}
}

func NewDoxxier(params ...DoxxierOption) *Doxxier {
	doxxier := &Doxxier{
		Id:        uuid.New().String(),
		CreatedAt: time.Now(),
	}
	for _, opt := range params {
		opt(doxxier)
	}
	return doxxier
}

func (d *Doxxier) AddPart(part DoxxierPart) {
	d.Parts = append(d.Parts, part)
}

func (d *Doxxier) GetPart(id string) *DoxxierPart {
	for _, part := range d.Parts {
		if part.Id == id {
			return &part
		}
	}
	return nil
}

func (d *Doxxier) ToJson() (string, error) {
	jsonData, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (d *Doxxier) UpdateFromDoxxier(doxxier Doxxier) error {
	if doxxier.Id == "" && doxxier.Description == "" && len(doxxier.Parts) == 0 && doxxier.Recipient == "" && doxxier.CreatedAt.IsZero() && doxxier.TransmissionStart.IsZero() && doxxier.TransmissionEnd.IsZero() {
		return errors.New("doxxier cannot be nil")
	}
	d.Description = doxxier.Description
	d.Recipient = doxxier.Recipient
	return nil
}

// Cloneable interface implementation
func (d *Doxxier) Clone() Doxxier {
	// Marshal the original struct to JSON
	data, err := json.Marshal(d)
	if err != nil {
		errors.WithMessage(err, "Error marshalling doxxier")
	}

	// Unmarshal the JSON back into a new struct
	var copy Doxxier
	err = json.Unmarshal(data, &copy)
	if err != nil {
		errors.WithMessage(err, "Error unmarshalling doxxier")
	}

	return copy
}
