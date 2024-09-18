package transformers

import (
	"doxxier.tech/doxxier/pkg/models"
)

type transformer interface {
	Transform(ctx *models.DoxxierPart) error
}
