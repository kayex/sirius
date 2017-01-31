package extension

import (
	"github.com/kayex/sirius/model"
)

type Extension interface {
	Run(model.Message) (error, []Transformation)
}
