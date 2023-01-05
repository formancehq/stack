package wallet

import (
	"context"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/go-libs/metadata"
	"github.com/formancehq/wallets/pkg/core"
	"github.com/pkg/errors"
)

type ListResponse[T any] struct {
	Data           []T
	Next, Previous string
	HasMore        bool
}

type ListQuery[T any] struct {
	Payload         T
	Limit           int
	PaginationToken string
}

func newListResponse[SRC any, DST any](cursor interface {
	GetData() []SRC
	GetNext() string
	GetPrevious() string
	GetHasMore() bool
}, mapper func(src SRC) DST,
) *ListResponse[DST] {
	ret := make([]DST, 0)
	for _, item := range cursor.GetData() {
		ret = append(ret, mapper(item))
	}

	return &ListResponse[DST]{
		Data:     ret,
		Next:     cursor.GetNext(),
		Previous: cursor.GetPrevious(),
		HasMore:  cursor.GetHasMore(),
	}
}

type ListHolds struct {
	WalletID string
	Metadata map[string]any
}

type ListTransactions struct {
	WalletID string
}

type ListWallets struct {
	Metadata metadata.Metadata
	Name     string
}

type Repository struct {
	ledgerName string
	chart      *core.Chart
	client     Ledger
}

func NewRepository(
	ledgerName string,
	client Ledger,
	chart *core.Chart,
) *Repository {
	return &Repository{
		ledgerName: ledgerName,
		chart:      chart,
		client:     client,
	}
}

type Data struct {
	Metadata metadata.Metadata `json:"metadata"`
	Name     string            `json:"name"`
}

func (r *Repository) CreateWallet(ctx context.Context, data *Data) (*core.Wallet, error) {
	wallet := core.NewWallet(data.Name, data.Metadata)

	if err := r.client.AddMetadataToAccount(
		ctx,
		r.ledgerName,
		r.chart.GetMainAccount(wallet.ID),
		wallet.LedgerMetadata(),
	); err != nil {
		return nil, errors.Wrap(err, "adding metadata to account")
	}

	return &wallet, nil
}

func (r *Repository) UpdateWallet(ctx context.Context, id string, data *Data) error {
	account, err := r.client.GetAccount(ctx, r.ledgerName, r.chart.GetMainAccount(id))
	if err != nil {
		return ErrWalletNotFound
	}

	if !core.IsPrimary(account) {
		return ErrWalletNotFound
	}

	newCustomMetadata := metadata.Metadata{}
	existingCustomMetadata := core.GetMetadata(account, core.MetadataKeyWalletCustomData)
	if existingCustomMetadata != nil {
		newCustomMetadata = newCustomMetadata.Merge(existingCustomMetadata.(map[string]any))
	}
	newCustomMetadata = newCustomMetadata.Merge(data.Metadata)

	meta := account.GetMetadata()
	meta[core.MetadataKeyWalletCustomData] = newCustomMetadata

	if err := r.client.AddMetadataToAccount(ctx, r.ledgerName, r.chart.GetMainAccount(id), meta); err != nil {
		return errors.Wrap(err, "adding metadata to account")
	}

	return nil
}

func (r *Repository) ListWallets(ctx context.Context, query ListQuery[ListWallets]) (*ListResponse[core.Wallet], error) {
	var (
		response *sdk.ListAccounts200ResponseCursor
		err      error
	)
	if query.PaginationToken == "" {
		metadata := map[string]interface{}{
			core.MetadataKeyWalletSpecType: core.PrimaryWallet,
		}
		if query.Payload.Metadata != nil && len(query.Payload.Metadata) > 0 {
			for k, v := range query.Payload.Metadata {
				metadata[core.MetadataKeyWalletCustomData+"."+k] = v
			}
		}
		if query.Payload.Name != "" {
			metadata[core.MetadataKeyWalletName] = query.Payload.Name
		}
		response, err = r.client.ListAccounts(ctx, r.ledgerName, ListAccountsQuery{
			Limit:    query.Limit,
			Metadata: metadata,
		})
	} else {
		response, err = r.client.ListAccounts(ctx, r.ledgerName, ListAccountsQuery{
			PaginationToken: query.PaginationToken,
		})
	}
	if err != nil {
		return nil, err
	}

	return newListResponse[sdk.Account, core.Wallet](response, func(account sdk.Account) core.Wallet {
		return core.WalletFromAccount(&account)
	}), nil
}

func (r *Repository) GetWallet(ctx context.Context, id string) (*core.WalletWithBalances, error) {
	account, err := r.client.GetAccount(
		ctx,
		r.ledgerName,
		r.chart.GetMainAccount(id),
	)
	if err != nil {
		return nil, errors.Wrap(err, "getting account")
	}

	if !core.IsPrimary(account) {
		return nil, ErrWalletNotFound
	}

	w := core.WalletWithBalancesFromAccount(account)

	return &w, nil
}

func (r *Repository) ListHolds(ctx context.Context, query ListQuery[ListHolds]) (*ListResponse[core.DebitHold], error) {
	var (
		response *sdk.ListAccounts200ResponseCursor
		err      error
	)
	if query.PaginationToken == "" {
		metadata := metadata.Metadata{
			core.MetadataKeyWalletSpecType: core.HoldWallet,
		}
		if query.Payload.WalletID != "" {
			metadata[core.MetadataKeyHoldWalletID] = query.Payload.WalletID
		}
		if query.Payload.Metadata != nil && len(query.Payload.Metadata) > 0 {
			for k, v := range query.Payload.Metadata {
				metadata[core.MetadataKeyWalletCustomData+"."+k] = v
			}
		}
		response, err = r.client.ListAccounts(ctx, r.ledgerName, ListAccountsQuery{
			Limit:    query.Limit,
			Metadata: metadata,
		})
	} else {
		response, err = r.client.ListAccounts(ctx, r.ledgerName, ListAccountsQuery{
			PaginationToken: query.PaginationToken,
		})
	}
	if err != nil {
		return nil, errors.Wrap(err, "listing accounts")
	}

	return newListResponse[sdk.Account, core.DebitHold](response, func(account sdk.Account) core.DebitHold {
		return core.DebitHoldFromLedgerAccount(&account)
	}), nil
}

func (r *Repository) ListTransactions(ctx context.Context, query ListQuery[ListTransactions]) (*ListResponse[sdk.Transaction], error) {
	var (
		response *sdk.ListTransactions200ResponseCursor
		err      error
	)
	if query.PaginationToken == "" {
		response, err = r.client.ListTransactions(ctx, r.ledgerName, ListTransactionsQuery{
			Limit: query.Limit,
			Account: func() string {
				if query.Payload.WalletID != "" {
					return r.chart.GetMainAccount(query.Payload.WalletID)
				}
				return ""
			}(),
			Metadata: core.WalletTransactionBaseMetadataFilter(),
		})
	} else {
		response, err = r.client.ListTransactions(ctx, r.ledgerName, ListTransactionsQuery{
			PaginationToken: query.PaginationToken,
		})
	}
	if err != nil {
		return nil, errors.Wrap(err, "listing transactions")
	}

	return newListResponse[sdk.Transaction, sdk.Transaction](response, func(tx sdk.Transaction) sdk.Transaction {
		return tx
	}), nil
}

func (r *Repository) GetHold(ctx context.Context, id string) (*core.ExpandedDebitHold, error) {
	account, err := r.client.GetAccount(ctx, r.ledgerName, r.chart.GetHoldAccount(id))
	if err != nil {
		return nil, err
	}

	hold := core.ExpandedDebitHoldFromLedgerAccount(account)

	return &hold, nil
}
