package sirius

import (
	"github.com/kayex/sirius/slack"
	"testing"
)

func TestMessage_Command(t *testing.T) {
	cases := []struct {
		msg    Message
		search string
		match  bool
		cmd    *Command
	}{
		{
			msg:    NewMessage(slack.UserID{UserID: "123", TeamID: "abc"}, "!foo bar", "#channel", "0"),
			search: "foo",
			match:  true,
			cmd: &Command{
				Name: "foo",
				Args: []string{
					"bar",
				},
			},
		},
		{
			msg:    NewMessage(slack.UserID{UserID: "123", TeamID: "abc"}, "!foo", "#channel", "0"),
			search: "foo",
			match:  true,
			cmd: &Command{
				Name: "foo",
				Args: []string{},
			},
		},
		{
			msg:    NewMessage(slack.UserID{UserID: "123", TeamID: "abc"}, "!foo bar baz", "#channel", "0"),
			search: "foo",
			match:  true,
			cmd: &Command{
				Name: "foo",
				Args: []string{
					"bar",
					"baz",
				},
			},
		},
		{
			msg:    NewMessage(slack.UserID{UserID: "123", TeamID: "abc"}, "!foo", "#channel", "0"),
			search: "bar",
			match:  false,
			cmd:    nil,
		},
	}

	for _, c := range cases {
		cmd, match := c.msg.Command(c.search)

		if match != c.match {
			t.Fatalf("Expected Command(%q) to return (%#v, %v), got (%#v, %v)", c.msg.Text, c.cmd, c.match, cmd, match)
		}

		if c.cmd == nil {
			return
		}

		if len(cmd.Args) != len(c.cmd.Args) {
			t.Fatalf("Expected len(Command(%q).Args) to be %d, got %d", c.msg.Text, len(c.cmd.Args), len(cmd.Args))
		}

		for i := 0; i < len(c.cmd.Args); i++ {
			expArg := c.cmd.Args[i]
			actArg := cmd.Args[i]

			if actArg != expArg {
				t.Errorf("Expected Command(%q).Args[%d] to be %q, got %q", c.msg.Text, i, expArg, actArg)
			}
		}
	}
}
