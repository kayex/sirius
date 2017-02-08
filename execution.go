package sirius

import (
	"time"
)

type ExtensionRunner interface {
	Run([]Execution, chan<- ExecutionResult, time.Duration)
}

type Execution struct {
	Extension Extension
	Message   Message
	Config    ExtensionConfig
}

type ExecutionResult struct {
	Error  error
	Action MessageAction
}

type AsyncRunner struct {
	Actions chan MessageAction
}

func NewExecution(x Extension, m Message, cfg ExtensionConfig) *Execution {
	return &Execution{
		Extension: x,
		Message:   m,
		Config:    cfg,
	}
}

func NewAsyncRunner() *AsyncRunner {
	return &AsyncRunner{
		Actions: make(chan MessageAction),
	}
}

func (r *AsyncRunner) Run(exe []Execution, res chan<- ExecutionResult, timeout time.Duration) {
	act := make(chan ExecutionResult, len(exe))

	for _, e := range exe {
		r.execute(e, act)
	}

ActionReceive:
	for range exe {
		select {
		case res <- <-act:
		case <-time.After(timeout):
			break ActionReceive
		}
	}

	close(res)
}

func (r *AsyncRunner) execute(e Execution, res chan<- ExecutionResult) {
	go func() {
		err, a := e.Extension.Run(e.Message, e.Config)

		r := ExecutionResult{
			Error:  err,
			Action: a,
		}

		res <- r
	}()
}
