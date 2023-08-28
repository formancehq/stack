package wallets

import (
	"flag"
	"fmt"
	"math/big"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDebit        = "debit <amount> <asset>"
	shortDebit      = "Debit a wallet"
	pendingFlag     = "pending"
	descriptionFlag = "description"
	balanceFlag     = "balance"
	destinationFlag = "destination"
)

type DebitWalletStore struct {
	HoldID  *string `json:"holdId"`
	Success bool    `json:"success"`
}
type DebitController struct {
	store  *DebitWalletStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*DebitWalletStore] = (*DebitController)(nil)

func NewDebitStore() *DebitWalletStore {
	return &DebitWalletStore{
		HoldID:  nil,
		Success: false,
	}
}
func NewDebitConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useDebit, flag.ExitOnError)
	flags.String(descriptionFlag, "", "Debit description")
	flags.String(pendingFlag, "", "Create a pending debit")
	flags.String(balanceFlag, "", "Balance to debit")
	flags.String(destinationFlag, "",
		`Use --destination account=<account> | --destination wallet=id:<wallet-id>[/<balance>] | --destination wallet=name:<wallet-name>[/<balance>]`)
	fctl.WithMetadataFlag(flags)
	fctl.WithConfirmFlag(flags)

	internal.WithTargetingWalletByName(flags)
	internal.WithTargetingWalletByID(flags)

	c := fctl.NewControllerConfig(
		useDebit,
		shortDebit,
		shortDebit,
		[]string{
			"deb",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)

	return c
}
func NewDebitController(config *fctl.ControllerConfig) *DebitController {
	return &DebitController{
		store:  NewDebitStore(),
		config: config,
	}
}

func (c *DebitController) GetStore() *DebitWalletStore {
	return c.store
}

func (c *DebitController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *DebitController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	cfg, err := fctl.GetConfig(flags)
	out := c.config.GetOut()
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are about to debit a wallets") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	pending := fctl.GetBool(flags, pendingFlag)

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(flags, fctl.MetadataFlag))
	if err != nil {
		return nil, err
	}

	if len(c.config.GetArgs()) < 2 {
		return nil, errors.New("missing amount and asset")
	}

	amountStr := c.config.GetArgs()[0]
	asset := c.config.GetArgs()[1]
	walletID, err := internal.RequireWalletID(flags, ctx, client)
	if err != nil {
		return nil, err
	}

	description := fctl.GetString(flags, descriptionFlag)

	amount, ok := big.NewInt(0).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse '%s' as big int", amountStr)
	}

	var destination *shared.Subject
	if destinationStr := fctl.GetString(flags, destinationFlag); destinationStr != "" {
		destination, err = internal.ParseSubject(destinationStr, flags, ctx, client)
		if err != nil {
			return nil, err
		}
	}

	response, err := client.Wallets.DebitWallet(ctx, operations.DebitWalletRequest{
		DebitWalletRequest: &shared.DebitWalletRequest{
			Amount: shared.Monetary{
				Asset:  asset,
				Amount: amount,
			},
			Pending:     &pending,
			Metadata:    metadata,
			Description: &description,
			Destination: destination,
			Balances:    fctl.GetStringSlice(flags, balanceFlag),
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

func (c *DebitController) Render() error {
	if c.store.HoldID != nil && *c.store.HoldID != "" {
		pterm.Success.WithWriter(c.config.GetOut()).Printfln("Wallet debited successfully with hold id '%s'!", *c.store.HoldID)
	} else {
		pterm.Success.WithWriter(c.config.GetOut()).Printfln("Wallet debited successfully!")
	}

	return nil

}

func NewDebitWalletCommand() *cobra.Command {
	c := NewDebitConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.RangeArgs(2, 3)),
		fctl.WithController[*DebitWalletStore](NewDebitController(c)),
	)
}
