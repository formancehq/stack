package internal

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	formance "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	walletNameFlag = "name"
	walletIDFlag   = "id"
)

var (
	ErrUndefinedName = errors.New("missing wallet name")
)

func WithTargetingWalletByName() fctl.CommandOptionFn {
	return fctl.WithStringFlag(walletNameFlag, "", "Wallet name to use")
}

func WithTargetingWalletByID() fctl.CommandOptionFn {
	return fctl.WithStringFlag(walletIDFlag, "", "Wallet ID to use")
}

func DiscoverWalletIDFromName(cmd *cobra.Command, client *formance.Formance, walletName string) (string, error) {
	request := operations.ListWalletsRequest{
		Name: &walletName,
	}
	wallets, err := client.Wallets.ListWallets(cmd.Context(), request)
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

func RetrieveWalletIDFromName(cmd *cobra.Command, client *formance.Formance) (string, error) {
	walletName := fctl.GetString(cmd, walletNameFlag)
	if walletName == "" {
		return "", ErrUndefinedName
	}
	return DiscoverWalletIDFromName(cmd, client, walletName)
}

func RetrieveWalletID(cmd *cobra.Command, client *formance.Formance) (string, error) {
	walletID, err := RetrieveWalletIDFromName(cmd, client)
	if err != nil && err != ErrUndefinedName {
		return "", err
	}
	if err == ErrUndefinedName {
		return fctl.GetString(cmd, walletIDFlag), nil
	}
	return walletID, nil
}

func RequireWalletID(cmd *cobra.Command, client *formance.Formance) (string, error) {
	walletID, err := RetrieveWalletID(cmd, client)
	if err != nil {
		return "", err
	}
	if walletID == "" {
		return "", errors.New("You need to specify wallet id using --id or --name flags")
	}
	return walletID, nil
}
