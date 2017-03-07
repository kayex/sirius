package text

import "strings"

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

	for {
		i := s.Search.Match(text)

		if i < 0 {
			return text
		}

		beginning := text[:i]
		end := text[i+len(s.Search.W):]

		text = beginning + s.Sub + end
	}

}

func (am *Append) Apply(text string) string {
	if len(am.Appendix) == 0 {
		return text
	}

	return text + am.Appendix
}

func (pm *Prepend) Apply(text string) string {
	if len(pm.Prefix) == 0 {
		return text
	}

	return pm.Prefix + text
}
