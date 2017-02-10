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
	for _, user := range users {
		cl := NewClient(&user, s.loader)

		go cl.Start(ctx)
	}

	select {
	case <-ctx.Done():
		break
	}
}
