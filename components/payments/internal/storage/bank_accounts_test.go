package storage

import (
	"context"
	"testing"
	"time"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/go-libs/query"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	defaultBankAccount = models.BankAccount{
		ID:            uuid.New(),
		CreatedAt:     now.Add(-60 * time.Minute).UTC().Time,
		Name:          "test1",
		AccountNumber: pointer.For("12345678"),
		Country:       pointer.For("US"),
		Metadata: map[string]string{
			"foo": "bar",
		},
	}

	bcID2               = uuid.New()
	defaultBankAccount2 = models.BankAccount{
		ID:           bcID2,
		CreatedAt:    now.Add(-30 * time.Minute).UTC().Time,
		Name:         "test2",
		IBAN:         pointer.For("DE89370400440532013000"),
		SwiftBicCode: pointer.For("COBADEFFXXX"),
		Country:      pointer.For("DE"),
		Metadata: map[string]string{
			"foo2": "bar2",
		},
		RelatedAccounts: []models.BankAccountRelatedAccount{
			{
				BankAccountID: bcID2,
				AccountID:     defaultAccounts[0].ID,
				ConnectorID:   defaultConnector.ID,
				CreatedAt:     now.Add(-30 * time.Minute).UTC().Time,
			},
		},
	}

	// No metadata
	defaultBankAccount3 = models.BankAccount{
		ID:            uuid.New(),
		CreatedAt:     now.Add(-55 * time.Minute).UTC().Time,
		Name:          "test1",
		AccountNumber: pointer.For("12345678"),
		Country:       pointer.For("US"),
	}
)

func upsertBankAccount(t *testing.T, ctx context.Context, storage Storage, bankAccounts models.BankAccount) {
	require.NoError(t, storage.BankAccountsUpsert(ctx, bankAccounts))
}

func TestBankAccountsUpsert(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertBankAccount(t, ctx, store, defaultBankAccount)
	upsertBankAccount(t, ctx, store, defaultBankAccount2)

	t.Run("upsert with same id", func(t *testing.T) {
		ba := models.BankAccount{
			ID:            defaultBankAccount.ID,
			CreatedAt:     now.UTC().Time,
			Name:          "changed",
			AccountNumber: pointer.For("987654321"),
			Country:       pointer.For("CA"),
			Metadata: map[string]string{
				"changed": "changed",
			},
		}

		require.NoError(t, store.BankAccountsUpsert(ctx, ba))

		actual, err := store.BankAccountsGet(ctx, ba.ID, true)
		require.NoError(t, err)
		// Should not update the bank account
		compareBankAccounts(t, defaultBankAccount, *actual)
	})

	t.Run("unknown connector id", func(t *testing.T) {
		ba := models.BankAccount{
			ID:            uuid.New(),
			CreatedAt:     now.UTC().Time,
			Name:          "foo",
			AccountNumber: pointer.For("12345678"),
			Country:       pointer.For("US"),
			Metadata: map[string]string{
				"foo": "bar",
			},
			RelatedAccounts: []models.BankAccountRelatedAccount{
				{
					BankAccountID: uuid.New(),
					AccountID:     defaultAccounts[0].ID,
					ConnectorID: models.ConnectorID{
						Reference: uuid.New(),
						Provider:  "unknown",
					},
					CreatedAt: now.UTC().Time,
				},
			},
		}

		require.Error(t, store.BankAccountsUpsert(ctx, ba))
		b, err := store.BankAccountsGet(ctx, ba.ID, true)
		require.Error(t, err)
		require.Nil(t, b)
	})
}

func TestBankAccountsUpdateMetadata(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertBankAccount(t, ctx, store, defaultBankAccount)
	upsertBankAccount(t, ctx, store, defaultBankAccount2)
	upsertBankAccount(t, ctx, store, defaultBankAccount3)

	t.Run("update metadata", func(t *testing.T) {
		metadata := map[string]string{
			"test1": "test2",
			"test3": "test4",
		}

		// redeclare it in order to not update the map of global variable
		acc := models.BankAccount{
			ID:            defaultBankAccount.ID,
			CreatedAt:     now.Add(-60 * time.Minute).UTC().Time,
			Name:          "test1",
			AccountNumber: pointer.For("12345678"),
			Country:       pointer.For("US"),
			Metadata: map[string]string{
				"foo": "bar",
			},
		}
		for k, v := range metadata {
			acc.Metadata[k] = v
		}

		require.NoError(t, store.BankAccountsUpdateMetadata(ctx, defaultBankAccount.ID, metadata))

		actual, err := store.BankAccountsGet(ctx, defaultBankAccount.ID, true)
		require.NoError(t, err)
		compareBankAccounts(t, acc, *actual)
	})

	t.Run("update same metadata", func(t *testing.T) {
		metadata := map[string]string{
			"foo2": "bar3",
		}

		acc := models.BankAccount{
			ID:           bcID2,
			CreatedAt:    now.Add(-30 * time.Minute).UTC().Time,
			Name:         "test2",
			IBAN:         pointer.For("DE89370400440532013000"),
			SwiftBicCode: pointer.For("COBADEFFXXX"),
			Country:      pointer.For("DE"),
			Metadata: map[string]string{
				"foo2": "bar2",
			},
			RelatedAccounts: []models.BankAccountRelatedAccount{
				{
					BankAccountID: bcID2,
					AccountID:     defaultAccounts[0].ID,
					ConnectorID:   defaultConnector.ID,
					CreatedAt:     now.Add(-30 * time.Minute).UTC().Time,
				},
			},
		}
		for k, v := range metadata {
			acc.Metadata[k] = v
		}

		require.NoError(t, store.BankAccountsUpdateMetadata(ctx, defaultBankAccount2.ID, metadata))

		actual, err := store.BankAccountsGet(ctx, defaultBankAccount2.ID, true)
		require.NoError(t, err)
		compareBankAccounts(t, acc, *actual)
	})

	t.Run("update metadata of bank accounts with nil map", func(t *testing.T) {
		metadata := map[string]string{
			"test1": "test2",
			"test3": "test4",
		}

		// redeclare it in order to not update the map of global variable
		acc := models.BankAccount{
			ID:            defaultBankAccount3.ID,
			CreatedAt:     now.Add(-55 * time.Minute).UTC().Time,
			Name:          "test1",
			AccountNumber: pointer.For("12345678"),
			Country:       pointer.For("US"),
		}
		acc.Metadata = make(map[string]string)
		for k, v := range metadata {
			acc.Metadata[k] = v
		}

		require.NoError(t, store.BankAccountsUpdateMetadata(ctx, defaultBankAccount3.ID, metadata))

		actual, err := store.BankAccountsGet(ctx, defaultBankAccount3.ID, true)
		require.NoError(t, err)
		compareBankAccounts(t, acc, *actual)
	})
}

func TestBankAccountsGet(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertBankAccount(t, ctx, store, defaultBankAccount)
	upsertBankAccount(t, ctx, store, defaultBankAccount2)
	upsertBankAccount(t, ctx, store, defaultBankAccount3)

	t.Run("get bank account without related accounts", func(t *testing.T) {
		actual, err := store.BankAccountsGet(ctx, defaultBankAccount.ID, true)
		require.NoError(t, err)
		compareBankAccounts(t, defaultBankAccount, *actual)
	})

	t.Run("get bank account without metadata", func(t *testing.T) {
		actual, err := store.BankAccountsGet(ctx, defaultBankAccount3.ID, true)
		require.NoError(t, err)
		compareBankAccounts(t, defaultBankAccount3, *actual)
	})

	t.Run("get bank account with related accounts", func(t *testing.T) {
		actual, err := store.BankAccountsGet(ctx, defaultBankAccount2.ID, true)
		require.NoError(t, err)
		compareBankAccounts(t, defaultBankAccount2, *actual)
	})

	t.Run("get unknown bank account", func(t *testing.T) {
		actual, err := store.BankAccountsGet(ctx, uuid.New(), true)
		require.Error(t, err)
		require.Nil(t, actual)
	})

	t.Run("get bank account with expand to false", func(t *testing.T) {
		acc := models.BankAccount{
			ID:        defaultBankAccount.ID,
			CreatedAt: now.Add(-60 * time.Minute).UTC().Time,
			Name:      "test1",
			Country:   pointer.For("US"),
			Metadata: map[string]string{
				"foo": "bar",
			},
		}

		actual, err := store.BankAccountsGet(ctx, defaultBankAccount.ID, false)
		require.NoError(t, err)
		compareBankAccounts(t, acc, *actual)
	})

	t.Run("get bank account with expand to false 2", func(t *testing.T) {
		acc := models.BankAccount{
			ID:        bcID2,
			CreatedAt: now.Add(-30 * time.Minute).UTC().Time,
			Name:      "test2",
			Country:   pointer.For("DE"),
			Metadata: map[string]string{
				"foo2": "bar2",
			},
			RelatedAccounts: []models.BankAccountRelatedAccount{
				{
					BankAccountID: bcID2,
					AccountID:     defaultAccounts[0].ID,
					ConnectorID:   defaultConnector.ID,
					CreatedAt:     now.Add(-30 * time.Minute).UTC().Time,
				},
			},
		}

		actual, err := store.BankAccountsGet(ctx, bcID2, false)
		require.NoError(t, err)
		compareBankAccounts(t, acc, *actual)
	})
}

func TestBankAccountsList(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	d1 := models.BankAccount{
		ID:        defaultBankAccount.ID,
		CreatedAt: defaultBankAccount.CreatedAt,
		Name:      defaultBankAccount.Name,
		Country:   defaultBankAccount.Country,
		Metadata:  defaultBankAccount.Metadata,
	}

	d2 := models.BankAccount{
		ID:              defaultBankAccount2.ID,
		CreatedAt:       defaultBankAccount2.CreatedAt,
		Name:            defaultBankAccount2.Name,
		Country:         defaultBankAccount2.Country,
		Metadata:        defaultBankAccount2.Metadata,
		RelatedAccounts: defaultBankAccount2.RelatedAccounts,
	}
	_ = d2

	d3 := models.BankAccount{
		ID:        defaultBankAccount3.ID,
		CreatedAt: defaultBankAccount3.CreatedAt,
		Name:      defaultBankAccount3.Name,
		Country:   defaultBankAccount3.Country,
	}

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertBankAccount(t, ctx, store, defaultBankAccount)
	upsertBankAccount(t, ctx, store, defaultBankAccount2)
	upsertBankAccount(t, ctx, store, defaultBankAccount3)

	t.Run("list bank accounts by name", func(t *testing.T) {
		q := NewListBankAccountsQuery(
			bunpaginate.NewPaginatedQueryOptions(BankAccountQuery{}).
				WithPageSize(15).
				WithQueryBuilder(query.Match("name", "test1")),
		)

		cursor, err := store.BankAccountsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
		compareBankAccounts(t, d3, cursor.Data[0])
		compareBankAccounts(t, d1, cursor.Data[1])
	})

	t.Run("list bank accounts by name 2", func(t *testing.T) {
		q := NewListBankAccountsQuery(
			bunpaginate.NewPaginatedQueryOptions(BankAccountQuery{}).
				WithPageSize(15).
				WithQueryBuilder(query.Match("name", "test2")),
		)

		cursor, err := store.BankAccountsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
		compareBankAccounts(t, d2, cursor.Data[0])
	})

	t.Run("list bank accounts by unknown name", func(t *testing.T) {
		q := NewListBankAccountsQuery(
			bunpaginate.NewPaginatedQueryOptions(BankAccountQuery{}).
				WithPageSize(15).
				WithQueryBuilder(query.Match("name", "unknown")),
		)

		cursor, err := store.BankAccountsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 0)
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
	})

	t.Run("list bank accounts by country", func(t *testing.T) {
		q := NewListBankAccountsQuery(
			bunpaginate.NewPaginatedQueryOptions(BankAccountQuery{}).
				WithPageSize(15).
				WithQueryBuilder(query.Match("country", "US")),
		)

		cursor, err := store.BankAccountsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
		compareBankAccounts(t, d3, cursor.Data[0])
		compareBankAccounts(t, d1, cursor.Data[1])
	})

	t.Run("list bank accounts by country 2", func(t *testing.T) {
		q := NewListBankAccountsQuery(
			bunpaginate.NewPaginatedQueryOptions(BankAccountQuery{}).
				WithPageSize(15).
				WithQueryBuilder(query.Match("country", "DE")),
		)

		cursor, err := store.BankAccountsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
		compareBankAccounts(t, d2, cursor.Data[0])
	})

	t.Run("list bank accounts by unknown country", func(t *testing.T) {
		q := NewListBankAccountsQuery(
			bunpaginate.NewPaginatedQueryOptions(BankAccountQuery{}).
				WithPageSize(15).
				WithQueryBuilder(query.Match("country", "unknown")),
		)

		cursor, err := store.BankAccountsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 0)
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
	})

	t.Run("list bank accounts by metadata", func(t *testing.T) {
		q := NewListBankAccountsQuery(
			bunpaginate.NewPaginatedQueryOptions(BankAccountQuery{}).
				WithPageSize(15).
				WithQueryBuilder(query.Match("metadata[foo]", "bar")),
		)

		cursor, err := store.BankAccountsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
		compareBankAccounts(t, d1, cursor.Data[0])
	})

	t.Run("list bank accounts by unknown metadata", func(t *testing.T) {
		q := NewListBankAccountsQuery(
			bunpaginate.NewPaginatedQueryOptions(BankAccountQuery{}).
				WithPageSize(15).
				WithQueryBuilder(query.Match("metadata[unknown]", "bar")),
		)

		cursor, err := store.BankAccountsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 0)
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
	})

	t.Run("list bank accounts test cursor", func(t *testing.T) {
		q := NewListBankAccountsQuery(
			bunpaginate.NewPaginatedQueryOptions(BankAccountQuery{}).
				WithPageSize(1),
		)

		cursor, err := store.BankAccountsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.NotEmpty(t, cursor.Next)
		compareBankAccounts(t, d2, cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Next, &q)
		require.NoError(t, err)
		cursor, err = store.BankAccountsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		require.NotEmpty(t, cursor.Previous)
		require.NotEmpty(t, cursor.Next)
		compareBankAccounts(t, d3, cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Next, &q)
		require.NoError(t, err)
		cursor, err = store.BankAccountsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.False(t, cursor.HasMore)
		require.NotEmpty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
		compareBankAccounts(t, d1, cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Previous, &q)
		require.NoError(t, err)
		cursor, err = store.BankAccountsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		require.NotEmpty(t, cursor.Previous)
		require.NotEmpty(t, cursor.Next)
		compareBankAccounts(t, d3, cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Previous, &q)
		require.NoError(t, err)
		cursor, err = store.BankAccountsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.NotEmpty(t, cursor.Next)
		compareBankAccounts(t, d2, cursor.Data[0])
	})
}

func TestBankAccountsAddRelatedAccount(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertBankAccount(t, ctx, store, defaultBankAccount)
	upsertBankAccount(t, ctx, store, defaultBankAccount2)
	upsertBankAccount(t, ctx, store, defaultBankAccount3)

	t.Run("add related account when empty", func(t *testing.T) {
		acc := models.BankAccountRelatedAccount{
			BankAccountID: defaultBankAccount.ID,
			AccountID:     defaultAccounts[0].ID,
			ConnectorID:   defaultConnector.ID,
			CreatedAt:     now.UTC().Time,
		}

		ba := defaultBankAccount
		ba.RelatedAccounts = append(ba.RelatedAccounts, acc)

		require.NoError(t, store.BankAccountsAddRelatedAccount(ctx, acc))

		actual, err := store.BankAccountsGet(ctx, defaultBankAccount.ID, true)
		require.NoError(t, err)
		compareBankAccounts(t, ba, *actual)
	})

	t.Run("add related account when not empty", func(t *testing.T) {
		acc := models.BankAccountRelatedAccount{
			BankAccountID: defaultBankAccount2.ID,
			AccountID:     defaultAccounts[1].ID,
			ConnectorID:   defaultConnector.ID,
			CreatedAt:     now.UTC().Time,
		}

		ba := defaultBankAccount2
		ba.RelatedAccounts = append(ba.RelatedAccounts, acc)

		require.NoError(t, store.BankAccountsAddRelatedAccount(ctx, acc))

		actual, err := store.BankAccountsGet(ctx, defaultBankAccount2.ID, true)
		require.NoError(t, err)
		compareBankAccounts(t, ba, *actual)
	})

	t.Run("add related account with unknown bank account", func(t *testing.T) {
		acc := models.BankAccountRelatedAccount{
			BankAccountID: uuid.New(),
			AccountID:     defaultAccounts[1].ID,
			ConnectorID:   defaultConnector.ID,
			CreatedAt:     now.UTC().Time,
		}

		require.Error(t, store.BankAccountsAddRelatedAccount(ctx, acc))
	})

	t.Run("add related account with unknown account", func(t *testing.T) {
		acc := models.BankAccountRelatedAccount{
			BankAccountID: defaultBankAccount2.ID,
			AccountID: models.AccountID{
				Reference:   "unknown",
				ConnectorID: defaultConnector.ID,
			},
			ConnectorID: defaultConnector.ID,
			CreatedAt:   now.UTC().Time,
		}

		require.Error(t, store.BankAccountsAddRelatedAccount(ctx, acc))
	})

	t.Run("add related account with unknown connector", func(t *testing.T) {
		acc := models.BankAccountRelatedAccount{
			BankAccountID: defaultBankAccount2.ID,
			AccountID:     defaultAccounts[2].ID,
			ConnectorID: models.ConnectorID{
				Reference: uuid.New(),
				Provider:  "unknown",
			},
			CreatedAt: now.UTC().Time,
		}

		require.Error(t, store.BankAccountsAddRelatedAccount(ctx, acc))
	})

	t.Run("add related account with existing related account", func(t *testing.T) {
		acc := models.BankAccountRelatedAccount{
			BankAccountID: defaultBankAccount3.ID,
			AccountID:     defaultAccounts[0].ID,
			ConnectorID:   defaultConnector.ID,
			CreatedAt:     now.Add(-30 * time.Minute).UTC().Time,
		}

		ba := defaultBankAccount3
		ba.RelatedAccounts = append(ba.RelatedAccounts, acc)

		require.NoError(t, store.BankAccountsAddRelatedAccount(ctx, acc))

		actual, err := store.BankAccountsGet(ctx, defaultBankAccount3.ID, true)
		require.NoError(t, err)
		compareBankAccounts(t, ba, *actual)

		require.NoError(t, store.BankAccountsAddRelatedAccount(ctx, acc))

		actual, err = store.BankAccountsGet(ctx, defaultBankAccount3.ID, true)
		require.NoError(t, err)
		compareBankAccounts(t, ba, *actual)
	})
}

func TestBankAccountsDeleteRelatedAccountFromConnectorID(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertBankAccount(t, ctx, store, defaultBankAccount)
	upsertBankAccount(t, ctx, store, defaultBankAccount2)
	upsertBankAccount(t, ctx, store, defaultBankAccount3)

	t.Run("delete related account with unknown connector", func(t *testing.T) {
		require.NoError(t, store.BankAccountsDeleteRelatedAccountFromConnectorID(ctx, models.ConnectorID{
			Reference: uuid.New(),
			Provider:  "unknown",
		}))

		actual, err := store.BankAccountsGet(ctx, defaultBankAccount2.ID, true)
		require.NoError(t, err)
		compareBankAccounts(t, defaultBankAccount2, *actual)
	})

	t.Run("delete related account with another connector id", func(t *testing.T) {
		require.NoError(t, store.BankAccountsDeleteRelatedAccountFromConnectorID(ctx, defaultConnector2.ID))

		actual, err := store.BankAccountsGet(ctx, defaultBankAccount2.ID, true)
		require.NoError(t, err)
		compareBankAccounts(t, defaultBankAccount2, *actual)
	})

	t.Run("delete related account", func(t *testing.T) {
		require.NoError(t, store.BankAccountsDeleteRelatedAccountFromConnectorID(ctx, defaultConnector.ID))

		ba := defaultBankAccount2
		ba.RelatedAccounts = nil

		actual, err := store.BankAccountsGet(ctx, defaultBankAccount2.ID, true)
		require.NoError(t, err)
		compareBankAccounts(t, ba, *actual)
	})
}

func compareBankAccounts(t *testing.T, expected, actual models.BankAccount) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.CreatedAt, actual.CreatedAt)
	require.Equal(t, expected.Name, actual.Name)

	require.Equal(t, len(expected.Metadata), len(actual.Metadata))
	for k, v := range expected.Metadata {
		require.Equal(t, v, actual.Metadata[k])
	}
	for k, v := range actual.Metadata {
		require.Equal(t, v, expected.Metadata[k])
	}

	switch {
	case expected.AccountNumber != nil && actual.AccountNumber != nil:
		require.Equal(t, *expected.AccountNumber, *actual.AccountNumber)
	case expected.AccountNumber == nil && actual.AccountNumber == nil:
		// Nothing to do
	default:
		require.Fail(t, "AccountNumber mismatch")
	}

	switch {
	case expected.IBAN != nil && actual.IBAN != nil:
		require.Equal(t, *expected.IBAN, *actual.IBAN)
	case expected.IBAN == nil && actual.IBAN == nil:
		// Nothing to do
	default:
		require.Fail(t, "IBAN mismatch")
	}

	switch {
	case expected.SwiftBicCode != nil && actual.SwiftBicCode != nil:
		require.Equal(t, *expected.SwiftBicCode, *actual.SwiftBicCode)
	case expected.SwiftBicCode == nil && actual.SwiftBicCode == nil:
		// Nothing to do
	default:
		require.Fail(t, "SwiftBicCode mismatch")
	}

	switch {
	case expected.Country != nil && actual.Country != nil:
		require.Equal(t, *expected.Country, *actual.Country)
	case expected.Country == nil && actual.Country == nil:
		// Nothing to do
	default:
		require.Fail(t, "Country mismatch")
	}

	require.Equal(t, len(expected.RelatedAccounts), len(actual.RelatedAccounts))
	for i := range expected.RelatedAccounts {
		require.Equal(t, expected.RelatedAccounts[i], actual.RelatedAccounts[i])
	}
}
