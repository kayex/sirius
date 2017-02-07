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
	fullCommandName := commandPrefix + c.name + " "
	if !strings.HasPrefix(m.Text, fullCommandName) {
		return "", false
	}

	return strings.TrimPrefix(m.Text, fullCommandName), true
}
