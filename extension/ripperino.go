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

func (*Ripperino) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	if !strings.HasPrefix(m.Text, base) {
		return sirius.NoAction(), nil
	}

	edit := m.EditText()

	// 1 in 10 times, go full Grino
	if rand.Int()%10 == 1 {
		edit.Substitute(base, fmt.Sprintf("~%s~ RAPPER GRINO", base))
	} else {
		edit.Substitute(base, fmt.Sprintf("%v %s", base, getRandomEnding()))
	}

	return edit, nil
}

func getRandomEnding() string {
	e := rand.Int() % len(endings)
	return endings[e]
}
