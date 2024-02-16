package backend

import (
	"context"

	manager "github.com/formancehq/payments/cmd/connectors/internal/api/connectors_manager"
	"github.com/formancehq/payments/cmd/connectors/internal/api/service"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

//go:generate mockgen -source backend.go -destination backend_generated.go -package backend . Service
type Service interface {
	Ping() error
	CreateBankAccount(ctx context.Context, req *service.CreateBankAccountRequest) (*models.BankAccount, error)
	ForwardBankAccountToConnector(ctx context.Context, id string, req *service.ForwardBankAccountToConnectorRequest) (*models.BankAccount, error)
	UpdateBankAccountMetadata(ctx context.Context, id string, req *service.UpdateBankAccountMetadataRequest) error
	ListConnectors(ctx context.Context) ([]*models.Connector, error)
	CreateTransferInitiation(ctx context.Context, req *service.CreateTransferInitiationRequest) (*models.TransferInitiation, error)
	UpdateTransferInitiationStatus(ctx context.Context, transferID string, req *service.UpdateTransferInitiationStatusRequest) error
	RetryTransferInitiation(ctx context.Context, id string) error
	DeleteTransferInitiation(ctx context.Context, id string) error
	ReverseTransferInitiation(ctx context.Context, transferID string, req *service.ReverseTransferInitiationRequest) (*models.TransferReversal, error)
}

//go:generate mockgen -source backend.go -destination backend_generated.go -package backend . Manager
type Manager[ConnectorConfig models.ConnectorConfigObject] interface {
	IsInstalled(ctx context.Context, connectorID models.ConnectorID) (bool, error)
	Connectors() map[string]*manager.ConnectorManager
	ReadConfig(ctx context.Context, connectorID models.ConnectorID) (ConnectorConfig, error)
	UpdateConfig(ctx context.Context, connectorID models.ConnectorID, config ConnectorConfig) error
	ListTasksStates(ctx context.Context, connectorID models.ConnectorID, q storage.ListTasksQuery) (*api.Cursor[models.Task], error)
	CreateWebhookAndContext(ctx context.Context, webhook *models.Webhook) (context.Context, error)
	ReadTaskState(ctx context.Context, connectorID models.ConnectorID, taskID uuid.UUID) (*models.Task, error)
	Install(ctx context.Context, name string, config ConnectorConfig) (models.ConnectorID, error)
	Reset(ctx context.Context, connectorID models.ConnectorID) error
	Uninstall(ctx context.Context, connectorID models.ConnectorID) error
}

type ServiceBackend interface {
	GetService() Service
}

type DefaultServiceBackend struct {
	service Service
}

func (d DefaultServiceBackend) GetService() Service {
	return d.service
}

func NewDefaultBackend(service Service) ServiceBackend {
	return &DefaultServiceBackend{
		service: service,
	}
}

type ManagerBackend[ConnectorConfig models.ConnectorConfigObject] interface {
	GetManager() Manager[ConnectorConfig]
}

type DefaultManagerBackend[ConnectorConfig models.ConnectorConfigObject] struct {
	manager Manager[ConnectorConfig]
}

func (m DefaultManagerBackend[ConnectorConfig]) GetManager() Manager[ConnectorConfig] {
	return m.manager
}

func NewDefaultManagerBackend[ConnectorConfig models.ConnectorConfigObject](manager Manager[ConnectorConfig]) ManagerBackend[ConnectorConfig] {
	return DefaultManagerBackend[ConnectorConfig]{
		manager: manager,
	}
}
