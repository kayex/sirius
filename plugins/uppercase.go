package plugins

import (
	"github.com/kayex/sirius/core"
	"strings"
)

type uppercase_plugin struct {
}

func (u *uppercase_plugin) Run(msg core.Message) string {
	return strings.ToUpper(msg.Text)
}

func NewUppercasePlugin() uppercase_plugin {
	return uppercase_plugin{}
}
