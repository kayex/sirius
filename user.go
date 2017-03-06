package sirius

import (
	"github.com/kayex/sirius/slack"
)

type User struct {
	ID             slack.ID
	Token          string
	Configurations []*Configuration
}

type Configuration struct {
	EID EID
	Cfg ExtensionConfig
}

func NewConfiguration(eid EID) Configuration {
	return Configuration{
		EID: eid,
	}
}

func NewUser(token string) *User {
	return &User{
		ID:    nil,
		Token: token,
	}
}
