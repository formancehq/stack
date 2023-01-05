package wallet

import (
	"context"
	"strings"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/go-libs/metadata"
	"github.com/formancehq/wallets/pkg/core"
	"github.com/formancehq/wallets/pkg/wallet/numscript"
	"github.com/pkg/errors"
)

const (
	DefaultCreditSource = "world"
	DefaultDebitDest    = "world"
)

type FundingService struct {
	client     Ledger
	chart      *core.Chart
	ledgerName string
}

func NewFundingService(
	ledgerName string,
	client Ledger,
	chart *core.Chart,
) *FundingService {
	return &FundingService{
		client:     client,
		chart:      chart,
		ledgerName: ledgerName,
	}
}

type Debit struct {
	WalletID    string            `json:"walletID"`
	Amount      core.Monetary     `json:"amount"`
	Destination string            `json:"destination"`
	Reference   string            `json:"reference"`
	Pending     bool              `json:"pending"`
	Metadata    metadata.Metadata `json:"metadata"`
	Description string            `json:"description"`
}

type ConfirmHold struct {
	HoldID    string `json:"holdID"`
	Amount    core.MonetaryInt
	Reference string
	Final     bool
}

type VoidHold struct {
	HoldID string `json:"holdID"`
}

type Credit struct {
	WalletID  string            `json:"walletID"`
	Source    string            `json:"source"`
	Amount    core.Monetary     `json:"amount"`
	Reference string            `json:"reference"`
	Metadata  metadata.Metadata `json:"metadata"`
}

func (s *FundingService) Debit(ctx context.Context, debit Debit) (*core.DebitHold, error) {
	dest := DefaultDebitDest
	if debit.Destination != "" {
		dest = debit.Destination
	}

	var hold *core.DebitHold
	if debit.Pending {
		md := debit.Metadata
		if md == nil {
			md = metadata.Metadata{}
		}
		newHold := core.NewDebitHold(debit.WalletID, dest, debit.Amount.Asset, debit.Description, md)
		hold = &newHold

		holdAccount := s.chart.GetHoldAccount(hold.ID)
		if err := s.client.AddMetadataToAccount(ctx, s.ledgerName, holdAccount,
			newHold.LedgerMetadata(s.chart)); err != nil {
			return nil, errors.Wrap(err, "adding metadata to account")
		}

		dest = holdAccount
	}

	transaction := sdk.TransactionData{
		Postings: []sdk.Posting{
			{
				// @todo: upgrade this to proper int after sdk is updated
				Amount:      int32(debit.Amount.Amount.Uint64()),
				Asset:       debit.Amount.Asset,
				Source:      s.chart.GetMainAccount(debit.WalletID),
				Destination: dest,
			},
		},
		Metadata: core.WalletTransactionBaseMetadata(),
	}

	if debit.Reference != "" {
		transaction.Reference = &debit.Reference
	}

	if err := s.client.CreateTransaction(ctx, s.ledgerName, transaction); err != nil {
		return nil, handleCreateTransactionError(err)
	}

	return hold, nil
}

func (s *FundingService) runScript(ctx context.Context, script sdk.Script) error {
	ret, err := s.client.RunScript(ctx, s.ledgerName, script)
	if err != nil {
		return err
	}
	if ret.ErrorCode == nil {
		return nil
	}
	if *ret.ErrorCode == string(sdk.INSUFFICIENT_FUND) {
		return ErrInsufficientFundError
	}
	if ret.ErrorMessage != nil {
		return errors.New(*ret.ErrorMessage)
	}
	return errors.New(*ret.ErrorCode)
}

func (s *FundingService) ConfirmHold(ctx context.Context, debit ConfirmHold) error {
	holdAccount := s.chart.GetHoldAccount(debit.HoldID)

	account, err := s.client.GetAccount(ctx, s.ledgerName, holdAccount)
	if err != nil {
		return errors.Wrap(err, "getting account")
	}

	if !core.IsHold(account) {
		return newErrMismatchType(core.HoldWallet, core.SpecType(account))
	}

	hold := core.ExpandedDebitHoldFromLedgerAccount(account)

	if hold.Remaining.Uint64() == 0 {
		return ErrClosedHold
	}

	amount := hold.Remaining.Uint64()
	if debit.Amount.Uint64() != 0 {
		if debit.Amount.Uint64() > amount {
			return ErrInsufficientFundError
		}
		amount = debit.Amount.Uint64()
	}

	return s.runScript(
		ctx,
		sdk.Script{
			Plain: numscript.BuildConfirmHoldScript(debit.Final, hold.Asset),
			Vars: map[string]interface{}{
				"hold": s.chart.GetHoldAccount(debit.HoldID),
				"amount": map[string]any{
					"amount": amount,
					"asset":  hold.Asset,
				},
			},
			Metadata: core.WalletTransactionBaseMetadata(),
		},
	)
}

func (s *FundingService) VoidHold(ctx context.Context, void VoidHold) error {
	account, err := s.client.GetAccount(ctx, s.ledgerName, s.chart.GetHoldAccount(void.HoldID))
	if err != nil {
		return errors.Wrap(err, "getting account")
	}

	hold := core.ExpandedDebitHoldFromLedgerAccount(account)
	if hold.Remaining.Uint64() == 0 {
		return ErrClosedHold
	}

	return s.runScript(ctx, sdk.Script{
		Plain: strings.ReplaceAll(numscript.CancelHold, "ASSET", hold.Asset),
		Vars: map[string]interface{}{
			"hold": s.chart.GetHoldAccount(void.HoldID),
		},
		Metadata: core.WalletTransactionBaseMetadata(),
	})
}

func (s *FundingService) Credit(ctx context.Context, credit Credit) error {
	source := DefaultCreditSource
	if credit.Source != "" {
		source = credit.Source
	}

	transaction := sdk.TransactionData{
		Postings: []sdk.Posting{
			{
				// @todo: upgrade this to proper int after sdk is updated
				Amount:      int32(credit.Amount.Amount.Uint64()),
				Asset:       credit.Amount.Asset,
				Source:      source,
				Destination: s.chart.GetMainAccount(credit.WalletID),
			},
		},
		Metadata: core.WalletTransactionBaseMetadata().Merge(metadata.Metadata{
			core.MetadataKeyWalletCustomData: credit.Metadata,
		}),
	}

	if credit.Reference != "" {
		transaction.Reference = &credit.Reference
	}

	if err := s.client.CreateTransaction(ctx, s.ledgerName, transaction); err != nil {
		return errors.Wrap(err, "creating transaction")
	}

	return nil
}

func handleCreateTransactionError(err error) error {
	//nolint:errorlint
	if err, ok := err.(interface {
		error
		Model() interface{}
	}); ok {
		if err, ok := err.(interface {
			GetErrorCode() sdk.ErrorCode
		}); ok {
			if err.GetErrorCode() == sdk.INSUFFICIENT_FUND {
				return ErrInsufficientFundError
			}
		}
	}
	return errors.Wrap(err, "running script")
}
