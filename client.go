package sirius

import (
	"errors"
	"golang.org/x/net/context"
	"strings"
	"time"
)

type Client struct {
	user       *User
	conn       Connection
	extensions []Extension
	runner     ExtensionRunner
}

func NewClient(user *User, ext []Extension) *Client {
	conn := NewRTMConnection(user.Token)

	return &Client{
		conn:       conn,
		user:       user,
		extensions: ext,
		runner:     NewAsyncRunner(),
	}
}

func (c *Client) Start(ctx context.Context) {
	go c.conn.Listen()

	err := c.authenticate()

	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.conn.Messages():
			c.handleMessage(&msg)
		}
	}
}

func (c *Client) authenticate() error {
	for c.user.ID.Empty() {
		select {
		case id := <-c.conn.Auth():
			c.user.ID = id.Secure()
			return nil
		case <-time.After(time.Second * 3):
			return errors.New("Dynamic client authentication timed out (<-c.conn.Auth())")
		}
	}

	return nil
}

func (c *Client) handleMessage(msg *Message) {
	// We only care about outgoing messages
	if !c.sender(msg) {
		return
	}

	if msg.escaped() {
		edit := msg.EditText().ReplaceWith(trimEscape(msg.Text))
		msg.perform(edit)

		c.conn.Update(msg)
		return
	}

	c.run(msg)
}

func (c *Client) run(m *Message) {
	var exe []Execution
	var act []MessageAction

	for _, x := range c.extensions {
		exe = append(exe, *NewExecution(x, *m, ExtensionConfig{}))
	}

	res := make(chan ExecutionResult, len(c.extensions))

	c.runner.Run(exe, res, time.Second*2)

	for {
		if r, running := <-res; running {
			if r.Error != nil {
				panic(r.Error)
			}

			act = append(act, r.Action)
			continue
		}

		break
	}

	updated := c.performActions(act, m)

	if updated {
		c.conn.Update(m)
	}
}

func (c *Client) performActions(act []MessageAction, msg *Message) bool {
	var update bool

	for _, a := range act {
		err, modified := msg.perform(a)

		if err != nil {
			panic(err)
		}

		update = update || modified
	}

	return update
}

func (c *Client) sender(msg *Message) bool {
	return c.user.ID.Equals(msg.UserID)
}

func (m *Message) escaped() bool {
	return strings.HasPrefix(m.Text, `\`)
}

func trimEscape(text string) string {
	return strings.TrimPrefix(text, `\`)
}
