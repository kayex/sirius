package db

import (
	"github.com/kayex/sirius"
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

func (db *Db) GetUser(id string) *sirius.User {
	user := sirius.User{Id: id}

	exec(db.conn.Select(&user))

	return &user
}

func (db *Db) SaveUser(usr *sirius.User) {
	exec(db.conn.Insert(usr))
}

func (db *Db) SaveConfiguration(cfg *sirius.Configuration) {
	exec(db.conn.Insert(cfg))
}

func (db *Db) UpdateConfiguration(cfg *sirius.Configuration) {
	_, err := db.conn.Model(cfg).Column("config").Update()

	if err != nil {
		panic(err)
	}
}

func (db *Db) GetUsers() *[]sirius.User {
	var users []sirius.User

	exec(db.conn.Model(&users).Select())

	return &users
}

func (db *Db) GetExtensions(usr *sirius.User) *[]sirius.Configuration {
	var exts []sirius.Configuration

	exec(db.conn.Model(&exts).
		Where("user_id = ?", usr.Id).
		Select())

	return &exts
}

func (db *Db) GetConfiguration(usr *sirius.User, pid string) *sirius.Configuration {
	cfg := &sirius.Configuration{
		User: usr,
	}

	db.conn.Model(cfg).
		Where("ext_guid = ?", pid)

	return cfg
}

func exec(err error) {
	if err != nil {
		panic(err)
	}
}
