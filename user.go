package sirius

import (
	"time"
)

type User struct {
	ID             string
	Token          string
	CreatedAt      time.Time
	Configurations []*Configuration
}

func NewUser(token string) User {
	return User{
		Token: token,
	}
}

func (u *User) AddConfiguration(cfg *Configuration) {
	u.Configurations = append(u.Configurations, cfg)
}
