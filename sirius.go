package sirius

import (
	"context"
	"fmt"
	"time"

	"github.com/kayex/sirius/slack"
	"github.com/kayex/sirius/text"
)

const EMOJI = "⚡" // The high voltage/lightning bolt emoji (:zap: in Slack)

type Service struct {
	loader  ExtensionLoader
	clients map[string]*CancelClient
	ctx     context.Context
}

func NewService(l ExtensionLoader) *Service {
	return &Service{
		loader:  l,
		clients: make(map[string]*CancelClient),
	}
}

func (s *Service) Start(ctx context.Context, users []User) error {
	s.ctx = ctx

	for _, u := range users {
		u := u
		err := s.AddUser(&u, false)
		if err != nil {
			return err
		}
	}

	select {
		case <-s.ctx.Done():
			return nil
	}
}

func (s *Service) AddUser(u *User, notify bool) error {
	stt := time.Now()

	cl := s.createClient(u)
	s.addClient(cl)
	err := cl.Start()
	if err != nil {
		return err
	}

	if notify {
		cl.notify(stt)
	}

	return nil
}

func (s *Service) DropUser(id slack.ID) {
	if cl, ok := s.clients[id.String()]; ok {
		cl.cancel()
		delete(s.clients, id.String())
	}
}

func (s *Service) addClient(cl *CancelClient) error {
	u := cl.user

	if u.ID == nil {
		id, err := cl.api.GetUserID(cl.user.Token)

		if err != nil {
			panic(err)
		}

		u.ID = id
	}

	if _, exists := s.clients[u.ID.String()]; exists {
		return fmt.Errorf("client with ID %v is already registered with service", u.ID)
	}

	s.clients[cl.user.ID.String()] = cl

	return nil
}

func (s *Service) createClient(u *User) *CancelClient {
	ctx, cancel := context.WithCancel(s.ctx)
	return &CancelClient{
		Client: NewClient(ClientConfig{
			user: u,
			loader: s.loader,
			runner: &AsyncRunner{},
		}),
		ctx: ctx,
		cancel: cancel,
	}
}

func (c *Client) notify(st time.Time) {
	et := time.Now()
	tt := et.Sub(st)

	// Display load time in seconds, with three decimals.
	conf := EMOJI + " " + text.Bold(fmt.Sprintf(
		"%d extensions loaded in %.3f seconds.",
		len(c.user.Settings),
		float64(tt.Nanoseconds())/float64(1e9)))

	if len(c.user.Settings) == 0 {
		conf += "\n" + text.Quote(text.Italic("No extensions activated."))
	} else {
		for _, cfg := range c.user.Settings {
			conf += "\n" + text.Quote(string(cfg.EID))
		}
	}

	c.conn.Send(&Message{
		Text:    conf,
		Channel: c.conn.Details().SelfChan,
	})
}
