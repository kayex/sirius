package plugins

import (
	"github.com/kayex/sirius/model"
	"math/rand"
	"strings"
)

var endings = []string{
	"casino",
	"linguini",
	"bambino",
	"ripperoni",
}

type Ripperino struct{}

func (r *Ripperino) Run(m model.Message) []Transformation {
	if !strings.HasPrefix(m.Text, "ripperino") {
		return NoTransformation()
	}

	// 1 in 10 times, go full Grino
	if rand.Int() % 10 == 1 {
		return []Transformation{rapperGrino()}
	}

	return []Transformation{Append(" " + getRandomEnding())}
}

func rapperGrino() Transformation {
	return Substitute("ripperino", "~ripperino~ RAPPER GRINO")
}

func getRandomEnding() string {
	e := rand.Int() % len(endings)
	return endings[e]
}
