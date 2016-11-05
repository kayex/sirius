package main

import (
	"github.com/kayex/sirius/config"
	"github.com/kayex/sirius/core"
	"github.com/kayex/sirius/store/db"
)

func main() {
	config := initConfig()
	store := initDbStore(&config)
	user := store.GetUser("d915025d-406f-49ac-8b74-0a9497fee4e7")

	cl := core.NewClient(user)
	cl.Start()

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

