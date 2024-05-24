package models

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

type DoxxierContext struct {
	Content         []byte
	Context         string
	Descriptiors    map[string]string
	Priority        Priority
	PrivacyLevel    PrivacyLevel
	Transformations []string
}
