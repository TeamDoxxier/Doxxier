package Doxexier

import (
	"doxxier.tech/doxxier/models"
)

type DoxxierManager struct {
	Doxxier Doxxier
}

func InitialiseDoxxierManager() DoxxierManager {
	return DoxxierManager{
		Doxxier: models.InitialiseDoxxier(),
	}
}
