package text

import (
	"testing"
)

func TestWordQuery_Match(t *testing.T) {
	cases := []struct {
		s   string
		q   WordQuery
		exp int
	}{
		{
			s:   "foo",
			q:   Word("foo"),
			exp: 0,
		},
		{
			s:   "foo bar",
			q:   Word("bar"),
			exp: 4,
		},
		{
			s:   "foo barbaz",
			q:   Word("bar"),
			exp: -1,
		},
		{
			s:   "foobar baz",
			q:   Word("bar"),
			exp: -1,
		},
		{
			s:   "foo bar",
			q:   Word("FOO"),
			exp: -1,
		},
		{
			s:   "FOO BAR",
			q:   Word("foo"),
			exp: -1,
		},
		{
			s:   "foo\tbar",
			q:   Word("bar"),
			exp: 4,
		},
		{
			s:   "foo\nbar",
			q:   Word("bar"),
			exp: 4,
		},
		{
			s:   "foo\vbar",
			q:   Word("bar"),
			exp: 4,
		},
		{
			s:   "foo\fbar",
			q:   Word("bar"),
			exp: 4,
		},
		{
			s:   "foo\rbar",
			q:   Word("bar"),
			exp: 4,
		},
		{
			s:   "åäö bar",
			q:   Word("bar"),
			exp: 4,
		},
	}

	for _, c := range cases {
		act := c.q.Match(c.s)

		if act != c.exp {
			t.Errorf("Expected Word(%q).Match(%q) to return %v, got %v", c.q.W, c.s, c.exp, act)
		}
	}
}

func TestCaseInsensitiveWordQuery_Match(t *testing.T) {
	cases := []struct {
		s   string
		q   CaseInsensitiveWordQuery
		exp int
	}{
		{
			s:   "foo",
			q:   IWord("foo"),
			exp: 0,
		},
		{
			s:   "FOO",
			q:   IWord("foo"),
			exp: 0,
		},
	}

	for _, c := range cases {
		act := c.q.Match(c.s)

		if act != c.exp {
			t.Errorf("Expected Lower(%#v).Match(%q) to return %v, got %v", c.q, c.s, c.exp, act)
		}
	}
}

// BenchmarkWord_MatchNotExist6_587 benchmarks a single Word query of length
// 6 against a search text of length 587, where the sought string does not
// exist in the search text.
//
// This benchmark gives a good indication of the average performance of an
// unsuccessful search.
func BenchmarkWord_MatchNotExist6_587(b *testing.B) {
	w := Word("foobar")

	txt := `Lorem ipsum dolor sit amet, an cum vero soleat concludaturque, te purto vero reprimique vis.
	Ignota mediocritatem ut sea. Cetero deserunt pericula te vel. Omnis legendos no per.
	Sale illum pertinax no sed, est posse putent minimum no. Pri et vitae mentitum eligendi,
	no ius reque fugit libris, eos ad quaeque pericula mediocrem. Habemus corpora an mea,
	inermis partiendo per et, at nemore dolorem iudicabit eos. At est mucius docendi. Sed et nisl facilisi.
	Idque suavitate argumentum eu eam, vis putant insolens dissentiunt id. Dictas labitur in mei, duo omnium assentior scripserit cu omnium`

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		w.Match(txt)
	}
}

// BenchmarkWord_MatchPartials5_587 benchmarks a single Word query of length
// 6 against a search text of length 587 with partials of length 5.
//
// This benchmark gives a good indication of the worst case performance
// of an unsuccessful search.
func BenchmarkWord_MatchPartials5_587(b *testing.B) {
	w := Word("foobar")

	txt := `fooab fooba fooba fooba fooba fooba fooba fooba fooba fooba
		fooba fooba fooba fooba fooba fooba fooba fooba fooba fooba
		fooba fooba fooba fooba fooba fooba fooba fooba fooba fooba
		fooba fooba fooba fooba fooba fooba fooba fooba fooba fooba
		fooba fooba fooba fooba fooba fooba fooba fooba fooba fooba
		fooba fooba fooba fooba fooba fooba fooba fooba fooba fooba
		fooba fooba fooba fooba fooba fooba fooba fooba fooba fooba
		fooba fooba fooba fooba fooba fooba fooba fooba fooba fooba
		fooba fooba fooba fooba fooba fooba fooba fooba fooba fooba
		fooba fooba fooba fooba fooba fooba fooba fooba`

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		w.Match(txt)
	}
}

// BenchmarkWord_MatchExist6_587 benchmarks a single Word query of length
// 6 against a search text of length 587, where the sought string is at the
// very end of the search text.
func BenchmarkWord_MatchExist6_587(b *testing.B) {
	w := Word("foobar")

	txt := `Lorem ipsum dolor sit amet, an cum vero soleat concludaturque, te purto vero reprimique vis.
	Ignota mediocritatem ut sea. Cetero deserunt pericula te vel. Omnis legendos no per.
	Sale illum pertinax no sed, est posse putent minimum no. Pri et vitae mentitum eligendi,
	no ius reque fugit libris, eos ad quaeque pericula mediocrem. Habemus corpora an mea,
	inermis partiendo per et, at nemore dolorem iudicabit eos. At est mucius docendi. Sed et nisl facilisi.
	Idque suavitate argumentum eu eam, vis putant insolens dissentiunt id. Dictas labitur in mei, duo omnium assentior scripserit cu foobar`

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		w.Match(txt)
	}
}
