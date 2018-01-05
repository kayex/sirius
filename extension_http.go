package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type HttpExtension struct {
	client *http.Client
	Host   string
}

func NewHttpExtension(host string, client *http.Client) *HttpExtension {
	x := &HttpExtension{
		Host:   host,
		client: client,
	}

	if x.client == nil {
		x.client = &http.Client{}
	}

	return x
}

type HttpExecutionResult struct {
	Err     error
	Actions []HttpMessageAction
}

type HttpMessageAction map[string]interface{}

type HttpExecution struct {
	ext     *HttpExtension
	Message string          `json:"message"`
	Config  ExtensionConfig `json:"config"`
}

func (x *HttpExtension) Run(m Message, cfg ExtensionConfig) (MessageAction, error) {
	res, err := runHttpExecution(&HttpExecution{
		ext:     x,
		Message: m.Text,
		Config:  cfg,
	})

	if err != nil {
		return nil, err
	}

	return res.ExecutionResult().Action, nil
}

func (her *HttpExecutionResult) ExecutionResult() *ExecutionResult {
	er := &ExecutionResult{
		Err: her.Err,
	}

	var edit TextEditAction

	for _, ha := range her.Actions {
		a := ha.ToMessageAction()

		if e, ok := a.(*TextEditAction); ok {
			edit.mutations = append(edit.mutations, e.mutations...)
		}
	}
	er.Action = &edit

	if er.Action == nil {
		er.Action = NoAction()
	}
	return er
}

func runHttpExecution(exe *HttpExecution) (er *HttpExecutionResult, err error) {
	pl := new(bytes.Buffer)
	json.NewEncoder(pl).Encode(exe)

	req, err := http.NewRequest("POST", exe.ext.Host, pl)
	req.Header.Set("Content-Type", "application/json")
	res, err := exe.ext.client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return
	}

	return er, json.NewDecoder(res.Body).Decode(&er)
}

func (act HttpMessageAction) ToMessageAction() MessageAction {
	switch act["type"] {
	case "substitute":
		edit := &TextEditAction{}

		search, a := act["search"].(string)
		sub, b := act["sub"].(string)
		if !(a && b) {
			return NoAction()
		}
		return edit.Substitute(search, sub)
	case "append":
		edit := &TextEditAction{}

		a, ok := act["suffix"].(string)
		if !ok {
			return NoAction()
		}
		return edit.Append(a)
	case "prepend":
		edit := &TextEditAction{}

		p, ok := act["prefix"].(string)
		if !ok {
			return NoAction()
		}
		return edit.Prepend(p)
	case "replace_with":
		edit := &TextEditAction{}

		t, ok := act["text"].(string)
		if !ok {
			return NoAction()
		}
		return edit.Set(t)
	}
	return NoAction()
}
