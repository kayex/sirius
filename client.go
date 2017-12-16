package sirius

import (
	"strings"
	"time"
	"context"
	"github.com/kayex/sirius/sync"
	"fmt"
)

type Client struct {
	user    *User
	api     *SlackAPI
	conn    Connection
	exe *Executor
	timeout time.Duration
	ctrl    *sync.Control
}

type CancelClient struct {
	*Client
	ctx context.Context
	cancel context.CancelFunc
}

func (c *CancelClient) Start() error {
	return c.Client.Start(c.ctx)
}

type ClientConfig struct {
	user    *User
	loader  ExtensionLoader
	timeout time.Duration
}

func NewClient(cfg ClientConfig) *Client {
	api := NewSlackAPI(cfg.user.Token)
	conn := api.NewRTMConnection()

	exe := &Executor{
		loader: cfg.loader,
	}

	cl := &Client{
		api:  api,
		conn: conn,
		user: cfg.user,
		exe:  exe,
		ctrl: sync.NewControl(),
	}
	return cl
}

func (c *Client) Start(ctx context.Context) error {
	err := c.conn.Start(ctx)
	if err != nil {
		return err
	}

	err = c.exe.Load(c.user.Settings)
	if err != nil {
		return fmt.Errorf("error loading user extensions: %v", err)
	}

	go c.run(ctx)

	return nil
}

func (c *Client) run(ctx context.Context) {
	defer c.ctrl.Finish(nil)

	for {
		select {
		case msg := <-c.conn.Messages():
			err := c.handle(&msg)
			if err != nil {
				c.ctrl.Finish(err)
			}
		case <-ctx.Done():
			return
		case err := <-c.conn.Closed():
			if err != nil {
				c.ctrl.Finish(err)
			}
			return
		}
	}
}

func (c *Client) handle(msg *Message) error {
	if !msg.sentBy(c.user) {
		return nil
	}

	m, mod := c.process(*msg)

	if !mod {
		return nil
	}

	return c.conn.Update(&m)
}

// process processes a message using the loaded Executor.
//
// Returns the processed message, and a bool indicating whether the message
// text property was modified from its original value. This is useful
// for determining if a message u
func (c *Client) process(msg Message) (Message, bool) {
	// If the message is escaped, we'll strip the escape character(s) and
	// return the message immediately.
	if msg.escaped() {
		edit := msg.EditText().Set(trimEscape(msg.Text))
		msg.alter(edit)

		return msg, true
	}

	var act []MessageAction

	res := c.exe.RunExtensions(msg)
	for r := range res {
		if r.Err != nil {
			fmt.Println(r.Err)
			continue
		}

		act = append(act, r.Action)
	}

	modified, err := msg.alterAll(act)
	if err != nil {
		panic(err)
	}

	return msg, modified
}

func trimEscape(text string) string {
	return strings.TrimPrefix(text, `\`)
}
