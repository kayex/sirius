package plugins

import (
	"github.com/kayex/sirius/core"
)

type plugin interface {
	Run(core.Message) string
}
