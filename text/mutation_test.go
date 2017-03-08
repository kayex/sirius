package text

import "testing"

func TestSubWord_Apply(t *testing.T) {
	cases := []struct {
		str string
		sw  SubWord
		exp string
	}{
		{
			str: "foo bar",
			sw:  SubWord{Word("bar"), "boo"},
			exp: "foo boo",
		},
		{
			str: "foo åäö",
			sw:  SubWord{Word("foo"), "bar"},
			exp: "bar åäö",
		},
		{
			str: "foo åäö",
			sw:  SubWord{Word("åäö"), "bar"},
			exp: "foo bar",
		},
		{
			str: "foo barbaz",
			sw:  SubWord{Word("foo"), "long replacement"},
			exp: "long replacement barbaz",
		},
	}

	for _, c := range cases {
		act := c.sw.Apply(c.str)

		if act != c.exp {
			t.Errorf("Expected SubWord{%q, %q}.Apply(%q) to return %q, got %q", c.sw.Search.W, c.sw.Sub, c.str, c.exp, act)
		}
	}
}
