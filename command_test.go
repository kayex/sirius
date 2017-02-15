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
			t.Fatalf("Expected Command(%v) for message \"%v\" to return (%#v, %v), got (%#v, %v)", c.search, c.msg.Text, c.cmd, c.match, cmd, match)
		}

		if c.cmd == nil {
			return
		}

		if len(cmd.Args) != len(c.cmd.Args) {
			t.Fatalf("Expected %d arguments from Command(%v) for message \"%v\", got %d", len(c.cmd.Args), c.search, c.msg.Text, len(cmd.Args))
		}

		for i := 0; i < len(c.cmd.Args); i++ {
			expArg := c.cmd.Args[i]
			actArg := cmd.Args[i]

			if actArg != expArg {
				t.Fatalf("Expected \"%s\" as argument number %v for message \"%v\", got \"%s\"", expArg, i, c.msg.Text, actArg)
			}
		}
	}
}
