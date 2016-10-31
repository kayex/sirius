package sirius

import (
	"github.com/Epoch2/slack-sirius/store"
	"github.com/Epoch2/slack-sirius/config"
)

type Sirius struct {
	cfg *config.Config
	store *store.Store
}

func New(cfg *config.Config, store *store.Store) {
	return Sirius{
		cfg: cfg,
		store: store,
	}
}

func Boot(s *Sirius) {

}
