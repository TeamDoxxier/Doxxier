package network

import "doxxier.tech/doxxier/pkg/models"

type Transport interface {
	Initialise() (string, error)
	Connect(callbackChan chan string, recipient string) error
	Send(contact string, doxxier models.Doxxier) error
	SendMessage(message string) error
	Disconnect() error
}
