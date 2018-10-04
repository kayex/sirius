package main

import (
	"context"
	"encoding/json"
	"github.com/kayex/sirius"
	"github.com/kayex/sirius/config"
	"github.com/kayex/sirius/extension"
	"io/ioutil"
)

var extensions = []string{
	"geocode",
	"ip_lookup",
	"quotes",
	"replacer",
	"ripperino",
	"thumbs_up",
}

func main() {
	cfg := config.FromEnv()
	users := createUsers(getTokensFromJSON())

	l := extension.NewStaticLoader(cfg)
	s := sirius.NewService(l)

	s.Start(context.Background(), users)
}

func createUsers(tokens []string) []sirius.User {
	var users []sirius.User

	for _, t := range tokens {
		u := sirius.NewUser(t)
		configure(u)

		users = append(users, *u)
	}

	return users
}

func configure(u *sirius.User) {
	m := make(map[string]interface{})

	for _, eid := range extensions {
		m[eid] = nil
	}

	u.Profile = append(u.Profile, sirius.FromConfigurationMap(m)...)
}

func getTokensFromJSON() []string {
	file, err := ioutil.ReadFile("./users.json")
	if err != nil {
		panic(err)
	}

	var tokens []string

	json.Unmarshal(file, &tokens)

	return tokens
}
