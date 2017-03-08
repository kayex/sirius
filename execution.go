package sirius

import (
	"time"
)

type Execution struct {
	Ext Extension
	Msg Message
	Cfg ExtensionConfig
}

type ExecutionResult struct {
	Err    error
	Action MessageAction
}

type ExtensionRunner interface {
	Run(exe []Execution, res chan<- ExecutionResult, timeout time.Duration)
}

type AsyncRunner struct{}

func NewExecution(x Extension, m Message, cfg ExtensionConfig) *Execution {
	return &Execution{
		Ext: x,
		Msg: m,
		Cfg: cfg,
	}
}

func NewAsyncRunner() *AsyncRunner {
	return &AsyncRunner{}
}

// Run executes all extensions in exe and returns all ExecutionResults that
// are received before timeout has elapsed.
func (r *AsyncRunner) Run(exe []Execution, res chan<- ExecutionResult, timeout time.Duration) {
	er := make(chan ExecutionResult, len(exe))

	for i := range exe {
		go func(ex *Execution) {
			a, err := ex.Ext.Run(ex.Msg, ex.Cfg)

			er <- ExecutionResult{
				Err:    err,
				Action: a,
			}
		}(&exe[i])
	}

Execution:
	for range exe {
		select {
		case <-time.After(timeout):
			break Execution
		case res <- <-er:
		}
	}

	close(res)
}
