package sirius

type Message struct {
	Text      string
	Modified  bool
	UserID    string
	Channel   string
	Timestamp string
}

func NewMessage(text string, user string, channel string, timestamp string) Message {
	return Message{
		Text:      text,
		UserID:    user,
		Channel:   channel,
		Timestamp: timestamp,
	}
}
