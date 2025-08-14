package flam

import (
	"io"
	"reflect"
	"slices"
	"strings"
	"sync"
)

type Factory[R Resource] interface {
	Close() error
	List() []string
	Has(id string) bool
	Get(id string) (R, error)
	Generate(id string) (R, error)
	Add(id string, value R) error
}

type factory[R Resource] struct {
	locker          *sync.Mutex
	creators        []ResourceCreator[R]
	configPath      string
	config          FactoryConfig
	configValidator FactoryConfigValidator
	entries         map[string]R
}

func NewFactory[R Resource](
	creators []ResourceCreator[R],
	configPath string,
	config FactoryConfig,
	configValidator FactoryConfigValidator,
) (Factory[R], error) {
	if config == nil {
		return nil, newErrNilReference("config")
	}

	return &factory[R]{
		locker:          &sync.Mutex{},
		creators:        creators,
		configPath:      configPath,
		config:          config,
		configValidator: configValidator,
		entries:         map[string]R{},
	}, nil
}

func (factory *factory[R]) Close() error {
	factory.locker.Lock()
	defer factory.locker.Unlock()

	for _, entry := range factory.entries {
		if closer, ok := any(entry).(io.Closer); ok {
			if e := closer.Close(); e != nil {
				return e
			}
		}
	}

	return nil
}

func (factory *factory[R]) List() []string {
	factory.locker.Lock()
	defer factory.locker.Unlock()

	factoryConfig := factory.config.Get(factory.configPath)
	ids := factoryConfig.Entries()

	for k := range factory.entries {
		if !slices.Contains(ids, k) {
			ids = append(ids, k)
		}
	}

	slices.SortFunc(ids, strings.Compare)

	return ids
}

func (factory *factory[R]) Has(
	id string,
) bool {
	factory.locker.Lock()
	defer factory.locker.Unlock()

	if _, ok := factory.entries[id]; ok {
		return true
	}

	factoryConfig := factory.config.Get(factory.configPath)

	return factoryConfig.Bag(id) != nil
}

func (factory *factory[R]) Get(
	id string,
) (R, error) {
	factory.locker.Lock()
	if entry, ok := factory.entries[id]; ok {
		factory.locker.Unlock()
		return entry, nil
	}
	factory.locker.Unlock()

	entry, e := factory.Generate(id)
	if e != nil {
		return entry, e
	}

	factory.locker.Lock()
	factory.entries[id] = entry
	factory.locker.Unlock()

	return entry, nil
}

func (factory *factory[R]) Generate(
	id string,
) (R, error) {
	factory.locker.Lock()
	defer factory.locker.Unlock()

	var zero R

	factoryConfig := factory.config.Get(factory.configPath)
	config := factoryConfig.Bag(id)
	if config == nil {
		return zero, newErrUnknownResource(
			reflect.TypeFor[R]().Name(),
			id)
	}
	_ = config.Set("id", id)

	if factory.configValidator != nil {
		if e := factory.configValidator(config); e != nil {
			return zero, e
		}
	}

	for _, creator := range factory.creators {
		if creator.Accept(config) {
			return creator.Create(config)
		}
	}

	return zero, newErrInvalidResourceConfig(
		reflect.TypeFor[R]().Name(),
		id,
		config)
}

func (factory *factory[R]) Add(
	id string,
	value R,
) error {
	switch {
	case any(value) == nil:
		return newErrNilReference("value")
	case factory.Has(id):
		return newErrDuplicateResource(id)
	}

	factory.locker.Lock()
	defer factory.locker.Unlock()

	factory.entries[id] = value

	return nil
}
