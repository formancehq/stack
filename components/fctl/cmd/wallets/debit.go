package wallets

import (
	"fmt"
	"strconv"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DebitWalletStore struct {
	HoldID  *string `json:"hold_id"`
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

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to debit a wallets") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	pending := fctl.GetBool(cmd, c.pendingFlag)

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, c.metadataFlag))
	if err != nil {
		return nil, err
	}

	amountStr := args[0]
	asset := args[1]
	walletID, err := internal.RequireWalletID(cmd, client)
	if err != nil {
		return nil, err
	}

	description := fctl.GetString(cmd, c.descriptionFlag)

	amount, err := strconv.ParseInt(amountStr, 10, 32)
	if err != nil {
		return nil, errors.Wrap(err, "parsing amount")
	}

	var destination *shared.Subject
	if destinationStr := fctl.GetString(cmd, c.destinationFlag); destinationStr != "" {
		destination, err = internal.ParseSubject(destinationStr, cmd, client)
		if err != nil {
			return nil, err
		}
	}

	response, err := client.Wallets.DebitWallet(cmd.Context(), operations.DebitWalletRequest{
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

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
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
