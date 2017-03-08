package text

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type Query interface {
	Match(string) int
}

// word matches the first occurrence of W in a search text, where W is a string
// not immediately preceded or followed by any characters that do not satisfy
// isWordDelimiter.
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

	// Make sure that any preceding or following characters are valid
	// word delimiters.
	prev := ir - 1
	next := ir + q.Length()
	if prev > 0 && !isWordDelimiter(sr[prev]) ||
		next <= len(sr)-1 && !isWordDelimiter(sr[next]) {
		return -1
	}

	return ir
}

func (q word) Length() int {
	return utf8.RuneCountInString(q.W)
}

// isWordDelimiter indicates if r is a word delimiter.
func isWordDelimiter(r rune) bool {
	if unicode.IsSpace(r) {
		return true
	}

	switch r {
	case ',', '.':
		return true
	}

	return false
}
