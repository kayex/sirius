package sirius

import (
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"strings"
)

type Connection struct {
	ID       SlackID
	Rtm      *slack.RTM
	Incoming chan Message
	ready    bool
	client   *slack.Client
	token    string
}

func CInit(logger *log.Logger) {
	slack.SetLogger(logger)
}

func NewConnection(token string) Connection {
	client := slack.New(token)

	rtm := client.NewRTM()
	inc := make(chan Message)

	return Connection{
		Rtm:      rtm,
		Incoming: inc,
		client:   client,
		token:    token,
	}
}

func (conn *Connection) Listen() {
	go conn.Rtm.ManageConnection()

	for !conn.ready {
		select {
		case ev := <-conn.Rtm.IncomingEvents:
			switch msg := ev.Data.(type) {
			case *slack.InvalidAuthEvent:
				panic(msg)
			case *slack.ConnectedEvent:
				conn.authenticate(msg)
				conn.handleIncomingEvent(ev)
			}
		}
	}

	for {
		select {
		case ev := <-conn.Rtm.IncomingEvents:
			conn.handleIncomingEvent(ev)
		}
	}
}

func (conn *Connection) SendMessage(msg *Message) {
	omsg := conn.Rtm.NewOutgoingMessage(msg.Text, msg.Channel)
	conn.Rtm.SendMessage(omsg)
}

func (conn *Connection) Update(msg *Message) error {
	_, _, _, err := conn.Rtm.UpdateMessage(msg.Channel, msg.Timestamp, msg.Text)
	return err
}

func (conn *Connection) authenticate(e *slack.ConnectedEvent) bool {
	info := e.Info
	conn.ID = NewSlackID(info.User.ID, info.Team.ID)
	conn.ready = true

	return true
}

func (conn *Connection) handleIncomingEvent(ev slack.RTMEvent) {
	switch msg := ev.Data.(type) {
	case *slack.MessageEvent:
		conn.handleIncomingMessage(msg)

	case *slack.RTMError:
		fmt.Printf("Error: %s\n", msg.Error())
		panic(msg)
	}
}

func (conn *Connection) handleIncomingMessage(ev *slack.MessageEvent) {
	text := removeEscapeCharacters(ev.Text)
	msg := NewMessage(text, ev.User, ev.Team, ev.Channel, ev.Timestamp)
	conn.Incoming <- msg
}

var escapeCharacters map[string]string = map[string]string{
	"&lt;":  "<",
	"&gt;":  ">",
	"&amp;": "&",
}

func removeEscapeCharacters(msg string) string {
	for html, unicode := range escapeCharacters {
		msg = strings.Replace(msg, html, unicode, -1)
	}

	return msg
}
