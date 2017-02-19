package text

import "strings"

type Query interface {
	Match(string) bool
}

// Word matches only complete words, i.e. strings that
// are not sub-strings of other words.
type Word struct {
	W string
}

func (q Word) Match(s string) bool {
	if s == q.W {
		return true
	}
	// "W lorem ipsum"
	//  ^^
	if strings.HasPrefix(s, q.W+" ") {
		return true
	}
	// "lorem ipsum W"
	//             ^^
	if strings.HasSuffix(s, " "+q.W) {
		return true
	}
	// "lorem W ipsum"
	//       ^^^
	if strings.Contains(s, " "+q.W+" ") {
		return true
	}

	return false
}
