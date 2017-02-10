package extension

import (
	"github.com/kayex/sirius"
	"strings"
)

type Quotes struct{}

func (*Quotes) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	if strings.HasPrefix(m.Text, ">") {
		edit := m.EditText().Substitute("\n", "\n>") // Don't break quote on line breaks

		return edit, nil
	}

	return sirius.NoAction(), nil
}
