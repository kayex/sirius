package extension

import (
	"errors"
	"fmt"
	"github.com/kayex/sirius"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

type Geocode struct {
	APIKey string
}

func (gc *Geocode) Run(m sirius.Message, cfg sirius.ExtensionConfig) (error, sirius.MessageAction) {
	address, run := sirius.NewCommand("address").Match(&m)

	if !run {
		return nil, sirius.NoAction()
	}

	c, err := maps.NewClient(maps.WithAPIKey(gc.APIKey))

	if err != nil {
		return errors.New(fmt.Sprintf("Maps client error: %v", err.Error())), nil
	}

	r := &maps.GeocodingRequest{
		Address: address,
	}

	res, err := c.Geocode(context.Background(), r)

	if err != nil {
		return errors.New(fmt.Sprintf("Geocoding error: %v", err.Error())), nil
	}

	pos := res[0]
	location := pos.Geometry.Location
	formatted := fmt.Sprintf("*%v*\n`(%.6f, %.6f)`", pos.FormattedAddress, location.Lat, location.Lng)

	edit := m.EditText().ReplaceWith(formatted)

	return nil, edit
}
