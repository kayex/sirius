package extension

import (
	"github.com/kayex/sirius"
	"github.com/kayex/sirius/text"
)

const slackThumb = ":+1:"

type ThumbsUp struct{}

func (*ThumbsUp) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	return m.EditText().SubstituteQuery(text.IWord("(y)"), slackThumb), nil
}
