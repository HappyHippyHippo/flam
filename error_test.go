package flam

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewError(t *testing.T) {
	t.Run("should create error with a simple message", func(t *testing.T) {
		msg := "test error"
		e := NewError(msg)
		assert.Equal(t, msg, e.Error())
		assert.NotNil(t, e.Context())
	})

	t.Run("should create error with a message and context", func(t *testing.T) {
		msg := "test error"
		ctx := Bag{"key": "value"}
		e := NewError(msg, ctx)
		assert.Equal(t, msg, e.Error())
		assert.Equal(t, ctx, *e.Context())
	})
}

func Test_NewErrorFrom(t *testing.T) {
	baseErr := errors.New("base error")

	t.Run("should wrapping an error with a message", func(t *testing.T) {
		msg := "additional context"
		e := NewErrorFrom(baseErr, msg)
		assert.ErrorIs(t, e, baseErr)
		assert.Equal(t, "base error: additional context", e.Error())
	})

	t.Run("should wrapping an error with a message and context", func(t *testing.T) {
		msg := "context"
		ctx := Bag{"key": "value"}
		e := NewErrorFrom(baseErr, msg, ctx)
		assert.ErrorIs(t, e, baseErr)
		assert.Equal(t, "base error: context", e.Error())
		assert.Equal(t, ctx, *e.Context())
	})
}

func Test_FlamError_GetCode(t *testing.T) {
	e := NewError("test").SetCode(404)
	assert.Equal(t, 404, e.GetCode())
}

func Test_FlamError_SetCode(t *testing.T) {
	e := NewError("test")
	assert.Same(t, e, e.SetCode(500))
	assert.Equal(t, 500, e.GetCode())
}

func Test_FlamError_Context(t *testing.T) {
	assert.Equal(t, &Bag{"detail": "not found"}, NewError("test").Set("detail", "not found").Context())
}

func Test_FlamError_Set(t *testing.T) {
	e := NewError("test")
	assert.Same(t, e, e.Set("user.id", 123))
	assert.Equal(t, 123, e.Get("user.id"))
}

func Test_FlamError_Get(t *testing.T) {
	e := NewError("test").Set("detail", "found")

	t.Run("should retrieve stored valid", func(t *testing.T) {
		assert.Equal(t, "found", e.Get("detail"))
	})

	t.Run("should return nil if the value is not present", func(t *testing.T) {
		assert.Nil(t, e.Get("non.existent"))
	})

	t.Run("should return given default if the value is not present", func(t *testing.T) {
		assert.Equal(t, "default", e.Get("non.existent", "default"))
	})
}

func Test_FlamError_Unwrap(t *testing.T) {
	base := errors.New("base")

	assert.ErrorIs(t, NewErrorFrom(base, "wrapped").Unwrap(), base)
}
