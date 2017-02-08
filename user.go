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

func (u *User) AddConfiguration(cfg *Configuration) {
	u.Configurations = append(u.Configurations, cfg)
}
