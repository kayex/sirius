package sirius

type Message struct {
	Text      string
	Modified  bool
	UserID    SlackID
	Channel   string
	Timestamp string
}

func NewMessage(userID SlackID, text, channel, timestamp string) Message {
	return Message{
		Text:      text,
		UserID:    userID,
		Channel:   channel,
		Timestamp: timestamp,
	}
}
