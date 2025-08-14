package flam

import (
	"go.uber.org/dig"
)

type RunnableProvider interface {
	Run(container *dig.Container) error
}
