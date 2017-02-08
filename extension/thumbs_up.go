package extension

import "github.com/kayex/sirius"

const slackThumb string = ":+1:"

type ThumbsUp struct{}

func (*ThumbsUp) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	edit := m.EditText()
	edit.Substitute("(y)", slackThumb)
	edit.Substitute("(Y)", slackThumb)

	return edit, nil
}
