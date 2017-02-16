package slack

import "testing"

func TestFormattingItalic(t *testing.T) {
	s := "Hello"
	exp := "_Hello_"
	act := Italic(s)

	if act != exp {
		t.Errorf("Expected Italic(%q) to be %q, got %q", s, exp, act)
	}
}

func TestFormattingBold(t *testing.T) {
	s := "Hello"
	exp := "*Hello*"
	act := Bold(s)

	if act != exp {
		t.Errorf("Expected Bold(%q) to be %q, got %q", s, exp, act)
	}
}

func TestFormattingStrike(t *testing.T) {
	s := "Hello"
	exp := "~Hello~"
	act := Strike(s)

	if act != exp {
		t.Errorf("Expected Strike(%q) to be %q, got %q", s, exp, act)
	}
}

func TestFormattingCode(t *testing.T) {
	s := "Hello"
	exp := "`Hello`"
	act := Code(s)

	if act != exp {
		t.Errorf("Expected Code(%q) to be %q, got %q", s, exp, act)
	}
}

func TestFormattingQuote(t *testing.T) {
	s := "Hello"
	exp := ">Hello"
	act := Quote(s)

	if act != exp {
		t.Errorf("Expected Quote(%q) to be %q, got %q", s, exp, act)
	}
}
