package network

import "doxxier.tech/doxxier/pkg/models"

type Transport interface {
	Connect() error
	Send(models.Doxxier) error
}
