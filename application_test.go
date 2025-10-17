package flam

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Application_NewApplication(t *testing.T) {
	t.Run("should return a valid application", func(t *testing.T) {
		app := NewApplication()
		assert.NotNil(t, app)
		assert.NotNil(t, app.Container())
	})
}

func Test_Application_Container(t *testing.T) {
	t.Run("should return the container", func(t *testing.T) {
		app := NewApplication()
		container := app.Container()
		assert.NotNil(t, container)
		assert.Same(t, container, app.Container())
	})
}

func Test_Application_Register(t *testing.T) {
	t.Run("should return ErrNilReference when provider is nil", func(t *testing.T) {
		app := NewApplication()
		assert.ErrorIs(t, app.Register(nil), ErrNilReference)
	})

	t.Run("should return ErrDuplicateProvider when provider is already registered", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := NewProviderMock(ctrl)
		provider.EXPECT().Id().Return("provider").AnyTimes()
		provider.EXPECT().Register(gomock.Any()).Return(nil).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(provider))
		assert.ErrorIs(t, app.Register(provider), ErrDuplicateProvider)
	})

	t.Run("should return an error if provider registration fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("registration error")
		provider := NewProviderMock(ctrl)
		provider.EXPECT().Register(gomock.Any()).Return(expectedErr).Times(1)

		app := NewApplication()
		assert.ErrorIs(t, app.Register(provider), expectedErr)
	})

	t.Run("should register a provider successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := NewProviderMock(ctrl)
		provider.EXPECT().Register(gomock.Any()).Return(nil).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(provider))
	})
}

func Test_Application_Boot(t *testing.T) {
	t.Run("should do nothing if already booted", func(t *testing.T) {
		app := NewApplication()
		assert.NoError(t, app.Boot())
		assert.NoError(t, app.Boot())
	})

	t.Run("should return an error if a bootable provider fails to boot", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("boot error")
		provider := NewBootableProviderMock(ctrl)
		provider.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		provider.EXPECT().Boot(gomock.Any()).Return(expectedErr).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(provider))
		assert.ErrorIs(t, app.Boot(), expectedErr)
	})

	t.Run("should boot all providers successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider1 := NewProviderMock(ctrl)
		provider1.EXPECT().Id().Return("provider1").AnyTimes()
		provider1.EXPECT().Register(gomock.Any()).Return(nil).Times(1)

		provider2 := NewBootableProviderMock(ctrl)
		provider2.EXPECT().Id().Return("provider2").AnyTimes()
		provider2.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		provider2.EXPECT().Boot(gomock.Any()).Return(nil).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(provider1))
		assert.NoError(t, app.Register(provider2))
		assert.NoError(t, app.Boot())
	})
}

func Test_Application_Run(t *testing.T) {
	t.Run("should boot if not already booted", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := NewBootableProviderMock(ctrl)
		provider.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		provider.EXPECT().Boot(gomock.Any()).Return(nil).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(provider))
		assert.NoError(t, app.Run())
	})

	t.Run("should return an error if boot fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("boot error")
		provider := NewBootableProviderMock(ctrl)
		provider.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		provider.EXPECT().Boot(gomock.Any()).Return(expectedErr).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(provider))
		assert.ErrorIs(t, app.Run(), expectedErr)
	})

	t.Run("should return an error if a runnable provider fails to run", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("run error")
		provider := NewRunnableProviderMock(ctrl)
		provider.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		provider.EXPECT().Run(gomock.Any()).Return(expectedErr).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(provider))
		assert.ErrorIs(t, app.Run(), expectedErr)
	})

	t.Run("should run all providers successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider1 := NewProviderMock(ctrl)
		provider1.EXPECT().Id().Return("provider1").Times(1)
		provider1.EXPECT().Register(gomock.Any()).Return(nil).Times(1)

		provider2 := NewRunnableProviderMock(ctrl)
		provider2.EXPECT().Id().Return("provider2").Times(1)
		provider2.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		provider2.EXPECT().Run(gomock.Any()).Return(nil).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(provider1))
		assert.NoError(t, app.Register(provider2))
		assert.NoError(t, app.Run())
	})
}

func Test_Application_Close(t *testing.T) {
	t.Run("should return the first error from closable providers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("close error")
		provider1 := NewClosableProviderMock(ctrl)
		provider1.EXPECT().Id().Return("closable1").Times(1)
		provider1.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		provider1.EXPECT().Close(gomock.Any()).Return(expectedErr).Times(1)

		provider2 := NewClosableProviderMock(ctrl)
		provider2.EXPECT().Id().Return("closable2").Times(1)
		provider2.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		provider2.EXPECT().Close(gomock.Any()).Return(nil).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(provider2))
		assert.NoError(t, app.Register(provider1))
		assert.ErrorIs(t, app.Close(), expectedErr)
	})

	t.Run("should close all providers successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider1 := NewProviderMock(ctrl)
		provider1.EXPECT().Id().Return("provider1").Times(1)
		provider1.EXPECT().Register(gomock.Any()).Return(nil).Times(1)

		provider2 := NewClosableProviderMock(ctrl)
		provider2.EXPECT().Id().Return("provider2").Times(1)
		provider2.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		provider2.EXPECT().Close(gomock.Any()).Return(nil).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(provider1))
		assert.NoError(t, app.Register(provider2))
		assert.NoError(t, app.Close())
	})
}
