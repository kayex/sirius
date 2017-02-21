package sirius

import (
	"github.com/kayex/sirius/slack"
)

type User struct {
	ID             slack.ID
	Token          string
	Configurations []*Configuration
}

func NewUser(token string) *User {
	return &User{
		ID:    slack.SecureID{},
		Token: token,
	}
}
