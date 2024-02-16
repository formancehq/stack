package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type BankAccountQuery struct{}

type ListBankAccountQuery bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[BankAccountQuery]]

func NewListBankAccountQuery(opts PaginatedQueryOptions[BankAccountQuery]) ListBankAccountQuery {
	return ListBankAccountQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}

func (s *Storage) ListBankAccounts(ctx context.Context, q ListBankAccountQuery) (*api.Cursor[models.BankAccount], error) {
	return PaginateWithOffset[PaginatedQueryOptions[BankAccountQuery], models.BankAccount](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[BankAccountQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.
				Relation("RelatedAccounts").
				Order("created_at DESC")

			if q.Options.Sorter != nil {
				query = q.Options.Sorter.Apply(query)
			}

			return query
		},
	)
}

func (s *Storage) GetBankAccount(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error) {
	var account models.BankAccount
	query := s.db.NewSelect().
		Model(&account).
		Column("id", "created_at", "name", "created_at", "country", "metadata").
		Relation("RelatedAccounts")

	if expand {
		query = query.ColumnExpr("pgp_sym_decrypt(account_number, ?, ?) AS decrypted_account_number", s.configEncryptionKey, encryptionOptions).
			ColumnExpr("pgp_sym_decrypt(iban, ?, ?) AS decrypted_iban", s.configEncryptionKey, encryptionOptions).
			ColumnExpr("pgp_sym_decrypt(swift_bic_code, ?, ?) AS decrypted_swift_bic_code", s.configEncryptionKey, encryptionOptions)
	}

	err := query.
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, e("get bank account", err)
	}

	return &account, nil
}
