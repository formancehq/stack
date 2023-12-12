package send

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/orchestration/internal/workflow/activities"
	"github.com/formancehq/orchestration/internal/workflow/stages/internal"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	internalLedger         = "orchestration-000-internal"
	moveToLedgerMetadata   = "orchestration/move-to-ledger"
	moveFromLedgerMetadata = "orchestration/move-from-ledger"
)

func extractFormanceAccountID[V any](metadataKey string, metadata map[string]V) (string, error) {
	formanceAccountID, ok := metadata[metadataKey]
	if !ok {
		return "", fmt.Errorf("expected '%s' metadata containing formance account ID", metadataKey)
	}
	if reflect.ValueOf(formanceAccountID).IsZero() {
		return "", errors.New("formance account ID empty")
	}
	return fmt.Sprint(formanceAccountID), nil
}

func justError[T any](v T, err error) error {
	return err
}

func RunSend(ctx workflow.Context, send Send) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.WithStack(fmt.Errorf("%s", e))
		}
	}()
	amount := send.Amount
	switch {
	case send.Source.Account != nil && send.Destination.Account != nil:
		return runAccountToAccount(ctx, send.Source.Account, send.Destination.Account, amount)
	case send.Source.Account != nil && send.Destination.Payment != nil:
		return runAccountToPayment(ctx, send.ConnectorID, send.Source.Account, send.Destination.Payment, amount)
	case send.Source.Account != nil && send.Destination.Wallet != nil:
		return runAccountToWallet(ctx, send.Source.Account, send.Destination.Wallet, amount)
	case send.Source.Wallet != nil && send.Destination.Account != nil:
		return runWalletToAccount(ctx, send.Source.Wallet, send.Destination.Account, amount)
	case send.Source.Wallet != nil && send.Destination.Payment != nil:
		return runWalletToPayment(ctx, send.ConnectorID, send.Source.Wallet, send.Destination.Payment, amount)
	case send.Source.Wallet != nil && send.Destination.Wallet != nil:
		return runWalletToWallet(ctx, send.Source.Wallet, send.Destination.Wallet, amount)
	case send.Source.Payment != nil && send.Destination.Account != nil:
		return runPaymentToAccount(ctx, send.Source.Payment, send.Destination.Account, amount)
	case send.Source.Payment != nil && send.Destination.Wallet != nil:
		return runPaymentToWallet(ctx, send.Source.Payment, send.Destination.Wallet, amount)
	case send.Source.Payment != nil && send.Destination.Payment != nil:
		return errors.New("send from payment to payment is not supported")
	}
	panic("should not happen")
}

func runPaymentToWallet(ctx workflow.Context, source *PaymentSource, destination *WalletSource, amount *shared.Monetary) error {
	payment, err := savePayment(ctx, source.ID)
	if err != nil {
		return err
	}
	if amount == nil {
		amount = &shared.Monetary{
			Amount: payment.InitialAmount,
			Asset:  payment.Asset,
		}
	}
	return runAccountToWallet(ctx, &LedgerAccountSource{
		ID:     paymentAccountName(source.ID),
		Ledger: internalLedger,
	}, destination, amount)
}

func paymentAccountName(paymentID string) string {
	paymentID = strings.ReplaceAll(paymentID, "-", "_")
	return fmt.Sprintf("payment:%s", paymentID)
}

func savePayment(ctx workflow.Context, paymentID string) (*shared.Payment, error) {
	payment, err := activities.GetPayment(internal.InfiniteRetryContext(ctx), paymentID)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving payment: %s", paymentID)
	}
	reference := paymentAccountName(paymentID)
	_, err = activities.CreateTransaction(internal.InfiniteRetryContext(ctx), internalLedger, shared.PostTransaction{
		Postings: []shared.Posting{{
			Amount:      payment.InitialAmount,
			Asset:       payment.Asset,
			Destination: paymentAccountName(paymentID),
			Source:      "world",
		}},
		Metadata:  map[string]interface{}{},
		Reference: &reference,
	})
	if err != nil {
		applicationError := &temporal.ApplicationError{}
		if errors.As(err, &applicationError) {
			if applicationError.Type() != "CONFLICT" {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return payment, nil
}

func runPaymentToAccount(ctx workflow.Context, source *PaymentSource, destination *LedgerAccountDestination, amount *shared.Monetary) error {
	payment, err := savePayment(ctx, source.ID)
	if err != nil {
		return err
	}
	if amount == nil {
		amount = &shared.Monetary{
			Amount: payment.InitialAmount,
			Asset:  payment.Asset,
		}
	}
	return runAccountToAccount(ctx, &LedgerAccountSource{
		ID:     paymentAccountName(source.ID),
		Ledger: internalLedger,
	}, destination, amount)
}

func runWalletToWallet(ctx workflow.Context, source *WalletSource, destination *WalletDestination, amount *shared.Monetary) error {
	if amount == nil {
		return errors.New("amount must be specified")
	}
	walletSource, err := activities.GetWallet(internal.InfiniteRetryContext(ctx), source.ID)
	if err != nil {
		return err
	}
	walletDestination, err := activities.GetWallet(internal.InfiniteRetryContext(ctx), destination.ID)
	if err != nil {
		return err
	}
	if walletSource.Ledger == walletDestination.Ledger {
		mainBalance := "main"
		sourceSubject := shared.WalletSubject{
			Balance:    &mainBalance,
			Identifier: source.ID,
			Type:       "WALLET",
		}
		return activities.CreditWallet(internal.InfiniteRetryContext(ctx), destination.ID, &shared.CreditWalletRequest{
			Amount:   *amount,
			Balance:  &destination.Balance,
			Metadata: map[string]string{},
			Sources:  []shared.Subject{{WalletSubject: &sourceSubject}},
		})
	}

	if err := justError(activities.DebitWallet(internal.InfiniteRetryContext(ctx), source.ID, &shared.DebitWalletRequest{
		Amount:   *amount,
		Balances: []string{source.Balance},
		Metadata: map[string]string{
			moveToLedgerMetadata: walletDestination.Ledger,
		},
	})); err != nil {
		return err
	}

	return activities.CreditWallet(internal.InfiniteRetryContext(ctx), destination.ID, &shared.CreditWalletRequest{
		Amount:  *amount,
		Balance: &destination.Balance,
		Metadata: map[string]string{
			moveFromLedgerMetadata: walletSource.Ledger,
		},
	})
}

func runWalletToPayment(ctx workflow.Context, connectorID *string, source *WalletSource, destination *PaymentDestination, amount *shared.Monetary) error {
	if amount == nil {
		return errors.New("amount must be specified")
	}
	if destination.PSP != "stripe" {
		return errors.New("only stripe actually supported")
	}
	wallet, err := activities.GetWallet(internal.InfiniteRetryContext(ctx), source.ID)
	if err != nil {
		return errors.Wrapf(err, "reading account: %s", source.ID)
	}

	formanceAccountID, err := extractFormanceAccountID(destination.Metadata, wallet.Metadata)
	if err != nil {
		return err
	}

	if err := activities.StripeTransfer(internal.InfiniteRetryContext(ctx), shared.ActivityStripeTransfer{
		Amount:            amount.Amount,
		Asset:             &amount.Asset,
		Destination:       &formanceAccountID,
		WaitingValidation: &destination.WaitingValidation,
		ConnectorID:       connectorID,
	}); err != nil {
		return err
	}
	return justError(activities.DebitWallet(internal.InfiniteRetryContext(ctx), source.ID, &shared.DebitWalletRequest{
		Amount:   *amount,
		Balances: []string{source.Balance},
	}))
}

func runWalletToAccount(ctx workflow.Context, source *WalletSource, destination *LedgerAccountDestination, amount *shared.Monetary) error {
	if amount == nil {
		return errors.New("amount must be specified")
	}
	wallet, err := activities.GetWallet(internal.InfiniteRetryContext(ctx), source.ID)
	if err != nil {
		return err
	}
	if wallet.Ledger == destination.Ledger {

		return justError(activities.DebitWallet(internal.InfiniteRetryContext(ctx), source.ID, &shared.DebitWalletRequest{
			Amount: *amount,
			Destination: &shared.Subject{
				LedgerAccountSubject: &shared.LedgerAccountSubject{
					Identifier: destination.ID,
					Type:       "ACCOUNT",
				},
			},
			Balances: []string{source.Balance},
		}))
	}

	if err := justError(activities.DebitWallet(internal.InfiniteRetryContext(ctx), source.ID, &shared.DebitWalletRequest{
		Amount:   *amount,
		Balances: []string{source.Balance},
		Metadata: map[string]string{
			moveToLedgerMetadata: destination.Ledger,
		},
	})); err != nil {
		return err
	}

	return justError(activities.CreateTransaction(internal.InfiniteRetryContext(ctx), destination.Ledger, shared.PostTransaction{
		Postings: []shared.Posting{{
			Amount:      amount.Amount,
			Asset:       amount.Asset,
			Destination: destination.ID,
			Source:      "world",
		}},
		Metadata: map[string]any{
			moveFromLedgerMetadata: wallet.Ledger,
		},
	}))
}

func runAccountToWallet(ctx workflow.Context, source *LedgerAccountSource, destination *WalletDestination, amount *shared.Monetary) error {
	if amount == nil {
		return errors.New("amount must be specified")
	}
	wallet, err := activities.GetWallet(internal.InfiniteRetryContext(ctx), destination.ID)
	if err != nil {
		return err
	}
	if wallet.Ledger == source.Ledger {
		return activities.CreditWallet(internal.InfiniteRetryContext(ctx), destination.ID, &shared.CreditWalletRequest{
			Amount: *amount,
			Sources: []shared.Subject{{
				LedgerAccountSubject: &shared.LedgerAccountSubject{
					Identifier: source.ID,
					Type:       "ACCOUNT",
				},
			}},
			Balance: &destination.Balance,
		})
	}

	if err := justError(activities.CreateTransaction(internal.InfiniteRetryContext(ctx), source.Ledger, shared.PostTransaction{
		Postings: []shared.Posting{{
			Amount:      amount.Amount,
			Asset:       amount.Asset,
			Destination: "world",
			Source:      source.ID,
		}},
		Metadata: map[string]any{
			moveToLedgerMetadata: wallet.Ledger,
		},
	})); err != nil {
		return err
	}

	return activities.CreditWallet(internal.InfiniteRetryContext(ctx), destination.ID, &shared.CreditWalletRequest{
		Amount: *amount,
		Sources: []shared.Subject{{
			LedgerAccountSubject: &shared.LedgerAccountSubject{
				Identifier: "world",
				Type:       "ACCOUNT",
			},
		}},
		Balance: &destination.Balance,
		Metadata: map[string]string{
			moveFromLedgerMetadata: source.Ledger,
		},
	})
}

func runAccountToAccount(ctx workflow.Context, source *LedgerAccountSource, destination *LedgerAccountDestination, amount *shared.Monetary) error {
	if amount == nil {
		return errors.New("amount must be specified")
	}
	if source.Ledger == destination.Ledger {
		return justError(activities.CreateTransaction(internal.InfiniteRetryContext(ctx), destination.Ledger, shared.PostTransaction{
			Postings: []shared.Posting{{
				Amount:      amount.Amount,
				Asset:       amount.Asset,
				Destination: destination.ID,
				Source:      source.ID,
			}},
			Metadata: map[string]interface{}{},
		}))
	}
	if err := justError(activities.CreateTransaction(internal.InfiniteRetryContext(ctx), source.Ledger, shared.PostTransaction{
		Postings: []shared.Posting{{
			Amount:      amount.Amount,
			Asset:       amount.Asset,
			Destination: "world",
			Source:      source.ID,
		}},
		Metadata: map[string]any{
			moveToLedgerMetadata: destination.Ledger,
		},
	})); err != nil {
		return err
	}
	return justError(activities.CreateTransaction(internal.InfiniteRetryContext(ctx), destination.Ledger, shared.PostTransaction{
		Postings: []shared.Posting{{
			Amount:      amount.Amount,
			Asset:       amount.Asset,
			Destination: destination.ID,
			Source:      "world",
		}},
		Metadata: map[string]any{
			moveFromLedgerMetadata: source.Ledger,
		},
	}))
}

func runAccountToPayment(ctx workflow.Context, connectorID *string, source *LedgerAccountSource, destination *PaymentDestination, amount *shared.Monetary) error {
	if amount == nil {
		return errors.New("amount must be specified")
	}
	if destination.PSP != "stripe" {
		return errors.New("only stripe actually supported")
	}
	account, err := activities.GetAccount(internal.InfiniteRetryContext(ctx), source.Ledger, source.ID)
	if err != nil {
		return errors.Wrapf(err, "reading account: %s", source.ID)
	}
	formanceAccountID, err := extractFormanceAccountID(destination.Metadata, account.Metadata)
	if err != nil {
		return err
	}

	if err := activities.StripeTransfer(internal.InfiniteRetryContext(ctx), shared.ActivityStripeTransfer{
		Amount:            amount.Amount,
		Asset:             &amount.Asset,
		Destination:       &formanceAccountID,
		WaitingValidation: &destination.WaitingValidation,
		ConnectorID:       connectorID,
	}); err != nil {
		return err
	}
	return justError(activities.CreateTransaction(internal.InfiniteRetryContext(ctx), source.Ledger, shared.PostTransaction{
		Postings: []shared.Posting{{
			Amount:      amount.Amount,
			Asset:       amount.Asset,
			Destination: "world",
			Source:      source.ID,
		}},
		Metadata: map[string]any{},
	}))
}
