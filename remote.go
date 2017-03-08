package sirius

import (
	"encoding/json"
	"net/http"

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
	Extensions     interface{} `json:"config"`
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

func (ru *RemoteUser) parseExtensionList(extl interface{}) []Configuration {
	var cfgs []Configuration

	switch ext := extl.(type) {
	case map[string]interface{}:
		cfgs = FromConfigurationMap(ext)
	case []interface{}:
		var m map[string]interface{}

		for _, v := range ext {
			if k, ok := v.(string); ok {
				m[k] = nil
			}
		}

		cfgs = FromConfigurationMap(m)
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
