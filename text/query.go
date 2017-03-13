package text

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type Query interface {
	Match(string) int
}

// Token represents any type of text token in a search text.
type Token interface {
	// Length returns the length of the token in runes.
	Length() int
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
	prev := previous(sr, ir)
	next := next(sr, ir, &q)
	if prev != nil && !isWordDelimiter(*prev) ||
		next != nil && !isWordDelimiter(*next) {
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

// previous returns a pointer to the rune immediately preceding index i
// in the search text s, or nil if it does not exist.
func previous(s []rune, i int) *rune {
	prev := i - 1
	if prev < 0 {
		return nil
	}

	r := s[prev]
	return &r
}

// next returns a pointer to the rune immediately following the token t in the
// search string s, or nil if it doesn't exist.
func next(s []rune, i int, t Token) *rune {
	next := i + t.Length()
	if next > len(s)-1 {
		return nil
	}

	r := s[next]
	return &r
}
