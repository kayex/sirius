package core

import (
	"github.com/kayex/sirius/api"
	"github.com/kayex/sirius/model"
	"github.com/kayex/sirius/extension"
	"log"
	"os"
	"strings"
)

type Client struct {
	api  *api.Api
	user *model.User
}

func NewClient(user *model.User) Client {
	api.Init(createLogger())
	ap := api.New(user.Token)

	return Client{
		api:  &ap,
		user: user,
	}
}

func (c *Client) Start() {
	go c.api.Listen()

	for {
		select {
		case msg := <-c.api.Incoming:
			c.handleMessage(&msg)
		}
	}
}

func (c *Client) handleMessage(msg *model.Message) {
	if !c.isSender(msg) {
		return
	}

	if shouldEscape(msg) {
		trimEscape(msg)
		return
	}

	c.runUserExtensions(msg)
}

func (c *Client) runUserExtensions(msg *model.Message) {
	trans := []extension.Transformation{}

	for _, cfg := range c.user.Configurations {
		ext := getExtensionForEID(cfg.ExtensionGUID)

		err, res := ext.Run(*msg)

		if err != nil {
			panic(err)
		}

		trans = append(trans, res...)
	}

	text := msg.Text

	for _, t := range trans {
		text = t.Apply(text)
	}

	c.updateText(msg, text)
}

func (c *Client) updateText(msg *model.Message, text string) {
	if text == msg.Text {
		return
	}

	msg.Text = text
	c.sendUpdate(msg)
}

func (c *Client) sendUpdate(msg *model.Message) {
	err := c.api.SendUpdatedMessage(msg)

	if err != nil {
		panic(err)
	}
}

func (c *Client) isSender(msg *model.Message) bool {
	return c.api.UserId == msg.UserID
}

func getExtensionForEID(eid string) extension.Extension {
	switch eid {
	case "thumbs_up":
		return &extension.ThumbsUp{}
	case "ripperino":
		return &extension.Ripperino{}
	case "replacer":
		return &extension.Replacer{}
	}

	panic("Invalid eid: " + eid)
}

func createLogger() *log.Logger {
	return log.New(os.Stdout, "sirius: ", log.Lshortfile|log.LstdFlags)
}

func shouldEscape(msg *model.Message) bool {
	return strings.HasPrefix(msg.Text, `\`)
}

func trimEscape(msg *model.Message) {
	msg.Text = strings.TrimPrefix(msg.Text, `\`)
}
