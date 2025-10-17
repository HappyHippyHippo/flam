package flam

import (
	"reflect"

	"github.com/golang/mock/gomock"
	"go.uber.org/dig"
)

type BootableProviderMock struct {
	ctrl     *gomock.Controller
	recorder *BootableProviderMockRecorder
}

type BootableProviderMockRecorder struct {
	mock *BootableProviderMock
}

func NewBootableProviderMock(ctrl *gomock.Controller) *BootableProviderMock {
	mock := &BootableProviderMock{ctrl: ctrl}
	mock.recorder = &BootableProviderMockRecorder{mock}
	return mock
}

func (m *BootableProviderMock) EXPECT() *BootableProviderMockRecorder {
	return m.recorder
}

func (m *BootableProviderMock) Id() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Id")
	ret0, _ := ret[0].(string)
	return ret0
}

func (mr *BootableProviderMockRecorder) Id() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Id", reflect.TypeOf((*BootableProviderMock)(nil).Id))
}

func (m *BootableProviderMock) Register(container *dig.Container) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", container)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *BootableProviderMockRecorder) Register(container interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*BootableProviderMock)(nil).Register), container)
}

func (m *BootableProviderMock) Boot(container *dig.Container) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Boot", container)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *BootableProviderMockRecorder) Boot(container interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Boot", reflect.TypeOf((*BootableProviderMock)(nil).Boot), container)
}
