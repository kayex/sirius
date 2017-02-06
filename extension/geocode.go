package extension

import (
	"github.com/kayex/sirius"
	"googlemaps.github.io/maps"
	"golang.org/x/net/context"
	"errors"
	"strings"
	"fmt"
)

type Geocode struct {
	APIKey string
}

func (gc *Geocode) Run(m sirius.Message, cfg sirius.ExtensionConfig) (error, sirius.MessageAction) {
	if !strings.HasPrefix(m.Text, "<address") {
		return nil, sirius.NoAction()
	}

	address := strings.TrimPrefix(m.Text, "<address ")

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
