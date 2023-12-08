// Code generated by MockGen. DO NOT EDIT.
// Source: backend.go
//
// Generated by this command:
//
//	mockgen -source backend.go -destination backend_test.go -package analytics . Ledger
//
// Package analytics is a generated GoMock package.
package analytics

import (
	context "context"
	reflect "reflect"

	systemstore "github.com/formancehq/ledger/internal/storage/systemstore"
	api "github.com/formancehq/stack/libs/go-libs/api"
	gomock "go.uber.org/mock/gomock"
)

// MockLedger is a mock of Ledger interface.
type MockLedger struct {
	ctrl     *gomock.Controller
	recorder *MockLedgerMockRecorder
}

// MockLedgerMockRecorder is the mock recorder for MockLedger.
type MockLedgerMockRecorder struct {
	mock *MockLedger
}

// NewMockLedger creates a new mock instance.
func NewMockLedger(ctrl *gomock.Controller) *MockLedger {
	mock := &MockLedger{ctrl: ctrl}
	mock.recorder = &MockLedgerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLedger) EXPECT() *MockLedgerMockRecorder {
	return m.recorder
}

// CountAccounts mocks base method.
func (m *MockLedger) CountAccounts(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountAccounts", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountAccounts indicates an expected call of CountAccounts.
func (mr *MockLedgerMockRecorder) CountAccounts(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountAccounts", reflect.TypeOf((*MockLedger)(nil).CountAccounts), ctx)
}

// CountTransactions mocks base method.
func (m *MockLedger) CountTransactions(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountTransactions", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountTransactions indicates an expected call of CountTransactions.
func (mr *MockLedgerMockRecorder) CountTransactions(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountTransactions", reflect.TypeOf((*MockLedger)(nil).CountTransactions), ctx)
}

// MockBackend is a mock of Backend interface.
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend.
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance.
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return m.recorder
}

// AppID mocks base method.
func (m *MockBackend) AppID(ctx context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppID", ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AppID indicates an expected call of AppID.
func (mr *MockBackendMockRecorder) AppID(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppID", reflect.TypeOf((*MockBackend)(nil).AppID), ctx)
}

// GetLedgerStore mocks base method.
func (m *MockBackend) GetLedgerStore(ctx context.Context, l string) (Ledger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLedgerStore", ctx, l)
	ret0, _ := ret[0].(Ledger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLedgerStore indicates an expected call of GetLedgerStore.
func (mr *MockBackendMockRecorder) GetLedgerStore(ctx, l any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLedgerStore", reflect.TypeOf((*MockBackend)(nil).GetLedgerStore), ctx, l)
}

// ListLedgers mocks base method.
func (m *MockBackend) ListLedgers(ctx context.Context, query systemstore.ListLedgersQuery) (*api.Cursor[systemstore.Ledger], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListLedgers", ctx, query)
	ret0, _ := ret[0].(*api.Cursor[systemstore.Ledger])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListLedgers indicates an expected call of ListLedgers.
func (mr *MockBackendMockRecorder) ListLedgers(ctx, query any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListLedgers", reflect.TypeOf((*MockBackend)(nil).ListLedgers), ctx, query)
}