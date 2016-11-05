package core

type Router struct {
	runner *Runner
}

func NewRouter(runner *Runner) Router {
	return Router{
		runner: runner,
	}
}
