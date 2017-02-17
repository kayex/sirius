package sirius

import (
	"github.com/kayex/sirius/slack"
	"testing"
)

func TestFullWordQuery_Query(t *testing.T) {
	cases := []struct {
		msg Message
		q   MessageQuery
		exp bool
	}{
		{
			msg: NewMessage(slack.UserID{"123", "abc"}, "Alligators eat mattresses", "#channel", "0"),
			q:   FullWordQuery{"Alligators"},
			exp: true,
		},
		{
			msg: NewMessage(slack.UserID{"123", "abc"}, "Alligators eat mattresses", "#channel", "0"),
			q:   FullWordQuery{"mattresses"},
			exp: true,
		},
		{
			msg: NewMessage(slack.UserID{"123", "abc"}, "Alligators eat mattresses", "#channel", "0"),
			q:   FullWordQuery{"gators"},
			exp: false,
		},
		{
			msg: NewMessage(slack.UserID{"123", "abc"}, "Alligators eat meat", "#channel", "0"),
			q:   FullWordQuery{"eat"},
			exp: true,
		},
	}

	for _, c := range cases {
		act := c.msg.Query(c.q)

		if act != c.exp {
			t.Errorf("Expected FullWordQuery(%q) for message %q to return %v, got %v", c.q, c.msg.Text, c.exp, act)
		}
	}
}
