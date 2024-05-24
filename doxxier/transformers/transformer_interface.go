package transformers

import (
	"doxxier.tech/doxxier/models"
)

type transformer interface {
	Transform(ctx *models.DoxxierContext) error
}
