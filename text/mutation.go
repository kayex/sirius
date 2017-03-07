package text

import (
	"strings"
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

	var tr []rune
	for i := s.Search.Match(text); i >= 0; i = s.Search.Match(text) {
		if tr == nil {
			tr = []rune(text)
		}

		beginning := tr[:i]
		end := tr[i+s.Search.Length():]

		tr = append(append(beginning, []rune(s.Sub)...), end...)

		text = string(tr)
	}

	return text
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
