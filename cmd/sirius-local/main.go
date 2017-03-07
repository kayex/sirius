package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kayex/sirius"
	"github.com/kayex/sirius/config"
	"github.com/kayex/sirius/extension"
	"io/ioutil"
)

var extensions []string = []string{
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

	s.Start(context.TODO(), users)
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
	var m map[string]interface{}

	for _, eid := range extensions {
		m[eid] = nil
	}

	u.Configurations = append(u.Configurations, sirius.FromConfigurationMap(m)...)
}

func getTokensFromJSON() []string {
	file, err := ioutil.ReadFile("./users.json")
	if err != nil {
		fmt.Println(err.Error())
	}

	var tokens []string

	json.Unmarshal(file, &tokens)

	return tokens
}
