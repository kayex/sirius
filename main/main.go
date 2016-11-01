package main

import (
	"github.com/kayex/sirius/store/db"
)

func main() {
	store := db.Connect("root", "eefca45AF")
	store.CreateSchema()
}
