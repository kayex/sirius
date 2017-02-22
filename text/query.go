package text

import "unicode"

type Query interface {
	Match(string) bool
}

// Word matches complete words only.
// The "complete words" of a string s is defined as the result of
// splitting the string on every single Unicode whitespace character.
type Word struct {
	W string
}

func (q Word) Match(s string) bool {
	if s == q.W {
		return true
	}

	sr := []rune(s)
	qr := []rune(q.W)

	if len(sr) < len(qr) {
		return false
	}

	var nMatch int

	for i := 0; i < len(sr); i++ {
		if nMatch < len(qr) && sr[i] != qr[nMatch] {
			nMatch = 0
			continue
		}

		nMatch++

		if nMatch == len(qr) {
			// Check that any immediately preceding or following
			// characters are spaces.

			next := i + 1
			prev := i - nMatch
			hasNext := len(sr) > next
			hasPrev := i-nMatch >= 0

			if hasNext && !unicode.IsSpace(sr[next]) || hasPrev && !unicode.IsSpace(sr[prev]) {
				nMatch = 0
				continue
			}

			return true
		}
	}

	return false
}
