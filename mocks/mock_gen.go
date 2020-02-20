// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sevigo/hokan/pkg/core (interfaces: DirectoryStore,EventCreator)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	core "github.com/sevigo/hokan/pkg/core"
	reflect "reflect"
)

// MockDirectoryStore is a mock of DirectoryStore interface
type MockDirectoryStore struct {
	ctrl     *gomock.Controller
	recorder *MockDirectoryStoreMockRecorder
}

// MockDirectoryStoreMockRecorder is the mock recorder for MockDirectoryStore
type MockDirectoryStoreMockRecorder struct {
	mock *MockDirectoryStore
}

// NewMockDirectoryStore creates a new mock instance
func NewMockDirectoryStore(ctrl *gomock.Controller) *MockDirectoryStore {
	mock := &MockDirectoryStore{ctrl: ctrl}
	mock.recorder = &MockDirectoryStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDirectoryStore) EXPECT() *MockDirectoryStoreMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockDirectoryStore) Create(arg0 context.Context, arg1 *core.Directory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockDirectoryStoreMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDirectoryStore)(nil).Create), arg0, arg1)
}

// Delete mocks base method
func (m *MockDirectoryStore) Delete(arg0 context.Context, arg1 *core.Directory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockDirectoryStoreMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDirectoryStore)(nil).Delete), arg0, arg1)
}

// Find mocks base method
func (m *MockDirectoryStore) Find(arg0 context.Context, arg1 int64) (*core.Directory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.Directory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockDirectoryStoreMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockDirectoryStore)(nil).Find), arg0, arg1)
}

// FindName mocks base method
func (m *MockDirectoryStore) FindName(arg0 context.Context, arg1 string) (*core.Directory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindName", arg0, arg1)
	ret0, _ := ret[0].(*core.Directory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindName indicates an expected call of FindName
func (mr *MockDirectoryStoreMockRecorder) FindName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindName", reflect.TypeOf((*MockDirectoryStore)(nil).FindName), arg0, arg1)
}

// List mocks base method
func (m *MockDirectoryStore) List(arg0 context.Context) ([]*core.Directory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]*core.Directory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockDirectoryStoreMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDirectoryStore)(nil).List), arg0)
}

// Update mocks base method
func (m *MockDirectoryStore) Update(arg0 context.Context, arg1 *core.Directory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockDirectoryStoreMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDirectoryStore)(nil).Update), arg0, arg1)
}

// MockEventCreator is a mock of EventCreator interface
type MockEventCreator struct {
	ctrl     *gomock.Controller
	recorder *MockEventCreatorMockRecorder
}

// MockEventCreatorMockRecorder is the mock recorder for MockEventCreator
type MockEventCreatorMockRecorder struct {
	mock *MockEventCreator
}

// NewMockEventCreator creates a new mock instance
func NewMockEventCreator(ctrl *gomock.Controller) *MockEventCreator {
	mock := &MockEventCreator{ctrl: ctrl}
	mock.recorder = &MockEventCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEventCreator) EXPECT() *MockEventCreatorMockRecorder {
	return m.recorder
}

// Publish mocks base method
func (m *MockEventCreator) Publish(arg0 context.Context, arg1 *core.EventData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish
func (mr *MockEventCreatorMockRecorder) Publish(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockEventCreator)(nil).Publish), arg0, arg1)
}
