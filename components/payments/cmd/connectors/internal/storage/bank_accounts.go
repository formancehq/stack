package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

func (s *Storage) CreateBankAccount(ctx context.Context, bankAccount *models.BankAccount) error {
	account := models.BankAccount{
		CreatedAt: bankAccount.CreatedAt,
		Country:   bankAccount.Country,
		Name:      bankAccount.Name,
		Metadata:  bankAccount.Metadata,
	}

	var id uuid.UUID
	err := s.db.NewInsert().Model(&account).Returning("id").Scan(ctx, &id)
	if err != nil {
		return e("install connector", err)
	}
	bankAccount.ID = id

	return s.updateBankAccountInformation(ctx, id, bankAccount.AccountNumber, bankAccount.IBAN, bankAccount.SwiftBicCode)
}

func (s *Storage) AddBankAccountRelatedAccount(ctx context.Context, relatedAccount *models.BankAccountRelatedAccount) error {
	_, err := s.db.NewInsert().Model(relatedAccount).Exec(ctx)
	if err != nil {
		return e("add bank account related account", err)
	}

	return nil
}

func (s *Storage) updateBankAccountInformation(ctx context.Context, id uuid.UUID, accountNumber, iban, swiftBicCode string) error {
	_, err := s.db.NewUpdate().
		Model(&models.BankAccount{}).
		Set("account_number = pgp_sym_encrypt(?::TEXT, ?, ?)", accountNumber, s.configEncryptionKey, encryptionOptions).
		Set("iban = pgp_sym_encrypt(?::TEXT, ?, ?)", iban, s.configEncryptionKey, encryptionOptions).
		Set("swift_bic_code = pgp_sym_encrypt(?::TEXT, ?, ?)", swiftBicCode, s.configEncryptionKey, encryptionOptions).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return e("update bank account information", err)
	}

	return nil
}

func (s *Storage) UpdateBankAccountMetadata(ctx context.Context, id uuid.UUID, metadata map[string]string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return e("update bank account metadata", err)
	}
	defer tx.Rollback()

	var account models.BankAccount
	err = tx.NewSelect().
		Model(&account).
		Column("id", "metadata").
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return e("update bank account metadata", err)
	}

	if account.Metadata == nil {
		account.Metadata = make(map[string]string)
	}

	for k, v := range metadata {
		account.Metadata[k] = v
	}

	_, err = s.db.NewUpdate().
		Model(&account).
		Column("metadata").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return e("update bank account metadata", err)
	}

	return e("failed to commit transaction", tx.Commit())
}

func (s *Storage) LinkBankAccountWithAccount(ctx context.Context, id uuid.UUID, accountID *models.AccountID) error {
	relatedAccount := &models.BankAccountRelatedAccount{
		ID:            uuid.New(),
		BankAccountID: id,
		ConnectorID:   accountID.ConnectorID,
		AccountID:     *accountID,
	}

	return s.AddBankAccountRelatedAccount(ctx, relatedAccount)
}

func (s *Storage) GetBankAccount(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error) {
	var account models.BankAccount
	query := s.db.NewSelect().
		Model(&account).
		Relation("RelatedAccounts").
		Column("id", "created_at", "name", "created_at", "country", "metadata")

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

func (s *Storage) GetBankAccountRelatedAccounts(ctx context.Context, id uuid.UUID) ([]*models.BankAccountRelatedAccount, error) {
	var relatedAccounts []*models.BankAccountRelatedAccount
	err := s.db.NewSelect().
		Model(&relatedAccounts).
		Where("bank_account_id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, e("get bank account related accounts", err)
	}

	return relatedAccounts, nil
}
