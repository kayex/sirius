package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type HttpExtension struct {
	Host string
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
	exe := &HttpExecution{
		ext:     x,
		Message: m.Text,
		Config:  cfg,
	}

	res, err := runHttpExecution(exe)

	if err != nil {
		return nil, err
	}

	er := res.ExecutionResult()

	return er.Action, nil
}

func (her *HttpExecutionResult) ExecutionResult() *ExecutionResult {
	er := &ExecutionResult{
		Err:    her.Err,
		Action: NoAction(),
	}

	var edit TextEditAction

	for _, ha := range her.Actions {
		a := ha.ToMessageAction()

		if e, ok := a.(*TextEditAction); ok {
			edit.mutations = append(edit.mutations, e.mutations...)
		}
	}

	er.Action = &edit

	return er
}

func runHttpExecution(exe *HttpExecution) (*HttpExecutionResult, error) {
	cl := &http.Client{}

	pl := new(bytes.Buffer)
	json.NewEncoder(pl).Encode(exe)

	req, err := http.NewRequest("POST", exe.ext.Host, pl)
	req.Header.Set("Content-Type", "application/json")

	res, err := cl.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var her HttpExecutionResult

	if err = json.NewDecoder(res.Body).Decode(&her); err != nil {
		return nil, err
	}

	return &her, nil
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

		edit.Substitute(search, sub)

		return edit
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

		return edit.ReplaceWith(t)
	}

	return NoAction()
}
