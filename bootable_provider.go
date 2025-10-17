package flam

import (
	"go.uber.org/dig"
)

type BootableProvider interface {
	Boot(container *dig.Container) error
}
