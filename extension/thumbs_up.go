package extension

import "github.com/kayex/sirius"

const slackThumb string = ":+1:"

type ThumbsUp struct{}

func (*ThumbsUp) Run(m sirius.Message) (error, sirius.MessageAction) {
	edit := m.EditText()
	edit.Substitute("(y)", slackThumb)
	edit.Substitute("(Y)", slackThumb)

	return nil, edit
}
