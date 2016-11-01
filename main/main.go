package main

import (
	"reflect"
	"fmt"
	"github.com/kayex/sirius/config"
	"github.com/kayex/sirius/store/db"
)

func main() {
	config := initConfig()
	store := initDbStore(&config)

	user := store.GetUser("d915025d-406f-49ac-8b74-0a9497fee4e7")
	cfg := store.GetConfiguration(user, "my-plugin")
	fmt.Println(reflect.TypeOf(cfg.Config))

	//cfg.SetConfig(map[string]interface{}{
	//	"setting-1": "value",
	//})

	store.UpdateConfiguration(cfg)

	fmt.Println(cfg.User)

	//cfg := model.NewConfiguration(user, "my-plugin")
	//store.SaveConfiguration(&cfg)

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
