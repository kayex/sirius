package sirius

import (
	"fmt"
	"github.com/nlopes/slack"
	"strings"
)

type SlackID struct {
	UserID string
	TeamID string
}

type Connection interface {
	Listen()
	Auth() chan SlackID
	Messages() chan Message
	Update(*Message) error
}

type RTMConnection struct {
	token    string
	rtm      *slack.RTM
	client   *slack.Client
	auth     chan SlackID
	messages chan Message
}

func NewSlackID(userID, teamID string) SlackID {
	return SlackID{
		UserID: userID,
		TeamID: teamID,
	}
}

func (id SlackID) Missing() bool {
	return id.UserID == "" && id.TeamID == ""
}

/*
Notice that user IDs are not guaranteed to be globally unique across all Slack users.
The combination of user ID and team ID, on the other hand, is guaranteed to be globally unique.

- Slack API documentation
*/
func (id SlackID) Equals(o SlackID) bool {
	return id.UserID == o.UserID && id.TeamID == o.TeamID
}

func NewRTMConnection(token string) *RTMConnection {
	client := slack.New(token)

	rtm := client.NewRTM()
	auth := make(chan SlackID)
	msg := make(chan Message)

	return &RTMConnection{
		rtm:      rtm,
		auth:     auth,
		messages: msg,
		client:   client,
		token:    token,
	}
}

func (conn *RTMConnection) Listen() {
	go conn.rtm.ManageConnection()

	for {
		select {
		case ev := <-conn.rtm.IncomingEvents:
			conn.handleIncomingEvent(ev)
		}
	}
}

func (conn *RTMConnection) Auth() chan SlackID {
	return conn.auth
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

func (conn *RTMConnection) authenticate(e *slack.ConnectedEvent) {
	info := e.Info
	id := NewSlackID(info.User.ID, info.Team.ID)
	conn.auth <- id
}

func (conn *RTMConnection) handleIncomingEvent(ev slack.RTMEvent) {
	switch msg := ev.Data.(type) {
	case *slack.ConnectedEvent:
		conn.authenticate(msg)
	case *slack.InvalidAuthEvent:
		panic(msg)
	case *slack.MessageEvent:
		conn.handleIncomingMessage(msg)

	case *slack.RTMError:
		fmt.Printf("Error: %s\n", msg.Error())
		panic(msg)
	}
}

func (conn *RTMConnection) handleIncomingMessage(ev *slack.MessageEvent) {
	// Drop messages with incomplete data
	if ev.User == "" || ev.Team == "" {
		return
	}

	text := removeEscapeCharacters(ev.Text)
	msg := NewMessage(SlackID{ev.User, ev.Team}, text, ev.Channel, ev.Timestamp)

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
