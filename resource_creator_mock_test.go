package flam

import (
	"reflect"

	"github.com/golang/mock/gomock"
)

type ResourceCreatorMock[R Resource] struct {
	ctrl     *gomock.Controller
	recorder *ResourceCreatorMockRecorder[R]
}

type ResourceCreatorMockRecorder[R Resource] struct {
	mock *ResourceCreatorMock[R]
}

func NewResourceCreatorMock[R Resource](ctrl *gomock.Controller) *ResourceCreatorMock[R] {
	mock := &ResourceCreatorMock[R]{ctrl: ctrl}
	mock.recorder = &ResourceCreatorMockRecorder[R]{mock}
	return mock
}

func (m *ResourceCreatorMock[R]) EXPECT() *ResourceCreatorMockRecorder[R] {
	return m.recorder
}

func (m *ResourceCreatorMock[R]) Accept(config Bag) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", config)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (mr *ResourceCreatorMockRecorder[R]) Accept(config any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*ResourceCreatorMock[R])(nil).Accept), config)
}

func (m *ResourceCreatorMock[R]) Create(config Bag) (R, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", config)
	ret0, _ := ret[0].(R)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *ResourceCreatorMockRecorder[R]) Create(config any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*ResourceCreatorMock[R])(nil).Create), config)
}
