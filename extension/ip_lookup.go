package extension

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kayex/sirius"
	"github.com/kayex/sirius/text"
	"net/http"
)

type IPLookup struct{}

func (*IPLookup) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	cmd, match := m.Command("ip")

	if !match {
		return sirius.NoAction(), nil
	}

	ip := cmd.Arg(0)

	if ip == "" {
		return sirius.NoAction(), nil
	}

	err, info := ipLookup(ip)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("IP Lookup error: %v", err))
	}

	if !info.complete() {
		return sirius.NoAction(), nil
	}

	// Output format:
	//
	// `IP`
	// City, Country (`COUNTRY CODE`)
	// ISP
	output := fmt.Sprintf("%v\n"+
		"%v, %v (%v)\n"+
		"%v",
		text.Code(ip),
		info.City,
		info.Country,
		text.Code(info.CountryCode),
		info.ISP)
	edit := m.EditText().Set(output)

	return edit, nil
}

func ipLookup(ip string) (error, *iPInfo) {
	url := "http://ip-api.com/json/" + ip

	c := &http.Client{}
	r, err := c.Get(url)

	if err != nil {
		return err, nil
	}
	defer r.Body.Close()

	var info iPInfo
	err = json.NewDecoder(r.Body).Decode(&info)

	return nil, &info
}

type iPInfo struct {
	City        string `json:"city"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	ISP         string `json:"isp"`
}

// complete indicates whether i has all data fields filled.
func (i *iPInfo) complete() bool {
	return !(i.City == "" || i.Country == "" || i.CountryCode == "" || i.ISP == "")
}
