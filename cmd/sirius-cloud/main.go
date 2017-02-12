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

	ld := extension.NewStaticLoader(cfg)
	sync := sirius.NewMQTTSync(rmt, cfg.MQTT.Config, cfg.MQTT.Topic)

	s := sirius.NewService(ld).WithSync(sync)

	s.Start(context.Background(), users)
}
