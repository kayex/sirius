package main

import (
	"fmt"
	"github.com/kayex/sirius/store/db"
)

func main() {
	store := db.Connect("jv", "")
	users := store.GetUsers()

	for _, u := range *users {
		fmt.Printf("User: %v\n", u)
	}
}
