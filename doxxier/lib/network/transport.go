package network

import "doxxier.tech/doxxier/lib/models"

type Transport interface {
	Connect() error
	Send(models.Doxxier) error
}
