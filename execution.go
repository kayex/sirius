package sirius

import (
	"time"
)

type ExecutionResult struct {
	Err    error
	Action MessageAction
}

type Executor struct {
	exs     []ConfigExtension
	timeout time.Duration
	loader  ExtensionLoader
}

func NewExecutor(loader ExtensionLoader, timeout time.Duration) *Executor {
	if timeout == 0 {
		timeout = time.Second * 2
	}

	return &Executor{
		loader: loader,
		timeout: timeout,
	}
}

func (e *Executor) FromSettings(s Settings) error {
	exs, err := LoadFromSettings(e.loader, s)
	if err != nil {
		return err
	}
	e.exs = exs

	return nil
}

func (e *Executor) RunExtensions(msg Message) <-chan ExecutionResult {
	res := make(chan ExecutionResult, len(e.exs))
	e.Run(msg, e.exs, res)

	return res
}

// Run executes all extensions in exs and returns all ExecutionResults that
// are received before timeout has elapsed.
func (e *Executor) Run(msg Message, exe []ConfigExtension, res chan<- ExecutionResult) {
	defer close(res)
	er := make(chan ExecutionResult, len(exe))

	for i := range exe {
		go func(msg Message, ex *ConfigExtension) {
			a, err := ex.Run(msg)

			er <- ExecutionResult{
				Err:    err,
				Action: a,
			}
		}(msg, &exe[i])
	}

Execution:
	for range exe {
		select {
		case <-time.After(e.timeout):
			break Execution
		case res<- <-er:
		}
	}
}
