package wallet

import (
	"context"
	"math/big"
	"sort"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/pkg/errors"
)

type ListResponse[T any] struct {
	Data           []T
	Next, Previous string
	HasMore        bool
}

type Pagination struct {
	Limit           int
	PaginationToken string
}

type ListQuery[T any] struct {
	Pagination
	Payload T
}

type mapper[SRC any, DST any] func(src SRC) DST

func newListResponse[SRC any, DST any](cursor interface {
	GetData() []SRC
	GetNext() string
	GetPrevious() string
	GetHasMore() bool
}, mapper mapper[SRC, DST],
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
	Metadata metadata.Metadata
}

type ListBalances struct {
	WalletID string
	Metadata metadata.Metadata
}

type ListTransactions struct {
	WalletID string
}

func BalancesMetadataFilter(walletID string) metadata.Metadata {
	return metadata.Metadata{
		MetadataKeyWalletBalance: TrueValue,
		MetadataKeyWalletID:      walletID,
	}
}

type Manager struct {
	client     Ledger
	chart      *Chart
	ledgerName string
}

func NewManager(
	ledgerName string,
	client Ledger,
	chart *Chart,
) *Manager {
	return &Manager{
		client:     client,
		chart:      chart,
		ledgerName: ledgerName,
	}
}

//nolint:cyclop
func (m *Manager) Debit(ctx context.Context, debit Debit) (*DebitHold, error) {
	if err := debit.Validate(); err != nil {
		return nil, err
	}

	dest := debit.getDestination()

	var hold *DebitHold
	if debit.Pending {
		hold = Ptr(debit.newHold())
		holdAccount := m.chart.GetHoldAccount(hold.ID)
		if err := m.client.AddMetadataToAccount(ctx, m.ledgerName, holdAccount, hold.LedgerMetadata(m.chart)); err != nil {
			return nil, errors.Wrap(err, "adding metadata to account")
		}

		dest = NewLedgerAccountSubject(holdAccount)
	}

	sources := make([]string, 0)
	switch {
	case len(debit.Balances) == 0:
		sources = append(sources, m.chart.GetMainBalanceAccount(debit.WalletID))
	case len(debit.Balances) == 1 && debit.Balances[0] == "*":
		balancesRaw, err := fetchAndMapAllAccounts[Balance](ctx, m, BalancesMetadataFilter(debit.WalletID), BalanceFromAccount)
		if err != nil {
			return nil, err
		}
		balances := Balances(balancesRaw)
		sort.Stable(balances)

		// Filter expired and generate sources
		for _, balance := range balances {
			if balance.ExpiresAt != nil && balance.ExpiresAt.Before(time.Now()) {
				continue
			}
			sources = append(sources, m.chart.GetBalanceAccount(debit.WalletID, balance.Name))
		}
	default:
		for _, balance := range debit.Balances {
			if balance == "*" {
				return nil, ErrInvalidBalanceSpecified
			}
			sources = append(sources, m.chart.GetBalanceAccount(debit.WalletID, balance))
		}
	}

	postTransaction := PostTransaction{
		Script: &PostTransactionScript{
			Plain: BuildDebitWalletScript(sources...),
			Vars: map[string]interface{}{
				"destination": dest.getAccount(m.chart),
				"amount": map[string]any{
					// @todo: upgrade this to proper int after sdk is updated
					"amount": debit.Amount.Amount.Uint64(),
					"asset":  debit.Amount.Asset,
				},
			},
		},
		Metadata: TransactionMetadata(debit.Metadata),
		//nolint:godox
		// TODO: Add set account metadata for hold when released on ledger (v1.9)
	}

	if debit.Reference != "" {
		postTransaction.Reference = &debit.Reference
	}

	if err := m.CreateTransaction(ctx, postTransaction); err != nil {
		return nil, err
	}

	return hold, nil
}

func (m *Manager) ConfirmHold(ctx context.Context, debit ConfirmHold) error {
	account, err := m.client.GetAccount(ctx, m.ledgerName, m.chart.GetHoldAccount(debit.HoldID))
	if err != nil {
		return errors.Wrap(err, "getting account")
	}
	if !IsHold(account) {
		return ErrHoldNotFound
	}

	hold := ExpandedDebitHoldFromLedgerAccount(*account)
	if hold.Remaining.Uint64() == 0 {
		return ErrClosedHold
	}

	amount, err := debit.resolveAmount(hold)
	if err != nil {
		return err
	}

	postTransaction := PostTransaction{
		Script: &PostTransactionScript{
			Plain: BuildConfirmHoldScript(debit.Final, hold.Asset),
			Vars: map[string]interface{}{
				"hold": m.chart.GetHoldAccount(debit.HoldID),
				"amount": map[string]any{
					"amount": amount,
					"asset":  hold.Asset,
				},
			},
		},
		Metadata: TransactionMetadata(metadata.Metadata{}),
	}

	if err := m.CreateTransaction(ctx, postTransaction); err != nil {
		return err
	}

	return nil
}

func (m *Manager) VoidHold(ctx context.Context, void VoidHold) error {
	account, err := m.client.GetAccount(ctx, m.ledgerName, m.chart.GetHoldAccount(void.HoldID))
	if err != nil {
		return errors.Wrap(err, "getting account")
	}

	hold := ExpandedDebitHoldFromLedgerAccount(*account)
	if hold.IsClosed() {
		return ErrClosedHold
	}

	postTransaction := PostTransaction{
		Script: &PostTransactionScript{
			Plain: BuildCancelHoldScript(hold.Asset),
			Vars: map[string]interface{}{
				"hold": m.chart.GetHoldAccount(void.HoldID),
			},
		},
		Metadata: TransactionMetadata(metadata.Metadata{}),
	}

	if err := m.CreateTransaction(ctx, postTransaction); err != nil {
		return err
	}

	return nil
}

func (m *Manager) Credit(ctx context.Context, credit Credit) error {
	if err := credit.Validate(); err != nil {
		return err
	}

	if credit.Balance != "" {
		if _, err := m.GetBalance(ctx, credit.WalletID, credit.Balance); err != nil {
			return err
		}
	}

	postTransaction := PostTransaction{
		Script: &PostTransactionScript{
			Plain: BuildCreditWalletScript(credit.Sources.ResolveAccounts(m.chart)...),
			Vars: map[string]interface{}{
				"destination": credit.destinationAccount(m.chart),
				"amount": map[string]any{
					// @todo: upgrade this to proper int after sdk is updated
					"amount": credit.Amount.Amount.Uint64(),
					"asset":  credit.Amount.Asset,
				},
			},
		},
		Metadata: TransactionMetadata(credit.Metadata),
	}
	if credit.Reference != "" {
		postTransaction.Reference = &credit.Reference
	}

	if err := m.CreateTransaction(ctx, postTransaction); err != nil {
		return err
	}

	return nil
}

func (m *Manager) CreateTransaction(ctx context.Context, postTransaction PostTransaction) error {
	if _, err := m.client.CreateTransaction(ctx, m.ledgerName, postTransaction); err != nil {
		apiErr, ok := err.(GenericOpenAPIError)
		if ok {
			respErr, ok := apiErr.Model().(shared.ErrorResponse)
			if ok {
				switch respErr.ErrorCode {
				case shared.ErrorsEnumInsufficientFund:
					return ErrInsufficientFundError
				}
			}
		}

		return errors.Wrap(err, "creating transaction")
	}

	return nil
}

func (m *Manager) ListWallets(ctx context.Context, query ListQuery[ListWallets]) (*ListResponse[Wallet], error) {
	return mapAccountList(ctx, m, mapAccountListQuery{
		Pagination: query.Pagination,
		Metadata: func() metadata.Metadata {
			metadata := metadata.Metadata{
				MetadataKeyWalletSpecType: PrimaryWallet,
			}
			if query.Payload.Metadata != nil && len(query.Payload.Metadata) > 0 {
				for k, v := range query.Payload.Metadata {
					metadata[MetadataKeyWalletCustomDataPrefix+k] = v
				}
			}
			if query.Payload.Name != "" {
				metadata[MetadataKeyWalletName] = query.Payload.Name
			}
			return metadata
		},
	}, func(account Account) Wallet {
		return FromAccount(m.ledgerName, account)
	})
}

func (m *Manager) ListHolds(ctx context.Context, query ListQuery[ListHolds]) (*ListResponse[DebitHold], error) {
	return mapAccountList(ctx, m, mapAccountListQuery{
		Pagination: query.Pagination,
		Metadata: func() metadata.Metadata {
			metadata := metadata.Metadata{
				MetadataKeyWalletSpecType: HoldWallet,
			}
			if query.Payload.WalletID != "" {
				metadata[MetadataKeyHoldWalletID] = query.Payload.WalletID
			}
			if query.Payload.Metadata != nil && len(query.Payload.Metadata) > 0 {
				for k, v := range query.Payload.Metadata {
					metadata[MetadataKeyWalletCustomDataPrefix+k] = v
				}
			}
			return metadata
		},
	}, DebitHoldFromLedgerAccount)
}

func (m *Manager) ListBalances(ctx context.Context, query ListQuery[ListBalances]) (*ListResponse[Balance], error) {
	return mapAccountList(ctx, m, mapAccountListQuery{
		Metadata: func() metadata.Metadata {
			metadata := BalancesMetadataFilter(query.Payload.WalletID)
			if query.Payload.Metadata != nil && len(query.Payload.Metadata) > 0 {
				for k, v := range query.Payload.Metadata {
					metadata[MetadataKeyWalletCustomDataPrefix+k] = v
				}
			}
			return metadata
		},
		Pagination: query.Pagination,
	}, BalanceFromAccount)
}

func (m *Manager) ListTransactions(ctx context.Context, query ListQuery[ListTransactions]) (*ListResponse[Transaction], error) {
	var (
		response *TransactionsCursorResponseCursor
		err      error
	)
	if query.PaginationToken == "" {
		response, err = m.client.ListTransactions(ctx, m.ledgerName, ListTransactionsQuery{
			Limit: query.Limit,
			Account: func() string {
				if query.Payload.WalletID != "" {
					return m.chart.GetMainBalanceAccount(query.Payload.WalletID)
				}
				return ""
			}(),
			Metadata: TransactionBaseMetadataFilter(),
		})
	} else {
		response, err = m.client.ListTransactions(ctx, m.ledgerName, ListTransactionsQuery{
			Cursor: query.PaginationToken,
		})
	}
	if err != nil {
		return nil, errors.Wrap(err, "listing transactions")
	}

	return newListResponse[ExpandedTransaction, Transaction](response, func(tx ExpandedTransaction) Transaction {
		return Transaction{
			ExpandedTransaction: tx,
			Ledger:              m.ledgerName,
		}
	}), nil
}

func (m *Manager) CreateWallet(ctx context.Context, data *CreateRequest) (*Wallet, error) {
	wallet := NewWallet(data.Name, m.ledgerName, data.Metadata)

	if err := m.client.AddMetadataToAccount(
		ctx,
		m.ledgerName,
		m.chart.GetMainBalanceAccount(wallet.ID),
		wallet.LedgerMetadata(),
	); err != nil {
		return nil, errors.Wrap(err, "adding metadata to account")
	}

	return &wallet, nil
}

func (m *Manager) UpdateWallet(ctx context.Context, id string, data *PatchRequest) error {
	account, err := m.client.GetAccount(ctx, m.ledgerName, m.chart.GetMainBalanceAccount(id))
	if err != nil {
		return ErrWalletNotFound
	}

	if !IsPrimary(account) {
		return ErrWalletNotFound
	}

	newCustomMetadata := metadata.Metadata{}
	newCustomMetadata = newCustomMetadata.Merge(ExtractCustomMetadata(account))
	newCustomMetadata = newCustomMetadata.Merge(data.Metadata)

	meta := account.GetMetadata()
	meta = meta.Merge(EncodeCustomMetadata(newCustomMetadata))

	if err := m.client.AddMetadataToAccount(ctx, m.ledgerName, m.chart.GetMainBalanceAccount(id), metadata.Metadata(meta)); err != nil {
		return errors.Wrap(err, "adding metadata to account")
	}

	return nil
}

func (m *Manager) GetWallet(ctx context.Context, id string) (*WithBalances, error) {
	account, err := m.client.GetAccount(
		ctx,
		m.ledgerName,
		m.chart.GetMainBalanceAccount(id),
	)
	if err != nil {
		return nil, errors.Wrap(err, "getting account")
	}

	if !IsPrimary(account) {
		return nil, ErrWalletNotFound
	}

	return Ptr(WithBalancesFromAccount(m.ledgerName, *account)), nil
}

type Summary struct {
	Balances       []ExpandedBalance   `json:"balances"`
	AvailableFunds map[string]*big.Int `json:"availableFunds"`
	ExpiredFunds   map[string]*big.Int `json:"expiredFunds"`
	ExpirableFunds map[string]*big.Int `json:"expirableFunds"`
	HoldFunds      map[string]*big.Int `json:"holdFunds"`
}

func (m *Manager) GetWalletSummary(ctx context.Context, id string) (*Summary, error) {
	balances, err := fetchAndMapAllAccounts(ctx, m, metadata.Metadata{
		MetadataKeyWalletID: id,
	}, func(src Account) ExpandedBalance {
		account, err := m.client.GetAccount(ctx, m.ledgerName, src.GetAddress())
		if err != nil {
			// TODO: refine error handling
			panic(errors.Wrap(err, "getting account"))
		}
		return ExpandedBalanceFromAccount(*account)
	})
	if err != nil {
		return nil, err
	}

	s := &Summary{
		Balances:       balances,
		AvailableFunds: map[string]*big.Int{},
		ExpiredFunds:   map[string]*big.Int{},
		ExpirableFunds: map[string]*big.Int{},
		HoldFunds:      map[string]*big.Int{},
	}

	for _, balance := range balances {
		for asset, amount := range balance.Assets {
			switch {
			case balance.ExpiresAt != nil && balance.ExpiresAt.Before(time.Now()):
				if s.ExpiredFunds[asset] == nil {
					s.ExpiredFunds[asset] = new(big.Int)
				}
				s.ExpiredFunds[asset].Add(s.ExpiredFunds[asset], amount)
			case balance.ExpiresAt != nil && !balance.ExpiresAt.Before(time.Now()):
				if s.ExpirableFunds[asset] == nil {
					s.ExpirableFunds[asset] = new(big.Int)
				}
				s.ExpirableFunds[asset].Add(s.ExpirableFunds[asset], amount)
				if s.AvailableFunds[asset] == nil {
					s.AvailableFunds[asset] = new(big.Int)
				}
				s.AvailableFunds[asset].Add(s.AvailableFunds[asset], amount)
			case balance.ExpiresAt == nil:
				if s.AvailableFunds[asset] == nil {
					s.AvailableFunds[asset] = new(big.Int)
				}
				s.AvailableFunds[asset].Add(s.AvailableFunds[asset], amount)
			}
		}
	}

	holds, err := fetchAndMapAllAccounts(ctx, m, metadata.Metadata{
		MetadataKeyHoldWalletID: id,
	}, func(src Account) ExpandedDebitHold {
		account, err := m.client.GetAccount(ctx, m.ledgerName, src.GetAddress())
		if err != nil {
			// TODO: refine error handling
			panic(errors.Wrap(err, "getting account"))
		}

		return ExpandedDebitHoldFromLedgerAccount(*account)
	})
	if err != nil {
		return nil, err
	}

	for _, hold := range holds {
		if s.HoldFunds[hold.Asset] == nil {
			s.HoldFunds[hold.Asset] = new(big.Int)
		}
		s.HoldFunds[hold.Asset].Add(s.HoldFunds[hold.Asset], hold.Remaining)
	}

	return s, nil
}

func (m *Manager) GetHold(ctx context.Context, id string) (*ExpandedDebitHold, error) {
	account, err := m.client.GetAccount(ctx, m.ledgerName, m.chart.GetHoldAccount(id))
	if err != nil {
		return nil, err
	}

	return Ptr(ExpandedDebitHoldFromLedgerAccount(*account)), nil
}

func (m *Manager) CreateBalance(ctx context.Context, data *CreateBalance) (*Balance, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}
	ret, err := m.client.GetAccount(ctx, m.ledgerName, m.chart.GetBalanceAccount(data.WalletID, data.Name))
	if err != nil {
		return nil, err
	}
	if ret.Metadata != nil && ret.Metadata[MetadataKeyWalletBalance] == TrueValue {
		return nil, ErrBalanceAlreadyExists
	}

	balance := NewBalance(data.Name, data.ExpiresAt)

	if err := m.client.AddMetadataToAccount(
		ctx,
		m.ledgerName,
		m.chart.GetBalanceAccount(data.WalletID, balance.Name),
		balance.LedgerMetadata(data.WalletID),
	); err != nil {
		return nil, errors.Wrap(err, "adding metadata to account")
	}

	return &balance, nil
}

func (m *Manager) GetBalance(ctx context.Context, walletID string, balanceName string) (*ExpandedBalance, error) {
	account, err := m.client.GetAccount(ctx, m.ledgerName, m.chart.GetBalanceAccount(walletID, balanceName))
	if err != nil {
		return nil, err
	}
	if account.Metadata[MetadataKeyWalletBalance] != TrueValue {
		return nil, ErrBalanceNotExists
	}

	return Ptr(ExpandedBalanceFromAccount(*account)), nil
}

type mapAccountListQuery struct {
	Pagination
	Metadata func() metadata.Metadata
}

func mapAccountList[TO any](ctx context.Context, r *Manager, query mapAccountListQuery, mapper mapper[Account, TO]) (*ListResponse[TO], error) {
	var (
		response *AccountsCursorResponseCursor
		err      error
	)
	if query.PaginationToken == "" {
		response, err = r.client.ListAccounts(ctx, r.ledgerName, ListAccountsQuery{
			Limit:    query.Limit,
			Metadata: query.Metadata(),
		})
	} else {
		response, err = r.client.ListAccounts(ctx, r.ledgerName, ListAccountsQuery{
			Cursor: query.PaginationToken,
		})
	}
	if err != nil {
		return nil, err
	}

	return newListResponse[Account, TO](response, func(item Account) TO {
		return mapper(item)
	}), nil
}

const maxPageSize = 100

func fetchAndMapAllAccounts[TO any](ctx context.Context, r *Manager, md metadata.Metadata, mapper mapper[Account, TO]) ([]TO, error) {
	ret := make([]TO, 0)
	query := mapAccountListQuery{
		Metadata: func() metadata.Metadata {
			return md
		},
		Pagination: Pagination{
			Limit: maxPageSize,
		},
	}
	for {
		listResponse, err := mapAccountList(ctx, r, query, mapper)
		if err != nil {
			return nil, err
		}
		ret = append(ret, listResponse.Data...)
		if listResponse.Next == "" {
			return ret, nil
		}
		query = mapAccountListQuery{
			Pagination: Pagination{
				PaginationToken: listResponse.Next,
			},
		}
	}
}
