package sirius

import (
	"errors"
	"strings"
	"time"
)

type Client struct {
	ExtensionLoader
	ExtensionRunner
	Ready   chan bool
	user    *User
	conn    Connection
	timeout time.Duration
	stop    chan bool
}

type ClientConfig struct {
	user    *User
	loader  ExtensionLoader
	runner  ExtensionRunner
	timeout time.Duration
}

const CLIENT_DEFAULT_TIMEOUT = time.Second * 2

func NewClient(cfg ClientConfig) *Client {
	cl := &Client{
		ExtensionLoader: cfg.loader,
		ExtensionRunner: cfg.runner,
		conn:            NewRTMConnection(cfg.user.Token),
		user:            cfg.user,
		timeout:         cfg.timeout,
		Ready:           make(chan bool, 1),
		stop:            make(chan bool),
	}
	if cl.ExtensionRunner == nil {
		cl.ExtensionRunner = NewAsyncRunner()
	}
	if cl.timeout == 0 {
		cl.timeout = CLIENT_DEFAULT_TIMEOUT
	}
	return cl
}

func (c *Client) Start() {
	go c.conn.Listen()

	err := c.authenticate(c.conn)
	if err != nil {
		panic(err)
	}

	c.Ready <- true

	for {
		select {
		// Make sure we always check for termination signal first
		case <-c.stop:
			c.conn.Close()
			return
		case <-c.conn.Finished():
			return
		default:
		}

		select {
		case msg := <-c.conn.Messages():
			c.handle(&msg)
		}
	}
}

func (c *Client) Stop() {
	c.stop <- true
}

func (c *Client) authenticate(conn Connection) error {
	auth := conn.Auth()
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
	var exe []Execution
	for _, cfg := range c.user.Configurations {
		exe = append(exe, *c.createExecution(m, &cfg))
	}

	act := c.runExecutions(exe)

	if performActions(act, m) {
		c.conn.Update(m)
	}
}

func (c *Client) runExecutions(exe []Execution) []MessageAction {
	var act []MessageAction
	res := make(chan ExecutionResult, len(c.user.Configurations))

	c.Run(exe, res, c.timeout)

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

func (c *Client) createExecution(m *Message, cfg *Configuration) *Execution {
	var x Extension

	// Check for HTTP extensions
	if cfg.URL != "" {
		x = NewHttpExtension(cfg.URL, nil)
	} else {
		ex, err := c.Load(cfg.EID)

		if err != nil {
			panic(err)
		}

		x = ex
	}

	return NewExecution(x, *m, cfg.Cfg)
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
