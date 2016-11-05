package store

import "github.com/kayex/sirius/model"

type Store interface {
	GetPlugins(usr *model.User) *[]model.Configuration
	GetUsers() *[]model.User
	GetUser(uid string) *model.User
	SaveUser(usr *model.User)
	SaveConfiguration(cfg *model.Configuration)
	UpdateConfiguration(cfg *model.Configuration)
}
