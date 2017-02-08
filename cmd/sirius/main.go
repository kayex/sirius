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

	for _, user := range users {
		var ext []sirius.Extension

		for _, c := range user.Configurations {
			err, x := loader.Load(c.EID)

			if err != nil {
				panic(err)
			}

			ext = append(ext, x)
		}

		cl := sirius.NewClient(&user, ext)

		go cl.Start(context.TODO())
	}

	select {}
}
