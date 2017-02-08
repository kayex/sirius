package extension

import (
	"errors"
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

func (l *StaticLoader) Load(eid sirius.EID) (error, sirius.Extension) {
	switch eid {
	case "thumbs_up":
		return nil, &ThumbsUp{}
	case "ripperino":
		return nil, &Ripperino{}
	case "replacer":
		return nil, &Replacer{}
	case "quotes":
		return nil, &Quotes{}
	case "ip_lookup":
		return nil, &IPLookup{}
	case "geocode":
		return nil, &Geocode{
			APIKey: l.cfg.Maps.APIKey,
		}
	}

	return errors.New(fmt.Sprintf("Invalid eid: %v", eid)), nil
}
