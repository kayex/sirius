package sirius

import "strings"

const prefix = `!`

type Command struct {
	name string
}

func NewCommand(name string) *Command {
	return &Command{
		name: name,
	}
}

func (c *Command) Match(m *Message) (string, bool) {
	command := prefix + c.name + " "
	if !strings.HasPrefix(m.Text, command) {
		return "", false
	}

	return strings.TrimPrefix(m.Text, command), true
}
