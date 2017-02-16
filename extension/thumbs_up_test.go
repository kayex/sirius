package extension

import (
	"github.com/kayex/sirius"
	"github.com/kayex/sirius/slack"
	"reflect"
	"testing"
)

func TestThumbsUp_Run(t *testing.T) {
	cases := []struct {
		msg sirius.Message
		exp sirius.MessageAction
	}{
		{
			msg: sirius.NewMessage(slack.UserID{UserID: "123", TeamID: "abc"}, "(y)", "#channel", "0"),
			exp: (&sirius.TextEditAction{}).Substitute("(y)", ":+1:").Substitute("(Y)", ":+1:"),
		},
		{
			msg: sirius.NewMessage(slack.UserID{UserID: "123", TeamID: "abc"}, "Hej!", "#channel", "0"),
			exp: (&sirius.TextEditAction{}).Substitute("(y)", ":+1:").Substitute("(Y)", ":+1:"),
		},
	}

	for _, c := range cases {
		tu := &ThumbsUp{}

		act, err := tu.Run(c.msg, sirius.ExtensionConfig{})

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(act, c.exp) {
			t.Fatalf("Expected %s but got %s for message %q", c.exp, act, c.msg.Text)
		}
	}
}
