package sirius

import (
	"strings"
)

type Replacer struct{}

var words = map[string]string{
	"overwatch": "abovelook",
}

func (r *Replacer) Run(m Message) (error, MessageAction) {
	edit := TextEdit()

	for s, r := range words {
		edit.Substitute(strings.ToLower(s), r)
	}

	return nil, edit
}
