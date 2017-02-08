package sirius

/*
The MessageAction interface represents an action
that an extension wishes to perform on the
current message after execution has finished.

A MessageAction may return an error if it could
not be performed for any reason.
*/
type MessageAction interface {
	Perform(*Message) error
}

type EmptyAction struct{}

func NoAction() *EmptyAction {
	return &EmptyAction{}
}

func (*EmptyAction) Perform(*Message) error {
	return nil
}

/*
Applies a MessageAction to a message, returning
a bool indicating whether the message was actually
modified or not.
*/
func (m *Message) perform(a MessageAction) (error, bool) {
	oldText := m.Text
	err := a.Perform(m)

	if err != nil {
		return err, false
	}

	return nil, m.Text != oldText
}
