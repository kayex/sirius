package sirius

import (
	"net/url"

	"github.com/kayex/sirius/slack"
)

type User struct {
	ID             slack.ID
	Token          string
	Configurations []Configuration
}

type Configuration struct {
	URL string
	EID EID
	Cfg ExtensionConfig
}

func FromConfigurationMap(cfg map[string]interface{}) []Configuration {
	var cfgs []Configuration

	for eid, settings := range cfg {
		var c Configuration

		// Check for HTTP extensions
		_, err := url.ParseRequestURI(eid)
		if err == nil {
			c = NewHTTPConfiguration(eid)
		} else {
			c = NewConfiguration(EID(eid))
		}

		switch ec := settings.(type) {
		case ExtensionConfig:
			c.Cfg = ec
		case map[string]interface{}:
			c.Cfg = ExtensionConfig(ec)
		}

		cfgs = append(cfgs, c)
	}

	return cfgs
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
