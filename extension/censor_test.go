package extension

import (
	"reflect"
	"testing"

	"github.com/kayex/sirius"
	"github.com/kayex/sirius/slack"
)

func TestCensor_Run(t *testing.T) {
	cases := []struct {
		msg sirius.Message
		cfg sirius.ExtensionConfig
		exp sirius.MessageAction
	}{
		{
			msg: sirius.NewMessage(slack.UserID{"123", "abc"}, "Voldemort", "#channel", "0"),
			cfg: sirius.ExtensionConfig{
				"phrases": []string{"Voldemort"},
				"strict":  false,
			},
			exp: (&sirius.TextEditAction{}).SubstituteWord("Voldemort", "`CENSORED`"),
		},
		{
			msg: sirius.NewMessage(slack.UserID{"123", "abc"}, "Voldemort", "#channel", "0"),
			cfg: sirius.ExtensionConfig{
				"phrases": []string{"Voldemort"},
				"strict":  true,
			},
			exp: (&sirius.TextEditAction{}).Set("`CENSORED`"),
		},
		{
			msg: sirius.NewMessage(slack.UserID{"123", "abc"}, "Rainbows", "#channel", "0"),
			cfg: sirius.ExtensionConfig{
				"phrases": []string{"Voldemort"},
				"strict":  true,
			},
			exp: &sirius.TextEditAction{},
		},
		{
			msg: sirius.NewMessage(slack.UserID{"123", "abc"}, "Dickens", "#channel", "0"),
			cfg: sirius.ExtensionConfig{
				"phrases": []string{"Dick"},
				"strict":  true,
			},
			exp: &sirius.TextEditAction{},
		},
	}

	for _, c := range cases {
		cs := &Censor{}

		act, err := cs.Run(c.msg, c.cfg)

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(act, c.exp) {
			t.Errorf("Expected %s but got %#v for message %q", c.exp, act, c.msg.Text)
		}
	}
}
