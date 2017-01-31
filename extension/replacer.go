package extension

import (
	"strings"
	"github.com/kayex/sirius"
)

type Replacer struct{}

var words = map[string]string{
	"overwatch": "abovelook",
}

func (r *Replacer) Run(m sirius.Message) (error, []Transformation) {
	trans := []Transformation{}

	for s, r := range words {
		trans = append(trans, Substitute(strings.ToLower(s), r))
	}

	return nil, trans
}
