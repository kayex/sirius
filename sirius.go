package sirius

import (
	"github.com/kayex/sirius/config"
	"golang.org/x/net/context"
)

type Service struct {
	loader ExtensionLoader
	config config.AppConfig
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
		var ext []Extension

		for _, c := range user.Configurations {
			err, x := s.loader.Load(c.EID)

			if err != nil {
				panic(err)
			}

			ext = append(ext, x)
		}

		cl := NewClient(&user, ext)

		go cl.Start(ctx)
	}

	select {}
}
