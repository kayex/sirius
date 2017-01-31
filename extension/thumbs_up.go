package extension

import "github.com/kayex/sirius/model"

const slackThumb string = ":+1:"

type ThumbsUp struct{}

func (tu *ThumbsUp) Run(m model.Message) (error, []Transformation) {
	return nil, []Transformation{
		Substitute("(y)", slackThumb),
		Substitute("(Y)", slackThumb),
	}
}
