package sirius

import (
	"golang.org/x/net/context"
)

type Service struct {
	loader  ExtensionLoader
	clients map[string]*CancelClient
	ctx     context.Context
}

type CancelClient struct {
	Client
	Cancel context.CancelFunc
}

func (c *Client) WithCancel(cancel context.CancelFunc) *CancelClient {
	return &CancelClient{
		Client: *c,
		Cancel: cancel,
	}
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
		s.startClient(&u)
	}

	select {
	case <-s.ctx.Done():
		break
	}
}

func (s *Service) AddUser(u *User) {
	s.startClient(u)
}

func (s *Service) DropUserWithToken(t string) bool {
	if cl, ok := s.clients[t]; ok {
		cl.Cancel()

		return true
	}

	return false
}

func (s *Service) startClient(u *User) {
	// Make sure we stop any existing client for the same user
	if ex, ok := s.clients[u.Token]; ok {
		ex.Cancel()
	}

	ctx, cancel := context.WithCancel(s.ctx)
	cl := NewClient(u, s.loader).WithCancel(cancel)

	go cl.Start(ctx)

	s.clients[u.Token] = cl
}
