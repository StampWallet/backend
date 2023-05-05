// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/StampWallet/backend/internal/managers (interfaces: AuthManager,BusinessManager,ItemDefinitionManager,LocalCardManager,VirtualCardManager,TransactionManager)

// Package mock_managers is a generated GoMock package.
package mock_managers

import (
	reflect "reflect"

	database "github.com/StampWallet/backend/internal/database"
	managers "github.com/StampWallet/backend/internal/managers"
	gomock "github.com/golang/mock/gomock"
)

// MockAuthManager is a mock of AuthManager interface.
type MockAuthManager struct {
	ctrl     *gomock.Controller
	recorder *MockAuthManagerMockRecorder
}

// MockAuthManagerMockRecorder is the mock recorder for MockAuthManager.
type MockAuthManagerMockRecorder struct {
	mock *MockAuthManager
}

// NewMockAuthManager creates a new mock instance.
func NewMockAuthManager(ctrl *gomock.Controller) *MockAuthManager {
	mock := &MockAuthManager{ctrl: ctrl}
	mock.recorder = &MockAuthManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthManager) EXPECT() *MockAuthManagerMockRecorder {
	return m.recorder
}

// ChangeEmail mocks base method.
func (m *MockAuthManager) ChangeEmail(arg0 database.User, arg1 string) (*database.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeEmail", arg0, arg1)
	ret0, _ := ret[0].(*database.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangeEmail indicates an expected call of ChangeEmail.
func (mr *MockAuthManagerMockRecorder) ChangeEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeEmail", reflect.TypeOf((*MockAuthManager)(nil).ChangeEmail), arg0, arg1)
}

// ChangePassword mocks base method.
func (m *MockAuthManager) ChangePassword(arg0 database.User, arg1, arg2 string) (*database.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePassword", arg0, arg1, arg2)
	ret0, _ := ret[0].(*database.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockAuthManagerMockRecorder) ChangePassword(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockAuthManager)(nil).ChangePassword), arg0, arg1, arg2)
}

// ConfirmEmail mocks base method.
func (m *MockAuthManager) ConfirmEmail(arg0, arg1 string) (*database.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfirmEmail", arg0, arg1)
	ret0, _ := ret[0].(*database.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConfirmEmail indicates an expected call of ConfirmEmail.
func (mr *MockAuthManagerMockRecorder) ConfirmEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfirmEmail", reflect.TypeOf((*MockAuthManager)(nil).ConfirmEmail), arg0, arg1)
}

// Create mocks base method.
func (m *MockAuthManager) Create(arg0 managers.UserDetails) (*database.User, *database.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(*database.User)
	ret1, _ := ret[1].(*database.Token)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockAuthManagerMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAuthManager)(nil).Create), arg0)
}

// Login mocks base method.
func (m *MockAuthManager) Login(arg0, arg1 string) (*database.User, *database.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(*database.User)
	ret1, _ := ret[1].(*database.Token)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Login indicates an expected call of Login.
func (mr *MockAuthManagerMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthManager)(nil).Login), arg0, arg1)
}

// Logout mocks base method.
func (m *MockAuthManager) Logout(arg0, arg1 string) (*database.User, *database.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", arg0, arg1)
	ret0, _ := ret[0].(*database.User)
	ret1, _ := ret[1].(*database.Token)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Logout indicates an expected call of Logout.
func (mr *MockAuthManagerMockRecorder) Logout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAuthManager)(nil).Logout), arg0, arg1)
}

// MockBusinessManager is a mock of BusinessManager interface.
type MockBusinessManager struct {
	ctrl     *gomock.Controller
	recorder *MockBusinessManagerMockRecorder
}

// MockBusinessManagerMockRecorder is the mock recorder for MockBusinessManager.
type MockBusinessManagerMockRecorder struct {
	mock *MockBusinessManager
}

// NewMockBusinessManager creates a new mock instance.
func NewMockBusinessManager(ctrl *gomock.Controller) *MockBusinessManager {
	mock := &MockBusinessManager{ctrl: ctrl}
	mock.recorder = &MockBusinessManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBusinessManager) EXPECT() *MockBusinessManagerMockRecorder {
	return m.recorder
}

// ChangeDetails mocks base method.
func (m *MockBusinessManager) ChangeDetails(arg0 *database.Business, arg1 *managers.BusinessDetails) (*database.Business, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeDetails", arg0, arg1)
	ret0, _ := ret[0].(*database.Business)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangeDetails indicates an expected call of ChangeDetails.
func (mr *MockBusinessManagerMockRecorder) ChangeDetails(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeDetails", reflect.TypeOf((*MockBusinessManager)(nil).ChangeDetails), arg0, arg1)
}

// Create mocks base method.
func (m *MockBusinessManager) Create(arg0 *database.User, arg1 *managers.BusinessDetails) (*database.Business, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*database.Business)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockBusinessManagerMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockBusinessManager)(nil).Create), arg0, arg1)
}

// Search mocks base method.
func (m *MockBusinessManager) Search(arg0, arg1 string, arg2, arg3, arg4 uint) ([]database.Business, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]database.Business)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockBusinessManagerMockRecorder) Search(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockBusinessManager)(nil).Search), arg0, arg1, arg2, arg3, arg4)
}

// MockItemDefinitionManager is a mock of ItemDefinitionManager interface.
type MockItemDefinitionManager struct {
	ctrl     *gomock.Controller
	recorder *MockItemDefinitionManagerMockRecorder
}

// MockItemDefinitionManagerMockRecorder is the mock recorder for MockItemDefinitionManager.
type MockItemDefinitionManagerMockRecorder struct {
	mock *MockItemDefinitionManager
}

// NewMockItemDefinitionManager creates a new mock instance.
func NewMockItemDefinitionManager(ctrl *gomock.Controller) *MockItemDefinitionManager {
	mock := &MockItemDefinitionManager{ctrl: ctrl}
	mock.recorder = &MockItemDefinitionManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockItemDefinitionManager) EXPECT() *MockItemDefinitionManagerMockRecorder {
	return m.recorder
}

// AddItem mocks base method.
func (m *MockItemDefinitionManager) AddItem(arg0 *database.Business, arg1 *managers.ItemDetails) (*database.ItemDefinition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddItem", arg0, arg1)
	ret0, _ := ret[0].(*database.ItemDefinition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddItem indicates an expected call of AddItem.
func (mr *MockItemDefinitionManagerMockRecorder) AddItem(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddItem", reflect.TypeOf((*MockItemDefinitionManager)(nil).AddItem), arg0, arg1)
}

// ChangeItemDetails mocks base method.
func (m *MockItemDefinitionManager) ChangeItemDetails(arg0 *database.ItemDefinition, arg1 *managers.ItemDetails) (*database.ItemDefinition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeItemDetails", arg0, arg1)
	ret0, _ := ret[0].(*database.ItemDefinition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangeItemDetails indicates an expected call of ChangeItemDetails.
func (mr *MockItemDefinitionManagerMockRecorder) ChangeItemDetails(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeItemDetails", reflect.TypeOf((*MockItemDefinitionManager)(nil).ChangeItemDetails), arg0, arg1)
}

// GetForBusiness mocks base method.
func (m *MockItemDefinitionManager) GetForBusiness(arg0 *database.Business) ([]database.ItemDefinition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForBusiness", arg0)
	ret0, _ := ret[0].([]database.ItemDefinition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetForBusiness indicates an expected call of GetForBusiness.
func (mr *MockItemDefinitionManagerMockRecorder) GetForBusiness(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForBusiness", reflect.TypeOf((*MockItemDefinitionManager)(nil).GetForBusiness), arg0)
}

// WithdrawItem mocks base method.
func (m *MockItemDefinitionManager) WithdrawItem(arg0 *database.ItemDefinition) (*database.ItemDefinition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithdrawItem", arg0)
	ret0, _ := ret[0].(*database.ItemDefinition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WithdrawItem indicates an expected call of WithdrawItem.
func (mr *MockItemDefinitionManagerMockRecorder) WithdrawItem(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithdrawItem", reflect.TypeOf((*MockItemDefinitionManager)(nil).WithdrawItem), arg0)
}

// MockLocalCardManager is a mock of LocalCardManager interface.
type MockLocalCardManager struct {
	ctrl     *gomock.Controller
	recorder *MockLocalCardManagerMockRecorder
}

// MockLocalCardManagerMockRecorder is the mock recorder for MockLocalCardManager.
type MockLocalCardManagerMockRecorder struct {
	mock *MockLocalCardManager
}

// NewMockLocalCardManager creates a new mock instance.
func NewMockLocalCardManager(ctrl *gomock.Controller) *MockLocalCardManager {
	mock := &MockLocalCardManager{ctrl: ctrl}
	mock.recorder = &MockLocalCardManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLocalCardManager) EXPECT() *MockLocalCardManagerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockLocalCardManager) Create(arg0 *database.User, arg1 *managers.LocalCardDetails) (database.LocalCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(database.LocalCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockLocalCardManagerMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockLocalCardManager)(nil).Create), arg0, arg1)
}

// Remove mocks base method.
func (m *MockLocalCardManager) Remove(arg0 *database.LocalCard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockLocalCardManagerMockRecorder) Remove(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockLocalCardManager)(nil).Remove), arg0)
}

// MockVirtualCardManager is a mock of VirtualCardManager interface.
type MockVirtualCardManager struct {
	ctrl     *gomock.Controller
	recorder *MockVirtualCardManagerMockRecorder
}

// MockVirtualCardManagerMockRecorder is the mock recorder for MockVirtualCardManager.
type MockVirtualCardManagerMockRecorder struct {
	mock *MockVirtualCardManager
}

// NewMockVirtualCardManager creates a new mock instance.
func NewMockVirtualCardManager(ctrl *gomock.Controller) *MockVirtualCardManager {
	mock := &MockVirtualCardManager{ctrl: ctrl}
	mock.recorder = &MockVirtualCardManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVirtualCardManager) EXPECT() *MockVirtualCardManagerMockRecorder {
	return m.recorder
}

// BuyItem mocks base method.
func (m *MockVirtualCardManager) BuyItem(arg0 *database.VirtualCard, arg1 *database.ItemDefinition) (database.OwnedItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuyItem", arg0, arg1)
	ret0, _ := ret[0].(database.OwnedItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuyItem indicates an expected call of BuyItem.
func (mr *MockVirtualCardManagerMockRecorder) BuyItem(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuyItem", reflect.TypeOf((*MockVirtualCardManager)(nil).BuyItem), arg0, arg1)
}

// Create mocks base method.
func (m *MockVirtualCardManager) Create(arg0 *database.User, arg1 *database.Business) (*database.VirtualCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*database.VirtualCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockVirtualCardManagerMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockVirtualCardManager)(nil).Create), arg0, arg1)
}

// FilterOwnedItems mocks base method.
func (m *MockVirtualCardManager) FilterOwnedItems(arg0 *database.VirtualCard, arg1 []string) ([]database.OwnedItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterOwnedItems", arg0, arg1)
	ret0, _ := ret[0].([]database.OwnedItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilterOwnedItems indicates an expected call of FilterOwnedItems.
func (mr *MockVirtualCardManagerMockRecorder) FilterOwnedItems(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterOwnedItems", reflect.TypeOf((*MockVirtualCardManager)(nil).FilterOwnedItems), arg0, arg1)
}

// GetForUser mocks base method.
func (m *MockVirtualCardManager) GetForUser(arg0 *database.User) ([]database.VirtualCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForUser", arg0)
	ret0, _ := ret[0].([]database.VirtualCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetForUser indicates an expected call of GetForUser.
func (mr *MockVirtualCardManagerMockRecorder) GetForUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForUser", reflect.TypeOf((*MockVirtualCardManager)(nil).GetForUser), arg0)
}

// GetOwnedItems mocks base method.
func (m *MockVirtualCardManager) GetOwnedItems(arg0 *database.VirtualCard) ([]database.OwnedItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOwnedItems", arg0)
	ret0, _ := ret[0].([]database.OwnedItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOwnedItems indicates an expected call of GetOwnedItems.
func (mr *MockVirtualCardManagerMockRecorder) GetOwnedItems(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOwnedItems", reflect.TypeOf((*MockVirtualCardManager)(nil).GetOwnedItems), arg0)
}

// Remove mocks base method.
func (m *MockVirtualCardManager) Remove(arg0 *database.VirtualCard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockVirtualCardManagerMockRecorder) Remove(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockVirtualCardManager)(nil).Remove), arg0)
}

// ReturnItem mocks base method.
func (m *MockVirtualCardManager) ReturnItem(arg0 *database.OwnedItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReturnItem", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReturnItem indicates an expected call of ReturnItem.
func (mr *MockVirtualCardManagerMockRecorder) ReturnItem(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReturnItem", reflect.TypeOf((*MockVirtualCardManager)(nil).ReturnItem), arg0)
}

// MockTransactionManager is a mock of TransactionManager interface.
type MockTransactionManager struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionManagerMockRecorder
}

// MockTransactionManagerMockRecorder is the mock recorder for MockTransactionManager.
type MockTransactionManagerMockRecorder struct {
	mock *MockTransactionManager
}

// NewMockTransactionManager creates a new mock instance.
func NewMockTransactionManager(ctrl *gomock.Controller) *MockTransactionManager {
	mock := &MockTransactionManager{ctrl: ctrl}
	mock.recorder = &MockTransactionManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionManager) EXPECT() *MockTransactionManagerMockRecorder {
	return m.recorder
}

// Finalize mocks base method.
func (m *MockTransactionManager) Finalize(arg0 *database.Transaction, arg1 []managers.ItemWithAction, arg2 uint64) (*database.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Finalize", arg0, arg1, arg2)
	ret0, _ := ret[0].(*database.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Finalize indicates an expected call of Finalize.
func (mr *MockTransactionManagerMockRecorder) Finalize(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Finalize", reflect.TypeOf((*MockTransactionManager)(nil).Finalize), arg0, arg1, arg2)
}

// Start mocks base method.
func (m *MockTransactionManager) Start(arg0 *database.VirtualCard, arg1 []database.OwnedItem) (*database.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0, arg1)
	ret0, _ := ret[0].(*database.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Start indicates an expected call of Start.
func (mr *MockTransactionManagerMockRecorder) Start(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockTransactionManager)(nil).Start), arg0, arg1)
}
