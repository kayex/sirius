package sirius

import (
	"testing"

	"github.com/kayex/sirius/slack"
)

func TestPerform(t *testing.T) {
	cases := []struct {
		act  MessageAction
		msg  Message
		mod  bool   // Expected value of mod.
		text string // The expected value of the msg.Text property after act has been performed.
	}{
		{
			act:  (&TextEditAction{}).Substitute("Foo", "Bar"),
			msg:  NewMessage(slack.UserID{UserID: "123", TeamID: "abc"}, "Foo", "#channel", "0"),
			mod:  true,
			text: "Bar",
		},
		{
			act:  (&TextEditAction{}).Substitute("Your reality", "My own"),
			msg:  NewMessage(slack.UserID{UserID: "123", TeamID: "abc"}, "Foo", "#channel", "0"),
			mod:  false,
			text: "Foo",
		},
		{
			act:  &TextEditAction{},
			msg:  NewMessage(slack.UserID{UserID: "123", TeamID: "abc"}, "Foo", "#channel", "0"),
			mod:  false,
			text: "Foo",
		},
		{
			act:  (&TextEditAction{}).Substitute("", ""),
			msg:  NewMessage(slack.UserID{UserID: "123", TeamID: "abc"}, "", "#channel", "0"),
			mod:  false,
			text: "",
		},
	}

	for _, c := range cases {
		oText := c.msg.Text

		err, mod := c.msg.perform(c.act)
		if err != nil {
			panic(err)
		}

		if mod != c.mod {
			t.Fatalf("Expected perform(Substitute(\"Foo\", \"Bar\")) to return (<nil>, true) for message %q, got (%#v, %v)", c.msg.Text, err, mod)
		}

		if c.msg.Text != c.text {
			t.Fatalf("Expected perform(Substitute(\"Foo\", \"Bar\")) to mutate message %q into %q, got %q", oText, c.text, c.msg.Text)
		}
	}
}
