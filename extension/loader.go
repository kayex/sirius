package extension

import "github.com/kayex/sirius"

type StaticExtensionLoader struct{}

func NewStaticExtensionLoader() *StaticExtensionLoader {
	return &StaticExtensionLoader{}
}

func (r *StaticExtensionLoader) Load(eid sirius.EID) sirius.Extension {
	switch eid {
	case "thumbs_up":
		return &ThumbsUp{}
	case "ripperino":
		return &Ripperino{}
	case "replacer":
		return &Replacer{}
	case "quotes":
		return &Quotes{}
	}

	panic("Invalid eid: " + eid)
}
