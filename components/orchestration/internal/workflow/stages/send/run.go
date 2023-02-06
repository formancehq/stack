package send

import (
	"fmt"
	"strings"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/orchestration/internal/workflow/activities"
	"github.com/formancehq/orchestration/internal/workflow/stages/internal"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/workflow"
)

const (
	internalLedger         = "orchestration-000-internal"
	moveToLedgerMetadata   = "orchestration/move-to-ledger"
	moveFromLedgerMetadata = "orchestration/move-from-ledger"
)

func extractStripeConnectID(metadataKey string, object interface {
	GetMetadata() map[string]any
}) (string, error) {
	stripeConnectIDAny, ok := object.GetMetadata()[metadataKey]
	if !ok {
		return "", fmt.Errorf("expected '%s' metadata containing connected account ID", metadataKey)
	}
	stripeConnectID, ok := stripeConnectIDAny.(string)
	if !ok {
		return "", fmt.Errorf("expected '%s' to be a string", metadataKey)
	}
	if stripeConnectID == "" {
		return "", errors.New("stripe connect ID empty")
	}
	return stripeConnectID, nil
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
	amount := *sdk.NewMonetary(send.Amount.Asset, send.Amount.Amount)
	switch {
	case send.Source.Account != nil && send.Destination.Account != nil:
		return runAccountToAccount(ctx, send.Source.Account, send.Destination.Account, amount)
	case send.Source.Account != nil && send.Destination.Payment != nil:
		return runAccountToPayment(ctx, send.Source.Account, send.Destination.Payment, amount)
	case send.Source.Account != nil && send.Destination.Wallet != nil:
		return runAccountToWallet(ctx, send.Source.Account, send.Destination.Wallet, amount)
	case send.Source.Wallet != nil && send.Destination.Account != nil:
		return runWalletToAccount(ctx, send.Source.Wallet, send.Destination.Account, amount)
	case send.Source.Wallet != nil && send.Destination.Payment != nil:
		return runWalletToPayment(ctx, send.Source.Wallet, send.Destination.Payment, amount)
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

func runPaymentToWallet(ctx workflow.Context, source *PaymentSource, destination *WalletSource, amount sdk.Monetary) error {
	if err := savePayment(ctx, source.ID); err != nil {
		return err
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

func savePayment(ctx workflow.Context, paymentID string) error {
	payment, err := activities.GetPayment(internal.SingleTryContext(ctx), paymentID)
	if err != nil {
		return errors.Wrapf(err, "retrieving payment: %s", paymentID)
	}
	_, err = activities.CreateTransaction(internal.SingleTryContext(ctx), internalLedger, sdk.PostTransaction{
		Postings: []sdk.Posting{{
			Amount:      payment.InitialAmount,
			Asset:       payment.Asset,
			Destination: paymentAccountName(paymentID),
			Source:      "world",
		}},
		Reference: sdk.PtrString(paymentAccountName(paymentID)),
	})
	if err != nil && err.Error() != "CONFLICT" {
		return err
	}
	return nil
}

func runPaymentToAccount(ctx workflow.Context, source *PaymentSource, destination *LedgerAccountDestination, amount sdk.Monetary) error {
	if err := savePayment(ctx, source.ID); err != nil {
		return err
	}
	return runAccountToAccount(ctx, &LedgerAccountSource{
		ID:     paymentAccountName(source.ID),
		Ledger: internalLedger,
	}, destination, amount)
}

func runWalletToWallet(ctx workflow.Context, source *WalletSource, destination *WalletDestination, amount sdk.Monetary) error {
	walletSource, err := activities.GetWallet(internal.SingleTryContext(ctx), source.ID)
	if err != nil {
		return err
	}
	walletDestination, err := activities.GetWallet(internal.SingleTryContext(ctx), destination.ID)
	if err != nil {
		return err
	}
	if walletSource.Ledger == walletDestination.Ledger {
		sourceSubject := sdk.NewWalletSubject("WALLET", source.ID)
		sourceSubject.SetBalance("main")
		return activities.CreditWallet(internal.SingleTryContext(ctx), destination.ID, sdk.CreditWalletRequest{
			Amount: *sdk.NewMonetary(amount.Asset, amount.Amount),
			Sources: []sdk.Subject{{
				WalletSubject: sourceSubject,
			}},
			Balance: sdk.PtrString(destination.Balance),
		})
	}

	if err := justError(activities.DebitWallet(internal.SingleTryContext(ctx), source.ID, sdk.DebitWalletRequest{
		Amount:   *sdk.NewMonetary(amount.Asset, amount.Amount),
		Balances: []string{source.Balance},
		Metadata: map[string]interface{}{
			moveToLedgerMetadata: walletDestination.Ledger,
		},
	})); err != nil {
		return err
	}

	return activities.CreditWallet(internal.InfiniteRetryContext(ctx), destination.ID, sdk.CreditWalletRequest{
		Amount:  *sdk.NewMonetary(amount.Asset, amount.Amount),
		Balance: sdk.PtrString(destination.Balance),
		Metadata: map[string]interface{}{
			moveFromLedgerMetadata: walletSource.Ledger,
		},
	})
}

func runWalletToPayment(ctx workflow.Context, source *WalletSource, destination *PaymentDestination, amount sdk.Monetary) error {
	if destination.PSP != "stripe" {
		return errors.New("only stripe actually supported")
	}
	wallet, err := activities.GetWallet(internal.SingleTryContext(ctx), source.ID)
	if err != nil {
		return errors.Wrapf(err, "reading account: %s", source.ID)
	}

	stripeConnectID, err := extractStripeConnectID(destination.Metadata, wallet)
	if err != nil {
		return err
	}
	if err := activities.StripeTransfer(internal.SingleTryContext(ctx), sdk.StripeTransferRequest{
		Amount:      sdk.PtrInt64(amount.Amount),
		Asset:       sdk.PtrString(amount.Asset),
		Destination: sdk.PtrString(stripeConnectID),
	}); err != nil {
		return err
	}
	return justError(activities.DebitWallet(internal.InfiniteRetryContext(ctx), source.ID, sdk.DebitWalletRequest{
		Amount:   *sdk.NewMonetary(amount.Asset, amount.Amount),
		Balances: []string{source.Balance},
	}))
}

func runWalletToAccount(ctx workflow.Context, source *WalletSource, destination *LedgerAccountDestination, amount sdk.Monetary) error {
	wallet, err := activities.GetWallet(internal.SingleTryContext(ctx), source.ID)
	if err != nil {
		return err
	}
	if wallet.Ledger == destination.Ledger {
		return justError(activities.DebitWallet(internal.SingleTryContext(ctx), source.ID, sdk.DebitWalletRequest{
			Amount: *sdk.NewMonetary(amount.Asset, amount.Amount),
			Destination: &sdk.Subject{
				LedgerAccountSubject: sdk.NewLedgerAccountSubject("ACCOUNT", destination.ID),
			},
			Balances: []string{source.Balance},
		}))
	}

	if err := justError(activities.DebitWallet(internal.SingleTryContext(ctx), source.ID, sdk.DebitWalletRequest{
		Amount:   *sdk.NewMonetary(amount.Asset, amount.Amount),
		Balances: []string{source.Balance},
		Metadata: map[string]interface{}{
			moveToLedgerMetadata: destination.Ledger,
		},
	})); err != nil {
		return err
	}

	return justError(activities.CreateTransaction(internal.InfiniteRetryContext(ctx), destination.Ledger, sdk.PostTransaction{
		Postings: []sdk.Posting{{
			Amount:      amount.Amount,
			Asset:       amount.Asset,
			Destination: destination.ID,
			Source:      "world",
		}},
		Metadata: map[string]interface{}{
			moveFromLedgerMetadata: wallet.Ledger,
		},
	}))
}

func runAccountToWallet(ctx workflow.Context, source *LedgerAccountSource, destination *WalletDestination, amount sdk.Monetary) error {
	wallet, err := activities.GetWallet(internal.SingleTryContext(ctx), destination.ID)
	if err != nil {
		return err
	}
	if wallet.Ledger == source.Ledger {
		return activities.CreditWallet(internal.SingleTryContext(ctx), destination.ID, sdk.CreditWalletRequest{
			Amount: *sdk.NewMonetary(amount.Asset, amount.Amount),
			Sources: []sdk.Subject{{
				LedgerAccountSubject: sdk.NewLedgerAccountSubject("ACCOUNT", source.ID),
			}},
			Balance: sdk.PtrString(destination.Balance),
		})
	}

	if err := justError(activities.CreateTransaction(internal.SingleTryContext(ctx), source.Ledger, sdk.PostTransaction{
		Postings: []sdk.Posting{{
			Amount:      amount.Amount,
			Asset:       amount.Asset,
			Destination: "world",
			Source:      source.ID,
		}},
		Metadata: map[string]interface{}{
			moveToLedgerMetadata: wallet.Ledger,
		},
	})); err != nil {
		return err
	}

	return activities.CreditWallet(internal.InfiniteRetryContext(ctx), destination.ID, sdk.CreditWalletRequest{
		Amount: *sdk.NewMonetary(amount.Asset, amount.Amount),
		Sources: []sdk.Subject{{
			LedgerAccountSubject: sdk.NewLedgerAccountSubject("ACCOUNT", "world"),
		}},
		Balance: sdk.PtrString(destination.Balance),
		Metadata: map[string]interface{}{
			moveFromLedgerMetadata: source.Ledger,
		},
	})
}

func runAccountToAccount(ctx workflow.Context, source *LedgerAccountSource, destination *LedgerAccountDestination, amount sdk.Monetary) error {
	if source.Ledger == destination.Ledger {
		return justError(activities.CreateTransaction(internal.SingleTryContext(ctx), destination.Ledger, sdk.PostTransaction{
			Postings: []sdk.Posting{{
				Amount:      amount.Amount,
				Asset:       amount.Asset,
				Destination: destination.ID,
				Source:      source.ID,
			}},
		}))
	}
	if err := justError(activities.CreateTransaction(internal.SingleTryContext(ctx), source.Ledger, sdk.PostTransaction{
		Postings: []sdk.Posting{{
			Amount:      amount.Amount,
			Asset:       amount.Asset,
			Destination: "world",
			Source:      source.ID,
		}},
		Metadata: map[string]interface{}{
			moveToLedgerMetadata: destination.Ledger,
		},
	})); err != nil {
		return err
	}
	return justError(activities.CreateTransaction(internal.InfiniteRetryContext(ctx), destination.Ledger, sdk.PostTransaction{
		Postings: []sdk.Posting{{
			Amount:      amount.Amount,
			Asset:       amount.Asset,
			Destination: destination.ID,
			Source:      "world",
		}},
		Metadata: map[string]interface{}{
			moveFromLedgerMetadata: source.Ledger,
		},
	}))
}

func runAccountToPayment(ctx workflow.Context, source *LedgerAccountSource, destination *PaymentDestination, amount sdk.Monetary) error {
	if destination.PSP != "stripe" {
		return errors.New("only stripe actually supported")
	}
	account, err := activities.GetAccount(internal.SingleTryContext(ctx), source.Ledger, source.ID)
	if err != nil {
		return errors.Wrapf(err, "reading account: %s", source.ID)
	}
	stripeConnectID, err := extractStripeConnectID(destination.Metadata, account)
	if err != nil {
		return err
	}
	if err := activities.StripeTransfer(internal.SingleTryContext(ctx), sdk.StripeTransferRequest{
		Amount:      sdk.PtrInt64(amount.Amount),
		Asset:       sdk.PtrString(amount.Asset),
		Destination: sdk.PtrString(stripeConnectID),
	}); err != nil {
		return err
	}
	return justError(activities.CreateTransaction(internal.InfiniteRetryContext(ctx), source.Ledger, sdk.PostTransaction{
		Postings: []sdk.Posting{{
			Amount:      amount.Amount,
			Asset:       amount.Asset,
			Destination: "world",
			Source:      source.ID,
		}},
	}))
}
