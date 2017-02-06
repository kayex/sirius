package extension

import "github.com/kayex/sirius"

const slackThumb string = ":+1:"

type ThumbsUp struct{}

func (*ThumbsUp) Run(sirius.Message, sirius.ExtensionConfig) (error, sirius.MessageAction) {
	edit := sirius.TextEdit()
	edit.Substitute("(y)", slackThumb)
	edit.Substitute("(Y)", slackThumb)

	return nil, edit
}
