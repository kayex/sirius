package sirius

import (
	"strings"

	"github.com/kayex/sirius/slack"
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

type MessageQuery interface {
	Query(*Message) bool
}

func (m *Message) Query(q MessageQuery) bool {
	return q.Query(m)
}

// WholeWordQuery matches only complete words, i.e. strings that
// are not sub-strings of other words.
type FullWordQuery struct {
	W string
}

func (q FullWordQuery) Query(m *Message) bool {
	if m.Text == q.W {
		return true
	}
	// "W lorem ipsum"
	//  ^^
	if strings.HasPrefix(m.Text, q.W+" ") {
		return true
	}
	// "lorem ipsum W"
	//             ^^
	if strings.HasSuffix(m.Text, " "+q.W) {
		return true
	}
	// "lorem W ipsum"
	//       ^^^
	if strings.Contains(m.Text, " "+q.W+" ") {
		return true
	}

	return false
}
