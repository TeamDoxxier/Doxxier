package main

import (
	"testing"

	"doxxier.tech/doxxier/lib"
	mockjs "doxxier.tech/wasm/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestOnEvent(t *testing.T) {
	// Create a mock global object
	mockGlobal := mockjs.ValueOf(map[string]interface{}{
		"console": mockjs.ValueOf(map[string]interface{}{
			"log": func(args ...mockjs.Value) mockjs.Value {
				assert.Equal(t, "test message", args[0].String())
				return mockjs.Undefined
			},
		}),
	})

	// Set the mock global to a package-level variable
	global = mockGlobal

	// Call the function being tested
	OnEvent("test message")
}

func TestOnMessage(t *testing.T) {
	// Create a mock global object
	mockGlobal := mockjs.ValueOf(map[string]interface{}{
		"console": mockjs.ValueOf(map[string]interface{}{
			"log": func(args ...mockjs.Value) mockjs.Value {
				assert.Equal(t, "test message", args[0].String())
				return mockjs.Undefined
			},
		}),
	})

	// Set the mock global to a package-level variable
	global = mockGlobal

	// Simulate the callback channel
	callbackChan = make(chan string, 1)
	callbackChan <- "test message"

	// Call the function being tested
	go onMessage()
}

func TestCreateDoxxier(t *testing.T) {
	// Create a mock global object
	mockGlobal := mockjs.ValueOf(map[string]interface{}{})
	global = mockGlobal

	// Mock the library function
	// lib.NewDoxxierManager = func(callback func(string)) *lib.DoxxierManager {
	// 	return &lib.DoxxierManager{}
	// }

	// Call the function being tested
	result := CreateDoxxier(mockjs.ValueOf(nil), nil)
	assert.NotNil(t, result)
}

func TestGetDoxxier(t *testing.T) {
	// Mock the DoxxierManager with a sample Doxxier
	doxxierManager := &lib.DoxxierManager{
		// Add mock methods if required
	}

	// Set up mock global object
	mockGlobal := mockjs.ValueOf(map[string]interface{}{})
	global = mockGlobal

	// Call the function being tested
	result := GetDoxxier(mockjs.ValueOf(nil), nil)
	assert.NotNil(t, result)
}

func TestUpdateDoxxier(t *testing.T) {
	// Mock the existing Doxxier
	mockDoxxier := &lib.Doxxier{
		Description: "",
		Recipient:   "",
	}

	// Mock DoxxierManager with a GetDoxxier method
	doxxierManager = lib.DoxxierManager{
		GetDoxxier: func() *lib.Doxxier {
			return mockDoxxier
		},
	}

	// Mock the arguments
	args := []mockjs.Value{
		mockjs.ValueOf(map[string]interface{}{
			"description": "Updated description",
			"recipient":   "Recipient Name",
		}),
	}

	// Call the function being tested
	UpdateDoxxier(mockjs.ValueOf(nil), args)

	// Assert the results
	assert.Equal(t, "Updated description", mockDoxxier.Description)
	assert.Equal(t, "Recipient Name", mockDoxxier.Recipient)
}
