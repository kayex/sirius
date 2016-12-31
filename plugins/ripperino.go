package plugins

import (
	"github.com/kayex/sirius/model"
	"math/rand"
	"strings"
	"fmt"
)

const base string = "ripperino"

var endings = []string{
	"casino",
	"linguini",
	"bambino",
	"ripperoni",
}

type Ripperino struct{}

func (r *Ripperino) Run(m model.Message) []Transformation {
	if !strings.HasPrefix(m.Text, base) {
		return NoTransformation()
	}

	// 1 in 10 times, go full Grino
	if rand.Int() % 10 == 1 {
		return []Transformation{rapperGrino()}
	}

	return []Transformation{Append(" " + getRandomEnding())}
}

func rapperGrino() Transformation {
	return Substitute(base, fmt.Sprintf("~%s~ RAPPER GRINO", base))
}

func getRandomEnding() string {
	e := rand.Int() % len(endings)
	return endings[e]
}
