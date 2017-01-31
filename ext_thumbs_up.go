package sirius

const slackThumb string = ":+1:"

type ThumbsUp struct{}

func (tu *ThumbsUp) Run(m Message) (error, MessageAction) {
	edit := TextEdit()
	edit.Substitute("(y)", slackThumb)
	edit.Substitute("(Y)", slackThumb)

	return nil, edit
}
