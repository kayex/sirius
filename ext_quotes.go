package sirius

import "strings"

type Quotes struct {}

func (*Quotes) Run(m Message) (error, MessageAction) {
	if strings.HasPrefix(m.Text, ">") {
		edit := TextEdit()

		edit.Substitute("\n", "\n>") // Don't break quote on line breaks

		return nil, edit
	}

	return nil, NoAction()
}
