package extension

import (
	"github.com/kayex/sirius"
	"github.com/kayex/sirius/text"
)

type Sin struct{}

func (*Sin) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	return m.EditText().SubstituteQuery(text.Word("010"), ":010:"), nil
}
