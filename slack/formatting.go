package slack

import "fmt"

func italic(s string) string {
	return fmt.Sprintf("_%v_", s)
}

func bold(s string) string {
	return fmt.Sprintf("*%v*", s)
}

func strike(s string) string {
	return fmt.Sprintf("~%v~", s)
}

func code(s string) string {
	return fmt.Sprintf("`%v`", s)
}
