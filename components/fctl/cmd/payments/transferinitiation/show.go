package transferinitiation

import (
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ShowStore struct {
	TransferInitiation *shared.TransferInitiation `json:"transferInitiation"`
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
	return fctl.NewCommand("get <transferID>",
		fctl.WithShortDescription("Get transfer initiation"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithAliases("sh", "s"),
		fctl.WithController[*ShowStore](NewShowController()),
	)
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	ledgerClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	response, err := ledgerClient.Payments.GetTransferInitiation(cmd.Context(), operations.GetTransferInitiationRequest{
		TransferID: args[0],
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.TransferInitiation = &response.TransferInitiationResponse.Data

	return c, nil
}

func (c *ShowController) Render(cmd *cobra.Command, args []string) error {
	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.TransferInitiation.ID})
	tableData = append(tableData, []string{pterm.LightCyan("CreatedAt"), c.store.TransferInitiation.CreatedAt.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("UpdatedAt"), c.store.TransferInitiation.UpdatedAt.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("ScheduledAt"), c.store.TransferInitiation.ScheduledAt.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("Description"), c.store.TransferInitiation.Description})
	tableData = append(tableData, []string{pterm.LightCyan("SourceAccountID"), c.store.TransferInitiation.SourceAccountID})
	tableData = append(tableData, []string{pterm.LightCyan("DestinationAccountID"), c.store.TransferInitiation.DestinationAccountID})
	tableData = append(tableData, []string{pterm.LightCyan("Provider"), string(c.store.TransferInitiation.Provider)})
	tableData = append(tableData, []string{pterm.LightCyan("Type"), string(c.store.TransferInitiation.Type)})
	tableData = append(tableData, []string{pterm.LightCyan("Amount"), fmt.Sprint(c.store.TransferInitiation.Amount)})
	tableData = append(tableData, []string{pterm.LightCyan("Asset"), c.store.TransferInitiation.Asset})
	tableData = append(tableData, []string{pterm.LightCyan("Status"), string(c.store.TransferInitiation.Status)})
	tableData = append(tableData, []string{pterm.LightCyan("Error"), c.store.TransferInitiation.Error})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return nil
}
