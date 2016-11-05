package model

type Message struct {
	Text string
	Channel string
}

func NewMessage(text string, channel string) Message {
	return Message{
		Text: text,
		Channel: channel,
	}
}
