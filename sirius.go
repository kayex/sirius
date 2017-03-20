package sirius

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kayex/sirius/slack"
	"github.com/kayex/sirius/text"
)

const EMOJI = "⚡" // The high voltage/lightning bolt emoji (:zap: in Slack)

type Service struct {
	loader  ExtensionLoader
	clients map[string]*Client
	ctx     context.Context
}

func NewService(l ExtensionLoader) *Service {
	return &Service{
		loader:  l,
		clients: make(map[string]*Client),
	}
}

func (s *Service) Start(ctx context.Context, users []User) {
	s.ctx = ctx

	for _, u := range users {
		u := u
		cl := s.createClient(&u)
		s.addClient(cl)

		go cl.Start()
	}

	select {
	case <-ctx.Done():
		for _, cl := range s.clients {
			cl.Stop()
		}
	}
}

func (s *Service) AddUser(u *User) {
	stt := time.Now()
	cl := s.createClient(u)
	s.addClient(cl)

	go cl.Start()

	<-cl.Ready
	cl.notify(stt)
}

func (s *Service) DropUser(id slack.ID) {
	if cl, ok := s.clients[id.String()]; ok {
		cl.Stop()
		delete(s.clients, id.String())
	}
}

func (s *Service) addClient(cl *Client) error {
	u := cl.user

	if u.ID == nil {
		id, err := cl.conn.GetUserID(cl.user.Token)

		if err != nil {
			panic(err)
		}

		u.ID = id
	}

	if _, exists := s.clients[cl.user.ID.String()]; exists {
		return errors.New("Client with ID %v is already registered with service.")
	}

	s.clients[cl.user.ID.String()] = cl

	return nil
}

func (s *Service) createClient(u *User) *Client {
	return NewClient(ClientConfig{
		user:   u,
		loader: s.loader,
	})
}

func (c *Client) notify(st time.Time) {
	et := time.Now()
	tt := et.Sub(st)

	// Display load time in seconds, with three decimals.
	conf := EMOJI + " " + text.Bold(fmt.Sprintf(
		"%d extensions loaded in %.3f seconds.",
		len(c.user.Configurations),
		float64(tt.Nanoseconds())/float64(1e9)))

	if len(c.user.Configurations) == 0 {
		conf += "\n" + text.Quote(text.Italic("No extensions activated."))
	} else {
		for _, cfg := range c.user.Configurations {
			conf += "\n" + text.Quote(string(cfg.EID))
		}
	}

	c.conn.Send(&Message{
		Text:    conf,
		Channel: c.conn.Details().SelfChan,
	})
}
