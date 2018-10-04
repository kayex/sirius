package sirius

import (
	"context"
	"fmt"
	"github.com/kayex/sirius/sync"
	"strings"
	"time"
)

type Client struct {
	user    *User
	api     *SlackAPI
	conn    Connection
	exe     *Executor
	timeout time.Duration
	ctrl    *sync.Control
}

type CancelClient struct {
	*Client
	ctx    context.Context
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

	exe := NewExecutor(cfg.loader, cfg.timeout)

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

	err = c.exe.Load(c.user.Profile)
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
		msg.apply(TextEdit().Set(trimEscape(msg.Text)))
		return c.conn.Update(msg)
	}

	m, mod := c.process(*msg)

	if !mod {
		return nil
	}

	return c.conn.Update(&m)
}

func (c *Client) process(msg Message) (Message, bool) {
	var act []MessageAction

	res := c.exe.RunExtensions(msg)
	for r := range res {
		if r.Err != nil {
			fmt.Printf("extension failed: %v\n", r.Err)
		}

		act = append(act, r.Action)
	}

	modified, err := msg.applyAll(act)
	if err != nil {
		fmt.Printf("error applying message action: %v\n", err)
	}

	return msg, modified
}

func trimEscape(text string) string {
	return strings.TrimPrefix(text, `\`)
}
