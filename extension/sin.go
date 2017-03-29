package extension

import "github.com/kayex/sirius"

type Sin struct{}

func (*Sin) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	return m.EditText().SubstituteWord("010", ":010:"), nil
}
