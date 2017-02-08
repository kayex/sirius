package sirius

import (
	"github.com/kayex/sirius/config"
	"golang.org/x/net/context"
)

type Service struct {
	loader  ExtensionLoader
	config  config.AppConfig
	clients []Client
}

func NewService(l ExtensionLoader, cfg config.AppConfig) *Service {
	return &Service{
		loader: l,
		config: cfg,
	}
}

func (s *Service) Start(ctx context.Context, users []User) {
	for _, user := range users {
		cl := NewClient(&user, s.loader)

		go cl.Start(ctx)
	}

	select {}
}
