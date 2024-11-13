package lib

import (
	"fmt"

	"doxxier.tech/doxxier/lib/network"
	"doxxier.tech/doxxier/pkg/models"
	"github.com/pkg/errors"
)

type DoxxierManager struct {
	doxxier       *models.Doxxier
	eventCallback func(string)
	transport     network.Transport
	callbackChan  chan string
	contact       string
}

func NewDoxxierManager(eventCallback func(string)) *DoxxierManager {
	orchestrator := network.TransportOrchestrator{}
	orchestrator.NewTransportOrchestrator()
	transport := orchestrator.GetTransportation(network.CONNECTION_PRIVACY_HIGH)

	callbackChan := make(chan string)

	dm := &DoxxierManager{
		doxxier:       models.NewDoxxier(),
		transport:     transport,
		callbackChan:  callbackChan,
		eventCallback: eventCallback,
	}

	go dm.listenForCallbacks()
	contact, err := transport.Connect(callbackChan)
	dm.contact = contact

	if err != nil {
		errors.WithMessage(err, "Error connecting to transport")
	}
	fmt.Println("Generated contact: ", dm.contact)

	return dm
}

func (dm *DoxxierManager) listenForCallbacks() {
	for message := range dm.callbackChan {
		dm.eventCallback(message)
	}
}

func (dm *DoxxierManager) Close() {
	dm.transport.Disconnect()
	close(dm.callbackChan)
}

func (dm *DoxxierManager) GetDoxxier() *models.Doxxier {
	return dm.doxxier
}

func (dm *DoxxierManager) AddPart() *models.DoxxierPart {
	part := models.NewDoxxierPart()
	dm.doxxier.AddPart(*part)
	return part
}

func (dm *DoxxierManager) GetPart(id string) *models.DoxxierPart {
	part := dm.doxxier.GetPart(id)
	return part
}

func (dm *DoxxierManager) SendDoxxier() error {
	if dm.doxxier.Recipient == "" {
		return errors.New("Recipient not set")
	}

	return nil
}

func (dm *DoxxierManager) EventCallback(recipient string) {
	dm.doxxier.Recipient = recipient
}

func (dm *DoxxierManager) Connect() error {
	contact, err := dm.transport.Connect(dm.callbackChan)
	if err != nil {
		return err
	}
	dm.contact = contact
	fmt.Println("Generated contact: ", dm.contact)
	return nil
}
