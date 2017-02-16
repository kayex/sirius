package sirius

import (
	"github.com/kayex/sirius/slack"
	"golang.org/x/net/context"
)

const EMOJI = "⚡" // The high voltage/lightning bolt emoji (:zap: in Slack)

type Service struct {
	loader  ExtensionLoader
	clients map[string]*CancelClient
	ctx     context.Context
}

type CancelClient struct {
	Client
	Cancel context.CancelFunc
	ctx    context.Context
}

func (c *Client) WithCancel(ctx context.Context, cancel context.CancelFunc) *CancelClient {
	return &CancelClient{
		Client: *c,
		Cancel: cancel,
		ctx:    ctx,
	}
}

func (c *CancelClient) Start() {
	c.Client.Start(c.ctx)
}

func NewService(l ExtensionLoader) *Service {
	return &Service{
		loader:  l,
		clients: make(map[string]*CancelClient),
	}
}

func (s *Service) Start(ctx context.Context, users []User) {
	s.ctx = ctx

	for _, u := range users {
		u := u
		cl := s.createClient(&u)
		s.clients[u.ID.HashSum] = cl

		go cl.Start()
	}

	select {
	case <-s.ctx.Done():
		break
	}
}

func (s *Service) AddUser(u *User) {
	s.DropUser(u.ID)

	cl := s.createClient(u)
	s.clients[u.ID.HashSum] = cl

	go cl.Start()

	<-cl.Ready
	s.notify(u)
}

func (s *Service) DropUser(id slack.SecureID) bool {
	if cl, ok := s.clients[id.HashSum]; ok {
		cl.Cancel()

		return true
	}

	return false
}

func (s *Service) stopClient(id slack.SecureID) {
	if ex, ok := s.clients[id.HashSum]; ok {
		ex.Cancel()
	}
}

func (s *Service) createClient(u *User) *CancelClient {
	return NewClient(u, s.loader).WithCancel(context.WithCancel(s.ctx))
}

func (s *Service) notify(u *User) {
	cl := s.clients[u.ID.HashSum]

	conf := EMOJI + " " + slack.Italic("Configuration loaded successfully.")

	if len(u.Configurations) == 0 {
		conf += slack.Quote("No extensions activated.")
	} else {
		for _, cfg := range u.Configurations {
			conf += "\n" + slack.Quote(string(cfg.EID))
		}
	}

	cl.conn.Send(&Message{
		Text:    conf,
		Channel: cl.conn.SelfChan(),
	})
}
