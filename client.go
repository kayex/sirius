package sirius

import (
	"golang.org/x/net/context"
	"time"
	"strings"
)

type Client struct {
	user *User
	conn  *Connection
}

func NewClient(user *User) *Client {
	conn := NewConnection(user.Token)

	return &Client{
		conn:  &conn,
		user: user,
	}
}

func (c *Client) Start(ctx context.Context) {
	go c.conn.Listen()

	for {
		select {
		case msg := <-c.conn.Incoming:
			c.handleMessage(&msg)
		}
	}
}

func (c *Client) handleMessage(msg *Message) {
	isSender := c.conn.UserId == msg.UserID

	if !isSender {
		return
	}

	if msg.escaped() {
		msg.Text = trimEscape(msg.Text)
		c.conn.Update(msg)
	}

	act := c.runExtensions(msg)
	c.applyActions(act, msg)
}

func (c *Client) runExtensions(msg *Message) []MessageAction {
	cfgs := c.user.Configurations
	act := make(chan MessageAction, len(cfgs))

	for _, cfg := range cfgs {
		ext := LoadExtension(cfg.EID)

		execute(ext, msg, act)
	}

	var actions []MessageAction

	ActionReceive:
	for range cfgs {
		select {
		case a := <-act:
			actions = append(actions, a)

		// Allow extensions max 200ms to execute and provide an actionable result
		case <-time.After(time.Millisecond * 200):
			break ActionReceive
		}
	}

	return actions
}

func (c *Client) applyActions(act []MessageAction, msg *Message) {
	oldText := msg.Text

	for _, a := range act {
		err := a.Perform(msg)

		if err != nil {
			panic(err)
		}
	}

	if msg.Text != oldText {
		c.conn.Update(msg)
	}
}

func (m *Message) escaped() bool {
	return strings.HasPrefix(m.Text, `\`)
}

func trimEscape(text string) string {
	return strings.TrimPrefix(text, `\`)
}

/*
Executes ext(msg) and passes the results onto act
 */
func execute(ext Extension, msg *Message, act chan<- MessageAction) {
	go func() {
		err, a := ext.Run(*msg)

		if err != nil {
			panic(err)
		}

		act <- a
	}()
}
