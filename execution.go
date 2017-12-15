package sirius

import (
	"time"
)

type Execution struct {
	Ext Extension
	Cfg ExtensionConfig
}

type ExecutionResult struct {
	Err    error
	Action MessageAction
}

type Executor struct {
	exe []Execution
	timeout time.Duration
	runner ExtensionRunner
	loader ExtensionLoader
}

func NewExecutor(loader ExtensionLoader, runner ExtensionRunner, timeout time.Duration) *Executor {
	if loader == nil {
		panic("executor requires an explicit loader")
	}
	if runner == nil {
		panic("executor requires an explicit runner")
	}
	if timeout == 0 {
		timeout = time.Second * 2
	}

	return &Executor{
		loader: loader,
		runner: runner,
		timeout: timeout,
	}
}

func (e *Executor) Run(msg Message) <-chan ExecutionResult {
	res := make(chan ExecutionResult, len(e.exe))
	e.runner.Run(msg, e.exe, res, e.timeout)

	return res
}

func (e *Executor) LoadFromSettings(s Settings) (error) {
	var exe []Execution

	for _, cfg := range s {
		ex, err := e.createExecution(&cfg)
		if err != nil {
			return err
		}

		exe = append(exe, *ex)
	}

	e.exe = exe

	return nil
}

func (e *Executor) createExecution(cfg *Configuration) (*Execution, error) {
	var x Extension

	// Check for HTTP extensions
	if cfg.URL != "" {
		x = NewHttpExtension(cfg.URL, nil)
	} else {
		ex, err := e.loader.Load(cfg.EID)

		if err != nil {
			panic(err)
		}

		x = ex
	}

	return NewExecution(x, cfg.Cfg), nil
}

type ExtensionRunner interface {
	Run(msg Message, exe []Execution, res chan<- ExecutionResult, timeout time.Duration)
}

type AsyncRunner struct{}

func NewExecution(x Extension, cfg ExtensionConfig) *Execution {
	return &Execution{
		Ext: x,
		Cfg: cfg,
	}
}

func NewAsyncRunner() *AsyncRunner {
	return &AsyncRunner{}
}

// Run executes all extensions in exe and returns all ExecutionResults that
// are received before timeout has elapsed.
func (r *AsyncRunner) Run(msg Message, exe []Execution, res chan<- ExecutionResult, timeout time.Duration) {
	er := make(chan ExecutionResult, len(exe))

	for i := range exe {
		go func(msg Message, ex *Execution) {
			a, err := ex.Ext.Run(msg, ex.Cfg)

			er <- ExecutionResult{
				Err:    err,
				Action: a,
			}
		}(msg, &exe[i])
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
