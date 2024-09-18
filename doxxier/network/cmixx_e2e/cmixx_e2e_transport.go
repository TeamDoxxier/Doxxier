package network

import "doxxier.tech/doxxier/pkg/models"

type CmixxE2eTransport struct {
}

func (c *CmixxE2eTransport) Connect() error {
	// Implementation for connecting to the transport
	return nil
}

func (c *CmixxE2eTransport) Send(doxxier models.Doxxier) error {
	// Implementation for sending the doxxier
	return nil
}
