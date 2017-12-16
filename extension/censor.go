package extension

import (
	"github.com/kayex/sirius"
	"github.com/kayex/sirius/text"
)

type Censor struct{}

func (*Censor) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	strict := cfg.Boolean("strict")
	phrases := cfg.List("phrases")

	edit := m.EditText()

	for _, p := range phrases {
		if !m.Query(text.IWord(p)) {
			continue
		}

		if strict {
			edit.Set(text.Code("CENSORED"))
			break
		}

		if cfg.Boolean("strike") {
			edit.SubstituteWord(p, text.Strike(p))
		} else {
			edit.SubstituteWord(p, text.Code("CENSORED"))
		}
	}

	return edit, nil
}
