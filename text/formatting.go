package text

import (
	"fmt"
	"strings"
)

func Italic(s string) string {
	return slackFmt("_%s_", strings.TrimSpace(s))
}

func Bold(s string) string {
	return slackFmt("*%s*", strings.TrimSpace(s))
}

func Strike(s string) string {
	return slackFmt("~%s~", strings.TrimSpace(s))
}

func Code(s string) string {
	return slackFmt("`%s`", strings.TrimSpace(s))
}

func Quote(s string) string {
	return slackFmt(">%s", s)
}

func slackFmt(f, s string) string {
	if s == "" {
		return s
	}
	return fmt.Sprintf(f, s)
}
