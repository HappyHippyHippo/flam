package flam

import (
	"go.uber.org/dig"
)

type ClosableProvider interface {
	Close(container *dig.Container) error
}
