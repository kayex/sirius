package extension

import (
	"github.com/kayex/sirius"
	"encoding/json"
	"net/http"
	"errors"
	"fmt"
)

type IPLookup struct {}

func (ipl *IPLookup) Run(m sirius.Message, cfg sirius.ExtensionConfig) (error, sirius.MessageAction) {
	ip, run := sirius.NewCommand("ip").Match(&m)

	if !run {
		return nil, sirius.NoAction()
	}

	var lookup map[string]interface{}

	err, lookup := ipLookup(ip)

	if err != nil {
		return errors.New(fmt.Sprintf("IP Lookup error: %v", err)), nil
	}

	output := fmt.Sprintf("`%v`\n%v, %v (`%v`)\n%v", ip, lookup["city"], lookup["country"], lookup["countryCode"], lookup["isp"])
	edit := m.EditText().ReplaceWith(output)

	return nil, edit
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
