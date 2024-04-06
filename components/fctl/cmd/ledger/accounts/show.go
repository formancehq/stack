package accounts

import (
	"fmt"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"

	internal "github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ShowStore struct {
	Account *shared.AccountWithVolumesAndBalances `json:"account"`
}
type ShowController struct {
	store *ShowStore
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

func NewDefaultShowStore() *ShowStore {
	return &ShowStore{}
}

func NewShowController() *ShowController {
	return &ShowController{
		store: NewDefaultShowStore(),
	}
}

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show <address>",
		fctl.WithShortDescription("Show account"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithAliases("sh", "s"),
		fctl.WithController[*ShowStore](NewShowController()),
	)
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := fctl.GetStackStore(cmd.Context())

	ledger := fctl.GetString(cmd, internal.LedgerFlag)
	response, err := store.Client().Ledger.GetAccount(cmd.Context(), operations.GetAccountRequest{
		Address: args[0],
		Ledger:  ledger,
	})
	if err != nil {
		return nil, err
	}

	c.store.Account = &response.AccountResponse.Data

	return c, nil
}

func (c *ShowController) Render(cmd *cobra.Command, args []string) error {
	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
	if c.store.Account.Volumes != nil && len(c.store.Account.Volumes) > 0 {
		tableData := pterm.TableData{}
		tableData = append(tableData, []string{"Asset", "Input", "Output"})
		for asset, volumes := range c.store.Account.Volumes {
			input := volumes.Input
			output := volumes.Output
			tableData = append(tableData, []string{pterm.LightCyan(asset), fmt.Sprint(input), fmt.Sprint(output)})
		}
		if err := pterm.DefaultTable.
			WithHasHeader(true).
			WithWriter(cmd.OutOrStdout()).
			WithData(tableData).
			Render(); err != nil {
			return err
		}
	} else {
		fctl.Println("No balances.")
	}

	fmt.Fprintln(cmd.OutOrStdout())

	return fctl.PrintMetadata(cmd.OutOrStdout(),
		collectionutils.ConvertMap(c.store.Account.Metadata, collectionutils.ToFmtString))
}
