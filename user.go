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
	URL string
	EID EID
	Cfg ExtensionConfig
}

func NewConfiguration(eid EID) Configuration {
	return Configuration{
		EID: eid,
	}
}

func NewHTTPConfiguration(url string) Configuration {
	return Configuration{
		URL: url,
	}
}

func NewUser(token string) *User {
	return &User{
		ID:    nil,
		Token: token,
	}
}
