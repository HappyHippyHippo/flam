package flam

import (
	"go.uber.org/dig"
)

type Provider interface {
	Id() string
	Register(container *dig.Container) error
}
