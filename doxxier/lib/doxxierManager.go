package lib

import (
	"errors"

	"doxxier.tech/doxxier/pkg/models"
)

type DoxxierManager struct {
	doxxier *models.Doxxier
}

func NewDoxxierManager() *DoxxierManager {
	return &DoxxierManager{
		doxxier: models.NewDoxxier(),
	}
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
