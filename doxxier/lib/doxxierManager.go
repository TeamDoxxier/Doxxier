package Foundation

import (
	"errors"

	"doxxier.tech/doxxier/pkg/models"
)

type DoxxierManager struct {
	Doxxier *models.Doxxier
}

func NewDoxxierManager() *DoxxierManager {
	return &DoxxierManager{
		Doxxier: models.NewDoxxier(),
	}
}

func (dm *DoxxierManager) AddPart(part models.DoxxierPart) *models.Doxxier {
	dm.Doxxier.Parts = append(dm.Doxxier.Parts, part)
	return dm.Doxxier
}

func (dm *DoxxierManager) SendDoxxier() error {
	if dm.Doxxier.Recipient == "" {
		return errors.New("Recipient not set")
	}
}
