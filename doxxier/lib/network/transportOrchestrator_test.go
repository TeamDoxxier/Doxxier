package network

import (
	"testing"

	"doxxier.tech/doxxier/lib/network/wasm/cmix_e2e"
	"github.com/stretchr/testify/assert"
)

func TestNewTransportOrchestrator(t *testing.T) {
	orchestrator := &TransportOrchestrator{
		TransportMap: make(map[PrivacyLevel]Transport),
	}

	orchestrator.Initialise()

	transport := orchestrator.GetTransportation(CONNECTION_PRIVACY_HIGH)
	if transport == nil {
		t.Errorf("Expected transport to be registered for CONNECTION_PRIVACY_HIGH, but it was nil")
	}

	// Additional checks can be added here to verify the transport configuration
	config := cmix_e2e.CMixE2eConfig{
		StoragePath: "assets/storage",
	}
	expectedTransport, err := cmix_e2e.NewCmixE2eTransport(config)
	if err != nil {
		t.Errorf("Error creating expected transport: %v", err)
	}
	assert.Equal(t, expectedTransport, transport)

}
