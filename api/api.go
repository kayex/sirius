package api

import (
	"fmt"
	"log"
	"github.com/nlopes/slack"
	"github.com/kayex/sirius/model"
)

type Api struct {
	Rtm *slack.RTM
	Incoming chan model.Message
}

func Init(logger *log.Logger) {
	slack.SetLogger(logger)
}

func New(token string) Api {
	api := slack.New(token)
	rtm := api.NewRTM()
	inc := make(chan model.Message)

	return Api{
		Rtm: rtm,
		Incoming: inc,
	}
}

func (api *Api) Listen() {
	go api.Rtm.ManageConnection()

	for {
		select {
		case ev := <-api.Rtm.IncomingEvents:
			api.handleIncomingEvent(ev)
		}
	}
}

func (api *Api) handleIncomingEvent(ev slack.RTMEvent) {
	switch msg := ev.Data.(type) {
	case *slack.MessageEvent:
		api.handleIncomingMessage(msg)

	case *slack.RTMError:
		fmt.Printf("Error: %s\n", ev.Error())

	case *slack.InvalidAuthEvent:
		fmt.Printf("Invalid credentials")
	}

}

func (api *Api) handleIncomingMessage(ev *slack.MessageEvent) {
	msg := model.NewMessage(ev.Text, ev.Channel)
	api.Incoming <- msg
}
