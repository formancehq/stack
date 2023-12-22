package send

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/metadata"

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

func getWalletFromReference(ctx workflow.Context, ref WalletReference) (*shared.Wallet, error) {
	if ref.ID != "" {
		walletSource, err := activities.GetWallet(internal.InfiniteRetryContext(ctx), ref.ID)
		if err != nil {
			return nil, err
		}
		return &shared.Wallet{
			CreatedAt: walletSource.CreatedAt,
			ID:        walletSource.ID,
			Ledger:    walletSource.Ledger,
			Metadata:  walletSource.Metadata,
			Name:      walletSource.Name,
		}, nil
	} else {
		wallets, err := activities.ListWallets(internal.InfiniteRetryContext(ctx), activities.ListWalletsRequest{
			Name: ref.Name,
		})
		if err != nil {
			return nil, err
		}
		switch len(wallets.Cursor.Data) {
		case 0:
			return nil, errors.New("wallet not found")
		case 1:
			return &wallets.Cursor.Data[0], nil
		default:
			return nil, errors.New("found multiple wallets with the same name")
		}
	}
}

func RunSend(ctx workflow.Context, send Send) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.WithStack(fmt.Errorf("%s", e))
		}
	}()
	amount := send.Amount
	metadata := send.Metadata
	if metadata == nil {
		metadata = map[string]string{}
	}
	switch {
	case send.Source.Account != nil && send.Destination.Account != nil:
		return runAccountToAccount(ctx, send.Source.Account, send.Destination.Account, amount, metadata)
	case send.Source.Account != nil && send.Destination.Payment != nil:
		return runAccountToPayment(ctx, send.Source.Account, send.Destination.Payment, amount, metadata)
	case send.Source.Account != nil && send.Destination.Wallet != nil:
		return runAccountToWallet(ctx, send.Source.Account, send.Destination.Wallet, amount, metadata)
	case send.Source.Wallet != nil && send.Destination.Account != nil:
		return runWalletToAccount(ctx, send.Source.Wallet, send.Destination.Account, amount, send.Metadata)
	case send.Source.Wallet != nil && send.Destination.Payment != nil:
		return runWalletToPayment(ctx, send.Source.Wallet, send.Destination.Payment, amount, metadata)
	case send.Source.Wallet != nil && send.Destination.Wallet != nil:
		return runWalletToWallet(ctx, send.Source.Wallet, send.Destination.Wallet, amount, metadata)
	case send.Source.Payment != nil && send.Destination.Account != nil:
		return runPaymentToAccount(ctx, send.Source.Payment, send.Destination.Account, amount, metadata)
	case send.Source.Payment != nil && send.Destination.Wallet != nil:
		return runPaymentToWallet(ctx, send.Source.Payment, send.Destination.Wallet, amount, metadata)
	case send.Source.Payment != nil && send.Destination.Payment != nil:
		return errors.New("send from payment to payment is not supported")
	}
	panic("should not happen")
}

func runPaymentToWallet(ctx workflow.Context, source *PaymentSource, destination *WalletSource, amount *shared.Monetary, m metadata.Metadata) error {
	payment, err := savePayment(ctx, source.ID, m)
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
	}, destination, amount, m)
}

func paymentAccountName(paymentID string) string {
	paymentID = strings.ReplaceAll(paymentID, "-", "_")
	return fmt.Sprintf("payment:%s", paymentID)
}

func savePayment(ctx workflow.Context, paymentID string, m metadata.Metadata) (*shared.Payment, error) {
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
		Metadata:  collectionutils.ConvertMap(m, collectionutils.ToAny[string]),
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

func runPaymentToAccount(ctx workflow.Context, source *PaymentSource, destination *LedgerAccountDestination, amount *shared.Monetary, m metadata.Metadata) error {
	payment, err := savePayment(ctx, source.ID, m)
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
	}, destination, amount, m)
}

func runWalletToWallet(ctx workflow.Context, source *WalletSource, destination *WalletDestination, amount *shared.Monetary, m metadata.Metadata) error {
	if amount == nil {
		return errors.New("amount must be specified")
	}
	sourceWallet, err := getWalletFromReference(ctx, source.WalletReference)
	if err != nil {
		return err
	}
	destinationWallet, err := getWalletFromReference(ctx, destination.WalletReference)
	if err != nil {
		return err
	}
	if sourceWallet.Ledger == destinationWallet.Ledger {
		mainBalance := "main"
		sourceSubject := shared.WalletSubject{
			Balance:    &mainBalance,
			Identifier: sourceWallet.ID,
			Type:       "WALLET",
		}
		return activities.CreditWallet(internal.InfiniteRetryContext(ctx), destinationWallet.ID, &shared.CreditWalletRequest{
			Amount:   *amount,
			Balance:  &destination.Balance,
			Metadata: m,
			Sources:  []shared.Subject{{WalletSubject: &sourceSubject}},
		})
	}

	if err := justError(activities.DebitWallet(internal.InfiniteRetryContext(ctx), sourceWallet.ID, &shared.DebitWalletRequest{
		Amount:   *amount,
		Balances: []string{source.Balance},
		Metadata: collectionutils.MergeMaps(m, map[string]string{
			moveToLedgerMetadata: destinationWallet.Ledger,
		}),
	})); err != nil {
		return err
	}

	return activities.CreditWallet(internal.InfiniteRetryContext(ctx), destinationWallet.ID, &shared.CreditWalletRequest{
		Amount:  *amount,
		Balance: &destination.Balance,
		Metadata: collectionutils.MergeMaps(m, map[string]string{
			moveFromLedgerMetadata: sourceWallet.Ledger,
		}),
	})
}

func runWalletToPayment(ctx workflow.Context, source *WalletSource, destination *PaymentDestination, amount *shared.Monetary, m metadata.Metadata) error {
	if amount == nil {
		return errors.New("amount must be specified")
	}
	if destination.PSP != "stripe" {
		return errors.New("only stripe actually supported")
	}
	sourceWallet, err := getWalletFromReference(ctx, source.WalletReference)
	if err != nil {
		return errors.Wrapf(err, "reading account: %s", source.ID)
	}

	formanceAccountID, err := extractFormanceAccountID(destination.Metadata, sourceWallet.Metadata)
	if err != nil {
		return err
	}

	if err := activities.StripeTransfer(internal.InfiniteRetryContext(ctx), activities.StripeTransferRequest{
		Amount:            amount.Amount,
		Asset:             &amount.Asset,
		Destination:       &formanceAccountID,
		WaitingValidation: &destination.WaitingValidation,
		ConnectorID:       destination.ConnectorID,
		Metadata:          m,
	}); err != nil {
		return err
	}

	return justError(activities.DebitWallet(internal.InfiniteRetryContext(ctx), sourceWallet.ID, &shared.DebitWalletRequest{
		Amount:   *amount,
		Balances: []string{source.Balance},
		Metadata: m,
	}))
}

func runWalletToAccount(ctx workflow.Context, source *WalletSource, destination *LedgerAccountDestination, amount *shared.Monetary, m metadata.Metadata) error {
	if amount == nil {
		return errors.New("amount must be specified")
	}
	sourceWallet, err := getWalletFromReference(ctx, source.WalletReference)
	if err != nil {
		return err
	}
	if sourceWallet.Ledger == destination.Ledger {
		return justError(activities.DebitWallet(internal.InfiniteRetryContext(ctx), sourceWallet.ID, &shared.DebitWalletRequest{
			Amount: *amount,
			Destination: &shared.Subject{
				LedgerAccountSubject: &shared.LedgerAccountSubject{
					Identifier: destination.ID,
					Type:       "ACCOUNT",
				},
			},
			Balances: []string{source.Balance},
			Metadata: m,
		}))
	}

	if err := justError(activities.DebitWallet(internal.InfiniteRetryContext(ctx), sourceWallet.ID, &shared.DebitWalletRequest{
		Amount:   *amount,
		Balances: []string{source.Balance},
		Metadata: collectionutils.MergeMaps(m, map[string]string{
			moveToLedgerMetadata: destination.Ledger,
		}),
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
		Metadata: collectionutils.MergeMaps(collectionutils.ConvertMap(m, collectionutils.ToAny[string]), map[string]any{
			moveFromLedgerMetadata: sourceWallet.Ledger,
		}),
	}))
}

func runAccountToWallet(ctx workflow.Context, source *LedgerAccountSource, destination *WalletDestination, amount *shared.Monetary, m metadata.Metadata) error {
	if amount == nil {
		return errors.New("amount must be specified")
	}
	destinationWallet, err := getWalletFromReference(ctx, destination.WalletReference)
	if err != nil {
		return err
	}
	if destinationWallet.Ledger == source.Ledger {
		return activities.CreditWallet(internal.InfiniteRetryContext(ctx), destinationWallet.ID, &shared.CreditWalletRequest{
			Amount: *amount,
			Sources: []shared.Subject{{
				LedgerAccountSubject: &shared.LedgerAccountSubject{
					Identifier: source.ID,
					Type:       "ACCOUNT",
				},
			}},
			Balance:  &destination.Balance,
			Metadata: m,
		})
	}

	if err := justError(activities.CreateTransaction(internal.InfiniteRetryContext(ctx), source.Ledger, shared.PostTransaction{
		Postings: []shared.Posting{{
			Amount:      amount.Amount,
			Asset:       amount.Asset,
			Destination: "world",
			Source:      source.ID,
		}},
		Metadata: collectionutils.MergeMaps(
			collectionutils.ConvertMap(m, collectionutils.ToAny[string]),
			map[string]any{
				moveToLedgerMetadata: destinationWallet.Ledger,
			},
		),
	})); err != nil {
		return err
	}

	return activities.CreditWallet(internal.InfiniteRetryContext(ctx), destinationWallet.ID, &shared.CreditWalletRequest{
		Amount: *amount,
		Sources: []shared.Subject{{
			LedgerAccountSubject: &shared.LedgerAccountSubject{
				Identifier: "world",
				Type:       "ACCOUNT",
			},
		}},
		Balance: &destination.Balance,
		Metadata: collectionutils.MergeMaps(m, map[string]string{
			moveFromLedgerMetadata: source.Ledger,
		}),
	})
}

func runAccountToAccount(ctx workflow.Context, source *LedgerAccountSource, destination *LedgerAccountDestination, amount *shared.Monetary, m metadata.Metadata) error {
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
			Metadata: collectionutils.ConvertMap(m, collectionutils.ToAny[string]),
		}))
	}
	if err := justError(activities.CreateTransaction(internal.InfiniteRetryContext(ctx), source.Ledger, shared.PostTransaction{
		Postings: []shared.Posting{{
			Amount:      amount.Amount,
			Asset:       amount.Asset,
			Destination: "world",
			Source:      source.ID,
		}},
		Metadata: collectionutils.MergeMaps(
			collectionutils.ConvertMap(m, collectionutils.ToAny[string]),
			map[string]any{
				moveToLedgerMetadata: destination.Ledger,
			},
		),
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
		Metadata: collectionutils.MergeMaps(
			collectionutils.ConvertMap(m, collectionutils.ToAny[string]),
			map[string]any{
				moveFromLedgerMetadata: source.Ledger,
			},
		),
	}))
}

func runAccountToPayment(ctx workflow.Context, source *LedgerAccountSource, destination *PaymentDestination, amount *shared.Monetary, m metadata.Metadata) error {
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

	if err := activities.StripeTransfer(internal.InfiniteRetryContext(ctx), activities.StripeTransferRequest{
		Amount:            amount.Amount,
		Asset:             &amount.Asset,
		Destination:       &formanceAccountID,
		WaitingValidation: &destination.WaitingValidation,
		ConnectorID:       destination.ConnectorID,
		Metadata:          m,
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
		Metadata: collectionutils.ConvertMap(m, collectionutils.ToAny[string]),
	}))
}
