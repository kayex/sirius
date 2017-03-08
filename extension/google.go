package extension

import (
	"bytes"
	"net/url"

	"github.com/kayex/sirius"
)

type Google struct{}

func (*Google) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	cmd, match := m.Command("g")

	if !match || len(cmd.Args) == 0 {
		return sirius.NoAction(), nil
	}

	var qbuf bytes.Buffer

	for i, a := range cmd.Args {
		if i != 0 {
			qbuf.WriteRune(' ')
		}
		qbuf.WriteString(a)
	}

	var urlb bytes.Buffer

	urlb.WriteString("https://www.google.com/search?q=")
	urlb.WriteString(url.QueryEscape(qbuf.String()))

	edit := m.EditText().Set(urlb.String())

	return edit, nil
}
