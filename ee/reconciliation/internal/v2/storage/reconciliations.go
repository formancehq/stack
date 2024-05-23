package storage

import (
	"context"
	"fmt"
	"math/big"
	"time"

	storageerrors "github.com/formancehq/reconciliation/internal/utils/storage/errors"
	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func (s *Storage) CreateReconciliation(ctx context.Context, reconciliation *models.Reconciliation) error {
	_, err := s.db.NewInsert().
		Model(reconciliation).
		Exec(ctx)
	if err != nil {
		return storageerrors.E("failed to create reconciliation", err)
	}
	return nil
}

func (s *Storage) CreateAccountBasedReconciliation(ctx context.Context, accountBasedReconciliation *models.ReconciliationAccountBased) error {
	_, err := s.db.NewInsert().
		Model(accountBasedReconciliation).
		Exec(ctx)
	if err != nil {
		return storageerrors.E("failed to create account based reconciliation", err)
	}
	return nil
}

func (s *Storage) CreateTransactionBasedReconciliation(ctx context.Context, transactionBasedReconciliation *models.ReconciliationTransactionBased) (bool, error) {
	res, err := s.db.NewInsert().
		Model(transactionBasedReconciliation).
		On("CONFLICT (payment_id, transaction_id, policy_id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return false, storageerrors.E("failed to create transaction based reconciliation", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, storageerrors.E("failed to get rows affected", err)
	}

	return rowsAffected > 0, nil
}

func (s *Storage) GetReconciliation(ctx context.Context, id uuid.UUID) (*models.Reconciliation, error) {
	var reconciliation models.Reconciliation
	err := s.db.NewSelect().
		Model(&reconciliation).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, storageerrors.E("failed to get reconciliation", err)
	}
	return &reconciliation, nil
}

func (s *Storage) GetAccountBasedReconciliation(ctx context.Context, id uuid.UUID) (*models.ReconciliationAccountBased, error) {
	var accountBasedReconciliation models.ReconciliationAccountBased
	err := s.db.NewSelect().
		Model(&accountBasedReconciliation).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, storageerrors.E("failed to get account based reconciliation", err)
	}
	return &accountBasedReconciliation, nil
}

func (s *Storage) GetTransactionBasedReconciliationByPaymentID(ctx context.Context, policyID uuid.UUID, id string) ([]*models.ReconciliationTransactionBased, error) {
	var transactionBasedReconciliation []*models.ReconciliationTransactionBased
	err := s.db.NewSelect().
		Model(&transactionBasedReconciliation).
		Where("payment_id = ?", id).
		Where("policy_id = ?", policyID).
		Scan(ctx)
	if err != nil {
		return nil, storageerrors.E("failed to get transaction based reconciliation", err)
	}
	return transactionBasedReconciliation, nil
}

func (s *Storage) GetTransactionBasedReconciliationByTransactionID(ctx context.Context, policyID uuid.UUID, id *big.Int) ([]*models.ReconciliationTransactionBased, error) {
	var transactionBasedReconciliation []*models.ReconciliationTransactionBased
	err := s.db.NewSelect().
		Model(&transactionBasedReconciliation).
		Where("transaction_id = ?", id).
		Where("policy_id = ?", policyID).
		Scan(ctx)
	if err != nil {
		return nil, storageerrors.E("failed to get transaction based reconciliation", err)
	}
	return transactionBasedReconciliation, nil
}

func (s *Storage) buildReconciliationListQuery(selectQuery *bun.SelectQuery, q ListReconciliationsQuery, where string, args []any) *bun.SelectQuery {
	selectQuery = selectQuery.
		Order("created_at DESC")

	if where != "" {
		return selectQuery.Where(where, args...)
	}

	return selectQuery
}

func (s *Storage) ListReconciliations(ctx context.Context, q ListReconciliationsQuery) (*bunpaginate.Cursor[models.Reconciliation], error) {
	var (
		where string
		args  []any
		err   error
	)

	if q.Options.QueryBuilder != nil {
		where, args, err = s.reconciliationsQueryContext(q.Options.QueryBuilder, q)
		if err != nil {
			return nil, err
		}
	}

	return paginateWithOffset[bunpaginate.PaginatedQueryOptions[ReconciliationFilters], models.Reconciliation](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[ReconciliationFilters]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return s.buildReconciliationListQuery(query, q, where, args)
		},
	)
}

func (s *Storage) reconciliationsQueryContext(qb query.Builder, q ListReconciliationsQuery) (string, []any, error) {
	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "createdAt":
			var timestamp time.Time
			switch t := value.(type) {
			case string:
				var err error
				timestamp, err = time.Parse(time.RFC3339, t)
				if err != nil {
					return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "invalid 'createdAt' value")
				}
			case time.Time:
				timestamp = t
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'createdAt' column can only be used with string or time.Time")
			}

			return fmt.Sprintf("created_at %s ?", query.DefaultComparisonOperatorsMapping[operator]), []any{timestamp}, nil

		case key == "name":
			if operator != "$match" {
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'name' column can only be used with $match")
			}
			switch name := value.(type) {
			case string:
				return "name = ?", []any{name}, nil
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'name' column can only be used with string")
			}

		case key == "policyID":
			if operator != "$match" {
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'policyID' column can only be used with $match")
			}
			switch policyID := value.(type) {
			case string:
				return "policy_id = ?", []any{policyID}, nil
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'policyID' column can only be used with string")
			}

		case key == "policyType":
			if operator != "$match" {
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'policyType' column can only be used with $match")
			}
			switch policyType := value.(type) {
			case string:
				return "policy_type = ?", []any{policyType}, nil
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'policyType' column can only be used with string")
			}

		case key == "status":
			if operator != "$match" {
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'status' column can only be used with $match")
			}
			switch reconciliationStatus := value.(type) {
			case string:
				return "status = ?", []any{reconciliationStatus}, nil
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'status' column can only be used with string")
			}

		default:
			return "", nil, errors.Wrapf(storageerrors.ErrInvalidQuery, "unknown key '%s' when building query", key)
		}
	}))
}

type ReconciliationFilters struct{}

type ListReconciliationsQuery bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[ReconciliationFilters]]

func NewListReconciliationsQuery(opts bunpaginate.PaginatedQueryOptions[ReconciliationFilters]) ListReconciliationsQuery {
	return ListReconciliationsQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}

func (s *Storage) buildTransactionBasedReconciliationListQuery(selectQuery *bun.SelectQuery, q ListTransactionBasedReconciliationsQuery, where string, args []any) *bun.SelectQuery {
	selectQuery = selectQuery.
		Order("created_at DESC")

	if where != "" {
		return selectQuery.Where(where, args...)
	}

	return selectQuery
}

func (s *Storage) ListTransactionBasedReconciliations(ctx context.Context, q ListTransactionBasedReconciliationsQuery) (*bunpaginate.Cursor[models.ReconciliationTransactionBased], error) {
	var (
		where string
		args  []any
		err   error
	)

	if q.Options.QueryBuilder != nil {
		where, args, err = s.transactionBasedReconciliationsQueryContext(q.Options.QueryBuilder, q)
		if err != nil {
			return nil, err
		}
	}

	return paginateWithOffset[bunpaginate.PaginatedQueryOptions[TransactionBasedReconciliationFilters], models.ReconciliationTransactionBased](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[TransactionBasedReconciliationFilters]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return s.buildTransactionBasedReconciliationListQuery(query, q, where, args)
		},
	)
}

func (s *Storage) transactionBasedReconciliationsQueryContext(qb query.Builder, q ListTransactionBasedReconciliationsQuery) (string, []any, error) {
	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "createdAt":
			var timestamp time.Time
			switch t := value.(type) {
			case string:
				var err error
				timestamp, err = time.Parse(time.RFC3339, t)
				if err != nil {
					return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "invalid 'createdAt' value")
				}
			case time.Time:
				timestamp = t
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'createdAt' column can only be used with string or time.Time")
			}

			return fmt.Sprintf("created_at %s ?", query.DefaultComparisonOperatorsMapping[operator]), []any{timestamp}, nil

		case key == "ruleID":
			if operator != "$match" {
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'ruleID' column can only be used with $match")
			}
			switch v := value.(type) {
			case string:
				ruleID, err := uuid.Parse(v)
				if err != nil {
					return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "invalid 'ruleID' value")
				}
				return "rule_id = ?", []any{ruleID}, nil
			case uuid.UUID:
				return "rule_id = ?", []any{v}, nil
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'ruleID' column can only be used with string")
			}

		case key == "policyID":
			if operator != "$match" {
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'policyID' column can only be used with $match")
			}
			switch policyID := value.(type) {
			case string:
				return "policy_id = ?", []any{policyID}, nil
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'policyID' column can only be used with string")
			}

		case key == "status":
			if operator != "$match" {
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'status' column can only be used with $match")
			}
			switch status := value.(type) {
			case string:
				return "status = ?", []any{status}, nil
			default:
				return "", nil, errors.Wrap(storageerrors.ErrInvalidQuery, "'status' column can only be used with string")
			}

		default:
			return "", nil, errors.Wrapf(storageerrors.ErrInvalidQuery, "unknown key '%s' when building query", key)
		}
	}))
}

type TransactionBasedReconciliationFilters struct{}

type ListTransactionBasedReconciliationsQuery bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[TransactionBasedReconciliationFilters]]

func NewListTransactionBasedReconciliationsQuery(opts bunpaginate.PaginatedQueryOptions[TransactionBasedReconciliationFilters]) ListTransactionBasedReconciliationsQuery {
	return ListTransactionBasedReconciliationsQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}
