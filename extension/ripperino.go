package extension

import (
	"fmt"
	"github.com/kayex/sirius"
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

func (*Ripperino) Run(m sirius.Message) (error, sirius.MessageAction) {
	if !strings.HasPrefix(m.Text, base) {
		return nil, sirius.NoAction()
	}

	edit := sirius.TextEdit()

	// 1 in 10 times, go full Grino
	if rand.Int()%10 == 1 {
		edit.Substitute(base, fmt.Sprintf("~%s~ RAPPER GRINO", base))
	} else {
		edit.Append(" " + getRandomEnding())
	}

	return nil, edit
}

func getRandomEnding() string {
	e := rand.Int() % len(endings)
	return endings[e]
}
