package system

import (
	"context"

	"github.com/formancehq/ledger/internal/tracing"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	ledger "github.com/formancehq/ledger/internal"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
)

//go:generate mockgen -source controller.go -destination controller_generated.go -package system . Controller

type Controller interface {
	GetLedgerController(ctx context.Context, name string) (ledgercontroller.Controller, error)
	GetLedger(ctx context.Context, name string) (*ledger.Ledger, error)
	ListLedgers(ctx context.Context, query ledgercontroller.ListLedgersQuery) (*bunpaginate.Cursor[ledger.Ledger], error)
	// CreateLedger can return following errors:
	//  * ErrLedgerAlreadyExists
	CreateLedger(ctx context.Context, name string, configuration ledger.Configuration) error
	UpdateLedgerMetadata(ctx context.Context, name string, m map[string]string) error
	DeleteLedgerMetadata(ctx context.Context, param string, key string) error
}

type DefaultController struct {
	resolver *ledgercontroller.Resolver
	store    Store
}

func (c *DefaultController) GetLedgerController(ctx context.Context, name string) (ledgercontroller.Controller, error) {
	return tracing.Trace(ctx, "GetLedgerController", func(ctx context.Context) (ledgercontroller.Controller, error) {
		return c.resolver.GetLedger(ctx, name)
	})
}

func (c *DefaultController) CreateLedger(ctx context.Context, name string, configuration ledger.Configuration) error {
	return tracing.SkipResult(tracing.Trace(ctx, "CreateLedger", tracing.NoResult(func(ctx context.Context) error {
		configuration.SetDefaults()
		return c.resolver.CreateLedger(ctx, name, configuration)
	})))
}

func (c *DefaultController) GetLedger(ctx context.Context, name string) (*ledger.Ledger, error) {
	return tracing.Trace(ctx, "GetLedger", func(ctx context.Context) (*ledger.Ledger, error) {
		return c.store.GetLedger(ctx, name)
	})
}

func (c *DefaultController) ListLedgers(ctx context.Context, query ledgercontroller.ListLedgersQuery) (*bunpaginate.Cursor[ledger.Ledger], error) {
	return tracing.Trace(ctx, "ListLedgers", func(ctx context.Context) (*bunpaginate.Cursor[ledger.Ledger], error) {
		return c.store.ListLedgers(ctx, query)
	})
}

func (c *DefaultController) UpdateLedgerMetadata(ctx context.Context, name string, m map[string]string) error {
	return tracing.SkipResult(tracing.Trace(ctx, "UpdateLedgerMetadata", tracing.NoResult(func(ctx context.Context) error {
		return c.store.UpdateLedgerMetadata(ctx, name, m)
	})))
}

func (c *DefaultController) DeleteLedgerMetadata(ctx context.Context, param string, key string) error {
	return tracing.SkipResult(tracing.Trace(ctx, "DeleteLedgerMetadata", tracing.NoResult(func(ctx context.Context) error {
		return c.store.DeleteLedgerMetadata(ctx, param, key)
	})))
}

func NewController(resolver *ledgercontroller.Resolver, store Store) *DefaultController {
	return &DefaultController{
		resolver: resolver,
		store:    store,
	}
}
