package network

import "doxxier.tech/doxxier/pkg/models"

type Transport interface {
	Connect(callbackChan chan string) (string, error)
	Send(string, models.Doxxier) error
	Disconnect() error
}
