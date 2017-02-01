package store

import "github.com/kayex/sirius"

type Store interface {
	GetExtensions(usr *sirius.User) *[]sirius.Configuration
	GetUsers() *[]sirius.User
	GetUser(uid string) *sirius.User
	SaveUser(usr *sirius.User)
	SaveConfiguration(cfg *sirius.Configuration)
	UpdateConfiguration(cfg *sirius.Configuration)
}
