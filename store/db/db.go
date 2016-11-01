package db

import (
	"github.com/kayex/sirius/model"
	"gopkg.in/pg.v5"
)

type Db struct {
	conn *pg.DB
}

func Connect(host string, port string, dbName string, user string, password string) Db {
	db := pg.Connect(&pg.Options{
		Addr:     host + `:` + port,
		Database: dbName,
		User:     user,
		Password: password,
	})

	return Db{
		conn: db,
	}
}

func (db *Db) GetUser(id string) *model.User {
	user := model.User{Id: id}

	exec(db.conn.Select(&user))

	return &user
}

func (db *Db) SaveUser(usr *model.User) {
	exec(db.conn.Insert(usr))
}

func (db *Db) SaveConfiguration(cfg *model.Configuration) {
	exec(db.conn.Insert(cfg))
}

func (db *Db) UpdateConfiguration(cfg *model.Configuration) {
	_, err := db.conn.Model(cfg).Column("config").Update()

	if (err != nil) {
		panic(err)
	}
}


func (db *Db) GetUsers() *[]model.User {
	var users []model.User

	exec(db.conn.Model(&users).Select())

	return &users
}

func (db *Db) GetPlugins(usr *model.User) *[]model.Configuration {
	var plugins []model.Configuration

	exec(db.conn.Model(&plugins).
		Where("user_id = ?", usr.Id).
		Select())

	return &plugins
}

func (db *Db) GetConfiguration(usr *model.User, pid string) *model.Configuration {
	cfg := &model.Configuration{
		User: usr,
	}

	db.conn.Model(cfg).
		Where("plugin_guid = ?", pid)

	return cfg;
}

func exec(err error) {
	if err != nil {
		panic(err)
	}
}
