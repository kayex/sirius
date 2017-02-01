package sirius

type Message struct {
	Text      string
	Modified  bool
	UserID    string
	TeamID    string
	Channel   string
	Timestamp string
}

func NewMessage(text, user, team, channel, timestamp string) Message {
	return Message{
		Text:      text,
		UserID:    user,
		TeamID:    team,
		Channel:   channel,
		Timestamp: timestamp,
	}
}
