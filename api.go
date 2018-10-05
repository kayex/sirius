package sirius

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kayex/sirius/slack"
	"github.com/kayex/sirius/sync"
	api "github.com/nlopes/slack"
)

// API is a connection to the Slack API.
type API interface {
	GetUserID(token string) (slack.ID, error)
}

// Connection is a connection to the Slack Messaging services.
type Connection interface {
	Start(ctx context.Context) error
	Send(*Message) error
	Update(*Message) error
	Messages() <-chan Message
	Details() ConnectionDetails

	// Closed returns a channel that can be used for detecting when the connection closes.
	// If the connection was closed due to an error, the error is returned. Otherwise, nil is returned.
	Closed() <-chan error
}

// ConnectionDetails are details that cannot be discerned until the client
// has performed authentication against Slack.
type ConnectionDetails struct {
	UserID   slack.UserID
	SelfChan string
}

// RTMConnection is a Connection that utilizes the Slack RTM API.
type RTMConnection struct {
	rtm      *api.RTM
	details  *ConnectionDetails
	messages chan Message
	ctrl     *sync.Control
}

type SlackAPI struct {
	client *api.Client
}

func NewSlackAPI(token string) *SlackAPI {
	client := api.New(token)

	return &SlackAPI{
		client: client,
	}
}

func (api *SlackAPI) GetUserID(token string) (slack.ID, error) {
	res, err := api.client.AuthTest()

	if err != nil {
		return nil, err
	}

	id := &slack.UserID{
		UserID: res.UserID,
		TeamID: res.TeamID,
	}

	return id, nil
}

func (api *SlackAPI) NewRTMConnection() *RTMConnection {
	return &RTMConnection{
		rtm:      api.client.NewRTM(),
		messages: make(chan Message),
		ctrl:     sync.NewControl(),
	}
}

func (conn *RTMConnection) Start(ctx context.Context) error {
	auth := make(chan ConnectionDetails)
	go conn.listen(ctx, auth)

	select {
	case d := <-auth:
		conn.details = &d
	case err := <-conn.Closed():
		return fmt.Errorf("RTM connection failed: %v", err)
	}

	return nil
}

// listen opens the RTM connection and listens for incoming events until
// ctx expires or an error occurs.
// Sends the connection details discerned from the initial authentication
// process on the auth channel.
func (conn *RTMConnection) listen(ctx context.Context, auth chan ConnectionDetails) {
	go conn.rtm.ManageConnection()
	defer conn.rtm.Disconnect()
	defer conn.ctrl.Finish(nil)

	d, err := conn.authenticate()
	if err != nil {
		conn.ctrl.Finish(err)
		close(auth)
		return
	}

	auth <- *d

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-conn.rtm.IncomingEvents:
			err := conn.handleEvent(ev)
			if err != nil {
				conn.ctrl.Finish(err)
				return
			}
		}
	}
}

func (conn *RTMConnection) handleEvent(ev api.RTMEvent) error {
	switch msg := ev.Data.(type) {
	case *api.MessageEvent:
		conn.handleIncomingMessage(msg)
		return nil
	case *api.DisconnectedEvent:
		return errors.New("received DisconnectedEvent")
	case *api.RTMError:
		return msg
	}

	return nil
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

func (conn *RTMConnection) Send(msg *Message) error {
	oMsg := conn.rtm.NewOutgoingMessage(msg.Text, msg.Channel)
	conn.rtm.SendMessage(oMsg)

	return nil
}

func (conn *RTMConnection) Update(msg *Message) error {
	_, _, _, err := conn.rtm.UpdateMessage(msg.Channel, msg.Timestamp, msg.Text)
	return err
}

func (conn *RTMConnection) Messages() <-chan Message {
	return conn.messages
}

func (conn *RTMConnection) Details() ConnectionDetails {
	return *conn.details
}

func (conn *RTMConnection) Closed() <-chan error {
	return conn.ctrl.Finished()
}

// rtmAuthTimeout is the maximum amount of time to wait for the first
// authentication event to arrive.
const rtmAuthTimeout = time.Second * 3

// authenticate authenticates the current connection by listening on the RTM
// connection for authentication events.
//
// Returns the connection details discerned during the authentication process.
func (conn *RTMConnection) authenticate() (*ConnectionDetails, error) {
	for {
		select {
		case <-time.After(rtmAuthTimeout):
			return nil, fmt.Errorf("authentication timed out after waiting %v seconds on first event", rtmAuthTimeout/time.Second)
		case ev := <-conn.rtm.IncomingEvents:
			switch msg := ev.Data.(type) {

			case *api.ConnectedEvent:
				id := slack.UserID{msg.Info.User.ID, msg.Info.Team.ID}
				selfChan, err := getSelfChan(id, msg)
				if err != nil {
					return nil, err
				}

				details := &ConnectionDetails{
					UserID:   id,
					SelfChan: selfChan,
				}

				return details, nil
			case *api.InvalidAuthEvent:
				return nil, errors.New("connection error \"invalid_auth\"")
			case *api.ConnectingEvent:
				continue
			default:
				return nil, fmt.Errorf("unexpected event during authentication: %+v", msg)
			}
		}
	}
}

func getSelfChan(id slack.UserID, ev *api.ConnectedEvent) (string, error) {
	for _, im := range ev.Info.IMs {
		if im.User == id.UserID {
			return im.ID, nil
		}
	}

	return "", fmt.Errorf("could not find self-channel for User(%v)", id.String())
}

var escapeCharacters = map[string]string{
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
