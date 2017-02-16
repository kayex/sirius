package slack

import (
	"fmt"
	"strings"
)

func Italic(s string) string {
	return fmt.Sprintf("_%s_", strings.TrimSpace(s))
}

func Bold(s string) string {
	return fmt.Sprintf("*%s*", s)
}

func Strike(s string) string {
	return fmt.Sprintf("~%s~", s)
}

func Code(s string) string {
	return fmt.Sprintf("`%s`", s)
}

func Quote(s string) string {
	return fmt.Sprintf(">%s", s)
}
