package sirius

import (
	"errors"
	"fmt"
	"github.com/nlopes/slack"
	"strings"
)

type SlackID struct {
	UserID string
	TeamID string
}

type Connection interface {
	ID() (error, SlackID)
	Listen()
	Messages() chan Message
	Update(*Message) error
}

type RTMConnection struct {
	rtm           *slack.RTM
	id            SlackID
	messages      chan Message
	authenticated bool
	client        *slack.Client
	token         string
}

func NewSlackID(userID, teamID string) SlackID {
	return SlackID{
		UserID: userID,
		TeamID: teamID,
	}
}

/*
Notice that user IDs are not guaranteed to be globally unique across all Slack users.
The combination of user ID and team ID, on the other hand, is guaranteed to be globally unique.

- Slack API documentation
*/
func (s *SlackID) Equals(o *SlackID) bool {
	return s.UserID == o.UserID && s.TeamID == o.TeamID
}

func NewRTMConnection(token string) *RTMConnection {
	client := slack.New(token)

	rtm := client.NewRTM()
	msg := make(chan Message)

	return &RTMConnection{
		rtm:      rtm,
		messages: msg,
		client:   client,
		token:    token,
	}
}

func (conn *RTMConnection) ID() (error, SlackID) {
	if !conn.authenticated {
		return errors.New("Cannot get UserID before authentication completes"), SlackID{}
	}

	return nil, conn.id
}

func (conn *RTMConnection) Listen() {
	go conn.rtm.ManageConnection()

	for !conn.authenticated {
		select {
		case ev := <-conn.rtm.IncomingEvents:
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
		case ev := <-conn.rtm.IncomingEvents:
			conn.handleIncomingEvent(ev)
		}
	}
}

func (conn *RTMConnection) Messages() chan Message {
	return conn.messages
}

func (conn *RTMConnection) SendMessage(msg *Message) {
	omsg := conn.rtm.NewOutgoingMessage(msg.Text, msg.Channel)
	conn.rtm.SendMessage(omsg)
}

func (conn *RTMConnection) Update(msg *Message) error {
	_, _, _, err := conn.rtm.UpdateMessage(msg.Channel, msg.Timestamp, msg.Text)
	return err
}

func (conn *RTMConnection) authenticate(e *slack.ConnectedEvent) bool {
	info := e.Info
	conn.id = NewSlackID(info.User.ID, info.Team.ID)
	conn.authenticated = true

	return true
}

func (conn *RTMConnection) handleIncomingEvent(ev slack.RTMEvent) {
	switch msg := ev.Data.(type) {
	case *slack.MessageEvent:
		conn.handleIncomingMessage(msg)

	case *slack.RTMError:
		fmt.Printf("Error: %s\n", msg.Error())
		panic(msg)
	}
}

func (conn *RTMConnection) handleIncomingMessage(ev *slack.MessageEvent) {
	text := removeEscapeCharacters(ev.Text)
	msg := NewMessage(text, ev.User, ev.Team, ev.Channel, ev.Timestamp)
	conn.messages <- msg
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
