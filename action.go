package sirius

// MessageAction represents an action that an extension
// wishes to perform on the current message after
// execution has finished.
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

// alter modifies a message by applying a MessageAction on it.
// Returns a bool indicating whether the message text property was modified
// by the action.
func (m *Message) alter(a MessageAction) (bool, error) {
	oldText := m.Text
	err := a.Perform(m)
	mod := m.Text != oldText

	return mod, err
}

// alterAll modifies a message by applying a series of MessageActions on it.
// Returns a bool indicating whether the message text property was modified
// by the action.
func (m *Message) alterAll(act []MessageAction) (bool, error) {
	var modified bool
	for _, a := range act {
		mod, err := m.alter(a)
		if err != nil {
			return modified, err
		}

		modified = modified || mod
	}

	return modified, nil
}
