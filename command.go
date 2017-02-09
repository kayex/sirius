package sirius

import "strings"

const prefix = `!`

type Command struct {
	Name string
	Args []string
}

func (m *Message) Command(name string) (*Command, bool) {
	cmd := prefix + name + " "

	if strings.HasPrefix(m.Text, cmd) {
		trim := strings.TrimPrefix(m.Text, cmd)
		args := strings.Split(trim, " ")

		return &Command{
			Name: name,
			Args: args,
		}, true
	}

	return nil, false
}
