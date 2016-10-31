package main

import (
	"fmt"
	"github.com/Epoch2/slack-sirius/store/db"
)

func main() {
	store := db.Connect("jv", "")
	users := store.GetUsers()

	for _, u := range *users {
		fmt.Printf("User: %v\n", u)
	}
}
