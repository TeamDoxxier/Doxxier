package network

import (
	"doxxier.tech/doxxier/lib/network/wasm/cmix_e2e"
)

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

func (orchestrator *TransportOrchestrator) Initialise() error {
	config := cmix_e2e.CMixE2eConfig{
		StoragePath: "xxstore",
	}
	println("Initializing CMixE2eTransport")
	transport, err := cmix_e2e.NewCmixE2eTransport(config)
	if err != nil {
		return err
	}
	orchestrator.TransportMap = make(map[PrivacyLevel]Transport)
	orchestrator.RegisterTransport(CONNECTION_PRIVACY_HIGH, transport)
	println("Registered CMixE2eTransport")
	return nil
}
