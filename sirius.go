package sirius

import (
	"github.com/kayex/sirius/config"
	"github.com/kayex/sirius/store"
)

type Sirius struct {
	cfg   *config.Config
	store *store.Store
}

func New(cfg *config.Config, store *store.Store) {
	return Sirius{
		cfg:   cfg,
		store: store,
	}
}

func Boot(s *Sirius) {

}
