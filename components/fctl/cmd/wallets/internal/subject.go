package internal

import (
	"context"
	"flag"
	"strings"

	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
)

func ParseSubject(subject string, flags *flag.FlagSet, ctx context.Context, client *formance.Formance) (*shared.Subject, error) {
	var err error
	switch {
	case strings.HasPrefix(subject, "wallet="):
		walletDefinition := strings.TrimPrefix(subject, "wallet=")
		parts := strings.SplitN(walletDefinition, "/", 2)
		balance := "main"
		if len(parts) > 1 {
			balance = parts[1]
		}

		var walletID string
		switch {
		case strings.HasPrefix(walletDefinition, "id:"):
			walletID = strings.TrimPrefix(parts[0], "id:")
		case strings.HasPrefix(walletDefinition, "name:"):
			walletID, err = DiscoverWalletIDFromName(flags, ctx, client, strings.TrimPrefix(parts[0], "name:"))
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("malformed wallet source definition")
		}
		subject := shared.CreateSubjectWallet(shared.WalletSubject{
			Identifier: walletID,
			Balance:    &balance,
		})
		return &subject, nil
	case strings.HasPrefix(subject, "account="):
		subject := shared.CreateSubjectAccount(shared.LedgerAccountSubject{
			Identifier: strings.TrimPrefix(subject, "account="),
		})
		return &subject, nil
	default:
		return nil, errors.New("malformed source definition")
	}
}
