package sirius

import (
	"encoding/json"
	"github.com/kayex/sirius/slack"
	"net/http"
)

type Remote struct {
	url    string
	token  string
	client *http.Client
}

type RemoteUser struct {
	IDHash string                 `json:"id_hash_sha256"`
	Token  string                 `json:"slack_token"`
	Config map[string]interface{} `json:"config"`
}

func NewRemote(url, token string) *Remote {
	return &Remote{
		url:    url,
		token:  token,
		client: &http.Client{},
	}
}

func (ru *RemoteUser) ToUser() *User {
	u := NewUser(ru.Token)
	u.ID = slack.SecureID{ru.IDHash}

	for eid, c := range ru.Config {
		cfg := NewConfiguration(EID(eid))

		if conf, ok := c.(ExtensionConfig); ok {
			cfg.Cfg = conf
		}

		u.Configurations = append(u.Configurations, &cfg)
	}

	return u
}

func (r *Remote) request(endpoint string) (*http.Response, error) {
	url := r.url + endpoint + "?token=" + r.token

	return r.client.Get(url)
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
