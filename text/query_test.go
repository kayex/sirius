package text

import (
	"testing"
)

func TestWord_Match(t *testing.T) {
	cases := []struct {
		s   string
		q   Word
		exp bool
	}{
		{
			s:   "Alligators eat mattresses",
			q:   Word{"Alligators"},
			exp: true,
		},
		{
			s:   "Alligators eat mattresses",
			q:   Word{"mattresses"},
			exp: true,
		},
		{
			s:   "Alligators eat mattresses",
			q:   Word{"gators"},
			exp: false,
		},
		{
			s:   "Alli\ngators eat mattresses",
			q:   Word{"gators"},
			exp: true,
		},
		{
			s:   "Alligators eat meat",
			q:   Word{"eat"},
			exp: true,
		},
	}

	for _, c := range cases {
		act := c.q.Match(c.s)

		if act != c.exp {
			t.Errorf("Expected Word(%q).Match(%q) to return %v, got %v", c.q.W, c.s, c.exp, act)
		}
	}
}

func BenchmarkWord_Match(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Word{"foo bar baz biz boz bem boos bick dale biz buul hum dirk hass tukk murr"}.Match("biz")
	}
}
