package extension

import (
	"strings"
)

type Transformation interface {
	Apply(text string) string
}

type SubstitutionTransformation struct {
	Search       string
	Substitution string
}

type AppendTransformation struct {
	Appendix string
}

func NoTransformation() []Transformation {
	return []Transformation{}
}

func Substitute(search string, substitution string) Transformation {
	return &SubstitutionTransformation{
		Search:       search,
		Substitution: substitution,
	}
}

func Append(appendix string) Transformation {
	return &AppendTransformation{
		Appendix: appendix,
	}
}

func (st *SubstitutionTransformation) Apply(text string) string {
	return strings.Replace(text, st.Search, st.Substitution, -1)
}

func (at *AppendTransformation) Apply(text string) string {
	if len(at.Appendix) == 0 {
		return text
	}

	return text + at.Appendix
}
