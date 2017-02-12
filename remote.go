package sirius

import (
	"encoding/json"
	"github.com/kayex/sirius/slack"
	"net/http"
)

type Remote struct {
	host   string
	token  string
	client *http.Client
}

type RemoteUser struct {
	IDHash string      `json:"id_hash_sha256"`
	Token  string      `json:"slack_token"`
	Config interface{} `json:"config"`
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

	switch cfg := ru.Config.(type) {
	case map[string]interface{}:
		for eid, c := range cfg {
			cfg := NewConfiguration(EID(eid))

			if conf, ok := c.(map[string]interface{}); ok {
				cfg.Cfg = conf
			}

			u.Configurations = append(u.Configurations, &cfg)
		}
	case []interface{}:
		for eid := range cfg {
			c := NewConfiguration(EID(eid))
			u.Configurations = append(u.Configurations, &c)
		}
	}

	return u
}

func (r *Remote) request(endpoint string) (*http.Response, error) {
	url := r.host + endpoint + "?token=" + r.token

	return r.client.Get(url)
}

func (r *Remote) GetUser(token string) (*User, error) {
	var ru RemoteUser

	res, err := r.request("/configs/" + token)

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
