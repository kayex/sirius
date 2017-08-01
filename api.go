package sirius

import (
	"fmt"
	"strings"

	"github.com/kayex/sirius/slack"
	api "github.com/nlopes/slack"
)

type API interface {
	GetUserID(token string) (slack.ID, error)
}

type MessageBroker interface {
	Send(*Message) error
	Update(*Message) error
	Messages() <-chan Message
}

type Connection interface {
	API
	MessageBroker
	Listen()
	Close()
	Finished() <-chan bool
	Auth() <-chan slack.UserID
	Details() ConnectionDetails
}

type ConnectionDetails struct {
	UserID   slack.UserID
	SelfChan string
}

type RTMConnection struct {
	token    string
	rtm      *api.RTM
	client   *api.Client
	details  ConnectionDetails
	auth     chan slack.UserID
	messages chan Message
	stop     chan bool
	finished chan bool
}

func NewRTMConnection(token string) *RTMConnection {
	client := api.New(token)

	return &RTMConnection{
		token:    token,
		rtm:      client.NewRTM(),
		client:   client,
		auth:     make(chan slack.UserID, 1),
		messages: make(chan Message),
		stop:     make(chan bool, 1),
	}
}

func (conn *RTMConnection) Listen() {
	go conn.rtm.ManageConnection()

	for {
		select {
		case <-conn.stop:
			return
		case ev := <-conn.rtm.IncomingEvents:
			conn.handleIncomingEvent(ev)
		}
	}
}

func (conn *RTMConnection) Close() {
	err := conn.rtm.Disconnect()

	if err != nil {
	}

	conn.stop <- true
}

func (conn *RTMConnection) Finished() <-chan bool {
	return conn.finished
}

func (conn *RTMConnection) Auth() <-chan slack.UserID {
	return conn.auth
}

func (conn *RTMConnection) Details() ConnectionDetails {
	return conn.details
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

func (conn *RTMConnection) GetUserID(token string) (slack.ID, error) {
	res, err := conn.client.AuthTest()

	if err != nil {
		return nil, err
	}

	id := &slack.UserID{
		UserID: res.UserID,
		TeamID: res.TeamID,
	}
	conn.details.UserID = *id

	return id, nil
}

func (conn *RTMConnection) authenticate(e *api.ConnectedEvent) {
	id := slack.UserID{e.Info.User.ID, e.Info.Team.ID}
	selfChan, err := conn.getSelfChan(id, e)
	if err != nil {
		panic(err)
	}
	conn.details.SelfChan = selfChan
	conn.details.UserID = id
	conn.auth <- id
}

func (conn *RTMConnection) getSelfChan(id slack.UserID, e *api.ConnectedEvent) (string, error) {
	for _, im := range e.Info.IMs {
		if im.User == id.UserID {
			return im.ID, nil
		}
	}

	return "", fmt.Errorf("Could not find self-channel for User(%v)", id.String())
}

func (conn *RTMConnection) handleIncomingEvent(ev api.RTMEvent) {
	switch msg := ev.Data.(type) {
	case *api.ConnectedEvent:
		conn.authenticate(msg)
	case *api.InvalidAuthEvent:
		panic(msg)
	case *api.MessageEvent:
		conn.handleIncomingMessage(msg)

	case *api.DisconnectedEvent:
		conn.Close()
	case *api.RTMError:
		panic(msg)
	}
}

func (conn *RTMConnection) handleIncomingMessage(ev *api.MessageEvent) {
	id := slack.UserID{ev.User, ev.Team}
	if !id.Valid() {
		return
	}

	text := stripEscapeCharacters(ev.Text)
	msg := NewMessage(id, text, ev.Channel, ev.Timestamp)

	conn.messages <- msg
}

var escapeCharacters map[string]string = map[string]string{
	"&lt;":  "<",
	"&gt;":  ">",
	"&amp;": "&",
}

func stripEscapeCharacters(msg string) string {
	for ec, unicode := range escapeCharacters {
		msg = strings.Replace(msg, ec, unicode, -1)
	}

	return msg
}
