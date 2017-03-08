package text

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type Query interface {
	Match(string) int
}

type word struct {
	W string
}

// Word returns a query that matches the first occurrence of w in a search text
// where w is not immediately preceded or followed by any characters that do not satisfy
// isWordDelimiter.
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
	ir := utf8.RuneCountInString(s[:i])

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
