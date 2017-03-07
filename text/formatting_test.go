package text

import "testing"

func TestFormattingItalic(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"Hello", "_Hello_"},
		{" Hello", "_Hello_"}, // Leading spaces break inline formatting within Slack.
		{"  Hello", "_Hello_"},
		{"Hello ", "_Hello_"}, // Trailing spaces too.
		{"Hello  ", "_Hello_"},
		{"  Hello  ", "_Hello_"},
		{"", ""},
	}

	for _, c := range cases {
		act := Italic(c.in)

		if act != c.out {
			t.Errorf("Expected Italic(%q) to be %q, got %q", c.in, c.out, act)
		}
	}
}

func TestFormattingBold(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"Hello", "*Hello*"},
		{" Hello", "*Hello*"},
		{"  Hello", "*Hello*"},
		{"Hello ", "*Hello*"},
		{"Hello  ", "*Hello*"},
		{"  Hello  ", "*Hello*"},
		{"", ""},
	}

	for _, c := range cases {
		act := Bold(c.in)

		if act != c.out {
			t.Errorf("Expected Bold(%q) to be %q, got %q", c.in, c.out, act)
		}
	}
}

func TestFormattingStrike(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"Hello", "~Hello~"},
		{" Hello", "~Hello~"},
		{"  Hello", "~Hello~"},
		{"Hello ", "~Hello~"},
		{"Hello  ", "~Hello~"},
		{"  Hello  ", "~Hello~"},
		{"", ""},
	}

	for _, c := range cases {
		act := Strike(c.in)

		if act != c.out {
			t.Errorf("Expected Strike(%q) to be %q, got %q", c.in, c.out, act)
		}
	}
}

func TestFormattingCode(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"Hello", "`Hello`"},
		{" Hello", "`Hello`"},
		{"  Hello", "`Hello`"},
		{"Hello ", "`Hello`"},
		{"Hello  ", "`Hello`"},
		{"  Hello  ", "`Hello`"},
		{"", ""},
	}

	for _, c := range cases {
		act := Code(c.in)

		if act != c.out {
			t.Errorf("Expected Code(%q) to be %q, got %q", c.in, c.out, act)
		}
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
