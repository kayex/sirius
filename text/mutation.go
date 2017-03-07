package text

import (
	"strings"
	"unicode/utf8"
)

type Mutation interface {
	Apply(text string) string
}

type Replace struct {
	Replacement string
}

type Sub struct {
	Search string
	Sub    string
}

type SubWord struct {
	Search Word
	Sub    string
}

type Append struct {
	Appendix string
}

type Prepend struct {
	Prefix string
}

func (r *Replace) Apply(text string) string {
	return r.Replacement
}

func (s *Sub) Apply(text string) string {
	return strings.Replace(text, s.Search, s.Sub, -1)
}

func (s *SubWord) Apply(text string) string {
	if text == s.Search.W {
		return s.Sub
	}

	var sr []rune
	for {
		i := s.Search.Match(text)

		if i < 0 {
			return text
		}

		if sr == nil {
			sr = []rune(text)
		}

		beginning := sr[:i]
		end := sr[i+utf8.RuneCountInString(s.Search.W):]

		sr = append(append(beginning, []rune(s.Sub)...), end...)

		text = string(sr)
	}

}

func (a *Append) Apply(text string) string {
	if len(a.Appendix) == 0 {
		return text
	}

	return text + a.Appendix
}

func (p *Prepend) Apply(text string) string {
	if len(p.Prefix) == 0 {
		return text
	}

	return p.Prefix + text
}
