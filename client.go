package sirius

import (
	"strings"
	"time"
	"context"
	"github.com/kayex/sirius/sync"
	"fmt"
	"log"
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
	runner  ExtensionRunner
	timeout time.Duration
}

func NewClient(cfg ClientConfig) *Client {
	api := NewSlackAPI(cfg.user.Token)
	conn := api.NewRTMConnection()

	exe := &Executor{
		runner: cfg.runner,
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

	err = c.exe.LoadFromSettings(c.user.Settings)
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

	if msg.escaped() {
		edit := msg.EditText().Set(trimEscape(msg.Text))
		msg.perform(edit)

		return c.conn.Update(msg)
	}

	m, mod := c.execute(*msg)

	if mod {
		err := c.conn.Update(&m)
		if err != nil {
			panic(err)
		}
		return nil
	}

	return nil
}

// execute runs the executions in exe on msg. Returns the new message, and a
// bool indicating if any changes were made to the message text.
func (c *Client) execute(msg Message) (Message, bool) {
	var act []MessageAction

	res := c.exe.Run(msg)
	for r := range res {
		if r.Err != nil {
			log.Println(r.Err)
		}

		if _, ok := r.Action.(*EmptyAction); ok {
			continue
		}

		act = append(act, r.Action)
	}

	modified := performActions(act, &msg)

	return msg, modified
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
