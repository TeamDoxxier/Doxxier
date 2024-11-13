package network

import (
	"testing"

	network "doxxier.tech/doxxier/lib/network/cmixx_e2e"
	"github.com/stretchr/testify/assert"
)

func TestNewTransportOrchestrator(t *testing.T) {
	orchestrator := &TransportOrchestrator{
		TransportMap: make(map[PrivacyLevel]Transport),
	}

	orchestrator.NewTransportOrchestrator()

	transport := orchestrator.GetTransportation(CONNECTION_PRIVACY_HIGH)
	if transport == nil {
		t.Errorf("Expected transport to be registered for CONNECTION_PRIVACY_HIGH, but it was nil")
	}

	// Additional checks can be added here to verify the transport configuration
	config := network.CMixxConfig{
		StoragePath: "assets/storage",
	}
	expectedTransport, _ := network.NewCMixxE2eTransport(config)
	assert.Equal(t, expectedTransport, transport)

}
