package text

import (
	"strings"
	"unicode/utf8"
)

type Query interface {
	Match(string) int
}

// word matches complete words only.
type word struct {
	W string
}

func Word(w string) word {
	if len(w) == 0 {
		panic("Cannot create word of length 0")
	}

	return word{w}
}

func (q word) Match(s string) int {
	if len(s) == 0 {
		return -1
	}

	if s == q.W {
		return 0
	}

	i := strings.Index(s, q.W)
	if i < 0 {
		return -1
	}

	sr := []rune(s)
	ir := len(sr[:i])

	// Check for any disallowed surrounding characters
	prev := ir - 1
	next := ir + q.Length()
	if prev > 0 && !isWordSurroundRune(sr[prev]) ||
		next <= len(sr)-1 && !isWordSurroundRune(sr[next]) {
		return -1
	}

	return ir
}

func (q word) Length() int {
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
