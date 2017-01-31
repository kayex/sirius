package main

import (
	"encoding/json"
	"fmt"
	"github.com/kayex/sirius/core"
	"github.com/kayex/sirius/model"
	"io/ioutil"
	"os"
)

func main() {
	tokens := getTokensFromJson()
	users := []model.User{}

	for _, token := range tokens {
		user := model.NewUser(token)

		tu := model.NewConfiguration(&user, "thumbs_up")
		rip := model.NewConfiguration(&user, "ripperino")
		rpl := model.NewConfiguration(&user, "replacer")

		user.AddConfiguration(&tu)
		user.AddConfiguration(&rip)
		user.AddConfiguration(&rpl)

		users = append(users, user)
	}

	for _, user := range users {
		cl := core.NewClient(&user)
		go cl.Start()
	}

	select {}
}

func getTokensFromJson() []string {
	file, err := ioutil.ReadFile("./users.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var users []string

	json.Unmarshal(file, &users)
	return users
}

