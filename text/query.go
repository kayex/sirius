package text

import (
	"strings"
	"unicode/utf8"
)

type Query interface {
	Match(string) int
}

// Word matches complete words only.
type Word struct {
	W string
}

func (q Word) Match(s string) int {
	if len(s) == 0 {
		return -1
	} else if s == q.W {
		return 0
	}

	i := strings.Index(s, q.W)
	if i < 0 {
		return -1
	}

	sr := []rune(s)
	ir := len(sr[:i])

	// Check for any disallowed surrounding characters
	if (ir > 0 && !isWordSurroundRune(sr[ir-1])) ||
		(i+len(q.W) <= len(s)-1 && !isWordSurroundRune(sr[ir+utf8.RuneCountInString(q.W)])) {
		return -1
	}

	return ir
}

func (q Word) Length() int {
	return utf8.RuneCountInString(q.W)
}

func isWordSurroundRune(r rune) bool {
	surr := []rune{'\t', '\n', '\v', '\f', '\r', ' ', ',', '.'}

	for _, v := range surr {
		if r == v {
			return true
		}
	}

	return false
}
