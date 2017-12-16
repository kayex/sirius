package extension

import (
	"github.com/kayex/sirius"
	"github.com/kayex/sirius/text"
)

type Replacer struct{}

var phrases = map[string]string{
	"overwatch": "abovelook",
	"cancer":    "Does this dress make me look fat?",
}

func (*Replacer) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	edit := m.EditText()

	for search, replace := range phrases {
		edit.SubstituteQuery(text.IWord(search), replace)
	}

	return edit, nil
}
