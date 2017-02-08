package sirius

import (
	"errors"
	"golang.org/x/net/context"
	"strings"
	"time"
)

type Client struct {
	user   *User
	conn   Connection
	loader ExtensionLoader
	runner ExtensionRunner
}

func NewClient(user *User, loader ExtensionLoader) *Client {
	conn := NewRTMConnection(user.Token)

	return &Client{
		conn:   conn,
		user:   user,
		loader: loader,
		runner: NewAsyncRunner(),
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
	if !msg.sentBy(c.user) {
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
	var act []MessageAction

	exe := c.loadExecutions(m)
	res := make(chan ExecutionResult, len(c.user.Configurations))

	c.runner.Run(exe, res, time.Second*2)

	for {
		if r, running := <-res; running {
			if r.Error != nil {
				panic(r.Error)
			}

			act = append(act, r.Action)
		} else {
			break
		}
	}

	updated := performActions(act, m)

	if updated {
		c.conn.Update(m)
	}
}

func (c *Client) loadExecutions(m *Message) []Execution {
	var exe []Execution

	for _, cf := range c.user.Configurations {
		x, err := c.loader.Load(cf.EID)

		if err != nil {
			panic(err)
		}

		exe = append(exe, *NewExecution(x, *m, cf.Config))
	}

	return exe
}

func (m *Message) sentBy(u *User) bool {
	return m.UserID.Secure().Equals(u.ID)
}

func (m *Message) escaped() bool {
	return strings.HasPrefix(m.Text, `\`)
}

func trimEscape(text string) string {
	return strings.TrimPrefix(text, `\`)
}

func performActions(act []MessageAction, msg *Message) bool {
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
