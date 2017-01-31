package main

import (
	"encoding/json"
	"fmt"
	"github.com/kayex/sirius"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
)

func main() {
	tokens := getTokensFromJson()
	users := []sirius.User{}

	for _, token := range tokens {
		user := sirius.NewUser(token)

		tu := sirius.NewConfiguration(&user, "thumbs_up")
		rip := sirius.NewConfiguration(&user, "ripperino")
		rpl := sirius.NewConfiguration(&user, "replacer")

		user.AddConfiguration(&tu)
		user.AddConfiguration(&rip)
		user.AddConfiguration(&rpl)

		users = append(users, user)
	}

	for _, user := range users {
		cl := sirius.NewClient(&user)
		go cl.Start(context.TODO())
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
