package sirius

import (
	"errors"

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
		s.addClient(cl)

		go cl.Start()
	}

	select {
	case <-s.ctx.Done():
		break
	}
}

func (s *Service) AddUser(u *User) {
	cl := s.createClient(u)
	s.addClient(cl)

	go cl.Start()

	<-cl.Ready
	cl.notify()
}

func (s *Service) DropUser(id slack.ID) bool {
	if cl, ok := s.clients[id.String()]; ok {
		cl.Cancel()

		return true
	}

	return false
}

func (s *Service) stopClient(id slack.ID) {
	if ex, ok := s.clients[id.String()]; ok {
		ex.Cancel()
		delete(s.clients, id.String())
	}
}

func (s *Service) addClient(cl *CancelClient) {
	u := cl.user

	if !u.ID.Valid() {
		id, err := cl.conn.GetUserID(cl.user.Token)

		if err != nil {
			panic(err)
		}

		u.ID = id
	}

	if _, exists := s.clients[cl.user.ID.String()]; exists {
		errors.New("Client with ID %v is already registered with service.")
	}

	s.clients[cl.user.ID.String()] = cl
}

func (s *Service) createClient(u *User) *CancelClient {
	return NewClient(ClientConfig{
		user:   u,
		loader: s.loader,
	}).WithCancel(context.WithCancel(s.ctx))
}

func (c *Client) notify() {
	conf := EMOJI + " Configuration loaded successfully."

	if len(c.user.Configurations) == 0 {
		conf += "\n" + slack.Quote(slack.Italic("No extensions activated."))
	} else {
		for _, cfg := range c.user.Configurations {
			conf += "\n" + slack.Quote(slack.Bold(string(cfg.EID)))
		}
	}

	c.conn.Send(&Message{
		Text:    conf,
		Channel: c.conn.Details().SelfChan,
	})
}
