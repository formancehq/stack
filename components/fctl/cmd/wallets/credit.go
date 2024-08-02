package wallets

import (
	"fmt"
	"math/big"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	formance "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type CreditWalletStore struct {
	Success bool `json:"success"`
}
type CreditWalletController struct {
	store        *CreditWalletStore
	metadataFlag string
	balanceFlag  string
	sourceFlag   string
	ikFlag       string
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
		ikFlag:       "ik",
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
		fctl.WithStringFlag(c.ikFlag, "", "Idempotency Key"),
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
	store := fctl.GetStackStore(cmd.Context())
	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to credit a wallets") {
		return nil, fctl.ErrMissingApproval
	}

	amountStr := args[0]
	asset := args[1]
	walletID, err := internal.RetrieveWalletID(cmd, store.Client())
	if err != nil {
		return nil, err
	}

	if walletID == "" {
		return nil, errors.New("You need to specify wallet id using --id or --name flags")
	}

	amount, ok := big.NewInt(0).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse '%s' as big int", amountStr)
	}

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, c.metadataFlag))
	if err != nil {
		return nil, err
	}

	sources := make([]shared.Subject, 0)
	for _, sourceStr := range fctl.GetStringSlice(cmd, c.sourceFlag) {
		source, err := internal.ParseSubject(sourceStr, cmd, store.Client())
		if err != nil {
			return nil, err
		}
		sources = append(sources, *source)
	}

	request := operations.CreditWalletRequest{
		IdempotencyKey: fctl.Ptr(fctl.GetString(cmd, c.ikFlag)),
		ID:             walletID,
		CreditWalletRequest: &shared.CreditWalletRequest{
			Amount: shared.Monetary{
				Asset:  asset,
				Amount: amount,
			},
			Metadata: metadata,
			Sources:  sources,
			Balance:  formance.String(fctl.GetString(cmd, c.balanceFlag)),
		},
	}
	_, err = store.Client().Wallets.CreditWallet(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "crediting wallet")
	}

	c.store.Success = true

	return c, nil
}

func (c *CreditWalletController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Wallet credited successfully!")
	return nil
}
