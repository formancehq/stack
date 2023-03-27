package vm

import (
	"context"

	"github.com/formancehq/ledger/pkg/core"
)

type Store interface {
	GetAccountWithVolumes(ctx context.Context, address string) (*core.AccountWithVolumes, error)
}
type StoreFn func(ctx context.Context, address string) (*core.AccountWithVolumes, error)

func (fn StoreFn) GetAccountWithVolumes(ctx context.Context, address string) (*core.AccountWithVolumes, error) {
	return fn(ctx, address)
}

var EmptyStore = StoreFn(func(ctx context.Context, address string) (*core.AccountWithVolumes, error) {
	return nil, nil
})
