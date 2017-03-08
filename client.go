package sirius

import (
	"errors"
	"strings"
	"time"

	"context"
)

type Client struct {
	Ready   chan bool
	user    *User
	conn    Connection
	loader  ExtensionLoader
	runner  ExtensionRunner
	timeout time.Duration
}

type ClientConfig struct {
	user    *User
	loader  ExtensionLoader
	runner  ExtensionRunner
	timeout time.Duration
}

func NewClient(cfg ClientConfig) *Client {
	cl := &Client{
		conn:    NewRTMConnection(cfg.user.Token),
		user:    cfg.user,
		loader:  cfg.loader,
		runner:  cfg.runner,
		timeout: cfg.timeout,
		Ready:   make(chan bool, 1),
	}
	if cl.runner == nil {
		cl.runner = NewAsyncRunner()
	}
	if cl.timeout == 0 {
		cl.timeout = time.Second * 2
	}
	return cl
}

func (c *Client) Start(ctx context.Context) {
	go c.conn.Listen(ctx)

	err := c.authenticate()
	if err != nil {
		panic(err)
	}

	c.Ready <- true

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.conn.Messages():
			c.handle(&msg)
		}
	}
}

func (c *Client) authenticate() error {
	auth := c.conn.Auth()

	for {
		select {
		case id := <-auth:
			c.user.ID = id
			return nil
		case <-time.After(time.Second * 3):
			return errors.New("Client authentication timed out (<-c.conn.Auth())")
		}
	}
}

func (c *Client) handle(msg *Message) {
	if !msg.sentBy(c.user) {
		return
	}

	if msg.escaped() {
		edit := msg.EditText().Set(trimEscape(msg.Text))
		msg.perform(edit)

		c.conn.Update(msg)
		return
	}

	c.run(msg)
}

func (c *Client) run(m *Message) {
	exe := c.loadExecutions(m)
	act := c.runExecutions(exe)

	if performActions(act, m) {
		c.conn.Update(m)
	}
}

func (c *Client) runExecutions(exe []Execution) []MessageAction {
	var act []MessageAction
	res := make(chan ExecutionResult, len(c.user.Configurations))

	c.runner.Run(exe, res, c.timeout)

	for r := range res {
		if r.Err != nil {
			panic(r.Err)
		}

		if _, ok := r.Action.(*EmptyAction); ok {
			continue
		}

		act = append(act, r.Action)
	}

	return act
}

func (c *Client) loadExecutions(m *Message) []Execution {
	var exe []Execution
	for _, cfg := range c.user.Configurations {
		var x Extension

		// Check for HTTP extensions
		if cfg.URL != "" {
			x = NewHttpExtension(cfg.URL, nil)
		} else {
			var err error
			x, err = c.loader.Load(cfg.EID)

			if err != nil {
				panic(err)
			}
		}

		exe = append(exe, *NewExecution(x, *m, cfg.Cfg))
	}

	return exe
}

func (m *Message) sentBy(u *User) bool {
	return u.ID.Equals(m.UserID)
}

func (m *Message) escaped() bool {
	return strings.HasPrefix(m.Text, `\`)
}

func trimEscape(text string) string {
	return strings.TrimPrefix(text, `\`)
}

func performActions(act []MessageAction, msg *Message) (modified bool) {
	for _, a := range act {
		err, mod := msg.perform(a)

		if err != nil {
			panic(err)
		}

		modified = modified || mod
	}

	return
}
