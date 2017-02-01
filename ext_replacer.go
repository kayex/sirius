package sirius

import (
	"strings"
)

type Replacer struct{}

var phrases = map[string]string{
	"overwatch": "abovelook",
	"cancer":    "Does this dress make me look fat?",
}

func (*Replacer) Run(Message) (error, MessageAction) {
	edit := TextEdit()

	for search, replace := range phrases {
		edit.Substitute(strings.ToLower(search), replace)
	}

	return nil, edit
}
