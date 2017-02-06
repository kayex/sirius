package extension

import (
	"github.com/kayex/sirius"
	"strings"
)

type Replacer struct{}

var phrases = map[string]string{
	"overwatch": "abovelook",
	"cancer":    "Does this dress make me look fat?",
}

func (*Replacer) Run(m sirius.Message) (error, sirius.MessageAction) {
	edit := m.EditText()

	for search, replace := range phrases {
		edit.Substitute(strings.ToLower(search), replace)
	}

	return nil, edit
}
