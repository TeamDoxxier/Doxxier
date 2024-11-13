package network

import network "doxxier.tech/doxxier/lib/network/cmixx_e2e"

type PrivacyLevel int

const (
	CONNECTION_PRIVACY_LOW    PrivacyLevel = 0
	CONNECTION_PRIVACY_MEDIUM PrivacyLevel = 1
	CONNECTION_PRIVACY_HIGH   PrivacyLevel = 2
)

type TransportOrchestrator struct {
	TransportMap map[PrivacyLevel]Transport
}

func (orchestrator *TransportOrchestrator) RegisterTransport(privacyLevel PrivacyLevel, transport Transport) {
	orchestrator.TransportMap[privacyLevel] = transport
}

func (orchestrator *TransportOrchestrator) GetTransportation(privacyLevel PrivacyLevel) Transport {
	return orchestrator.TransportMap[privacyLevel]
}

func (orchestrator *TransportOrchestrator) NewTransportOrchestrator() {
	config := network.CMixxConfig{
		StoragePath: "assets/storage",
	}
	transport, _ := network.NewCMixxE2eTransport(config)
	orchestrator.TransportMap = make(map[PrivacyLevel]Transport)
	orchestrator.RegisterTransport(CONNECTION_PRIVACY_HIGH, transport)
}
