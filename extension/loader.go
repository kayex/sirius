package extension

import (
	"fmt"
	"github.com/kayex/sirius"
	"github.com/kayex/sirius/config"
)

type StaticLoader struct {
	cfg config.AppConfig
}

func NewStaticLoader(cfg config.AppConfig) *StaticLoader {
	return &StaticLoader{
		cfg: cfg,
	}
}

func (l *StaticLoader) Load(eid sirius.EID) (sirius.Extension, error) {
	switch eid {
	case "thumbs_up":
		return &ThumbsUp{}, nil
	case "ripperino":
		return &Ripperino{}, nil
	case "replacer":
		return &Replacer{}, nil
	case "quotes":
		return &Quotes{}, nil
	case "ip_lookup":
		return &IPLookup{}, nil
	case "censor":
		return &Censor{}, nil
	case "google":
		return &Google{}, nil
	case "sin":
		return &Sin{}, nil
	case "geocode":
		return &Geocode{
			APIKey: l.cfg.Maps.APIKey,
		}, nil
	}

	return nil, fmt.Errorf("invalid eid: %v", eid)
}
