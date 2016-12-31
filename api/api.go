package api

import (
	"fmt"
	"github.com/kayex/sirius/model"
	"github.com/nlopes/slack"
	"log"
)

type Api struct {
	UserId   string
	Rtm      *slack.RTM
	Incoming chan model.Message
	client   *slack.Client
	token    string
}

func Init(logger *log.Logger) {
	slack.SetLogger(logger)
}

func New(token string) Api {
	client := slack.New(token)

	rtm := client.NewRTM()
	inc := make(chan model.Message)

	return Api{
		Rtm:      rtm,
		Incoming: inc,
		client:   client,
		token:    token,
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

func (api *Api) SendMessage(msg *model.Message) {
	omsg := api.Rtm.NewOutgoingMessage(msg.Text, msg.Channel)
	api.Rtm.SendMessage(omsg)
}

func (api *Api) SendUpdatedMessage(msg *model.Message) error {
	_, _, _, err := api.Rtm.UpdateMessage(msg.Channel, msg.Timestamp, msg.Text)
	return err
}

func (api *Api) handleIncomingEvent(ev slack.RTMEvent) {
	switch msg := ev.Data.(type) {
	case *slack.ConnectedEvent:
		api.UserId = msg.Info.User.ID

	case *slack.MessageEvent:
		api.handleIncomingMessage(msg)

	case *slack.RTMError:
		fmt.Printf("Error: %s\n", msg.Error())
		panic(msg)

	case *slack.InvalidAuthEvent:
		panic(msg)
	}

}

func (api *Api) handleIncomingMessage(ev *slack.MessageEvent) {
	msg := model.NewMessage(ev.Text, ev.User, ev.Channel, ev.Timestamp)
	api.Incoming <- msg
}
