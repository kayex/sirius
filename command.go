package sirius

import "strings"

const prefix = `!`

type Command struct {
	Name string
	Args []string
}

func (m *Message) Command(name string) (*Command, bool) {
	cmd := prefix + name

	if strings.HasPrefix(m.Text, cmd) {
		var args []string
		inv := strings.Split(m.Text, " ")

		if len(inv) >= 2 {
			args = append(args, inv[1:]...)
		}

		return &Command{
			Name: name,
			Args: args,
		}, true
	}

	return nil, false
}
