// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sevigo/hokan/pkg/core (interfaces: ConfigStore,DirectoryStore,EventCreator,FileStore,MinioWrapper,Notifier,TargetStorage,TargetRegister,UserStore,Watcher)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	minio "github.com/minio/minio-go"
	core "github.com/sevigo/hokan/pkg/core"
	core0 "github.com/sevigo/notify/core"
	event "github.com/sevigo/notify/event"
	reflect "reflect"
)

// MockConfigStore is a mock of ConfigStore interface
type MockConfigStore struct {
	ctrl     *gomock.Controller
	recorder *MockConfigStoreMockRecorder
}

// MockConfigStoreMockRecorder is the mock recorder for MockConfigStore
type MockConfigStoreMockRecorder struct {
	mock *MockConfigStore
}

// NewMockConfigStore creates a new mock instance
func NewMockConfigStore(ctrl *gomock.Controller) *MockConfigStore {
	mock := &MockConfigStore{ctrl: ctrl}
	mock.recorder = &MockConfigStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigStore) EXPECT() *MockConfigStoreMockRecorder {
	return m.recorder
}

// Find mocks base method
func (m *MockConfigStore) Find(arg0 context.Context, arg1 string) (*core.TargetConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.TargetConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockConfigStoreMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockConfigStore)(nil).Find), arg0, arg1)
}

// Save mocks base method
func (m *MockConfigStore) Save(arg0 context.Context, arg1 *core.TargetConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockConfigStoreMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockConfigStore)(nil).Save), arg0, arg1)
}

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

// Subscribe mocks base method
func (m *MockEventCreator) Subscribe(arg0 context.Context, arg1 core.EventType) <-chan *core.EventData {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", arg0, arg1)
	ret0, _ := ret[0].(<-chan *core.EventData)
	return ret0
}

// Subscribe indicates an expected call of Subscribe
func (mr *MockEventCreatorMockRecorder) Subscribe(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockEventCreator)(nil).Subscribe), arg0, arg1)
}

// MockFileStore is a mock of FileStore interface
type MockFileStore struct {
	ctrl     *gomock.Controller
	recorder *MockFileStoreMockRecorder
}

// MockFileStoreMockRecorder is the mock recorder for MockFileStore
type MockFileStoreMockRecorder struct {
	mock *MockFileStore
}

// NewMockFileStore creates a new mock instance
func NewMockFileStore(ctrl *gomock.Controller) *MockFileStore {
	mock := &MockFileStore{ctrl: ctrl}
	mock.recorder = &MockFileStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFileStore) EXPECT() *MockFileStoreMockRecorder {
	return m.recorder
}

// Delete mocks base method
func (m *MockFileStore) Delete(arg0 context.Context, arg1 string, arg2 *core.File) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockFileStoreMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockFileStore)(nil).Delete), arg0, arg1, arg2)
}

// Find mocks base method
func (m *MockFileStore) Find(arg0 context.Context, arg1, arg2 string) (*core.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1, arg2)
	ret0, _ := ret[0].(*core.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockFileStoreMockRecorder) Find(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockFileStore)(nil).Find), arg0, arg1, arg2)
}

// List mocks base method
func (m *MockFileStore) List(arg0 context.Context, arg1 string) ([]*core.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*core.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockFileStoreMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockFileStore)(nil).List), arg0, arg1)
}

// Save mocks base method
func (m *MockFileStore) Save(arg0 context.Context, arg1 string, arg2 *core.File) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockFileStoreMockRecorder) Save(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockFileStore)(nil).Save), arg0, arg1, arg2)
}

// MockMinioWrapper is a mock of MinioWrapper interface
type MockMinioWrapper struct {
	ctrl     *gomock.Controller
	recorder *MockMinioWrapperMockRecorder
}

// MockMinioWrapperMockRecorder is the mock recorder for MockMinioWrapper
type MockMinioWrapperMockRecorder struct {
	mock *MockMinioWrapper
}

// NewMockMinioWrapper creates a new mock instance
func NewMockMinioWrapper(ctrl *gomock.Controller) *MockMinioWrapper {
	mock := &MockMinioWrapper{ctrl: ctrl}
	mock.recorder = &MockMinioWrapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMinioWrapper) EXPECT() *MockMinioWrapperMockRecorder {
	return m.recorder
}

// BucketExists mocks base method
func (m *MockMinioWrapper) BucketExists(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BucketExists", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BucketExists indicates an expected call of BucketExists
func (mr *MockMinioWrapperMockRecorder) BucketExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BucketExists", reflect.TypeOf((*MockMinioWrapper)(nil).BucketExists), arg0)
}

// FPutObjectWithContext mocks base method
func (m *MockMinioWrapper) FPutObjectWithContext(arg0 context.Context, arg1, arg2, arg3 string, arg4 minio.PutObjectOptions) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FPutObjectWithContext", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FPutObjectWithContext indicates an expected call of FPutObjectWithContext
func (mr *MockMinioWrapperMockRecorder) FPutObjectWithContext(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FPutObjectWithContext", reflect.TypeOf((*MockMinioWrapper)(nil).FPutObjectWithContext), arg0, arg1, arg2, arg3, arg4)
}

// MockNotifier is a mock of Notifier interface
type MockNotifier struct {
	ctrl     *gomock.Controller
	recorder *MockNotifierMockRecorder
}

// MockNotifierMockRecorder is the mock recorder for MockNotifier
type MockNotifierMockRecorder struct {
	mock *MockNotifier
}

// NewMockNotifier creates a new mock instance
func NewMockNotifier(ctrl *gomock.Controller) *MockNotifier {
	mock := &MockNotifier{ctrl: ctrl}
	mock.recorder = &MockNotifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNotifier) EXPECT() *MockNotifierMockRecorder {
	return m.recorder
}

// Error mocks base method
func (m *MockNotifier) Error() chan event.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Error")
	ret0, _ := ret[0].(chan event.Error)
	return ret0
}

// Error indicates an expected call of Error
func (mr *MockNotifierMockRecorder) Error() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockNotifier)(nil).Error))
}

// Event mocks base method
func (m *MockNotifier) Event() chan event.Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Event")
	ret0, _ := ret[0].(chan event.Event)
	return ret0
}

// Event indicates an expected call of Event
func (mr *MockNotifierMockRecorder) Event() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Event", reflect.TypeOf((*MockNotifier)(nil).Event))
}

// RescanAll mocks base method
func (m *MockNotifier) RescanAll() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RescanAll")
}

// RescanAll indicates an expected call of RescanAll
func (mr *MockNotifierMockRecorder) RescanAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RescanAll", reflect.TypeOf((*MockNotifier)(nil).RescanAll))
}

// StartWatching mocks base method
func (m *MockNotifier) StartWatching(arg0 string, arg1 *core0.WatchingOptions) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartWatching", arg0, arg1)
}

// StartWatching indicates an expected call of StartWatching
func (mr *MockNotifierMockRecorder) StartWatching(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartWatching", reflect.TypeOf((*MockNotifier)(nil).StartWatching), arg0, arg1)
}

// StopWatching mocks base method
func (m *MockNotifier) StopWatching(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StopWatching", arg0)
}

// StopWatching indicates an expected call of StopWatching
func (mr *MockNotifierMockRecorder) StopWatching(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopWatching", reflect.TypeOf((*MockNotifier)(nil).StopWatching), arg0)
}

// MockTargetStorage is a mock of TargetStorage interface
type MockTargetStorage struct {
	ctrl     *gomock.Controller
	recorder *MockTargetStorageMockRecorder
}

// MockTargetStorageMockRecorder is the mock recorder for MockTargetStorage
type MockTargetStorageMockRecorder struct {
	mock *MockTargetStorage
}

// NewMockTargetStorage creates a new mock instance
func NewMockTargetStorage(ctrl *gomock.Controller) *MockTargetStorage {
	mock := &MockTargetStorage{ctrl: ctrl}
	mock.recorder = &MockTargetStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTargetStorage) EXPECT() *MockTargetStorageMockRecorder {
	return m.recorder
}

// Delete mocks base method
func (m *MockTargetStorage) Delete(arg0 context.Context, arg1 *core.File) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockTargetStorageMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTargetStorage)(nil).Delete), arg0, arg1)
}

// Find mocks base method
func (m *MockTargetStorage) Find(arg0 context.Context, arg1 string) (*core.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockTargetStorageMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockTargetStorage)(nil).Find), arg0, arg1)
}

// List mocks base method
func (m *MockTargetStorage) List(arg0 context.Context) ([]*core.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]*core.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockTargetStorageMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTargetStorage)(nil).List), arg0)
}

// Ping mocks base method
func (m *MockTargetStorage) Ping(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping
func (mr *MockTargetStorageMockRecorder) Ping(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockTargetStorage)(nil).Ping), arg0)
}

// Save mocks base method
func (m *MockTargetStorage) Save(arg0 context.Context, arg1 *core.File) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockTargetStorageMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockTargetStorage)(nil).Save), arg0, arg1)
}

// MockTargetRegister is a mock of TargetRegister interface
type MockTargetRegister struct {
	ctrl     *gomock.Controller
	recorder *MockTargetRegisterMockRecorder
}

// MockTargetRegisterMockRecorder is the mock recorder for MockTargetRegister
type MockTargetRegisterMockRecorder struct {
	mock *MockTargetRegister
}

// NewMockTargetRegister creates a new mock instance
func NewMockTargetRegister(ctrl *gomock.Controller) *MockTargetRegister {
	mock := &MockTargetRegister{ctrl: ctrl}
	mock.recorder = &MockTargetRegisterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTargetRegister) EXPECT() *MockTargetRegisterMockRecorder {
	return m.recorder
}

// AllConfigs mocks base method
func (m *MockTargetRegister) AllConfigs() map[string]*core.TargetConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllConfigs")
	ret0, _ := ret[0].(map[string]*core.TargetConfig)
	return ret0
}

// AllConfigs indicates an expected call of AllConfigs
func (mr *MockTargetRegisterMockRecorder) AllConfigs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllConfigs", reflect.TypeOf((*MockTargetRegister)(nil).AllConfigs))
}

// AllTargets mocks base method
func (m *MockTargetRegister) AllTargets() map[string]core.TargetFactory {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllTargets")
	ret0, _ := ret[0].(map[string]core.TargetFactory)
	return ret0
}

// AllTargets indicates an expected call of AllTargets
func (mr *MockTargetRegisterMockRecorder) AllTargets() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllTargets", reflect.TypeOf((*MockTargetRegister)(nil).AllTargets))
}

// GetConfig mocks base method
func (m *MockTargetRegister) GetConfig(arg0 context.Context, arg1 string) (*core.TargetConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfig", arg0, arg1)
	ret0, _ := ret[0].(*core.TargetConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfig indicates an expected call of GetConfig
func (mr *MockTargetRegisterMockRecorder) GetConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockTargetRegister)(nil).GetConfig), arg0, arg1)
}

// GetTarget mocks base method
func (m *MockTargetRegister) GetTarget(arg0 string) core.TargetStorage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTarget", arg0)
	ret0, _ := ret[0].(core.TargetStorage)
	return ret0
}

// GetTarget indicates an expected call of GetTarget
func (mr *MockTargetRegisterMockRecorder) GetTarget(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTarget", reflect.TypeOf((*MockTargetRegister)(nil).GetTarget), arg0)
}

// SetConfig mocks base method
func (m *MockTargetRegister) SetConfig(arg0 context.Context, arg1 *core.TargetConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetConfig", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetConfig indicates an expected call of SetConfig
func (mr *MockTargetRegisterMockRecorder) SetConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConfig", reflect.TypeOf((*MockTargetRegister)(nil).SetConfig), arg0, arg1)
}

// MockUserStore is a mock of UserStore interface
type MockUserStore struct {
	ctrl     *gomock.Controller
	recorder *MockUserStoreMockRecorder
}

// MockUserStoreMockRecorder is the mock recorder for MockUserStore
type MockUserStoreMockRecorder struct {
	mock *MockUserStore
}

// NewMockUserStore creates a new mock instance
func NewMockUserStore(ctrl *gomock.Controller) *MockUserStore {
	mock := &MockUserStore{ctrl: ctrl}
	mock.recorder = &MockUserStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserStore) EXPECT() *MockUserStoreMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockUserStore) Create(arg0 context.Context, arg1 *core.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockUserStoreMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserStore)(nil).Create), arg0, arg1)
}

// Delete mocks base method
func (m *MockUserStore) Delete(arg0 context.Context, arg1 *core.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockUserStoreMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUserStore)(nil).Delete), arg0, arg1)
}

// Find mocks base method
func (m *MockUserStore) Find(arg0 context.Context, arg1 int64) (*core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockUserStoreMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockUserStore)(nil).Find), arg0, arg1)
}

// Update mocks base method
func (m *MockUserStore) Update(arg0 context.Context, arg1 *core.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockUserStoreMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserStore)(nil).Update), arg0, arg1)
}

// MockWatcher is a mock of Watcher interface
type MockWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockWatcherMockRecorder
}

// MockWatcherMockRecorder is the mock recorder for MockWatcher
type MockWatcherMockRecorder struct {
	mock *MockWatcher
}

// NewMockWatcher creates a new mock instance
func NewMockWatcher(ctrl *gomock.Controller) *MockWatcher {
	mock := &MockWatcher{ctrl: ctrl}
	mock.recorder = &MockWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWatcher) EXPECT() *MockWatcherMockRecorder {
	return m.recorder
}

// GetDirsToWatch mocks base method
func (m *MockWatcher) GetDirsToWatch() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDirsToWatch")
	ret0, _ := ret[0].(error)
	return ret0
}

// GetDirsToWatch indicates an expected call of GetDirsToWatch
func (mr *MockWatcherMockRecorder) GetDirsToWatch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDirsToWatch", reflect.TypeOf((*MockWatcher)(nil).GetDirsToWatch))
}

// StartDirWatcher mocks base method
func (m *MockWatcher) StartDirWatcher() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartDirWatcher")
}

// StartDirWatcher indicates an expected call of StartDirWatcher
func (mr *MockWatcherMockRecorder) StartDirWatcher() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartDirWatcher", reflect.TypeOf((*MockWatcher)(nil).StartDirWatcher))
}

// StartFileWatcher mocks base method
func (m *MockWatcher) StartFileWatcher() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartFileWatcher")
}

// StartFileWatcher indicates an expected call of StartFileWatcher
func (mr *MockWatcherMockRecorder) StartFileWatcher() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartFileWatcher", reflect.TypeOf((*MockWatcher)(nil).StartFileWatcher))
}
