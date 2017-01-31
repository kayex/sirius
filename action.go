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

func (na *EmptyAction) Perform(*Message) error {
	return nil
}
