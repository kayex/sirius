package model

import "time"

type Extension struct {
	Id        int
	Name      string
	CreatedAt time.Time
	Users     []User `pg:"many2many:user_extensions"`
}
