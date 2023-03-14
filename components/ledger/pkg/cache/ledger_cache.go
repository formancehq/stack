package cache

import (
	"context"
	"strings"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/storage"
)

type Ledger struct {
	cache  *Cache
	ledger string
}

func (c *Ledger) GetAccountWithVolumes(ctx context.Context, address string) (*core.AccountWithVolumes, error) {

	address = strings.TrimPrefix(address, "@")

	rawAccount, err := c.cache.cache.Get(c.accountKey(address))
	if err != nil {
		store, _, err := c.cache.driver.GetLedgerStore(ctx, c.ledger, false)
		if err != nil && err != storage.ErrLedgerStoreNotFound {
			return nil, err
		}
		if err == storage.ErrLedgerStoreNotFound {
			return nil, nil
		}

		// TODO: Rename later ?
		account, err := store.ComputeAccount(ctx, address)
		if err != nil {
			return nil, err
		}

		if err := c.cache.cache.Set(c.accountKey(account.Address), account); err != nil {
			panic(err)
		}

		return account, nil
	}
	cp := rawAccount.(*core.AccountWithVolumes).Copy()

	return &cp, nil
}

func (c *Ledger) Update(ctx context.Context, tx *TxInfo, accounts core.AccountsAssetsVolumes) {
	c.cache.lastInsertedTransaction[c.ledger] = tx
	for address, volumes := range accounts {
		rawAccount, err := c.cache.cache.Get(c.accountKey(address))
		if err != nil {
			// Cannot update cache, item maybe evicted
			continue
		}
		account := rawAccount.(*core.AccountWithVolumes)
		account.Volumes = volumes
		account.Balances = volumes.Balances()
		if err := c.cache.cache.Set(c.accountKey(address), account); err != nil {
			panic(err)
		}
	}
}

func (c *Ledger) GetLastTransaction(ctx context.Context) (*TxInfo, error) {
	if c.cache.lastInsertedTransaction[c.ledger] == nil {
		store, _, err := c.cache.driver.GetLedgerStore(ctx, c.ledger, false)
		if err != nil && err != storage.ErrLedgerStoreNotFound {
			return nil, err
		}
		if err == storage.ErrLedgerStoreNotFound {
			return nil, nil
		}
		lastTx, err := store.GetLastTransaction(ctx)
		if err != nil {
			return nil, err
		}

		if lastTx != nil {
			c.cache.lastInsertedTransaction[c.ledger] = &TxInfo{
				Date: lastTx.Timestamp,
				ID:   lastTx.ID,
			}
		}
	}

	return c.cache.lastInsertedTransaction[c.ledger], nil
}

func (c *Ledger) UpdateAccountMetadata(ctx context.Context, address string, m core.Metadata) error {
	account, err := c.GetAccountWithVolumes(ctx, address)
	if err != nil {
		return err
	}
	account.Metadata = account.Metadata.Merge(m)
	_ = c.cache.cache.Set(c.accountKey(address), account)
	return nil
}

func (c *Ledger) accountKey(address string) string {
	return c.ledger + "-" + address
}

func newLedgerCache(cache *Cache, ledger string) *Ledger {
	return &Ledger{
		cache:  cache,
		ledger: ledger,
	}
}
