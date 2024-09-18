package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type DoxxierPart struct {
	Id              string            `json:"id"`
	Content         []byte            `json:"-"`
	Context         string            `json:"context"`
	Descriptors     map[string]string `json:"descriptors"`
	Priority        Priority          `json:"priority"`
	PrivacyLevel    PrivacyLevel      `json:"privacy_level"`
	Transformations []string          `json:"transformations"`
	Metadata        Metadata          `json:"metadata"`
}

type DoxxierPartOption func(*DoxxierPart)

func DoxxierPartWithId(id string) DoxxierPartOption {
	return func(dp *DoxxierPart) {
		dp.Id = id
	}
}

func DoxxierPartWithContent(content []byte) DoxxierPartOption {
	return func(dp *DoxxierPart) {
		dp.Content = content
	}
}

func NewDoxxierPart(params ...DoxxierPartOption) *DoxxierPart {
	if len(params) == 0 {
		return &DoxxierPart{
			Id: uuid.New().String(),
		}
	}
	dp := &DoxxierPart{}
	for _, opt := range params {
		opt(dp)
	}
	return dp
}

func (dp *DoxxierPart) ToJson() (string, error) {
	jsonData, err := json.Marshal(dp)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (dp *DoxxierPart) MarshalJSON() ([]byte, error) {
	type Alias DoxxierPart
	return json.Marshal(&struct {
		Size int `json:"size"`
		*Alias
	}{
		Size:  dp.Size(),
		Alias: (*Alias)(dp),
	})
}

func (dp *DoxxierPart) Size() int {
	if dp.Content == nil {
		return 0
	}
	return len(dp.Content)
}
