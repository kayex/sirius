package sirius

import (
	"golang.org/x/net/context"
)

type Service struct {
	loader  ExtensionLoader
	clients []Client
}

func NewService(l ExtensionLoader) *Service {
	return &Service{
		loader: l,
	}
}

func (s *Service) Start(ctx context.Context, users []User) {
	for _, u := range users {
		u := u
		cl := NewClient(&u, s.loader)
		s.clients = append(s.clients, *cl)

		go cl.Start(ctx)
	}

	select {
	case <-ctx.Done():
		break
	}
}
