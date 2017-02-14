package slack

import "testing"

func TestFormattingItalic(t *testing.T) {
	s := "Hello"
	exp := "_Hello_"
	act := italic(s)

	if act != exp {
		t.Fatalf("Expected italic(%s) to be (%s), got (%s)", s, exp, act)
	}
}

func TestFormattingBold(t *testing.T) {
	s := "Hello"
	exp := "*Hello*"
	act := bold(s)

	if act != exp {
		t.Fatalf("Expected bold(%s) to be (%s), got (%s)", s, exp, act)
	}
}

func TestFormattingStrike(t *testing.T) {
	s := "Hello"
	exp := "~Hello~"
	act := strike(s)

	if act != exp {
		t.Fatalf("Expected strike(%s) to be (%s), got (%s)", s, exp, act)
	}
}

func TestFormattingCode(t *testing.T) {
	s := "Hello"
	exp := "`Hello`"
	act := code(s)

	if act != exp {
		t.Fatalf("Expected code(%s) to be (%s), got (%)", s, exp, act)
	}
}
