package sirius

type EID string

func LoadExtension(eid EID) Extension {
	switch eid {
	case "thumbs_up":
		return &ThumbsUp{}
	case "ripperino":
		return &Ripperino{}
	case "replacer":
		return &Replacer{}
	}

	panic("Invalid eid: " + eid)
}
