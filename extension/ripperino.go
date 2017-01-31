package extension

import (
	"fmt"
	"github.com/kayex/sirius/model"
	"math/rand"
	"strings"
)

const base string = "ripperino"

var endings = []string{
	"casino",
	"linguini",
	"bambino",
	"ripperoni",
}

type Ripperino struct{}

func (r *Ripperino) Run(m model.Message) (error, []Transformation) {
	if !strings.HasPrefix(m.Text, base) {
		return nil, NoTransformation()
	}

	// 1 in 10 times, go full Grino
	if rand.Int()%10 == 1 {
		return nil, []Transformation{rapperGrino()}
	}

	return nil, []Transformation{Append(" " + getRandomEnding())}
}

func rapperGrino() Transformation {
	return Substitute(base, fmt.Sprintf("~%s~ RAPPER GRINO", base))
}

func getRandomEnding() string {
	e := rand.Int() % len(endings)
	return endings[e]
}
