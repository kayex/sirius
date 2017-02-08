package extension

import (
	"github.com/kayex/sirius/config"
	"github.com/kayex/sirius"
	"errors"
	"fmt"
)

type StaticLoader struct {
	config config.AppConfig
}

func NewStaticLoader(cfg config.AppConfig) *StaticLoader {
	return &StaticLoader{
		config: cfg,
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
	case "geocode":
		return &Geocode{
			APIKey: l.config.Maps.APIKey,
		}, nil
	}

	return nil, errors.New(fmt.Sprintf("Invalid eid: %v", eid))
}
