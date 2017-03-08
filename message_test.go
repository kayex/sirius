package sirius

import (
	"github.com/kayex/sirius/slack"
	"github.com/kayex/sirius/text"
	"testing"
)

func TestMessage_Query(t *testing.T) {
	cases := []struct {
		msg Message
		q   text.Query
		exp bool
	}{
		{
			msg: NewMessage(slack.UserID{"123", "abc"}, "Alligators eat mattresses", "#channel", "0"),
			q:   text.Word{"Alligators"},
			exp: true,
		},
		{
			msg: NewMessage(slack.UserID{"123", "abc"}, "Alligators eat mattresses", "#channel", "0"),
			q:   text.Word{"mattresses"},
			exp: true,
		},
		{
			msg: NewMessage(slack.UserID{"123", "abc"}, "Alligators eat mattresses", "#channel", "0"),
			q:   text.Word{"gators"},
			exp: false,
		},
		{
			msg: NewMessage(slack.UserID{"123", "abc"}, "Alligators eat meat", "#channel", "0"),
			q:   text.Word{"eat"},
			exp: true,
		},
	}

	for _, c := range cases {
		act := c.msg.Query(c.q)

		if act != c.exp {
			t.Errorf("Expected Message.Query(%#v) for message %q to return %v, got %v", c.q, c.msg.Text, c.exp, act)
		}
	}
}
