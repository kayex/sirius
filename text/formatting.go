package text

import (
	"fmt"
	"strings"
)

func Italic(s string) string {
	return format("_%s_", strings.TrimSpace(s))
}

func Bold(s string) string {
	return format("*%s*", strings.TrimSpace(s))
}

func Strike(s string) string {
	return format("~%s~", strings.TrimSpace(s))
}

func Code(s string) string {
	return format("`%s`", strings.TrimSpace(s))
}

func Quote(s string) string {
	return format(">%s", s)
}

func format(f, s string) string {
	if s == "" {
		return s
	}
	return fmt.Sprintf(f, s)
}
