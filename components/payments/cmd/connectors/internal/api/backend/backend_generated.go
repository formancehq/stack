// Code generated by MockGen. DO NOT EDIT.
// Source: backend.go

// Package backend is a generated GoMock package.
package backend

import (
	context "context"
	reflect "reflect"

	connectors_manager "github.com/formancehq/payments/cmd/connectors/internal/api/connectors_manager"
	service "github.com/formancehq/payments/cmd/connectors/internal/api/service"
	storage "github.com/formancehq/payments/cmd/connectors/internal/storage"
	models "github.com/formancehq/payments/internal/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateBankAccount mocks base method.
func (m *MockService) CreateBankAccount(ctx context.Context, req *service.CreateBankAccountRequest) (*models.BankAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBankAccount", ctx, req)
	ret0, _ := ret[0].(*models.BankAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBankAccount indicates an expected call of CreateBankAccount.
func (mr *MockServiceMockRecorder) CreateBankAccount(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBankAccount", reflect.TypeOf((*MockService)(nil).CreateBankAccount), ctx, req)
}

// CreateTransferInitiation mocks base method.
func (m *MockService) CreateTransferInitiation(ctx context.Context, req *service.CreateTransferInitiationRequest) (*models.TransferInitiation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransferInitiation", ctx, req)
	ret0, _ := ret[0].(*models.TransferInitiation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransferInitiation indicates an expected call of CreateTransferInitiation.
func (mr *MockServiceMockRecorder) CreateTransferInitiation(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransferInitiation", reflect.TypeOf((*MockService)(nil).CreateTransferInitiation), ctx, req)
}

// DeleteTransferInitiation mocks base method.
func (m *MockService) DeleteTransferInitiation(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTransferInitiation", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTransferInitiation indicates an expected call of DeleteTransferInitiation.
func (mr *MockServiceMockRecorder) DeleteTransferInitiation(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTransferInitiation", reflect.TypeOf((*MockService)(nil).DeleteTransferInitiation), ctx, id)
}

// ListConnectors mocks base method.
func (m *MockService) ListConnectors(ctx context.Context) ([]*models.Connector, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListConnectors", ctx)
	ret0, _ := ret[0].([]*models.Connector)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListConnectors indicates an expected call of ListConnectors.
func (mr *MockServiceMockRecorder) ListConnectors(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListConnectors", reflect.TypeOf((*MockService)(nil).ListConnectors), ctx)
}

// Ping mocks base method.
func (m *MockService) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockServiceMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockService)(nil).Ping))
}

// RetryTransferInitiation mocks base method.
func (m *MockService) RetryTransferInitiation(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetryTransferInitiation", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RetryTransferInitiation indicates an expected call of RetryTransferInitiation.
func (mr *MockServiceMockRecorder) RetryTransferInitiation(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetryTransferInitiation", reflect.TypeOf((*MockService)(nil).RetryTransferInitiation), ctx, id)
}

// UpdateTransferInitiationStatus mocks base method.
func (m *MockService) UpdateTransferInitiationStatus(ctx context.Context, transferID string, req *service.UpdateTransferInitiationStatusRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTransferInitiationStatus", ctx, transferID, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTransferInitiationStatus indicates an expected call of UpdateTransferInitiationStatus.
func (mr *MockServiceMockRecorder) UpdateTransferInitiationStatus(ctx, transferID, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransferInitiationStatus", reflect.TypeOf((*MockService)(nil).UpdateTransferInitiationStatus), ctx, transferID, req)
}

// MockManager is a mock of Manager interface.
type MockManager[ConnectorConfig models.ConnectorConfigObject] struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder[ConnectorConfig]
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder[ConnectorConfig models.ConnectorConfigObject] struct {
	mock *MockManager[ConnectorConfig]
}

// NewMockManager creates a new mock instance.
func NewMockManager[ConnectorConfig models.ConnectorConfigObject](ctrl *gomock.Controller) *MockManager[ConnectorConfig] {
	mock := &MockManager[ConnectorConfig]{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder[ConnectorConfig]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager[ConnectorConfig]) EXPECT() *MockManagerMockRecorder[ConnectorConfig] {
	return m.recorder
}

// Connectors mocks base method.
func (m *MockManager[ConnectorConfig]) Connectors() map[string]*connectors_manager.ConnectorManager {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connectors")
	ret0, _ := ret[0].(map[string]*connectors_manager.ConnectorManager)
	return ret0
}

// Connectors indicates an expected call of Connectors.
func (mr *MockManagerMockRecorder[ConnectorConfig]) Connectors() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connectors", reflect.TypeOf((*MockManager[ConnectorConfig])(nil).Connectors))
}

// CreateWebhookAndContext mocks base method.
func (m *MockManager[ConnectorConfig]) CreateWebhookAndContext(ctx context.Context, webhook *models.Webhook) (context.Context, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWebhookAndContext", ctx, webhook)
	ret0, _ := ret[0].(context.Context)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWebhookAndContext indicates an expected call of CreateWebhookAndContext.
func (mr *MockManagerMockRecorder[ConnectorConfig]) CreateWebhookAndContext(ctx, webhook interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWebhookAndContext", reflect.TypeOf((*MockManager[ConnectorConfig])(nil).CreateWebhookAndContext), ctx, webhook)
}

// Install mocks base method.
func (m *MockManager[ConnectorConfig]) Install(ctx context.Context, name string, config ConnectorConfig) (models.ConnectorID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Install", ctx, name, config)
	ret0, _ := ret[0].(models.ConnectorID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Install indicates an expected call of Install.
func (mr *MockManagerMockRecorder[ConnectorConfig]) Install(ctx, name, config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Install", reflect.TypeOf((*MockManager[ConnectorConfig])(nil).Install), ctx, name, config)
}

// IsInstalled mocks base method.
func (m *MockManager[ConnectorConfig]) IsInstalled(ctx context.Context, connectorID models.ConnectorID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsInstalled", ctx, connectorID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsInstalled indicates an expected call of IsInstalled.
func (mr *MockManagerMockRecorder[ConnectorConfig]) IsInstalled(ctx, connectorID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsInstalled", reflect.TypeOf((*MockManager[ConnectorConfig])(nil).IsInstalled), ctx, connectorID)
}

// ListTasksStates mocks base method.
func (m *MockManager[ConnectorConfig]) ListTasksStates(ctx context.Context, connectorID models.ConnectorID, pagination storage.PaginatorQuery) ([]*models.Task, storage.PaginationDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTasksStates", ctx, connectorID, pagination)
	ret0, _ := ret[0].([]*models.Task)
	ret1, _ := ret[1].(storage.PaginationDetails)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListTasksStates indicates an expected call of ListTasksStates.
func (mr *MockManagerMockRecorder[ConnectorConfig]) ListTasksStates(ctx, connectorID, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTasksStates", reflect.TypeOf((*MockManager[ConnectorConfig])(nil).ListTasksStates), ctx, connectorID, pagination)
}

// ReadConfig mocks base method.
func (m *MockManager[ConnectorConfig]) ReadConfig(ctx context.Context, connectorID models.ConnectorID) (ConnectorConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadConfig", ctx, connectorID)
	ret0, _ := ret[0].(ConnectorConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadConfig indicates an expected call of ReadConfig.
func (mr *MockManagerMockRecorder[ConnectorConfig]) ReadConfig(ctx, connectorID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadConfig", reflect.TypeOf((*MockManager[ConnectorConfig])(nil).ReadConfig), ctx, connectorID)
}

// ReadTaskState mocks base method.
func (m *MockManager[ConnectorConfig]) ReadTaskState(ctx context.Context, connectorID models.ConnectorID, taskID uuid.UUID) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadTaskState", ctx, connectorID, taskID)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadTaskState indicates an expected call of ReadTaskState.
func (mr *MockManagerMockRecorder[ConnectorConfig]) ReadTaskState(ctx, connectorID, taskID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadTaskState", reflect.TypeOf((*MockManager[ConnectorConfig])(nil).ReadTaskState), ctx, connectorID, taskID)
}

// Reset mocks base method.
func (m *MockManager[ConnectorConfig]) Reset(ctx context.Context, connectorID models.ConnectorID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reset", ctx, connectorID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reset indicates an expected call of Reset.
func (mr *MockManagerMockRecorder[ConnectorConfig]) Reset(ctx, connectorID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockManager[ConnectorConfig])(nil).Reset), ctx, connectorID)
}

// Uninstall mocks base method.
func (m *MockManager[ConnectorConfig]) Uninstall(ctx context.Context, connectorID models.ConnectorID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Uninstall", ctx, connectorID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Uninstall indicates an expected call of Uninstall.
func (mr *MockManagerMockRecorder[ConnectorConfig]) Uninstall(ctx, connectorID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Uninstall", reflect.TypeOf((*MockManager[ConnectorConfig])(nil).Uninstall), ctx, connectorID)
}

// UpdateConfig mocks base method.
func (m *MockManager[ConnectorConfig]) UpdateConfig(ctx context.Context, connectorID models.ConnectorID, config ConnectorConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateConfig", ctx, connectorID, config)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateConfig indicates an expected call of UpdateConfig.
func (mr *MockManagerMockRecorder[ConnectorConfig]) UpdateConfig(ctx, connectorID, config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConfig", reflect.TypeOf((*MockManager[ConnectorConfig])(nil).UpdateConfig), ctx, connectorID, config)
}

// MockServiceBackend is a mock of ServiceBackend interface.
type MockServiceBackend struct {
	ctrl     *gomock.Controller
	recorder *MockServiceBackendMockRecorder
}

// MockServiceBackendMockRecorder is the mock recorder for MockServiceBackend.
type MockServiceBackendMockRecorder struct {
	mock *MockServiceBackend
}

// NewMockServiceBackend creates a new mock instance.
func NewMockServiceBackend(ctrl *gomock.Controller) *MockServiceBackend {
	mock := &MockServiceBackend{ctrl: ctrl}
	mock.recorder = &MockServiceBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceBackend) EXPECT() *MockServiceBackendMockRecorder {
	return m.recorder
}

// GetService mocks base method.
func (m *MockServiceBackend) GetService() Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetService")
	ret0, _ := ret[0].(Service)
	return ret0
}

// GetService indicates an expected call of GetService.
func (mr *MockServiceBackendMockRecorder) GetService() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetService", reflect.TypeOf((*MockServiceBackend)(nil).GetService))
}

// MockManagerBackend is a mock of ManagerBackend interface.
type MockManagerBackend[ConnectorConfig models.ConnectorConfigObject] struct {
	ctrl     *gomock.Controller
	recorder *MockManagerBackendMockRecorder[ConnectorConfig]
}

// MockManagerBackendMockRecorder is the mock recorder for MockManagerBackend.
type MockManagerBackendMockRecorder[ConnectorConfig models.ConnectorConfigObject] struct {
	mock *MockManagerBackend[ConnectorConfig]
}

// NewMockManagerBackend creates a new mock instance.
func NewMockManagerBackend[ConnectorConfig models.ConnectorConfigObject](ctrl *gomock.Controller) *MockManagerBackend[ConnectorConfig] {
	mock := &MockManagerBackend[ConnectorConfig]{ctrl: ctrl}
	mock.recorder = &MockManagerBackendMockRecorder[ConnectorConfig]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManagerBackend[ConnectorConfig]) EXPECT() *MockManagerBackendMockRecorder[ConnectorConfig] {
	return m.recorder
}

// GetManager mocks base method.
func (m *MockManagerBackend[ConnectorConfig]) GetManager() Manager[ConnectorConfig] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManager")
	ret0, _ := ret[0].(Manager[ConnectorConfig])
	return ret0
}

// GetManager indicates an expected call of GetManager.
func (mr *MockManagerBackendMockRecorder[ConnectorConfig]) GetManager() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManager", reflect.TypeOf((*MockManagerBackend[ConnectorConfig])(nil).GetManager))
}
