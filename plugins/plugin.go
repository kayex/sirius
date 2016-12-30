package plugins

import (
	"github.com/kayex/sirius/model"
)

type Plugin interface {
	Run(model.Message) []Transformation
}
