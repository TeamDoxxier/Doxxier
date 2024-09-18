package wasm

import (
	"encoding/json"
	"testing"

	"doxxier.tech/doxxier/pkg/models"
)

// Mock implementation for testing
type MockJSGlobal struct {
	values map[string]interface{}
}

func (m *MockJSGlobal) Get(key string) interface{} {
	return m.values[key]
}

func (m *MockJSGlobal) Set(key string, value interface{}) {
	m.values[key] = value
}

var mockGlobal JSGlobal = &MockJSGlobal{values: make(map[string]interface{})}

func TestCreateDoxxier(t *testing.T) {
	mockGlobal := &MockJSGlobal{values: make(map[string]interface{})}
	global = mockGlobal

	CreateDoxxier(nil, nil)

	doxxierValue := mockGlobal.Get("doxxier")
	if doxxierValue == nil {
		t.Fatalf("Expected doxxier to be set in global, but it was not")
	}

	var doxxier models.Doxxier
	err := json.Unmarshal([]byte(doxxierValue.(string)), &doxxier)
	if err != nil {
		t.Fatalf("Error unmarshalling doxxier: %v", err)
	}

	if doxxier.Id == "" && len(doxxier.Parts) == 0 {
		t.Fatalf("Expected doxxier to be initialized, but it was empty")
	}
}

func TestAddPart(t *testing.T) {
	mockGlobal := &MockJSGlobal{values: make(map[string]interface{})}

	// Initialize doxxier in global
	doxxier := models.NewDoxxier()
	doxxierJSON, _ := json.Marshal(doxxier)
	mockGlobal.Set("doxxier", string(doxxierJSON))

	AddPart(nil, nil)

	doxxierValue := mockGlobal.Get("doxxier")
	if doxxierValue == nil {
		t.Fatalf("Expected doxxier to be set in global, but it was not")
	}

	var updatedDoxxier models.Doxxier
	err := json.Unmarshal([]byte(doxxierValue.(string)), &updatedDoxxier)
	if err != nil {
		t.Fatalf("Error unmarshalling doxxier: %v", err)
	}

	if len(updatedDoxxier.Parts) == 0 {
		t.Fatalf("Expected doxxier to have parts, but it was empty")
	}
}
