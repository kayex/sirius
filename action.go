package sirius

/*
MessageAction represents an action that an extension
wishes to perform on the current message after
execution has finished.
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

//perform applies a to m
//Returns a bool indicating whether m was actually modified
func (m *Message) perform(a MessageAction) (err error, mod bool) {
	oldText := m.Text
	err = a.Perform(m)
	mod = m.Text != oldText

	return
}
