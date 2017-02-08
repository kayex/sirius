package main

import (
	"github.com/kayex/sirius"
	"github.com/kayex/sirius/config"
	"github.com/kayex/sirius/extension"
	"golang.org/x/net/context"
)

func main() {
	cfg := config.FromEnv()

	rmt := sirius.NewRemote(cfg.Remote.URL, cfg.Remote.Token)
	users, err := rmt.GetUsers()

	if err != nil {
		panic(err)
	}

	loader := extension.NewStaticLoader(cfg)
	service := sirius.NewService(loader, cfg)

	service.Start(context.TODO(), users)
}
