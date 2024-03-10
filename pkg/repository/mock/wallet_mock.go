// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/repository/interface/wallet.go

// Package mock is a generated GoMock package.
package mock

import (
	models "WatchHive/pkg/utils/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockWalletRepository is a mock of WalletRepository interface.
type MockWalletRepository struct {
	ctrl     *gomock.Controller
	recorder *MockWalletRepositoryMockRecorder
}

// MockWalletRepositoryMockRecorder is the mock recorder for MockWalletRepository.
type MockWalletRepositoryMockRecorder struct {
	mock *MockWalletRepository
}

// NewMockWalletRepository creates a new mock instance.
func NewMockWalletRepository(ctrl *gomock.Controller) *MockWalletRepository {
	mock := &MockWalletRepository{ctrl: ctrl}
	mock.recorder = &MockWalletRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWalletRepository) EXPECT() *MockWalletRepositoryMockRecorder {
	return m.recorder
}

// AddToWallet mocks base method.
func (m *MockWalletRepository) AddToWallet(userID int, Amount float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToWallet", userID, Amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToWallet indicates an expected call of AddToWallet.
func (mr *MockWalletRepositoryMockRecorder) AddToWallet(userID, Amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToWallet", reflect.TypeOf((*MockWalletRepository)(nil).AddToWallet), userID, Amount)
}

// AddToWalletHistory mocks base method.
func (m *MockWalletRepository) AddToWalletHistory(wallet models.WalletHistory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToWalletHistory", wallet)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToWalletHistory indicates an expected call of AddToWalletHistory.
func (mr *MockWalletRepositoryMockRecorder) AddToWalletHistory(wallet interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToWalletHistory", reflect.TypeOf((*MockWalletRepository)(nil).AddToWalletHistory), wallet)
}

// CreateWallet mocks base method.
func (m *MockWalletRepository) CreateWallet(userID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWallet", userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateWallet indicates an expected call of CreateWallet.
func (mr *MockWalletRepositoryMockRecorder) CreateWallet(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWallet", reflect.TypeOf((*MockWalletRepository)(nil).CreateWallet), userID)
}

// DebitFromWallet mocks base method.
func (m *MockWalletRepository) DebitFromWallet(userID int, amount float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DebitFromWallet", userID, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// DebitFromWallet indicates an expected call of DebitFromWallet.
func (mr *MockWalletRepositoryMockRecorder) DebitFromWallet(userID, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DebitFromWallet", reflect.TypeOf((*MockWalletRepository)(nil).DebitFromWallet), userID, amount)
}

// GetWallet mocks base method.
func (m *MockWalletRepository) GetWallet(userID int) (models.WalletAmount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWallet", userID)
	ret0, _ := ret[0].(models.WalletAmount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWallet indicates an expected call of GetWallet.
func (mr *MockWalletRepositoryMockRecorder) GetWallet(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWallet", reflect.TypeOf((*MockWalletRepository)(nil).GetWallet), userID)
}

// GetWalletData mocks base method.
func (m *MockWalletRepository) GetWalletData(userID int) (models.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWalletData", userID)
	ret0, _ := ret[0].(models.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWalletData indicates an expected call of GetWalletData.
func (mr *MockWalletRepositoryMockRecorder) GetWalletData(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWalletData", reflect.TypeOf((*MockWalletRepository)(nil).GetWalletData), userID)
}

// GetWalletHistory mocks base method.
func (m *MockWalletRepository) GetWalletHistory(walletId int) ([]models.WalletHistoryResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWalletHistory", walletId)
	ret0, _ := ret[0].([]models.WalletHistoryResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWalletHistory indicates an expected call of GetWalletHistory.
func (mr *MockWalletRepositoryMockRecorder) GetWalletHistory(walletId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWalletHistory", reflect.TypeOf((*MockWalletRepository)(nil).GetWalletHistory), walletId)
}

// GetWalletHistoryAmount mocks base method.
func (m *MockWalletRepository) GetWalletHistoryAmount(orderId int) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWalletHistoryAmount", orderId)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWalletHistoryAmount indicates an expected call of GetWalletHistoryAmount.
func (mr *MockWalletRepositoryMockRecorder) GetWalletHistoryAmount(orderId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWalletHistoryAmount", reflect.TypeOf((*MockWalletRepository)(nil).GetWalletHistoryAmount), orderId)
}

// IsWalletExist mocks base method.
func (m *MockWalletRepository) IsWalletExist(userID int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsWalletExist", userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsWalletExist indicates an expected call of IsWalletExist.
func (mr *MockWalletRepositoryMockRecorder) IsWalletExist(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsWalletExist", reflect.TypeOf((*MockWalletRepository)(nil).IsWalletExist), userID)
}