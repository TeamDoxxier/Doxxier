package models

import "time"

type Doxxier struct {
	Id                string
	Parts             []DoxxierPart
	ClientId          string
	TransmissionStart time.Time
}

func InitialiseDoxxier() Doxxier {
	return Doxxier{}
}
