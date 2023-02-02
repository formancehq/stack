package send

import (
	"fmt"
	"runtime/debug"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/orchestration/internal/workflow/activities"
	"github.com/formancehq/orchestration/internal/workflow/stages/internal"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/workflow"
)

func justError[T any](v T, err error) error {
	return err
}

func RunSend(ctx workflow.Context, send Send) error {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
			debug.PrintStack()
		}
	}()
	switch {
	case send.Source.Account != nil && send.Destination.Account != nil:
		return runAccountToAccount(ctx, send)
	case send.Source.Account != nil && send.Destination.Payment != nil:
		return runAccountToPayment(ctx, send)
	case send.Source.Account != nil && send.Destination.Wallet != nil:
		return runAccountToWallet(ctx, send)
	case send.Source.Wallet != nil && send.Destination.Account != nil:
		return runWalletToAccount(ctx, send)
	case send.Source.Wallet != nil && send.Destination.Payment != nil:
		return runWalletToPayment(ctx, send)
	case send.Source.Wallet != nil && send.Destination.Wallet != nil:
		return runWalletToWallet(ctx, send)
	case send.Source.Payment != nil && send.Destination.Account != nil:
		return runPaymentToAccount(ctx, send)
	case send.Source.Payment != nil && send.Destination.Wallet != nil:
		return runPaymentToWallet(ctx, send)
	case send.Source.Payment != nil && send.Destination.Payment != nil:
		return errors.New("send from payment to payment is not supported")
	}
	panic("should not happen")
}

func runPaymentToWallet(ctx workflow.Context, send Send) error {
	payment, err := activities.GetPayment(internal.SingleTryContext(ctx), send.Source.Payment.ID)
	if err != nil {
		return errors.Wrapf(err, "retrieving payment: %s", send.Source.Payment.ID)
	}
	amount := *sdk.NewMonetary(send.Amount.Asset, send.Amount.Amount)
	if amount.Asset != payment.Asset || amount.Amount > payment.InitialAmount {
		return fmt.Errorf("payment amount invalid")
	}
	return activities.CreditWallet(internal.SingleTryContext(ctx), send.Destination.Wallet.ID, sdk.CreditWalletRequest{
		Amount: amount,
		Sources: []sdk.Subject{{
			LedgerAccountSubject: sdk.NewLedgerAccountSubject("ACCOUNT", "world"),
		}},
		Balance: sdk.PtrString(send.Destination.Wallet.Balance),
	})
}

func runPaymentToAccount(ctx workflow.Context, send Send) error {
	payment, err := activities.GetPayment(internal.SingleTryContext(ctx), send.Source.Payment.ID)
	if err != nil {
		return errors.Wrapf(err, "retrieving payment: %s", send.Source.Payment.ID)
	}
	amount := *sdk.NewMonetary(send.Amount.Asset, send.Amount.Amount)
	if amount.Asset != payment.Asset || amount.Amount > payment.InitialAmount {
		return fmt.Errorf("payment amount invalid")
	}
	return justError(activities.CreateTransaction(internal.SingleTryContext(ctx), send.Destination.Account.Ledger, sdk.PostTransaction{
		Postings: []sdk.Posting{{
			Amount:      amount.Amount,
			Asset:       amount.Asset,
			Destination: send.Destination.Account.ID,
			Source:      "world",
		}},
	}))
}

func runWalletToWallet(ctx workflow.Context, send Send) error {
	walletSource, err := activities.GetWallet(internal.SingleTryContext(ctx), send.Source.Wallet.ID)
	if err != nil {
		return err
	}
	walletDestination, err := activities.GetWallet(internal.SingleTryContext(ctx), send.Destination.Wallet.ID)
	if err != nil {
		return err
	}
	if walletSource.Ledger != walletDestination.Ledger {
		return errors.New("wallets not on the same ledger")
	}
	sourceSubject := sdk.NewWalletSubject("WALLET", send.Source.Wallet.ID)
	sourceSubject.SetBalance("main")
	return activities.CreditWallet(internal.SingleTryContext(ctx), send.Destination.Wallet.ID, sdk.CreditWalletRequest{
		Amount: *sdk.NewMonetary(send.Amount.Asset, send.Amount.Amount),
		Sources: []sdk.Subject{{
			WalletSubject: sourceSubject,
		}},
		Balance: sdk.PtrString(send.Destination.Wallet.Balance),
	})
}

func runWalletToPayment(ctx workflow.Context, send Send) error {
	if send.Destination.Payment.PSP != "stripe" {
		return errors.New("only stripe actually supported")
	}
	wallet, err := activities.GetWallet(internal.SingleTryContext(ctx), send.Source.Wallet.ID)
	if err != nil {
		return errors.Wrapf(err, "reading account: %s", send.Source.Account.ID)
	}

	stripeConnectID, err := extractStripeConnectID(wallet)
	if err != nil {
		return err
	}
	if err := activities.StripeTransfer(internal.SingleTryContext(ctx), sdk.StripeTransferRequest{
		Amount:      sdk.PtrInt64(send.Amount.Amount),
		Asset:       sdk.PtrString(send.Amount.Asset),
		Destination: sdk.PtrString(stripeConnectID),
	}); err != nil {
		return err
	}
	return justError(activities.DebitWallet(internal.InfiniteRetryContext(ctx), send.Source.Wallet.ID, sdk.DebitWalletRequest{
		Amount:   *sdk.NewMonetary(send.Amount.Asset, send.Amount.Amount),
		Balances: []string{send.Source.Wallet.Balance},
	}))
}

func runWalletToAccount(ctx workflow.Context, send Send) error {
	wallet, err := activities.GetWallet(internal.SingleTryContext(ctx), send.Source.Wallet.ID)
	if err != nil {
		return err
	}
	if wallet.Ledger != send.Destination.Account.Ledger {
		return errors.New("wallet not on the same ledger than the account")
	}
	return justError(activities.DebitWallet(internal.SingleTryContext(ctx), send.Source.Wallet.ID, sdk.DebitWalletRequest{
		Amount: *sdk.NewMonetary(send.Amount.Asset, send.Amount.Amount),
		Destination: &sdk.Subject{
			LedgerAccountSubject: sdk.NewLedgerAccountSubject("ACCOUNT", send.Destination.Account.ID),
		},
		Balances: []string{send.Source.Wallet.Balance},
	}))
}

func runAccountToWallet(ctx workflow.Context, send Send) error {
	wallet, err := activities.GetWallet(internal.SingleTryContext(ctx), send.Destination.Wallet.ID)
	if err != nil {
		return err
	}
	if wallet.Ledger == send.Source.Account.Ledger {
		return activities.CreditWallet(internal.SingleTryContext(ctx), send.Destination.Wallet.ID, sdk.CreditWalletRequest{
			Amount: *sdk.NewMonetary(send.Amount.Asset, send.Amount.Amount),
			Sources: []sdk.Subject{{
				LedgerAccountSubject: sdk.NewLedgerAccountSubject("ACCOUNT", send.Source.Account.ID),
			}},
			Balance: sdk.PtrString(send.Destination.Wallet.Balance),
		})
	}

	if err := justError(activities.CreateTransaction(internal.SingleTryContext(ctx), send.Source.Account.Ledger, sdk.PostTransaction{
		Postings: []sdk.Posting{{
			Amount:      send.Amount.Amount,
			Asset:       send.Amount.Asset,
			Destination: "world",
			Source:      send.Source.Account.ID,
		}},
		Metadata: map[string]interface{}{
			moveToLedgerMetadata: wallet.Ledger,
		},
	})); err != nil {
		return err
	}

	return activities.CreditWallet(internal.InfiniteRetryContext(ctx), send.Destination.Wallet.ID, sdk.CreditWalletRequest{
		Amount: *sdk.NewMonetary(send.Amount.Asset, send.Amount.Amount),
		Sources: []sdk.Subject{{
			LedgerAccountSubject: sdk.NewLedgerAccountSubject("ACCOUNT", "world"),
		}},
		Balance: sdk.PtrString(send.Destination.Wallet.Balance),
		Metadata: map[string]interface{}{
			moveFromLedgerMetadata: send.Source.Account.Ledger,
		},
	})
}

func extractStripeConnectID(object interface {
	GetMetadata() map[string]any
}) (string, error) {
	stripeConnectIDAny, ok := object.GetMetadata()["stripeConnectID"]
	if !ok {
		return "", errors.New("expected 'stripeConnectID' metadata containing connected account ID")
	}
	stripeConnectID, ok := stripeConnectIDAny.(string)
	if !ok {
		return "", errors.New("expected 'stripeConnectID' to be a string")
	}
	if stripeConnectID == "" {
		return "", errors.New("stripe connect ID empty")
	}
	return stripeConnectID, nil
}

const (
	moveToLedgerMetadata   = "orchestration/move-to-ledger"
	moveFromLedgerMetadata = "orchestration/move-from-ledger"
)

func runAccountToAccount(ctx workflow.Context, send Send) error {
	if send.Source.Account.Ledger == send.Destination.Account.Ledger {
		return justError(activities.CreateTransaction(internal.SingleTryContext(ctx), send.Destination.Account.Ledger, sdk.PostTransaction{
			Postings: []sdk.Posting{{
				Amount:      send.Amount.Amount,
				Asset:       send.Amount.Asset,
				Destination: send.Destination.Account.ID,
				Source:      send.Source.Account.ID,
			}},
		}))
	}
	if err := justError(activities.CreateTransaction(internal.SingleTryContext(ctx), send.Source.Account.Ledger, sdk.PostTransaction{
		Postings: []sdk.Posting{{
			Amount:      send.Amount.Amount,
			Asset:       send.Amount.Asset,
			Destination: "world",
			Source:      send.Source.Account.ID,
		}},
		Metadata: map[string]interface{}{
			moveToLedgerMetadata: send.Destination.Account.Ledger,
		},
	})); err != nil {
		return err
	}
	return justError(activities.CreateTransaction(internal.InfiniteRetryContext(ctx), send.Destination.Account.Ledger, sdk.PostTransaction{
		Postings: []sdk.Posting{{
			Amount:      send.Amount.Amount,
			Asset:       send.Amount.Asset,
			Destination: send.Destination.Account.ID,
			Source:      "world",
		}},
		Metadata: map[string]interface{}{
			moveFromLedgerMetadata: send.Source.Account.Ledger,
		},
	}))
}

func runAccountToPayment(ctx workflow.Context, send Send) error {
	if send.Destination.Payment.PSP != "stripe" {
		return errors.New("only stripe actually supported")
	}
	account, err := activities.GetAccount(internal.SingleTryContext(ctx), send.Source.Account.Ledger, send.Source.Account.ID)
	if err != nil {
		return errors.Wrapf(err, "reading account: %s", send.Source.Account.ID)
	}
	stripeConnectID, err := extractStripeConnectID(account)
	if err != nil {
		return err
	}
	if err := activities.StripeTransfer(internal.SingleTryContext(ctx), sdk.StripeTransferRequest{
		Amount:      sdk.PtrInt64(send.Amount.Amount),
		Asset:       sdk.PtrString(send.Amount.Asset),
		Destination: sdk.PtrString(stripeConnectID),
	}); err != nil {
		return err
	}
	return justError(activities.CreateTransaction(internal.InfiniteRetryContext(ctx), send.Source.Account.Ledger, sdk.PostTransaction{
		Postings: []sdk.Posting{{
			Amount:      send.Amount.Amount,
			Asset:       send.Amount.Asset,
			Destination: "world",
			Source:      send.Source.Account.ID,
		}},
	}))
}
