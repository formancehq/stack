package backend

import (
	"context"

	"github.com/formancehq/payments/cmd/api/internal/api/service"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

//go:generate mockgen -source backend.go -destination backend_generated.go -package backend . Service
type Service interface {
	Ping() error
	ListAccounts(ctx context.Context, q storage.ListAccountsQuery) (*api.Cursor[models.Account], error)
	GetAccount(ctx context.Context, id string) (*models.Account, error)
	ListBalances(ctx context.Context, q storage.ListBalancesQuery) (*api.Cursor[models.Balance], error)
	ListBankAccounts(ctx context.Context, a storage.ListBankAccountQuery) (*api.Cursor[models.BankAccount], error)
	GetBankAccount(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error)
	ListTransferInitiations(ctx context.Context, q storage.ListTransferInitiationsQuery) (*api.Cursor[models.TransferInitiation], error)
	ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error)
	CreatePayment(ctx context.Context, req *service.CreatePaymentRequest) (*models.Payment, error)
	ListPayments(ctx context.Context, q storage.ListPaymentsQuery) (*api.Cursor[models.Payment], error)
	GetPayment(ctx context.Context, id string) (*models.Payment, error)
	UpdatePaymentMetadata(ctx context.Context, paymentID models.PaymentID, metadata map[string]string) error
	CreatePool(ctx context.Context, req *service.CreatePoolRequest) (*models.Pool, error)
	AddAccountToPool(ctx context.Context, poolID string, req *service.AddAccountToPoolRequest) error
	RemoveAccountFromPool(ctx context.Context, poolID string, accountID string) error
	ListPools(ctx context.Context, q storage.ListPoolsQuery) (*api.Cursor[models.Pool], error)
	GetPool(ctx context.Context, poolID string) (*models.Pool, error)
	GetPoolBalance(ctx context.Context, poolID string, atTime string) (*service.GetPoolBalanceResponse, error)
	DeletePool(ctx context.Context, poolID string) error
}

type Backend interface {
	GetService() Service
}

type DefaultBackend struct {
	service Service
}

func (d DefaultBackend) GetService() Service {
	return d.service
}

func NewDefaultBackend(service Service) Backend {
	return &DefaultBackend{
		service: service,
	}
}
