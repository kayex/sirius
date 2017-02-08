package sirius

import (
	"github.com/kayex/sirius/slack"
)

type User struct {
	ID             slack.SecureID
	Token          string
	Configurations []*Configuration
}

func NewUser(token string) *User {
	return &User{
		Token: token,
	}
}
