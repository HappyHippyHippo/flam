package flam

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

func Test_Registerer_NewRegisterer(t *testing.T) {
	t.Run("should return a valid registerer", func(t *testing.T) {
		registerer := NewRegisterer()
		assert.NotNil(t, registerer)
		assert.Empty(t, registerer.entries)
	})
}

func Test_Registerer_Queue(t *testing.T) {
	t.Run("should queue a constructor", func(t *testing.T) {
		registerer := NewRegisterer()
		constructor := func() {}
		registerer.Queue(constructor)
		assert.Len(t, registerer.entries, 1)
		assert.NotNil(t, registerer.entries[0].constructor)
		assert.Empty(t, registerer.entries[0].opts)
	})

	t.Run("should queue a constructor with options", func(t *testing.T) {
		registerer := NewRegisterer()
		constructor := func() {}
		opts := []dig.ProvideOption{dig.Name("test")}
		registerer.Queue(constructor, opts...)
		assert.Len(t, registerer.entries, 1)
		assert.NotNil(t, registerer.entries[0].constructor)
		assert.Equal(t, opts, registerer.entries[0].opts)
	})

	t.Run("should allow chaining", func(t *testing.T) {
		registerer := NewRegisterer()
		constructor1 := func() {}
		constructor2 := func() {}
		registerer.Queue(constructor1).Queue(constructor2)
		assert.Len(t, registerer.entries, 2)
	})
}

func Test_Registerer_Run(t *testing.T) {
	t.Run("should run without error on an empty queue", func(t *testing.T) {
		registerer := NewRegisterer()
		container := dig.New()
		assert.NoError(t, registerer.Run(container))
	})

	t.Run("should return an error if providing a non-function constructor fails", func(t *testing.T) {
		registerer := NewRegisterer().Queue(struct{}{})
		container := dig.New()

		assert.Error(t, registerer.Run(container))
	})

	t.Run("should successfully provide a constructor", func(t *testing.T) {
		type Dep struct{}
		constructor := func() *Dep { return &Dep{} }

		registerer := NewRegisterer().Queue(constructor)
		container := dig.New()

		assert.NoError(t, registerer.Run(container))

		err := container.Invoke(func(d *Dep) {
			assert.NotNil(t, d)
		})
		assert.NoError(t, err)
	})

	t.Run("should register a fallible constructor without immediate error", func(t *testing.T) {
		type Dep struct{}
		expectedErr := errors.New("constructor error")
		constructor := func() (*Dep, error) {
			return nil, expectedErr
		}

		registerer := NewRegisterer().Queue(constructor)
		container := dig.New()

		assert.NoError(t, registerer.Run(container))

		err := container.Invoke(func(d *Dep) {})
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("should stop execution on the first Provide error", func(t *testing.T) {
		type Dep struct{}
		constructor1 := struct{}{}
		constructor2 := func() *Dep { return &Dep{} }

		registerer := NewRegisterer().Queue(constructor1).Queue(constructor2)
		container := dig.New()

		assert.Error(t, registerer.Run(container))

		// Check that the second constructor was not provided
		err := container.Invoke(func(d *Dep) {})
		assert.Error(t, err)
	})
}
