package wallets

import (
	"fmt"
	"math/big"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DebitWalletStore struct {
	HoldID  *string `json:"holdId"`
	Success bool    `json:"success"`
}
type DebitWalletController struct {
	store           *DebitWalletStore
	pendingFlag     string
	metadataFlag    string
	descriptionFlag string
	balanceFlag     string
	destinationFlag string
}

var _ fctl.Controller[*DebitWalletStore] = (*DebitWalletController)(nil)

func NewDefaultDebitWalletStore() *DebitWalletStore {
	return &DebitWalletStore{
		HoldID:  nil,
		Success: false,
	}
}

func NewDebitWalletController() *DebitWalletController {
	return &DebitWalletController{
		store:           NewDefaultDebitWalletStore(),
		pendingFlag:     "pending",
		metadataFlag:    "metadata",
		descriptionFlag: "description",
		balanceFlag:     "balance",
		destinationFlag: "destination",
	}
}

func NewDebitWalletCommand() *cobra.Command {
	c := NewDebitWalletController()
	return fctl.NewCommand("debit <amount> <asset>",
		fctl.WithShortDescription("Debit a wallet"),
		fctl.WithAliases("deb"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.RangeArgs(2, 3)),
		fctl.WithStringFlag(c.descriptionFlag, "", "Debit description"),
		fctl.WithBoolFlag(c.pendingFlag, false, "Create a pending debit"),
		fctl.WithStringSliceFlag(c.metadataFlag, []string{""}, "Metadata to use"),
		fctl.WithStringSliceFlag(c.balanceFlag, []string{""}, "Balance to debit"),
		fctl.WithStringFlag(c.destinationFlag, "",
			`Use --destination account=<account> | --destination wallet=id:<wallet-id>[/<balance>] | --destination wallet=name:<wallet-name>[/<balance>]`),
		internal.WithTargetingWalletByName(),
		internal.WithTargetingWalletByID(),
		fctl.WithController[*DebitWalletStore](c),
	)
}

func (c *DebitWalletController) GetStore() *DebitWalletStore {
	return c.store
}

func (c *DebitWalletController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())
	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to debit a wallets") {
		return nil, fctl.ErrMissingApproval
	}

	pending := fctl.GetBool(cmd, c.pendingFlag)

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, c.metadataFlag))
	if err != nil {
		return nil, err
	}

	amountStr := args[0]
	asset := args[1]
	walletID, err := internal.RequireWalletID(cmd, store.Client())
	if err != nil {
		return nil, err
	}

	description := fctl.GetString(cmd, c.descriptionFlag)

	amount, ok := big.NewInt(0).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse '%s' as big int", amountStr)
	}

	var destination *shared.Subject
	if destinationStr := fctl.GetString(cmd, c.destinationFlag); destinationStr != "" {
		destination, err = internal.ParseSubject(destinationStr, cmd, store.Client())
		if err != nil {
			return nil, err
		}
	}

	response, err := store.Client().Wallets.DebitWallet(cmd.Context(), operations.DebitWalletRequest{
		DebitWalletRequest: &shared.DebitWalletRequest{
			Amount: shared.Monetary{
				Asset:  asset,
				Amount: amount,
			},
			Pending:     &pending,
			Metadata:    metadata,
			Description: &description,
			Destination: destination,
			Balances:    fctl.GetStringSlice(cmd, c.balanceFlag),
		},
		ID: walletID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "debiting wallet")
	}

	if response.DebitWalletResponse != nil {
		c.store.HoldID = &response.DebitWalletResponse.Data.ID
	}

	c.store.Success = true

	return c, nil
}

func (c *DebitWalletController) Render(cmd *cobra.Command, args []string) error {
	if c.store.HoldID != nil && *c.store.HoldID != "" {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Wallet debited successfully with hold id '%s'!", *c.store.HoldID)
	} else {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Wallet debited successfully!")
	}

	return nil

}
