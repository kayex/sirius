package main

import (
	"context"
	"fmt"

	"github.com/kayex/sirius"
	"github.com/kayex/sirius/config"
	"github.com/kayex/sirius/extension"
)

func main() {
	cfg := config.FromEnv()
	rmt := sirius.NewRemote(cfg.Remote.Host, cfg.Remote.Token)

	users, err := rmt.GetUsers()

	if err != nil {
		panic(err)
	}

	ld := extension.NewStaticLoader(cfg)
	s := sirius.NewService(ld)

	printRunInfo(users, cfg)

	s.Start(context.Background(), users)
}

func printRunInfo(users []sirius.User, cfg config.AppConfig) {
	fmt.Println("Connecting to remote: " + cfg.Remote.Host)

	for _, u := range users {
		tks := ""
		for i, c := range u.Settings {
			if i != 0 {
				tks += ", "
			}
			tks += string(c.EID)
		}
		fmt.Printf("[%v] %v (%v)\n", u.ID.String(), len(u.Settings), tks)
	}
}
