package sync

type signal struct{}

type Control struct {
	quit     chan signal
	done     chan error
	finished bool
}

func NewControl() *Control {
	return &Control{
		quit: make(chan signal),
		done: make(chan error, 1),
	}
}

func (c *Control) Finish(err error) {
	if c.finished {
		return
	}

	if err != nil {
		c.done <- err
	}
	close(c.done)
	c.finished = true
}

func (c *Control) Finished() <-chan error {
	return c.done
}
