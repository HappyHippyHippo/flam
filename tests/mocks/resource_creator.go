package mocks

import (
	"reflect"

	"github.com/golang/mock/gomock"

	"github.com/happyhippyhippo/flam"
)

// ResourceCreator is a mock of ResourceCreator interface.
type ResourceCreator[Resource flam.Resource] struct {
	ctrl     *gomock.Controller
	recorder *ResourceCreatorRecorder[Resource]
}

// ResourceCreatorRecorder is the mock recorder for ResourceCreator.
type ResourceCreatorRecorder[Resource flam.Resource] struct {
	mock *ResourceCreator[Resource]
}

// NewResourceCreator creates a new mock instance.
func NewResourceCreator[Resource flam.Resource](ctrl *gomock.Controller) *ResourceCreator[Resource] {
	mock := &ResourceCreator[Resource]{ctrl: ctrl}
	mock.recorder = &ResourceCreatorRecorder[Resource]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *ResourceCreator[Resource]) EXPECT() *ResourceCreatorRecorder[Resource] {
	return m.recorder
}

// Accept mocks base method.
func (m *ResourceCreator[Resource]) Accept(config flam.Bag) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", config)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Accept indicates an expected call of Accept.
func (mr *ResourceCreatorRecorder[Resource]) Accept(config any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*ResourceCreator[Resource])(nil).Accept), config)
}

// Create mocks base method.
func (m *ResourceCreator[Resource]) Create(config flam.Bag) (Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", config)
	ret0, _ := ret[0].(Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *ResourceCreatorRecorder[Resource]) Create(config any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*ResourceCreator[Resource])(nil).Create), config)
}
