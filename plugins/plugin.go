package plugins

import (
	"github.com/Epoch2/slack-sirius/core"
)

type plugin interface {
	Run(core.Message) string
}
