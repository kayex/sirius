package extension

import "github.com/kayex/sirius"

const slackThumb string = ":+1:"

type ThumbsUp struct{}

func (tu *ThumbsUp) Run(m sirius.Message) (error, []Transformation) {
	return nil, []Transformation{
		Substitute("(y)", slackThumb),
		Substitute("(Y)", slackThumb),
	}
}
