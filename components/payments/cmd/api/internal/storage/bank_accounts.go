package storage

import (
	"context"
	"sort"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

func (s *Storage) ListBankAccounts(ctx context.Context, pagination PaginatorQuery) ([]*models.BankAccount, PaginationDetails, error) {
	var bankAccounts []*models.BankAccount

	query := s.db.NewSelect().
		Model(&bankAccounts).
		Relation("Adjustments")

	query = pagination.apply(query, "bank_account.created_at")

	err := query.Scan(ctx)
	if err != nil {
		return nil, PaginationDetails{}, e("failed to list bank accounts", err)
	}

	var (
		hasMore                       = len(bankAccounts) > pagination.pageSize
		hasPrevious                   bool
		firstReference, lastReference string
	)

	if hasMore {
		if pagination.cursor.Next || pagination.cursor.Reference == "" {
			bankAccounts = bankAccounts[:pagination.pageSize]
		} else {
			bankAccounts = bankAccounts[1:]
		}
	}

	sort.Slice(bankAccounts, func(i, j int) bool {
		return bankAccounts[i].CreatedAt.After(bankAccounts[j].CreatedAt)
	})

	if len(bankAccounts) > 0 {
		firstReference = bankAccounts[0].CreatedAt.Format(time.RFC3339Nano)
		lastReference = bankAccounts[len(bankAccounts)-1].CreatedAt.Format(time.RFC3339Nano)

		query = s.db.NewSelect().Model(&bankAccounts).Relation("Adjustments")

		hasPrevious, err = pagination.hasPrevious(ctx, query, "bank_account.created_at", firstReference)
		if err != nil {
			return nil, PaginationDetails{}, e("failed to check if there is a previous page", err)
		}
	}

	paginationDetails, err := pagination.paginationDetails(hasMore, hasPrevious, firstReference, lastReference)
	if err != nil {
		return nil, PaginationDetails{}, e("failed to get pagination details", err)
	}

	return bankAccounts, paginationDetails, nil
}

func (s *Storage) GetBankAccount(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error) {
	var account models.BankAccount
	query := s.db.NewSelect().
		Model(&account).
		Column("id", "created_at", "name", "created_at", "country", "metadata").
		Relation("Adjustments")

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
