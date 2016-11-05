package core

import (
	"os"
	"log"
	"fmt"
	"github.com/kayex/sirius/model"
	"github.com/kayex/sirius/api"
)

type Client struct {
	api *api.Api
	user *model.User
}

func NewClient(user *model.User) Client {
	api.Init(createLogger())
	api := api.New(user.Token)

	return Client{
		api: &api,
		user: user,
	}
}

func (c *Client) Start() {
	go c.api.Listen()

	for {
		select {
		case msg := <-c.api.Incoming:

		}
	}

	fmt.Println("Done.")
}

func (c *Client) handleMessage(msg *model.Message) {

}

func createLogger() *log.Logger {
	return log.New(os.Stdout, "sirius: ", log.Lshortfile|log.LstdFlags)
}