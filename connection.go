package sirius

import (
	"fmt"
	"github.com/nlopes/slack"
	"log"
)

type Connection struct {
	UserID   string
	TeamID   string
	Rtm      *slack.RTM
	Incoming chan Message
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

func (conn *Connection) handleIncomingEvent(ev slack.RTMEvent) {
	switch msg := ev.Data.(type) {
	case *slack.ConnectedEvent:
		conn.UserID = msg.Info.User.ID
		conn.TeamID = msg.Info.Team.ID

	case *slack.MessageEvent:
		conn.handleIncomingMessage(msg)

	case *slack.RTMError:
		fmt.Printf("Error: %s\n", msg.Error())
		panic(msg)

	case *slack.InvalidAuthEvent:
		panic(msg)
	}

}

func (conn *Connection) handleIncomingMessage(ev *slack.MessageEvent) {
	msg := NewMessage(ev.Text, ev.User, ev.Team, ev.Channel, ev.Timestamp)
	conn.Incoming <- msg
}
