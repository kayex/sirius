package extension

import (
	"fmt"
	"github.com/kayex/sirius"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

type Geocode struct {
	APIKey string
}

func (x *Geocode) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	cmd, match := m.Command("address")

	if !match {
		return sirius.NoAction(), nil
	}

	c, err := maps.NewClient(maps.WithAPIKey(x.APIKey))

	if err != nil {
		return nil, err
	}

	r := &maps.GeocodingRequest{
		Address: cmd.Args[0],
	}

	res, err := c.Geocode(context.Background(), r)

	if err != nil {
		return nil, err
	}

	pos := res[0]
	location := pos.Geometry.Location
	formatted := fmt.Sprintf("*%v*\n`(%.6f, %.6f)`", pos.FormattedAddress, location.Lat, location.Lng)

	edit := m.EditText().ReplaceWith(formatted)

	return edit, nil
}
