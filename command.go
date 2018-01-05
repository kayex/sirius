package sirius

import "strings"

const prefix = `!`

type Command struct {
	Name string
	Args []string
}

// Arg returns argument number a, or "" if there is
// no argument in that position.
func (c *Command) Arg(a int) string {
	if len(c.Args) > a {
		return c.Args[a]
	}

	return ""
}

func (c *Command) CollapseArgs() string {
	return strings.Join(c.Args, " ")
}

func (m *Message) Command(name string) (*Command, bool) {
	cmd := prefix + name

	if strings.HasPrefix(m.Text, cmd) {
		var args []string
		pieces := strings.Split(m.Text, " ")

		if len(pieces) > 1 {
			args = append(args, pieces[1:]...)
		}

		return &Command{
			Name: name,
			Args: args,
		}, true
	}

	return nil, false
}
