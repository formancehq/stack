package wallets

import (
	"fmt"
	"strconv"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	cobra "github.com/spf13/cobra"
)

type CreditWalletStore struct {
	Success bool `json:"success"`
}
type CreditWalletController struct {
	store        *CreditWalletStore
	metadataFlag string
	balanceFlag  string
	sourceFlag   string
}

var _ fctl.Controller[*CreditWalletStore] = (*CreditWalletController)(nil)

func NewDefaultCreditWalletStore() *CreditWalletStore {
	return &CreditWalletStore{}
}

func NewCreditWalletController() *CreditWalletController {
	return &CreditWalletController{
		store:        NewDefaultCreditWalletStore(),
		metadataFlag: "metadata",
		balanceFlag:  "balance",
		sourceFlag:   "source",
	}
}
func NewCreditWalletCommand() *cobra.Command {
	c := NewCreditWalletController()
	return fctl.NewCommand("credit <amount> <asset>",
		fctl.WithShortDescription("Credit a wallets"),
		fctl.WithAliases("cr"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithStringSliceFlag(c.metadataFlag, []string{""}, "Metadata to use"),
		fctl.WithStringFlag(c.balanceFlag, "", "Balance to credit"),
		fctl.WithStringSliceFlag(c.sourceFlag, []string{}, `Use --source account=<account> | --source wallet=id:<wallet-id>[/<balance>] | --source wallet=name:<wallet-name>[/<balance>]`),
		internal.WithTargetingWalletByName(),
		internal.WithTargetingWalletByID(),
		fctl.WithController[*CreditWalletStore](c),
	)
}

func (c *CreditWalletController) GetStore() *CreditWalletStore {
	return c.store
}

func (c *CreditWalletController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "reading config")
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to credit a wallets") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	amountStr := args[0]
	asset := args[1]
	walletID, err := internal.RetrieveWalletID(cmd, client)
	if err != nil {
		return nil, err
	}

	if walletID == "" {
		return nil, errors.New("You need to specify wallet id using --id or --name flags")
	}

	amount, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "parsing amount")
	}

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, c.metadataFlag))
	if err != nil {
		return nil, err
	}

	sources := make([]shared.Subject, 0)
	for _, sourceStr := range fctl.GetStringSlice(cmd, c.sourceFlag) {
		source, err := internal.ParseSubject(sourceStr, cmd, client)
		if err != nil {
			return nil, err
		}
		sources = append(sources, *source)
	}

	request := operations.CreditWalletRequest{
		ID: walletID,
		CreditWalletRequest: &shared.CreditWalletRequest{
			Amount: shared.Monetary{
				Asset:  asset,
				Amount: int64(amount),
			},
			Metadata: metadata,
			Sources:  sources,
			Balance:  formance.String(fctl.GetString(cmd, c.balanceFlag)),
		},
	}
	response, err := client.Wallets.CreditWallet(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "crediting wallet")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true

	return c, nil
}

func (c *CreditWalletController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Wallet credited successfully!")
	return nil
}
