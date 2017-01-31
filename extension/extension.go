package extension

import (
	"github.com/kayex/sirius"
)

type Extension interface {
	Run(sirius.Message) (error, []Transformation)
}
