package flam

import (
	"reflect"

	"github.com/golang/mock/gomock"
	"go.uber.org/dig"
)

type RunnableProviderMock struct {
	ctrl     *gomock.Controller
	recorder *RunnableProviderMockRecorder
}

type RunnableProviderMockRecorder struct {
	mock *RunnableProviderMock
}

func NewRunnableProviderMock(ctrl *gomock.Controller) *RunnableProviderMock {
	mock := &RunnableProviderMock{ctrl: ctrl}
	mock.recorder = &RunnableProviderMockRecorder{mock}
	return mock
}

func (m *RunnableProviderMock) EXPECT() *RunnableProviderMockRecorder {
	return m.recorder
}

func (m *RunnableProviderMock) Id() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Id")
	ret0, _ := ret[0].(string)
	return ret0
}

func (mr *RunnableProviderMockRecorder) Id() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Id", reflect.TypeOf((*RunnableProviderMock)(nil).Id))
}

func (m *RunnableProviderMock) Register(container *dig.Container) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", container)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *RunnableProviderMockRecorder) Register(container interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*RunnableProviderMock)(nil).Register), container)
}

func (m *RunnableProviderMock) Run(container *dig.Container) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", container)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *RunnableProviderMockRecorder) Run(container interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*RunnableProviderMock)(nil).Run), container)
}
