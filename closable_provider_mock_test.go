package flam

import (
	"reflect"

	"github.com/golang/mock/gomock"
	"go.uber.org/dig"
)

type ClosableProviderMock struct {
	ctrl     *gomock.Controller
	recorder *ClosableProviderMockRecorder
}

type ClosableProviderMockRecorder struct {
	mock *ClosableProviderMock
}

func NewClosableProviderMock(ctrl *gomock.Controller) *ClosableProviderMock {
	mock := &ClosableProviderMock{ctrl: ctrl}
	mock.recorder = &ClosableProviderMockRecorder{mock}
	return mock
}

func (m *ClosableProviderMock) EXPECT() *ClosableProviderMockRecorder {
	return m.recorder
}

func (m *ClosableProviderMock) Id() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Id")
	ret0, _ := ret[0].(string)
	return ret0
}

func (mr *ClosableProviderMockRecorder) Id() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Id", reflect.TypeOf((*ClosableProviderMock)(nil).Id))
}

func (m *ClosableProviderMock) Register(container *dig.Container) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", container)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *ClosableProviderMockRecorder) Register(container interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*ClosableProviderMock)(nil).Register), container)
}

func (m *ClosableProviderMock) Close(container *dig.Container) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", container)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *ClosableProviderMockRecorder) Close(container interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*ClosableProviderMock)(nil).Close), container)
}
