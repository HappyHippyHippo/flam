package flam

import (
	"reflect"

	"github.com/golang/mock/gomock"
)

type CloserMock struct {
	ctrl     *gomock.Controller
	recorder *CloserMockRecorder
}

type CloserMockRecorder struct {
	mock *CloserMock
}

func NewCloserMock(ctrl *gomock.Controller) *CloserMock {
	mock := &CloserMock{ctrl: ctrl}
	mock.recorder = &CloserMockRecorder{mock}
	return mock
}

func (m *CloserMock) EXPECT() *CloserMockRecorder {
	return m.recorder
}

func (m *CloserMock) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *CloserMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*CloserMock)(nil).Close))
}
