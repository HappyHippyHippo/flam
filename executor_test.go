package flam

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"
)

func Test_Executor_NewExecutor(t *testing.T) {
	t.Run("should return a valid executor", func(t *testing.T) {
		executor := NewExecutor()
		assert.NotNil(t, executor)
		assert.Empty(t, executor.entries)
	})
}

func Test_Executor_Queue(t *testing.T) {
	t.Run("should queue a callback", func(t *testing.T) {
		executor := NewExecutor()
		callback := func() {}
		executor.Queue(callback)
		assert.Len(t, executor.entries, 1)
		assert.NotNil(t, executor.entries[0].callback)
	})

	t.Run("should allow chaining", func(t *testing.T) {
		executor := NewExecutor()
		callback1 := func() {}
		callback2 := func() {}
		executor.Queue(callback1).Queue(callback2)
		assert.Len(t, executor.entries, 2)
	})
}

func Test_Executor_Run(t *testing.T) {
	t.Run("should run without error on an empty queue", func(t *testing.T) {
		executor := NewExecutor()
		container := dig.New()
		assert.NoError(t, executor.Run(container))
	})

	t.Run("should return an error if callback invocation fails", func(t *testing.T) {
		expectedErr := errors.New("invoke error")
		callback := func() error {
			return expectedErr
		}

		executor := NewExecutor().Queue(callback)
		container := dig.New()

		assert.ErrorIs(t, executor.Run(container), expectedErr)
	})

	t.Run("should successfully invoke a callback with dependencies", func(t *testing.T) {
		type Dep struct{}
		var wasCalled bool

		callback := func(d *Dep) {
			assert.NotNil(t, d)
			wasCalled = true
		}

		executor := NewExecutor().Queue(callback)
		container := dig.New()
		require.NoError(t, container.Provide(func() *Dep { return &Dep{} }))

		assert.NoError(t, executor.Run(container))
		assert.True(t, wasCalled)
	})

	t.Run("should stop execution on the first error", func(t *testing.T) {
		var firstCallbackCalled bool
		var secondCallbackCalled bool
		expectedErr := errors.New("invoke error")

		callback1 := func() error {
			firstCallbackCalled = true
			return expectedErr
		}
		callback2 := func() {
			secondCallbackCalled = true
		}

		executor := NewExecutor().Queue(callback1).Queue(callback2)
		container := dig.New()

		assert.ErrorIs(t, executor.Run(container), expectedErr)
		assert.True(t, firstCallbackCalled)
		assert.False(t, secondCallbackCalled)
	})
}
