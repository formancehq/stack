package accounts

import (
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ShowStore struct {
	Account *shared.PaymentsAccount `json:"account"`
}
type ShowController struct {
	store *ShowStore
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

func NewShowStore() *ShowStore {
	return &ShowStore{}
}

func NewShowController() *ShowController {
	return &ShowController{
		store: NewShowStore(),
	}
}

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("get <accountID>",
		fctl.WithShortDescription("Get account"),
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

	response, err := store.Client().Payments.V1.PaymentsgetAccount(cmd.Context(), operations.PaymentsgetAccountRequest{
		AccountID: args[0],
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Account = &response.PaymentsAccountResponse.Data

	return c, nil
}

func (c *ShowController) Render(cmd *cobra.Command, args []string) error {
	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.Account.ID})
	tableData = append(tableData, []string{pterm.LightCyan("AccountName"), c.store.Account.AccountName})
	tableData = append(tableData, []string{pterm.LightCyan("CreatedAt"), c.store.Account.CreatedAt.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("ConnectorID"), string(c.store.Account.ConnectorID)})
	tableData = append(tableData, []string{pterm.LightCyan("DefaultAsset"), c.store.Account.DefaultAsset})
	tableData = append(tableData, []string{pterm.LightCyan("DefaultCurrency"), c.store.Account.DefaultCurrency})
	tableData = append(tableData, []string{pterm.LightCyan("Reference"), c.store.Account.Reference})
	tableData = append(tableData, []string{pterm.LightCyan("Type"), string(c.store.Account.Type)})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return fctl.PrintMetadata(cmd.OutOrStdout(), c.store.Account.Metadata)
}
