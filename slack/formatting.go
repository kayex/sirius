package slack

import (
	"fmt"
	"strings"
)

func Italic(s string) string {
	if s == "" {
		return s
	}
	return fmt.Sprintf("_%s_", strings.TrimSpace(s))
}

func Bold(s string) string {
	if s == "" {
		return s
	}

	return fmt.Sprintf("*%s*", strings.TrimSpace(s))
}

func Strike(s string) string {
	if s == "" {
		return s
	}

	return fmt.Sprintf("~%s~", strings.TrimSpace(s))
}

func Code(s string) string {
	if s == "" {
		return s
	}

	return fmt.Sprintf("`%s`", strings.TrimSpace(s))
}

func Quote(s string) string {
	if s == "" {
		return s
	}

	return fmt.Sprintf(">%s", s)
}
