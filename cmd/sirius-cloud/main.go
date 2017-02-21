package main

import (
	"fmt"
	"github.com/kayex/sirius"
	"github.com/kayex/sirius/config"
	"github.com/kayex/sirius/extension"
	"golang.org/x/net/context"
)

func main() {
	cfg := config.FromEnv()
	rmt := sirius.NewRemote(cfg.Remote.Host, cfg.Remote.Token)

	users, err := rmt.GetUsers()

	if err != nil {
		panic(err)
	}

	ld := extension.NewStaticLoader(cfg)
	sync := sirius.NewMQTTSync(rmt, cfg.MQTT.Config, cfg.MQTT.Topic)

	s := sirius.NewService(ld).WithSync(sync)

	printRunInfo(users, cfg)

	s.Start(context.Background(), users)
}

func printRunInfo(users []sirius.User, cfg config.AppConfig) {
	fmt.Println("Connecting to remote: " + cfg.Remote.Host)
	fmt.Printf("Establishing MQTT sync: %v@%v:%v [%v]\n", cfg.MQTT.CID, cfg.MQTT.Host, cfg.MQTT.Port, cfg.MQTT.Topic)

	for _, u := range users {
		tks := ""
		for i, c := range u.Configurations {
			if i != 0 {
				tks += ", "
			}
			tks += string(c.EID)
		}
		fmt.Printf("[%v] %v (%v)\n", u.ID.String(), len(u.Configurations), tks)
	}
}
