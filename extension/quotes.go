package extension

import (
	"github.com/kayex/sirius"
	"strings"
)

type Quotes struct{}

func (*Quotes) Run(m sirius.Message) (error, sirius.MessageAction) {
	if strings.HasPrefix(m.Text, ">") {
		edit := sirius.TextEdit()

		edit.Substitute("\n", "\n>") // Don't break quote on line breaks

		return nil, edit
	}

	return nil, sirius.NoAction()
}
