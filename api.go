package sirius

import (
	"strings"

	"github.com/kayex/sirius/slack"
	api "github.com/nlopes/slack"
	"golang.org/x/net/context"
)

type Connection interface {
	Listen(context.Context)
	Auth() <-chan slack.UserID
	Send(*Message) error
	Update(*Message) error
	Messages() <-chan Message
	SelfChan() string
}

type RTMConnection struct {
	token    string
	rtm      *api.RTM
	client   *api.Client
	auth     chan slack.UserID
	messages chan Message
	selfChan string
}

func NewRTMConnection(token string) *RTMConnection {
	client := api.New(token)

	rtm := client.NewRTM()
	auth := make(chan slack.UserID, 1)
	msg := make(chan Message)

	return &RTMConnection{
		rtm:      rtm,
		auth:     auth,
		messages: msg,
		client:   client,
		token:    token,
	}
}

func (conn *RTMConnection) Listen(ctx context.Context) {
	go conn.rtm.ManageConnection()

	for {
		select {
		case <-ctx.Done():
			conn.rtm.Disconnect()
			return
		case ev := <-conn.rtm.IncomingEvents:
			conn.handleIncomingEvent(ev)
		}
	}
}

func (conn *RTMConnection) SelfChan() string {
	return conn.selfChan
}

func (conn *RTMConnection) Auth() <-chan slack.UserID {
	return conn.auth
}

func (conn *RTMConnection) Messages() <-chan Message {
	return conn.messages
}

func (conn *RTMConnection) Send(msg *Message) error {
	oMsg := conn.rtm.NewOutgoingMessage(msg.Text, msg.Channel)
	conn.rtm.SendMessage(oMsg)

	return nil
}

func (conn *RTMConnection) Update(msg *Message) error {
	_, _, _, err := conn.rtm.UpdateMessage(msg.Channel, msg.Timestamp, msg.Text)
	return err
}

func (conn *RTMConnection) authenticate(e *api.ConnectedEvent) {
	id := slack.UserID{e.Info.User.ID, e.Info.Team.ID}
	conn.selfChan = conn.getSelfChan(id, e)
	conn.auth <- id
}

func (conn *RTMConnection) getSelfChan(id slack.UserID, e *api.ConnectedEvent) string {
	for _, im := range e.Info.IMs {
		if im.User == id.UserID {
			return im.ID
		}
	}

	panic("")
}

func (conn *RTMConnection) handleIncomingEvent(ev api.RTMEvent) {
	switch msg := ev.Data.(type) {
	case *api.ConnectedEvent:
		conn.authenticate(msg)
	case *api.InvalidAuthEvent:
		panic(msg)
	case *api.MessageEvent:
		conn.handleIncomingMessage(msg)

	case *api.RTMError:
		panic(msg)
	}
}

func (conn *RTMConnection) handleIncomingMessage(ev *api.MessageEvent) {
	// Drop messages with incomplete data
	if ev.User == "" || ev.Team == "" {
		return
	}

	text := removeEscapeCharacters(ev.Text)
	msg := NewMessage(slack.UserID{ev.User, ev.Team}, text, ev.Channel, ev.Timestamp)

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
