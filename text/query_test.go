package text

import (
	"testing"
)

func TestWord_Match(t *testing.T) {
	cases := []struct {
		s   string
		q   word
		exp int
	}{
		{
			s:   "Alligators eat mattresses",
			q:   Word("Alligators"),
			exp: 0,
		},
		{
			s:   "Alligators eat mattresses",
			q:   Word("mattresses"),
			exp: 15,
		},
		{
			s:   "Alligators eat mattresses",
			q:   Word("gators"),
			exp: -1,
		},
		{
			s:   "Alli\ngators eat mattresses",
			q:   Word("gators"),
			exp: 5,
		},
		{
			s:   "Alligators eat meat",
			q:   Word("eat"),
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

func BenchmarkWord_MatchExist(b *testing.B) {
	w := Word("foobar")

	txt := `Lorem ipsum dolor sit amet, an cum vero soleat concludaturque, te purto vero reprimique vis.
	Ignota mediocritatem ut sea. Cetero deserunt pericula te vel. Omnis legendos no per.
	Sale illum pertinax no sed, est posse putent minimum foobar no. Pri et vitae mentitum eligendi,
	no ius reque fugit libris, eos ad quaeque pericula mediocrem. Habemus corpora an mea,
	inermis partiendo per et, at nemore dolorem iudicabit eos. At est mucius docendi. Sed et nisl facilisi.
	Idque suavitate argumentum eu eam, vis putant insolens dissentiunt id. Dictas labitur in mei, duo omnium assentior scripserit cu.`

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		w.Match(txt)
	}
}

func BenchmarkWord_MatchNotExist(b *testing.B) {
	w := Word("foobar")

	txt := `Lorem ipsum dolor sit amet, an cum vero soleat concludaturque, te purto vero reprimique vis.
	Ignota mediocritatem ut sea. Cetero deserunt pericula te vel. Omnis legendos no per.
	Sale illum pertinax no sed, est posse putent minimum no. Pri et vitae mentitum eligendi,
	no ius reque fugit libris, eos ad quaeque pericula mediocrem. Habemus corpora an mea,
	inermis partiendo per et, at nemore dolorem iudicabit eos. At est mucius docendi. Sed et nisl facilisi.
	Idque suavitate argumentum eu eam, vis putant insolens dissentiunt id. Dictas labitur in mei, duo omnium assentior scripserit cu.`

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		w.Match(txt)
	}
}
