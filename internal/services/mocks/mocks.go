// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/StampWallet/backend/internal/services (interfaces: TokenService,EmailService,FileStorageService)

// Package mock_services is a generated GoMock package.
package mock_services

import (
	os "os"
	reflect "reflect"
	time "time"

	database "github.com/StampWallet/backend/internal/database"
	gomock "github.com/golang/mock/gomock"
)

// MockTokenService is a mock of TokenService interface.
type MockTokenService struct {
	ctrl     *gomock.Controller
	recorder *MockTokenServiceMockRecorder
}

// MockTokenServiceMockRecorder is the mock recorder for MockTokenService.
type MockTokenServiceMockRecorder struct {
	mock *MockTokenService
}

// NewMockTokenService creates a new mock instance.
func NewMockTokenService(ctrl *gomock.Controller) *MockTokenService {
	mock := &MockTokenService{ctrl: ctrl}
	mock.recorder = &MockTokenServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenService) EXPECT() *MockTokenServiceMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockTokenService) Check(arg0, arg1 string) (*database.User, *database.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", arg0, arg1)
	ret0, _ := ret[0].(*database.User)
	ret1, _ := ret[1].(*database.Token)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Check indicates an expected call of Check.
func (mr *MockTokenServiceMockRecorder) Check(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockTokenService)(nil).Check), arg0, arg1)
}

// Create mocks base method.
func (m *MockTokenService) Create(arg0 database.User, arg1 database.TokenPurposeEnum, arg2 time.Time) (*database.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(*database.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTokenServiceMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTokenService)(nil).Create), arg0, arg1, arg2)
}

// Invalidate mocks base method.
func (m *MockTokenService) Invalidate(arg0 database.Token) (*database.User, *database.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Invalidate", arg0)
	ret0, _ := ret[0].(*database.User)
	ret1, _ := ret[1].(*database.Token)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Invalidate indicates an expected call of Invalidate.
func (mr *MockTokenServiceMockRecorder) Invalidate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Invalidate", reflect.TypeOf((*MockTokenService)(nil).Invalidate), arg0)
}

// MockEmailService is a mock of EmailService interface.
type MockEmailService struct {
	ctrl     *gomock.Controller
	recorder *MockEmailServiceMockRecorder
}

// MockEmailServiceMockRecorder is the mock recorder for MockEmailService.
type MockEmailServiceMockRecorder struct {
	mock *MockEmailService
}

// NewMockEmailService creates a new mock instance.
func NewMockEmailService(ctrl *gomock.Controller) *MockEmailService {
	mock := &MockEmailService{ctrl: ctrl}
	mock.recorder = &MockEmailServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailService) EXPECT() *MockEmailServiceMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockEmailService) Send(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockEmailServiceMockRecorder) Send(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockEmailService)(nil).Send), arg0, arg1, arg2)
}

// MockFileStorageService is a mock of FileStorageService interface.
type MockFileStorageService struct {
	ctrl     *gomock.Controller
	recorder *MockFileStorageServiceMockRecorder
}

// MockFileStorageServiceMockRecorder is the mock recorder for MockFileStorageService.
type MockFileStorageServiceMockRecorder struct {
	mock *MockFileStorageService
}

// NewMockFileStorageService creates a new mock instance.
func NewMockFileStorageService(ctrl *gomock.Controller) *MockFileStorageService {
	mock := &MockFileStorageService{ctrl: ctrl}
	mock.recorder = &MockFileStorageServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileStorageService) EXPECT() *MockFileStorageServiceMockRecorder {
	return m.recorder
}

// CreateStub mocks base method.
func (m *MockFileStorageService) CreateStub(arg0 *database.User) (database.FileMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStub", arg0)
	ret0, _ := ret[0].(database.FileMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateStub indicates an expected call of CreateStub.
func (mr *MockFileStorageServiceMockRecorder) CreateStub(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStub", reflect.TypeOf((*MockFileStorageService)(nil).CreateStub), arg0)
}

// GetData mocks base method.
func (m *MockFileStorageService) GetData(arg0 string) (*os.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetData", arg0)
	ret0, _ := ret[0].(*os.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetData indicates an expected call of GetData.
func (mr *MockFileStorageServiceMockRecorder) GetData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetData", reflect.TypeOf((*MockFileStorageService)(nil).GetData), arg0)
}

// Remove mocks base method.
func (m *MockFileStorageService) Remove(arg0 database.FileMetadata) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockFileStorageServiceMockRecorder) Remove(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockFileStorageService)(nil).Remove), arg0)
}

// Upload mocks base method.
func (m *MockFileStorageService) Upload(arg0 database.FileMetadata, arg1 os.File, arg2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upload indicates an expected call of Upload.
func (mr *MockFileStorageServiceMockRecorder) Upload(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockFileStorageService)(nil).Upload), arg0, arg1, arg2)
}