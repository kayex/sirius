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

// WordQuery matches the first occurrence of w in a search text where w is
// surrounded by word delimiters (as defined by isWordDelimiter).
type WordQuery struct {
	W string
}

func Word(w string) WordQuery {
	if len(w) == 0 {
		panic("Cannot create WordQuery of length 0")
	}

	return WordQuery{w}
}

func (q WordQuery) Match(s string) int {
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
	// WordQuery delimiters.
	prev, p := at(sr, ir-1)
	next, n := at(sr, ir+q.Length())
	if p && !isWordDelimiter(prev) || n && !isWordDelimiter(next) {
		return -1
	}

	return ir
}

func (q WordQuery) Length() int {
	return utf8.RuneCountInString(q.W)
}

// IgnoreCaseQuery provides case-insensitive matching by folding the search string
// before passing it to the wrapped query.
type IgnoreCaseQuery struct {
	Query
}

func (q IgnoreCaseQuery) Match(s string) int {
	sl := strings.ToLower(s)

	return q.Query.Match(sl)
}

func IgnoreCase(q Query) IgnoreCaseQuery {
	return IgnoreCaseQuery{q}
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

// at returns the rune at index i, and a bool indicating if i exists in s.
func at(s []rune, i int) (rune, bool) {
	if i < 0 || i > len(s)-1 {
		return 0, false
	}

	r := s[i]
	return r, true
}
