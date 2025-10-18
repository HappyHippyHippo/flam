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

		providerMock := NewProviderMock(ctrl)
		providerMock.EXPECT().Id().Return("provider").AnyTimes()
		providerMock.EXPECT().Register(gomock.Any()).Return(nil).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(providerMock))
		assert.ErrorIs(t, app.Register(providerMock), ErrDuplicateProvider)
	})

	t.Run("should return an error if provider registration fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("registration error")
		providerMock := NewProviderMock(ctrl)
		providerMock.EXPECT().Register(gomock.Any()).Return(expectedErr).Times(1)

		app := NewApplication()
		assert.ErrorIs(t, app.Register(providerMock), expectedErr)
	})

	t.Run("should register a provider successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		providerMock := NewProviderMock(ctrl)
		providerMock.EXPECT().Register(gomock.Any()).Return(nil).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(providerMock))
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
		providerMock := NewBootableProviderMock(ctrl)
		providerMock.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		providerMock.EXPECT().Boot(gomock.Any()).Return(expectedErr).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(providerMock))
		assert.ErrorIs(t, app.Boot(), expectedErr)
	})

	t.Run("should boot all providers successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		providerMock1 := NewProviderMock(ctrl)
		providerMock1.EXPECT().Id().Return("provider1").AnyTimes()
		providerMock1.EXPECT().Register(gomock.Any()).Return(nil).Times(1)

		providerMock2 := NewBootableProviderMock(ctrl)
		providerMock2.EXPECT().Id().Return("provider2").AnyTimes()
		providerMock2.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		providerMock2.EXPECT().Boot(gomock.Any()).Return(nil).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(providerMock1))
		assert.NoError(t, app.Register(providerMock2))
		assert.NoError(t, app.Boot())
	})
}

func Test_Application_Run(t *testing.T) {
	t.Run("should boot if not already booted", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		providerMock := NewBootableProviderMock(ctrl)
		providerMock.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		providerMock.EXPECT().Boot(gomock.Any()).Return(nil).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(providerMock))
		assert.NoError(t, app.Run())
	})

	t.Run("should return an error if boot fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("boot error")
		providerMock := NewBootableProviderMock(ctrl)
		providerMock.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		providerMock.EXPECT().Boot(gomock.Any()).Return(expectedErr).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(providerMock))
		assert.ErrorIs(t, app.Run(), expectedErr)
	})

	t.Run("should return an error if a runnable provider fails to run", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("run error")
		providerMock := NewRunnableProviderMock(ctrl)
		providerMock.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		providerMock.EXPECT().Run(gomock.Any()).Return(expectedErr).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(providerMock))
		assert.ErrorIs(t, app.Run(), expectedErr)
	})

	t.Run("should run all providers successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		providerMock1 := NewProviderMock(ctrl)
		providerMock1.EXPECT().Id().Return("provider1").Times(1)
		providerMock1.EXPECT().Register(gomock.Any()).Return(nil).Times(1)

		providerMock2 := NewRunnableProviderMock(ctrl)
		providerMock2.EXPECT().Id().Return("provider2").Times(1)
		providerMock2.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		providerMock2.EXPECT().Run(gomock.Any()).Return(nil).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(providerMock1))
		assert.NoError(t, app.Register(providerMock2))
		assert.NoError(t, app.Run())
	})
}

func Test_Application_Close(t *testing.T) {
	t.Run("should return the first error from closable providers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("close error")
		providerMock1 := NewClosableProviderMock(ctrl)
		providerMock1.EXPECT().Id().Return("closable1").Times(1)
		providerMock1.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		providerMock1.EXPECT().Close(gomock.Any()).Return(expectedErr).Times(1)

		providerMock2 := NewClosableProviderMock(ctrl)
		providerMock2.EXPECT().Id().Return("closable2").Times(1)
		providerMock2.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		providerMock2.EXPECT().Close(gomock.Any()).Return(nil).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(providerMock2))
		assert.NoError(t, app.Register(providerMock1))
		assert.ErrorIs(t, app.Close(), expectedErr)
	})

	t.Run("should close all providers successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		providerMock1 := NewProviderMock(ctrl)
		providerMock1.EXPECT().Id().Return("provider1").Times(1)
		providerMock1.EXPECT().Register(gomock.Any()).Return(nil).Times(1)

		providerMock2 := NewClosableProviderMock(ctrl)
		providerMock2.EXPECT().Id().Return("provider2").Times(1)
		providerMock2.EXPECT().Register(gomock.Any()).Return(nil).Times(1)
		providerMock2.EXPECT().Close(gomock.Any()).Return(nil).Times(1)

		app := NewApplication()
		assert.NoError(t, app.Register(providerMock1))
		assert.NoError(t, app.Register(providerMock2))
		assert.NoError(t, app.Close())
	})
}
