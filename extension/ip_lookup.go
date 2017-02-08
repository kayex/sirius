package extension

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kayex/sirius"
	"net/http"
)

type IPLookup struct{}

func (ipl *IPLookup) Run(m sirius.Message, cfg sirius.ExtensionConfig) (sirius.MessageAction, error) {
	ip, run := sirius.NewCommand("ip").Match(&m)

	if !run {
		return sirius.NoAction(), nil
	}

	var lookup map[string]interface{}

	err, lookup := ipLookup(ip)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("IP Lookup error: %v", err))
	}

	output := fmt.Sprintf("`%v`\n%v, %v (`%v`)\n%v", ip, lookup["city"], lookup["country"], lookup["countryCode"], lookup["isp"])
	edit := m.EditText().ReplaceWith(output)

	return edit, nil
}

func ipLookup(ip string) (error, map[string]interface{}) {
	url := "http://ip-api.com/json/" + ip

	c := &http.Client{}
	r, err := c.Get(url)

	if err != nil {
		return err, nil
	}
	defer r.Body.Close()

	var lookup map[string]interface{}
	json.NewDecoder(r.Body).Decode(&lookup)

	return nil, lookup
}
