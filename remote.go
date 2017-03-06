package sirius

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/kayex/sirius/slack"
)

type Remote struct {
	host   string
	token  string
	client *http.Client
}

type RemoteUser struct {
	IDHash         string      `json:"sirius_id"`
	Token          string      `json:"slack_token"`
	Extensions     interface{} `json:"extensions"`
	HttpExtensions interface{} `json:"http_extensions"`
}

func NewRemote(host, token string) *Remote {
	return &Remote{
		host:   host,
		token:  token,
		client: &http.Client{},
	}
}

func (ru *RemoteUser) ToUser() *User {
	u := NewUser(ru.Token)
	u.ID = slack.SecureID{ru.IDHash}

	u.Configurations = append(u.Configurations, ru.parseExtensionList(ru.Extensions)...)
	u.Configurations = append(u.Configurations, ru.parseExtensionList(ru.HttpExtensions)...)

	return u
}

func (ru *RemoteUser) parseExtensionList(extl interface{}) []*Configuration {
	var cfgs []*Configuration

	switch ext := extl.(type) {
	case map[string]interface{}:
		for eid, settings := range ext {
			var c Configuration

			// Check for HTTP extensions
			_, err := url.ParseRequestURI(eid)
			if err == nil {
				c = NewHTTPConfiguration(eid)
			} else {
				c = NewConfiguration(EID(eid))
			}

			if conf, ok := settings.(map[string]interface{}); ok {
				c.Cfg = ExtensionConfig(conf)
			}
			cfgs = append(cfgs, &c)
		}
	case []interface{}:
		for eid := range ext {
			c := NewConfiguration(EID(eid))
			cfgs = append(cfgs, &c)
		}
	}

	return cfgs
}

func (r *Remote) request(endpoint string) (*http.Response, error) {
	url := r.host + endpoint + "?token=" + r.token

	return r.client.Get(url)
}

func (r *Remote) GetUser(id slack.SecureID) (*User, error) {
	var ru RemoteUser

	res, err := r.request("/configs/" + id.HashSum)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&ru); err != nil {
		return nil, err
	}

	return ru.ToUser(), nil
}

func (r *Remote) GetUsers() ([]User, error) {
	var ru []RemoteUser

	res, err := r.request("/configs")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&ru); err != nil {
		return nil, err
	}

	var users []User

	for _, u := range ru {
		users = append(users, *u.ToUser())
	}
	return users, nil
}
