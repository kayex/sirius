package extension

import (
	"bytes"
	"net/url"

	"github.com/kayex/sirius"
)

type Google struct{}

func (*Google) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	cmd, match := m.Command("g")

	if !match {
		return sirius.NoAction(), nil
	}

	q := cmd.Arg(0)

	if q == "" {
		return sirius.NoAction(), nil
	}

	var urlb bytes.Buffer

	urlb.WriteString("https://www.google.com/search?q=")
	urlb.WriteString(url.QueryEscape(q))

	edit := m.EditText().Set(urlb.String())

	return edit, nil
}
