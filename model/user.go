package model

import (
	"time"
)

type User struct {
	Id             string
	Token          string
	CreatedAt      time.Time
	Configurations []*Configuration
}

func NewUser(token string) User {
	return User{
		Token: token,
	}
}
