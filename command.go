package sirius

import "strings"

const prefix = `!`

type Command struct {
	Name string
	Args []string
}

// Arg returns argument number a, or nil if there is no argument in that
// position.
func (c *Command) Arg(a int) string {
	if len(c.Args) > a {
		return c.Args[a]
	}

	return ""
}

func (m *Message) Command(name string) (*Command, bool) {
	cmd := prefix + name

	if strings.HasPrefix(m.Text, cmd) {
		var args []string
		pieces := strings.Split(m.Text, " ")

		if len(pieces) >= 2 {
			args = append(args, pieces[1:]...)
		}

		return &Command{
			Name: name,
			Args: args,
		}, true
	}

	return nil, false
}
