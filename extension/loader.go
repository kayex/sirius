package extension

import (
	"fmt"
	"errors"
	"github.com/kayex/sirius"
	"github.com/kayex/sirius/config"
)

type StaticExtensionLoader struct{
	cfg config.AppConfig
}

func NewStaticExtensionLoader(cfg config.AppConfig) *StaticExtensionLoader {
	return &StaticExtensionLoader{
		cfg: cfg,
	}
}

func (l *StaticExtensionLoader) Load(eid sirius.EID) (error, sirius.Extension) {
	switch eid {
	case "thumbs_up":
		return nil, &ThumbsUp{}
	case "ripperino":
		return nil, &Ripperino{}
	case "replacer":
		return nil, &Replacer{}
	case "quotes":
		return nil, &Quotes{}
	}

	return errors.New(fmt.Sprintf("Invalid eid: %v", eid)), nil
}
