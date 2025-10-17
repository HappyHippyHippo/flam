package flam

import "go.uber.org/dig"

type registererEntry struct {
	constructor any
	opts        []dig.ProvideOption
}

type Registerer struct {
	entries []registererEntry
}

func NewRegisterer() *Registerer {
	return &Registerer{
		entries: make([]registererEntry, 0),
	}
}

func (registerer *Registerer) Queue(
	constructor any,
	opts ...dig.ProvideOption,
) *Registerer {
	registerer.entries = append(
		registerer.entries,
		registererEntry{
			constructor: constructor,
			opts:        opts,
		})
	return registerer
}

func (registerer *Registerer) Run(
	container *dig.Container,
) error {
	for _, entry := range registerer.entries {
		if e := container.Provide(
			entry.constructor,
			entry.opts...,
		); e != nil {
			return e
		}
	}
	return nil
}
