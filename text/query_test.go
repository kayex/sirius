package text

import (
	"testing"
)

func TestWord_Match(t *testing.T) {
	cases := []struct {
		s   string
		q   Word
		exp int
	}{
		{
			s:   "Alligators eat mattresses",
			q:   Word{"Alligators"},
			exp: 0,
		},
		{
			s:   "Alligators eat mattresses",
			q:   Word{"mattresses"},
			exp: 15,
		},
		{
			s:   "Alligators eat mattresses",
			q:   Word{"gators"},
			exp: -1,
		},
		{
			s:   "Alli\ngators eat mattresses",
			q:   Word{"gators"},
			exp: 5,
		},
		{
			s:   "Alligators eat meat",
			q:   Word{"eat"},
			exp: 11,
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
	w := Word{"foo bar baz biz boz bem boos bick dale biz buul hum dirk hass tukk murr"}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		w.Match("biz")
	}
}
