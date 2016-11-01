package model

import (
	"time"
)

type User struct {
	Id        int
	CreatedAt time.Time
	Plugins   []Plugin `pg:",many2many:user_plugins"`
}
