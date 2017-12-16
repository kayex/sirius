package sirius

import (
	"github.com/kayex/sirius/slack"
	"github.com/kayex/sirius/text"
	"strings"
)

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

func (m *Message) Query(q text.Query) bool {
	i, _ := q.Match(m.Text)
	return i >= 0
}

func (m *Message) sentBy(u *User) bool {
	return u.ID.Equals(m.UserID)
}

func (m *Message) escaped() bool {
	return strings.HasPrefix(m.Text, `\`)
}
