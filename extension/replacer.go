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

func (*Replacer) Run(sirius.Message) (error, sirius.MessageAction) {
	edit := sirius.TextEdit()

	for search, replace := range phrases {
		edit.Substitute(strings.ToLower(search), replace)
	}

	return nil, edit
}
