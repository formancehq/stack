package command

import (
	"context"

	"github.com/formancehq/ledger/pkg/core"
)

type executionContext struct {
	context.Context
	cache    Cache
	ingested chan struct{}
	accounts map[string]*core.AccountWithVolumes
}

func (ec *executionContext) GetAccountWithVolumes(ctx context.Context, address string) (*core.AccountWithVolumes, error) {
	account, ok := ec.accounts[address]
	if ok {
		return account, nil
	}

	account, release, err := ec.cache.GetAccountWithVolumes(ctx, address)
	if err != nil {
		return nil, err
	}

	ec.accounts[address] = account
	go func() {
		<-ec.ingested
		release()
	}()

	return account, nil
}

func (ctx *executionContext) SetIngested() {
	close(ctx.ingested)
}

func (ec *executionContext) RetainAccount(address string) error {
	_, err := ec.GetAccountWithVolumes(ec.Context, address)
	return err
}

func newExecutionContext(ctx context.Context, cache Cache) *executionContext {
	return &executionContext{
		Context:  ctx,
		cache:    cache,
		ingested: make(chan struct{}),
		accounts: map[string]*core.AccountWithVolumes{},
	}
}
