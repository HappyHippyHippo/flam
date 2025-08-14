package flam

import (
	"slices"
	"sync"

	"go.uber.org/dig"
)

type Application interface {
	Container() *dig.Container
	Register(provider Provider) error
	Boot() error
	Run() error
	Close() error
}

type application struct {
	locker    sync.Locker
	container *dig.Container
	providers []Provider
	isBooted  bool
}

func NewApplication() Application {
	return &application{
		locker:    &sync.Mutex{},
		container: dig.New(),
		providers: []Provider{},
		isBooted:  false,
	}
}

func (app *application) Container() *dig.Container {
	return app.container
}

func (app *application) Register(
	provider Provider,
) error {
	if provider == nil {
		return newErrNilReference("provider")
	}

	app.locker.Lock()
	defer app.locker.Unlock()

	for _, registered := range app.providers {
		if registered.Id() == provider.Id() {
			return newErrDuplicateProvider(provider.Id())
		}
	}

	if e := provider.Register(app.container); e != nil {
		return e
	}
	app.providers = append(app.providers, provider)

	return nil
}

func (app *application) Boot() error {
	if app.isBooted {
		return nil
	}

	app.locker.Lock()
	defer app.locker.Unlock()

	for _, registered := range app.providers {
		if bootable, ok := registered.(BootableProvider); ok {
			if e := bootable.Boot(app.container); e != nil {
				return e
			}
		}
	}

	app.isBooted = true

	return nil
}

func (app *application) Run() error {
	if !app.isBooted {
		if e := app.Boot(); e != nil {
			return e
		}
	}

	app.locker.Lock()
	defer app.locker.Unlock()

	for _, registered := range app.providers {
		if runnable, ok := registered.(RunnableProvider); ok {
			if e := runnable.Run(app.container); e != nil {
				return e
			}
		}
	}

	return nil
}

func (app *application) Close() error {
	app.locker.Lock()
	defer app.locker.Unlock()

	slices.Reverse(app.providers)

	var e error
	for _, registered := range app.providers {
		if closable, ok := registered.(ClosableProvider); ok {
			if closeError := closable.Close(app.container); closeError != nil {
				e = closeError
			}
		}
	}

	return e
}
