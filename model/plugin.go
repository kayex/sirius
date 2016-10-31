package model

import "time"

type Plugin struct {
	Id int
	Name string
	CreatedAt time.Time
	Users []User `pg:"many2many:user_plugins"`
}
