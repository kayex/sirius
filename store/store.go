package store

import "github.com/Epoch2/slack-sirius/model"

type Store interface {
	GetPlugins(*model.User) *[]model.Plugin
	GetUsers() *[]model.User
}
