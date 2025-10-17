package flam

import "go.uber.org/dig"

type executorEntry struct {
	callback any
}

type Executor struct {
	entries []executorEntry
}

func NewExecutor() *Executor {
	return &Executor{
		entries: make([]executorEntry, 0),
	}
}

func (executor *Executor) Queue(
	callback any,
) *Executor {
	executor.entries = append(
		executor.entries,
		executorEntry{
			callback: callback,
		})
	return executor
}

func (executor *Executor) Run(
	container *dig.Container,
) error {
	for _, entry := range executor.entries {
		if e := container.Invoke(entry.callback); e != nil {
			return e
		}
	}
	return nil
}
