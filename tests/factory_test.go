package tests

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/happyhippyhippo/flam"
	"github.com/happyhippyhippo/flam/tests/mocks"
)

type testResource struct{}

func Test_Factory_NewFactory(t *testing.T) {
	t.Run("should return ErrNilReference when config is nil", func(t *testing.T) {
		factory, e := flam.NewFactory[flam.Resource](nil, "path", nil, nil)
		assert.Nil(t, factory)
		assert.ErrorIs(t, e, flam.ErrNilReference)
	})

	t.Run("should return valid factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factory, e := flam.NewFactory[flam.Resource](nil, "path", factoryConfig, nil)
		assert.NotNil(t, factory)
		assert.NoError(t, e)
	})
}

func Test_Factory_Close(t *testing.T) {
	t.Run("should correctly close closable resources", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		closer := mocks.NewCloser(ctrl)
		closer.EXPECT().Close().Return(nil).Times(1)

		config := flam.Bag{}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get("path").Return(config).Times(1)

		factory, e := flam.NewFactory[flam.Resource](nil, "path", factoryConfig, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		assert.NoError(t, factory.Add("resource1", closer))
		assert.NoError(t, factory.Close())
	})

	t.Run("should propagate the first error from closers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("close error")
		closer := mocks.NewCloser(ctrl)
		closer.EXPECT().Close().Return(expectedErr).Times(1)

		config := flam.Bag{}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get("path").Return(config).Times(1)

		factory, e := flam.NewFactory[flam.Resource](nil, "path", factoryConfig, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		assert.NoError(t, factory.Add("resource", closer))
		assert.ErrorIs(t, factory.Close(), expectedErr)
	})

	t.Run("should not fail with non-closable resources", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := flam.Bag{}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get("path").Return(config).Times(1)

		factory, e := flam.NewFactory[flam.Resource](nil, "path", factoryConfig, nil)
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

		config := flam.Bag{"zulu": flam.Bag{}, "alpha": flam.Bag{}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get("path").Return(config).Times(2)

		factory, e := flam.NewFactory[flam.Resource](nil, "path", factoryConfig, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		assert.NoError(t, factory.Add("charlie", &testResource{}))
		assert.Equal(t, []string{"alpha", "charlie", "zulu"}, factory.List())
	})

	t.Run("should return an empty list when there are no entries", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get("path").Return(flam.Bag{}).Times(1)

		factory, e := flam.NewFactory[flam.Resource](nil, "path", factoryConfig, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		assert.Empty(t, factory.List())
	})
}

func Test_Factory_Has(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := flam.Bag{"entry1": flam.Bag{}}
	factoryConfig := mocks.NewFactoryConfig(ctrl)
	factoryConfig.EXPECT().Get("path").Return(config).AnyTimes()

	factory, e := flam.NewFactory[flam.Resource](nil, "path", factoryConfig, nil)
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

		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get("path").Return(flam.Bag{}).Times(1)

		factory, e := flam.NewFactory[flam.Resource](nil, "path", factoryConfig, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		got, e := factory.Get("nonexistent")
		assert.Nil(t, got)
		assert.ErrorIs(t, e, flam.ErrUnknownResource)
	})
}

func Test_Factory_Generate(t *testing.T) {
	t.Run("should return ErrUnknownResource for an unknown id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get("path").Return(flam.Bag{}).Times(1)

		factory, e := flam.NewFactory[flam.Resource](nil, "path", factoryConfig, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		got, e := factory.Generate("unknown")
		assert.Nil(t, got)
		assert.ErrorIs(t, e, flam.ErrUnknownResource)
	})

	t.Run("should return an error if creator's Create method fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := flam.Bag{"default": flam.Bag{}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get("path").Return(config).Times(1)

		expectedErr := errors.New("creation failed")
		creator := mocks.NewResourceCreator[flam.Resource](ctrl)
		creator.EXPECT().Accept(flam.Bag{"id": "default"}).Return(true)
		creator.EXPECT().Create(flam.Bag{"id": "default"}).Return(nil, expectedErr)

		factory, e := flam.NewFactory([]flam.ResourceCreator[flam.Resource]{creator}, "path", factoryConfig, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		got, e := factory.Generate("default")
		assert.Nil(t, got)
		assert.ErrorIs(t, e, expectedErr)
	})

	t.Run("should return ErrInvalidResourceConfig if no creator accepts the config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := flam.Bag{"default": flam.Bag{}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get("path").Return(config).Times(1)

		creator := mocks.NewResourceCreator[flam.Resource](ctrl)
		creator.EXPECT().Accept(flam.Bag{"id": "default"}).Return(false)

		factory, e := flam.NewFactory([]flam.ResourceCreator[flam.Resource]{creator}, "path", factoryConfig, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		got, e := factory.Generate("default")
		assert.Nil(t, got)
		assert.ErrorIs(t, e, flam.ErrInvalidResourceConfig)
	})

	t.Run("should return the validation error if there is a validator and it return an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := flam.Bag{"default": flam.Bag{"driver": "test"}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get("path").Return(config).Times(1)

		expectedErr := errors.New("validation error")
		validator := func(config flam.Bag) error { return expectedErr }

		factory, e := flam.NewFactory[flam.Resource](nil, "path", factoryConfig, validator)
		require.NotNil(t, factory)
		require.NoError(t, e)

		got, e := factory.Generate("default")
		assert.Nil(t, got)
		assert.ErrorIs(t, e, expectedErr)
	})

	t.Run("should generate, cache, and return an entry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := flam.Bag{"default": flam.Bag{"driver": "test"}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get("path").Return(config).Times(1)

		resource := &testResource{}

		creator := mocks.NewResourceCreator[flam.Resource](ctrl)
		creator.EXPECT().Accept(flam.Bag{"id": "default", "driver": "test"}).Return(true).Times(1)
		creator.EXPECT().Create(flam.Bag{"id": "default", "driver": "test"}).Return(resource, nil).Times(1)

		validator := func(config flam.Bag) error { return nil }

		factory, e := flam.NewFactory([]flam.ResourceCreator[flam.Resource]{creator}, "path", factoryConfig, validator)
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

		factoryConfig := mocks.NewFactoryConfig(ctrl)

		factory, e := flam.NewFactory[flam.Resource](nil, "path", factoryConfig, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		assert.ErrorIs(t, factory.Add("nil_resource", nil), flam.ErrNilReference)
	})

	t.Run("should successfully add a resource and allow it to be retrieved", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := flam.Bag{}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get("path").Return(config).Times(1)

		factory, e := flam.NewFactory[flam.Resource](nil, "path", factoryConfig, nil)
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

		config := flam.Bag{}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get("path").Return(config).Times(1)

		factory, e := flam.NewFactory[flam.Resource](nil, "path", factoryConfig, nil)
		require.NotNil(t, factory)
		require.NoError(t, e)

		resource := &testResource{}
		assert.NoError(t, factory.Add("my_resource", resource))
		assert.ErrorIs(t, factory.Add("my_resource", resource), flam.ErrDuplicateResource)
	})
}
