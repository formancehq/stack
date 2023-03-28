package ledger_test

import (
	"context"
	"testing"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/ledgertesting"
	"github.com/formancehq/ledger/pkg/storage"
	"github.com/formancehq/ledger/pkg/storage/sqlstorage"
	"github.com/stretchr/testify/assert"
)

func TestAccounts(t *testing.T) {
	d := ledgertesting.StorageDriver(t)

	assert.NoError(t, d.Initialize(context.Background()))

	defer func(d *sqlstorage.Driver, ctx context.Context) {
		assert.NoError(t, d.Close(ctx))
	}(d, context.Background())

	store, _, err := d.GetLedgerStore(context.Background(), "foo", true)
	assert.NoError(t, err)

	_, err = store.Initialize(context.Background())
	assert.NoError(t, err)

	t.Run("success balance", func(t *testing.T) {
		q := storage.AccountsQuery{
			PageSize: 10,
			Filters: storage.AccountsQueryFilters{
				Balance: "50",
			},
		}

		_, err := store.GetAccounts(context.Background(), q)
		assert.NoError(t, err, "balance filter should not fail")
	})

	t.Run("panic invalid balance", func(t *testing.T) {
		q := storage.AccountsQuery{
			PageSize: 10,
			Filters: storage.AccountsQueryFilters{
				Balance: "TEST",
			},
		}

		assert.PanicsWithError(
			t, `invalid balance parameter: strconv.ParseInt: parsing "TEST": invalid syntax`,

			func() {
				_, _ = store.GetAccounts(context.Background(), q)
			}, "invalid balance in storage should panic")
	})

	t.Run("panic invalid balance operator", func(t *testing.T) {
		assert.PanicsWithValue(t, "invalid balance operator parameter", func() {
			q := storage.AccountsQuery{
				PageSize: 10,
				Filters: storage.AccountsQueryFilters{
					Balance:         "50",
					BalanceOperator: "TEST",
				},
			}

			_, _ = store.GetAccounts(context.Background(), q)
		}, "invalid balance operator in storage should panic")
	})

	t.Run("success balance operator", func(t *testing.T) {
		q := storage.AccountsQuery{
			PageSize: 10,
			Filters: storage.AccountsQueryFilters{
				Balance:         "50",
				BalanceOperator: storage.BalanceOperatorGte,
			},
		}

		_, err := store.GetAccounts(context.Background(), q)
		assert.NoError(t, err, "balance operator filter should not fail")
	})

	t.Run("success account insertion", func(t *testing.T) {
		addr := "test:account"
		metadata := core.Metadata(map[string]any{
			"foo": "bar",
		})

		err := store.UpdateAccountMetadata(context.Background(), addr, metadata)
		assert.NoError(t, err, "account insertion should not fail")

		account, err := store.GetAccount(context.Background(), addr)
		assert.NoError(t, err, "account retrieval should not fail")

		assert.Equal(t, addr, account.Address, "account address should match")
		assert.Equal(t, metadata, account.Metadata, "account metadata should match")
	})

	t.Run("success multiple account insertions", func(t *testing.T) {
		accounts := []core.Account{
			{
				Address:  "test:account1",
				Metadata: core.Metadata(map[string]any{"foo1": "bar1"}),
			},
			{
				Address:  "test:account2",
				Metadata: core.Metadata(map[string]any{"foo2": "bar2"}),
			},
			{
				Address:  "test:account3",
				Metadata: core.Metadata(map[string]any{"foo3": "bar3"}),
			},
		}

		err := store.UpdateAccountsMetadata(context.Background(), accounts)
		assert.NoError(t, err, "account insertion should not fail")

		for _, account := range accounts {
			acc, err := store.GetAccount(context.Background(), account.Address)
			assert.NoError(t, err, "account retrieval should not fail")

			assert.Equal(t, account.Address, acc.Address, "account address should match")
			assert.Equal(t, account.Metadata, acc.Metadata, "account metadata should match")
		}
	})
}
