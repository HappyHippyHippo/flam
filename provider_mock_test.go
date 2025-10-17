package flam

import (
	"reflect"

	"github.com/golang/mock/gomock"
	"go.uber.org/dig"
)

type ProviderMock struct {
	ctrl     *gomock.Controller
	recorder *ProviderMockRecorder
}

type ProviderMockRecorder struct {
	mock *ProviderMock
}

func NewProviderMock(ctrl *gomock.Controller) *ProviderMock {
	mock := &ProviderMock{ctrl: ctrl}
	mock.recorder = &ProviderMockRecorder{mock}
	return mock
}

func (m *ProviderMock) EXPECT() *ProviderMockRecorder {
	return m.recorder
}

func (m *ProviderMock) Id() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Id")
	ret0, _ := ret[0].(string)
	return ret0
}

func (mr *ProviderMockRecorder) Id() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Id", reflect.TypeOf((*ProviderMock)(nil).Id))
}

func (m *ProviderMock) Register(container *dig.Container) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", container)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *ProviderMockRecorder) Register(container interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*ProviderMock)(nil).Register), container)
}
