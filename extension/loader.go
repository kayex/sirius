package extension

import (
	"fmt"
	"errors"
	"github.com/kayex/sirius"
)

type StaticExtensionLoader struct{}

func NewStaticExtensionLoader() *StaticExtensionLoader {
	return &StaticExtensionLoader{}
}

func (r *StaticExtensionLoader) Load(eid sirius.EID) (error, sirius.Extension) {
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
