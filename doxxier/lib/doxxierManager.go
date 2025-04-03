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
	println("Initializing DoxxierManager")
	orchestrator := network.TransportOrchestrator{}
	if err := orchestrator.Initialise(); err != nil {
		fmt.Printf("Failed to initialize orchestrator: %v\n", err)
		return nil
	}

	transport := orchestrator.GetTransportation(network.CONNECTION_PRIVACY_HIGH)
	callbackChan := make(chan string)

	dm := &DoxxierManager{
		doxxier:       models.NewDoxxier(),
		transport:     transport,
		callbackChan:  callbackChan,
		eventCallback: eventCallback,
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				println(fmt.Printf("Recovered from panic in listenForCallbacks: %v\n", r))
			}
		}()
		dm.listenForCallbacks()
	}()
	jsonData, err := dm.doxxier.ToJson()
	if err != nil {
		println(fmt.Printf("Error converting Doxxier to JSON: %v\n", err))
	} else {
		println(fmt.Sprintf("Initialized DoxxierManager: %v", jsonData))
	}
	err = transport.Connect(callbackChan, "")
	if err != nil {
		fmt.Printf("Error connecting to transport: %v\n", err)
		close(callbackChan) // Close the channel on error
		return nil
	}
	// dm.contact = contact
	// fmt.Println("Generated contact: ", dm.contact)

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
