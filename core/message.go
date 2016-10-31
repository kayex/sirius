package core

type Message struct {
	Text string
}

func NewMessage(text string) Message {
	return Message{
		Text: text,
	}
}
