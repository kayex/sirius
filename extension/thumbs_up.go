package extension

import "github.com/kayex/sirius/model"

const slackThumb string = ":+1:"

type ThumbsUp struct{}

func (tu *ThumbsUp) Run(m model.Message) []Transformation {
	return []Transformation{
		Substitute("(y)", slackThumb),
		Substitute("(Y)", slackThumb),
	}
}
