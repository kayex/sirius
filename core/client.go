package core

import (
	"github.com/kayex/sirius/api"
	"github.com/kayex/sirius/model"
	"github.com/kayex/sirius/plugins"
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

	c.runUserPlugins(msg)
}

func (c *Client) runUserPlugins(msg *model.Message) {
	trans := []plugins.Transformation{}

	for _, cfg := range c.user.Configurations {
		pg := getPluginForPid(cfg.PluginGuid)
		trans = append(trans, pg.Run(*msg)...)
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

func getPluginForPid(pid string) plugins.Plugin {
	switch pid {
	case "thumbs_up":
		return &plugins.ThumbsUp{}
	case "ripperino":
		return &plugins.Ripperino{}
	}

	panic("Invalid pid: " + pid)
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
