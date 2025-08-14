package mocks

import (
	"reflect"

	"github.com/golang/mock/gomock"
)

// Closer is a mock of Closer interface.
type Closer struct {
	ctrl     *gomock.Controller
	recorder *CloserRecorder
}

// CloserRecorder is the mock recorder for Closer.
type CloserRecorder struct {
	mock *Closer
}

// NewCloser creates a new mock instance.
func NewCloser(ctrl *gomock.Controller) *Closer {
	mock := &Closer{ctrl: ctrl}
	mock.recorder = &CloserRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Closer) EXPECT() *CloserRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *Closer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *CloserRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*Closer)(nil).Close))
}
