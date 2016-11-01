package store

import "github.com/kayex/sirius/model"

type Store interface {
	GetPlugins(*model.User) *[]model.Plugin
	GetUsers() *[]model.User
}
