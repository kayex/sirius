package sirius

import (
	"github.com/kayex/sirius/config"
	"github.com/kayex/sirius/store"
)

type Sirius struct {
	cfg   *config.Config
	store *store.Store
}

func Init(cfg *config.Config, store *store.Store) {
	return Sirius{
		cfg:   cfg,
		store: store,
	}
}

func (s *Sirius) Boot() {
}
