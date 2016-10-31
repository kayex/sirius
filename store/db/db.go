package db

import (
	"gopkg.in/pg.v5"
	"github.com/Epoch2/slack-sirius/model"
)

type Db struct {
	conn *pg.DB
}

func Connect(user string, password string) Db {
	db := pg.Connect(&pg.Options{
		User: user,
		Password: password,
	})

	return Db{
		conn: db,
	}
}

func (db *Db) GetUsers() *[]model.User {
	var users []model.User

	exec(db.conn.Model(&users).Select())

	return &users
}

func (db *Db) GetPlugins(u *model.User) *[]model.Plugin {
	var plugins []model.Plugin

	exec(db.conn.Model(&plugins).Select())

	return &plugins
}

func exec(err error) {
	if (err != nil) {
		panic(err)
	}
}
