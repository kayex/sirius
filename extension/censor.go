package extension

import (
	"strings"

	"github.com/kayex/sirius"
	"github.com/kayex/sirius/slack"
)

type Censor struct{}

func (*Censor) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	strict := cfg.Boolean("strict")
	phrases := cfg.List("phrases")

	edit := m.EditText()

	for _, p := range phrases {
		if !strings.Contains(m.Text, p) {
			continue
		}

		if strict {
			return edit.ReplaceWith(slack.Code("CENSORED")), nil
		}

		edit.Substitute(p, slack.Code("CENSORED"))
	}

	return edit, nil
}
