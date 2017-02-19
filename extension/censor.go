package extension

import (
	"github.com/kayex/sirius"
	"github.com/kayex/sirius/slack"
	"github.com/kayex/sirius/text"
)

type Censor struct{}

func (*Censor) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	strict := cfg.Boolean("strict")
	phrases := cfg.List("phrases")

	edit := m.EditText()

	for _, p := range phrases {
		if !m.Query(text.Word{p}) {
			continue
		}

		if strict {
			return edit.ReplaceWith(slack.Code("CENSORED")), nil
		}

		edit.SubstituteWord(p, slack.Code("CENSORED"))
	}

	return edit, nil
}
