package models

import "time"

type PrivacyLevel int
type Priority int

const (
	PrivacyLow PrivacyLevel = iota
	PrivacyMedium
	PrivacyHigh
)

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
)

type Metadata struct {
	Gps              GpsInfo
	OriginalDateTime time.Time
	ModifiedDateTime time.Time
	CreationDateTime time.Time
}
