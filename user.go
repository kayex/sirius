package sirius

import (
	"github.com/kayex/sirius/slack"
)

type User struct {
	Profile
	ID    slack.ID
	Token string
}

type Profile struct {
	Configurations []Configuration
}

func NewUser(token string) *User {
	return &User{
		ID:    nil,
		Token: token,
	}
}
