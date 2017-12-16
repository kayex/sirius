package sirius

import (
	"github.com/kayex/sirius/slack"
)

type User struct {
	ID             slack.ID
	Token          string
	Settings Settings
}

func NewUser(token string) *User {
	return &User{
		ID:    nil,
		Token: token,
	}
}
