package internal

import (
	"context"
	"flag"
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
)

const (
	walletNameFlag = "name"
	walletIDFlag   = "id"
)

var (
	ErrUndefinedName = errors.New("missing wallet name")
)

func WithTargetingWalletByName(flag *flag.FlagSet) *flag.FlagSet {
	flag.String(walletNameFlag, "", "Wallet name to use")
	return flag
}

func WithTargetingWalletByID(flag *flag.FlagSet) *flag.FlagSet {
	flag.String(walletIDFlag, "", "Wallet ID to use")
	return flag
}

func DiscoverWalletIDFromName(flags *flag.FlagSet, ctx context.Context, client *formance.Formance, walletName string) (string, error) {
	request := operations.ListWalletsRequest{
		Name: &walletName,
	}
	wallets, err := client.Wallets.ListWallets(ctx, request)
	if err != nil {
		return "", errors.Wrap(err, "listing wallets to retrieve wallet by name")
	}

	if wallets.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status code: %d", wallets.StatusCode)
	}

	if len(wallets.ListWalletsResponse.Cursor.Data) > 1 {
		return "", fmt.Errorf("found multiple wallets with name: %s", walletName)
	}
	if len(wallets.ListWalletsResponse.Cursor.Data) == 0 {
		return "", fmt.Errorf("wallet with name '%s' not found", walletName)
	}
	return wallets.ListWalletsResponse.Cursor.Data[0].ID, nil
}

func RetrieveWalletIDFromName(flags *flag.FlagSet, ctx context.Context, client *formance.Formance) (string, error) {
	walletName := fctl.GetString(flags, walletNameFlag)
	if walletName == "" {
		return "", ErrUndefinedName
	}
	return DiscoverWalletIDFromName(flags, ctx, client, walletName)
}

func RetrieveWalletID(flags *flag.FlagSet, ctx context.Context, client *formance.Formance) (string, error) {
	walletID, err := RetrieveWalletIDFromName(flags, ctx, client)
	if err != nil && err != ErrUndefinedName {
		return "", err
	}
	if err == ErrUndefinedName {
		return fctl.GetString(flags, walletIDFlag), nil
	}
	return walletID, nil
}

func RequireWalletID(flags *flag.FlagSet, ctx context.Context, client *formance.Formance) (string, error) {
	walletID, err := RetrieveWalletID(flags, ctx, client)
	if err != nil {
		return "", err
	}
	if walletID == "" {
		return "", errors.New("You need to specify wallet id using --id or --name flags")
	}
	return walletID, nil
}
