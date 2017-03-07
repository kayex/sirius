package text

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type Query interface {
	Match(string) int
}

// Word matches complete words only.
// The "complete words" of a string s is defined as the result of
// splitting the string on every single Unicode whitespace character.
type Word struct {
	W string
}

func (q Word) Match(s string) int {
	if s == q.W || len(s) == 0 {
		return 0
	}

	sp := strings.FieldsFunc(s, func(r rune) bool {
		if unicode.IsSpace(r) {
			return true
		}

		switch r {
		case ',', '\n':
			return true
		default:
			return false
		}
	})

	var rCount int
	for i, w := range sp {
		if w == q.W {
			return rCount + i
		}
		rCount += utf8.RuneCountInString(w)
	}

	return -1
}

func (q Word) Length() int {
	return utf8.RuneCountInString(q.W)
}
