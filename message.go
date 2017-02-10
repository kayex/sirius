package sirius

import "github.com/kayex/sirius/slack"

type Message struct {
	Text      string
	UserID    slack.UserID
	Channel   string
	Timestamp string
}

func NewMessage(userID slack.UserID, text, channel, timestamp string) Message {
	return Message{
		Text:      text,
		UserID:    userID,
		Channel:   channel,
		Timestamp: timestamp,
	}
}
