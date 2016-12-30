package main

import (
	"github.com/kayex/sirius/config"
	"github.com/kayex/sirius/core"
	"github.com/kayex/sirius/model"
	"github.com/kayex/sirius/store/db"
)

func main() {
	tokens := map[string]string{
		"johan": "xoxp-14643781812-14649325041-100316954898-6bfd950500977c5b25ffa6636b818811",
		"ash":   "xoxp-14643781812-14671711527-104875993958-fa382e478a906f535ca80f2864e6b90f",
	}

	users := []model.User{}

	for _, token := range tokens {
		user := model.NewUser(token)

		tu := model.NewConfiguration(&user, "thumbs_up")
		rip := model.NewConfiguration(&user, "ripperino")

		user.AddConfiguration(&tu)
		user.AddConfiguration(&rip)

		users = append(users, user)
	}

	for _, user := range users {
		cl := core.NewClient(&user)
		go cl.Start()
	}

	select {}
}

func initConfig() config.Config {
	return config.FromEnv()
}

func initDbStore(cfg *config.Config) db.Db {
	host := cfg.DbHost
	port := cfg.DbPort
	dbName := cfg.DbDatabase
	user := cfg.DbUser
	pwd := cfg.DbPassword

	return db.Connect(host, port, dbName, user, pwd)
}
