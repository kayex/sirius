package sirius

import (
	"fmt"
	"github.com/kayex/sirius/core"
	"github.com/kayex/sirius/plugins"
	"github.com/kayex/sirius/store"
	"github.com/nlopes/slack"
	"log"
	"os"
)

func Route() {
	api := slack.New("xoxp-14643781812-14649325041-34618221140-986ad19416")
	logger := log.New(os.Stdout, "app-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connected: %v\n", ev)

			case *slack.MessageEvent:
				rtm.SendMessage(rtm.NewOutgoingMessage(newMsg.Text, "D0EJXP1C4"))

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:

				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
