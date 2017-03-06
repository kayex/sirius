package main

import (
	"encoding/json"
	"fmt"
	"github.com/kayex/sirius"
	"github.com/kayex/sirius/config"
	"github.com/kayex/sirius/extension"
	"context"
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
	for _, eid := range extensions {
		c := sirius.NewConfiguration(sirius.EID(eid))
		u.Configurations = append(u.Configurations, &c)
	}
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
