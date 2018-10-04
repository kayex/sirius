package sirius

import "time"

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
	e := &Executor{
		loader:  loader,
		timeout: timeout,
	}

	if e.timeout == 0 {
		e.timeout = time.Second * 2
	}

	return e
}

// Load sets up the executor by reading s and creating the appropriate extensions.
func (e *Executor) Load(s Profile) error {
	exs, err := LoadFromSettings(e.loader, s)
	if err != nil {
		return err
	}
	e.exs = exs

	return nil
}

func (e *Executor) Run(msg Message) <-chan ExecutionResult {
	res := make(chan ExecutionResult, len(e.exs))
	e.run(msg, e.exs, res)

	return res
}

// run executes all extensions in exs and returns all ExecutionResults that are received before timeout has elapsed.
func (e *Executor) run(msg Message, exs []ConfigExtension, res chan<- ExecutionResult) {
	defer close(res)
	er := make(chan ExecutionResult, len(exs))

	for i := range exs {
		go func(msg Message, ex *ConfigExtension) {
			a, err := ex.Run(msg)

			er <- ExecutionResult{
				Err:    err,
				Action: a,
			}
		}(msg, &exs[i])
	}

	for range exs {
		select {
		case <-time.After(e.timeout):
			return
		case res <- <-er:
		}
	}
}
