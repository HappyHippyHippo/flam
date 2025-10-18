package flam

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testResource struct{}

func Test_Factory_NewFactory(t *testing.T) {
	t.Run("should return ErrNilReference when config is nil", func(t *testing.T) {
		factory, e := NewFactory[Resource](nil, "path", nil, nil)
		assert.Nil(t, factory)
		assert.ErrorIs(t, e, ErrNilReference)
	})

	t.Run("should return valid factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factory, e := NewFactory[Resource](nil, "path", factoryConfigMock, nil)
		assert.NotNil(t, factory)
		assert.NoError(t, e)
	})
}

func Test_Factory_Close(t *testing.T) {
	t.Run("should correctly close closable resources", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		closerMock := NewCloserMock(ctrl)
		closerMock.EXPECT().Close().Return(nil).Times(1)

		config := Bag{}
		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.EXPECT().Get("path").Return(config).Times(1)

		factory, e := NewFactory[Resource](nil, "path", factoryConfigMock, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		assert.NoError(t, factory.Add("resource1", closerMock))
		assert.NoError(t, factory.Close())
	})

	t.Run("should propagate the first error from closers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("close error")
		closerMock := NewCloserMock(ctrl)
		closerMock.EXPECT().Close().Return(expectedErr).Times(1)

		config := Bag{}
		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.EXPECT().Get("path").Return(config).Times(1)

		factory, e := NewFactory[Resource](nil, "path", factoryConfigMock, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		assert.NoError(t, factory.Add("resource", closerMock))
		assert.ErrorIs(t, factory.Close(), expectedErr)
	})

	t.Run("should not fail with non-closable resources", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := Bag{}
		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.EXPECT().Get("path").Return(config).Times(1)

		factory, e := NewFactory[Resource](nil, "path", factoryConfigMock, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		assert.NoError(t, factory.Add("resource1", &testResource{}))
		assert.NoError(t, factory.Close())
	})
}

func Test_Factory_List(t *testing.T) {
	t.Run("should return a sorted list of ids from config and added entries", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := Bag{"zulu": Bag{}, "alpha": Bag{}}
		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.EXPECT().Get("path").Return(config).Times(2)

		factory, e := NewFactory[Resource](nil, "path", factoryConfigMock, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		assert.NoError(t, factory.Add("charlie", &testResource{}))
		assert.Equal(t, []string{"alpha", "charlie", "zulu"}, factory.List())
	})

	t.Run("should return an empty list when there are no entries", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.EXPECT().Get("path").Return(Bag{}).Times(1)

		factory, e := NewFactory[Resource](nil, "path", factoryConfigMock, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		assert.Empty(t, factory.List())
	})
}

func Test_Factory_Has(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := Bag{"entry1": Bag{}}
	factoryConfigMock := NewFactoryConfigMock(ctrl)
	factoryConfigMock.EXPECT().Get("path").Return(config).AnyTimes()

	factory, e := NewFactory[Resource](nil, "path", factoryConfigMock, nil)
	require.NotNil(t, factory)
	require.NoError(t, e)

	assert.NoError(t, factory.Add("entry2", &testResource{}))

	testCases := []struct {
		name     string
		id       string
		expected bool
	}{
		{"entry in config", "entry1", true},
		{"manually added entry", "entry2", true},
		{"non-existent entry", "nonexistent", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, factory.Has(tc.id))
		})
	}
}

func Test_Factory_Get(t *testing.T) {
	t.Run("should return an error if generation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.EXPECT().Get("path").Return(Bag{}).Times(1)

		factory, e := NewFactory[Resource](nil, "path", factoryConfigMock, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		got, e := factory.Get("nonexistent")
		assert.Nil(t, got)
		assert.ErrorIs(t, e, ErrUnknownResource)
	})
}

func Test_Factory_Generate(t *testing.T) {
	t.Run("should return ErrUnknownResource for an unknown id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.EXPECT().Get("path").Return(Bag{}).Times(1)

		factory, e := NewFactory[Resource](nil, "path", factoryConfigMock, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		got, e := factory.Generate("unknown")
		assert.Nil(t, got)
		assert.ErrorIs(t, e, ErrUnknownResource)
	})

	t.Run("should return an error if creator's Create method fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := Bag{"default": Bag{}}
		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.EXPECT().Get("path").Return(config).Times(1)

		expectedErr := errors.New("creation failed")
		creatorMock := NewResourceCreatorMock[Resource](ctrl)
		creatorMock.EXPECT().Accept(Bag{"id": "default"}).Return(true)
		creatorMock.EXPECT().Create(Bag{"id": "default"}).Return(nil, expectedErr)

		factory, e := NewFactory([]ResourceCreator[Resource]{creatorMock}, "path", factoryConfigMock, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		got, e := factory.Generate("default")
		assert.Nil(t, got)
		assert.ErrorIs(t, e, expectedErr)
	})

	t.Run("should return ErrInvalidResourceConfig if no creator accepts the config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := Bag{"default": Bag{}}
		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.EXPECT().Get("path").Return(config).Times(1)

		creatorMock := NewResourceCreatorMock[Resource](ctrl)
		creatorMock.EXPECT().Accept(Bag{"id": "default"}).Return(false)

		factory, e := NewFactory([]ResourceCreator[Resource]{creatorMock}, "path", factoryConfigMock, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		got, e := factory.Generate("default")
		assert.Nil(t, got)
		assert.ErrorIs(t, e, ErrInvalidResourceConfig)
	})

	t.Run("should return the validation error if there is a validator and it return an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := Bag{"default": Bag{"driver": "test"}}
		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.EXPECT().Get("path").Return(config).Times(1)

		expectedErr := errors.New("validation error")
		validator := func(config Bag) error { return expectedErr }

		factory, e := NewFactory[Resource](nil, "path", factoryConfigMock, validator)
		require.NotNil(t, factory)
		require.NoError(t, e)

		got, e := factory.Generate("default")
		assert.Nil(t, got)
		assert.ErrorIs(t, e, expectedErr)
	})

	t.Run("should generate, cache, and return an entry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := Bag{"default": Bag{"driver": "test"}}
		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.EXPECT().Get("path").Return(config).Times(1)

		resource := &testResource{}

		creatorMock := NewResourceCreatorMock[Resource](ctrl)
		creatorMock.EXPECT().Accept(Bag{"id": "default", "driver": "test"}).Return(true).Times(1)
		creatorMock.EXPECT().Create(Bag{"id": "default", "driver": "test"}).Return(resource, nil).Times(1)

		validator := func(config Bag) error { return nil }

		factory, e := NewFactory([]ResourceCreator[Resource]{creatorMock}, "path", factoryConfigMock, validator)
		require.NotNil(t, factory)
		require.NoError(t, e)

		entry1, e1 := factory.Get("default")
		assert.NoError(t, e1)
		assert.Same(t, resource, entry1)

		entry2, e2 := factory.Get("default")
		assert.NoError(t, e2)
		assert.Same(t, resource, entry2)
	})
}

func Test_Factory_Add(t *testing.T) {
	t.Run("should return ErrNilReference when adding a nil value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factoryConfigMock := NewFactoryConfigMock(ctrl)

		factory, e := NewFactory[Resource](nil, "path", factoryConfigMock, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		assert.ErrorIs(t, factory.Add("nil_resource", nil), ErrNilReference)
	})

	t.Run("should successfully add a resource and allow it to be retrieved", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := Bag{}
		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.EXPECT().Get("path").Return(config).Times(1)

		factory, e := NewFactory[Resource](nil, "path", factoryConfigMock, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		resource := &testResource{}
		assert.NoError(t, factory.Add("my_resource", resource))

		retrieved, e := factory.Get("my_resource")
		assert.NoError(t, e)
		assert.Same(t, resource, retrieved)
	})

	t.Run("should return ErrDuplicateResource when adding a duplicate resource", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := Bag{}
		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.EXPECT().Get("path").Return(config).Times(1)

		factory, e := NewFactory[Resource](nil, "path", factoryConfigMock, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		resource := &testResource{}
		assert.NoError(t, factory.Add("my_resource", resource))
		assert.ErrorIs(t, factory.Add("my_resource", resource), ErrDuplicateResource)
	})
}
