package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type bankAccount struct {
	bun.BaseModel `bun:"table:bank_accounts"`

	// Mandatory fields
	ID        uuid.UUID `bun:"id,pk,type:uuid,notnull"`
	CreatedAt time.Time `bun:"created_at,type:timestamp without time zone,notnull"`
	Name      string    `bun:"name,type:text,notnull"`

	// Field encrypted
	AccountNumber string `bun:"decrypted_account_number,scanonly"`
	IBAN          string `bun:"decrypted_iban,scanonly"`
	SwiftBicCode  string `bun:"decrypted_swift_bic_code,scanonly"`

	// Optional fields
	// c.f.: https://bun.uptrace.dev/guide/models.html#nulls
	Country *string `bun:"country,type:text,nullzero"`

	// Optional fields with default
	// c.f. https://bun.uptrace.dev/guide/models.html#default
	Metadata map[string]string `bun:"metadata,type:jsonb,nullzero,notnull,default:'{}'"`

	RelatedAccounts []*bankAccountRelatedAccount `bun:"rel:has-many,join:id=bank_account_id"`
}

func (s *store) BankAccountsUpsert(ctx context.Context, bankAccount models.BankAccount) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return e("begin transaction", err)
	}
	defer tx.Rollback()

	toInsert := fromBankAccountModels(bankAccount)
	// Insert or update the bank account
	_, err = tx.NewInsert().
		Model(&toInsert).
		On("CONFLICT (id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return e("insert bank account", err)
	}

	// Insert or update the related accounts
	_, err = tx.NewInsert().
		Model(&toInsert.RelatedAccounts).
		On("CONFLICT (bank_account_id, account_id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return e("insert related accounts", err)
	}

	return e("commit transaction", tx.Commit())
}

func (s *store) BankAccountsUpdateMetadata(ctx context.Context, id uuid.UUID, metadata map[string]string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return e("update bank account metadata", err)
	}
	defer tx.Rollback()

	var account bankAccount
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

	_, err = tx.NewUpdate().
		Model(&account).
		Column("metadata").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return e("update bank account metadata", err)
	}

	return e("commit transaction", tx.Commit())
}

func (s *store) BankAccountsGet(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error) {
	var account bankAccount
	query := s.db.NewSelect().
		Model(&account).
		Relation("RelatedAccounts")
	if !expand {
		query = query.Column("id", "created_at", "name", "country", "metadata")
	}
	err := query.Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, e("get bank account", err)
	}

	return pointer.For(toBankAccountModels(account)), nil
}

type BankAccountQuery struct{}

type ListBankAccountsQuery bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[BankAccountQuery]]

func NewListBankAccountsQuery(opts bunpaginate.PaginatedQueryOptions[BankAccountQuery]) ListBankAccountsQuery {
	return ListBankAccountsQuery{
		Order:    bunpaginate.OrderAsc,
		PageSize: opts.PageSize,
		Options:  opts,
	}
}

func (s *store) bankAccountsQueryContext(qb query.Builder) (string, []any, error) {
	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "name", key == "country":
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrValidation, fmt.Sprintf("'%s' column can only be used with $match", key))
			}
			return fmt.Sprintf("%s = ?", key), []any{value}, nil
		case metadataRegex.Match([]byte(key)):
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrValidation, "'metadata' column can only be used with $match")
			}
			match := metadataRegex.FindAllStringSubmatch(key, 3)

			key := "metadata"
			return key + " @> ?", []any{map[string]any{
				match[0][1]: value,
			}}, nil
		default:
			return "", nil, errors.Wrap(ErrValidation, fmt.Sprintf("unknown key '%s' when building query", key))
		}
	}))
}

func (s *store) BankAccountsList(ctx context.Context, q ListBankAccountsQuery) (*bunpaginate.Cursor[models.BankAccount], error) {
	var (
		where string
		args  []any
		err   error
	)
	if q.Options.QueryBuilder != nil {
		where, args, err = s.bankAccountsQueryContext(q.Options.QueryBuilder)
		if err != nil {
			return nil, err
		}
	}

	cursor, err := paginateWithOffset[bunpaginate.PaginatedQueryOptions[BankAccountQuery], bankAccount](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[BankAccountQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.Relation("RelatedAccounts")
			if where != "" {
				query = query.Where(where, args...)
			}

			query = query.Order("created_at DESC")

			return query
		},
	)
	if err != nil {
		return nil, e("failed to fetch accounts", err)
	}

	bankAccounts := make([]models.BankAccount, 0, len(cursor.Data))
	for _, a := range cursor.Data {
		bankAccounts = append(bankAccounts, toBankAccountModels(a))
	}

	return &bunpaginate.Cursor[models.BankAccount]{
		PageSize: cursor.PageSize,
		HasMore:  cursor.HasMore,
		Previous: cursor.Previous,
		Next:     cursor.Next,
		Data:     bankAccounts,
	}, nil
}

type bankAccountRelatedAccount struct {
	bun.BaseModel `bun:"table:bank_accounts_related_accounts"`

	// Mandatory fields
	BankAccountID uuid.UUID          `bun:"bank_account_id,pk,type:uuid,notnull"`
	AccountID     models.AccountID   `bun:"account_id,pk,type:character varying,notnull"`
	ConnectorID   models.ConnectorID `bun:"connector_id,type:character varying,notnull"`
	CreatedAt     time.Time          `bun:"created_at,type:timestamp without time zone,notnull"`
}

func (s *store) BankAccountsAddRelatedAccount(ctx context.Context, relatedAccount models.BankAccountRelatedAccount) error {
	toInsert := fromBankAccountRelatedAccountModels(relatedAccount)

	_, err := s.db.NewInsert().
		Model(&toInsert).
		On("CONFLICT (bank_account_id, account_id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return e("add bank account related account", err)
	}

	return nil
}

func (s *store) BankAccountsDeleteRelatedAccountFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error {
	_, err := s.db.NewDelete().
		Model((*bankAccountRelatedAccount)(nil)).
		Where("connector_id = ?", connectorID).
		Exec(ctx)
	if err != nil {
		return e("delete bank account related account", err)
	}

	return nil
}

func fromBankAccountModels(from models.BankAccount) bankAccount {
	ba := bankAccount{
		ID:        from.ID,
		CreatedAt: from.CreatedAt,
		Name:      from.Name,
		Country:   from.Country,
		Metadata:  from.Metadata,
	}

	if from.AccountNumber != nil {
		ba.AccountNumber = *from.AccountNumber
	}

	if from.IBAN != nil {
		ba.IBAN = *from.IBAN
	}

	if from.SwiftBicCode != nil {
		ba.SwiftBicCode = *from.SwiftBicCode
	}

	relatedAccounts := make([]*bankAccountRelatedAccount, 0, len(from.RelatedAccounts))
	for _, ra := range from.RelatedAccounts {
		relatedAccounts = append(relatedAccounts, pointer.For(fromBankAccountRelatedAccountModels(ra)))
	}
	ba.RelatedAccounts = relatedAccounts

	return ba
}

func toBankAccountModels(from bankAccount) models.BankAccount {
	ba := models.BankAccount{
		ID:        from.ID,
		CreatedAt: from.CreatedAt,
		Name:      from.Name,
		Country:   from.Country,
		Metadata:  from.Metadata,
	}

	if from.AccountNumber != "" {
		ba.AccountNumber = &from.AccountNumber
	}

	if from.IBAN != "" {
		ba.IBAN = &from.IBAN
	}

	if from.SwiftBicCode != "" {
		ba.SwiftBicCode = &from.SwiftBicCode
	}

	relatedAccounts := make([]models.BankAccountRelatedAccount, 0, len(from.RelatedAccounts))
	for _, ra := range from.RelatedAccounts {
		relatedAccounts = append(relatedAccounts, toBankAccountRelatedAccountModels(*ra))
	}
	ba.RelatedAccounts = relatedAccounts

	return ba
}

func fromBankAccountRelatedAccountModels(from models.BankAccountRelatedAccount) bankAccountRelatedAccount {
	return bankAccountRelatedAccount{
		BankAccountID: from.BankAccountID,
		AccountID:     from.AccountID,
		ConnectorID:   from.ConnectorID,
		CreatedAt:     from.CreatedAt,
	}
}

func toBankAccountRelatedAccountModels(from bankAccountRelatedAccount) models.BankAccountRelatedAccount {
	return models.BankAccountRelatedAccount{
		BankAccountID: from.BankAccountID,
		AccountID:     from.AccountID,
		ConnectorID:   from.ConnectorID,
		CreatedAt:     from.CreatedAt,
	}
}
