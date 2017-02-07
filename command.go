package sirius

import "strings"

const commandPrefix = `<`

type Command struct {
	name string
}

func NewCommand(name string) *Command {
	return &Command{
		name: name,
	}
}

func (c *Command) Match(m *Message) (string, bool) {
	if !strings.HasPrefix(m.Text, commandPrefix) {
		return "", false
	}

	return strings.TrimPrefix(m.Text, commandPrefix), true
}
