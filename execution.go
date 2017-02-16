package sirius

import (
	"time"
)

type ExtensionRunner interface {
	Run([]Execution, chan<- ExecutionResult, time.Duration)
}

type Execution struct {
	Ext Extension
	Msg Message
	Cfg ExtensionConfig
}

type ExecutionResult struct {
	Err    error
	Action MessageAction
}

type AsyncRunner struct {
	Actions chan MessageAction
}

func NewExecution(x Extension, m Message, cfg ExtensionConfig) *Execution {
	return &Execution{
		Ext: x,
		Msg: m,
		Cfg: cfg,
	}
}

func NewAsyncRunner() *AsyncRunner {
	return &AsyncRunner{
		Actions: make(chan MessageAction),
	}
}

func (r *AsyncRunner) Run(exe []Execution, res chan<- ExecutionResult, timeout time.Duration) {
	er := make(chan ExecutionResult, len(exe))

	for _, e := range exe {
		r.execute(e, er)
	}

ActionReceive:
	for range exe {
		select {
		case <-time.After(timeout):
			break ActionReceive
		case res <- <-er:
		}
	}

	close(res)
}

func (r *AsyncRunner) execute(e Execution, res chan<- ExecutionResult) {
	go func() {
		a, err := e.Ext.Run(e.Msg, e.Cfg)

		res <- ExecutionResult{
			Err:    err,
			Action: a,
		}
	}()
}
