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

	l := extension.NewStaticLoader(cfg)
	s := sirius.NewService(l)

	s.Start(context.TODO(), users)
}
