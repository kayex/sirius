package sirius

const slackThumb string = ":+1:"

type ThumbsUp struct{}

func (*ThumbsUp) Run(Message) (error, MessageAction) {
	edit := TextEdit()
	edit.Substitute("(y)", slackThumb)
	edit.Substitute("(Y)", slackThumb)

	return nil, edit
}
