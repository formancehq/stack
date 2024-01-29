package payments

import (
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ShowStore struct {
	Payment *shared.Payment `json:"payment"`
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
	c := NewShowController()
	return fctl.NewCommand("get <paymentID>",
		fctl.WithShortDescription("Get payment"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithAliases("sh", "s"),
		fctl.WithController[*ShowStore](c),
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

	response, err := ledgerClient.Payments.GetPayment(cmd.Context(), operations.GetPaymentRequest{
		PaymentID: args[0],
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Payment = &response.PaymentResponse.Data

	return c, nil
}

func (c *ShowController) Render(cmd *cobra.Command, args []string) error {
	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.Payment.ID})
	tableData = append(tableData, []string{pterm.LightCyan("Reference"), c.store.Payment.Reference})
	tableData = append(tableData, []string{pterm.LightCyan("CreatedAt"), c.store.Payment.CreatedAt.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("ConnectorID"), c.store.Payment.ConnectorID})
	tableData = append(tableData, []string{pterm.LightCyan("Asset"), c.store.Payment.Asset})
	tableData = append(tableData, []string{pterm.LightCyan("Amount"), c.store.Payment.Amount.String()})
	tableData = append(tableData, []string{pterm.LightCyan("InitialAmount"), c.store.Payment.InitialAmount.String()})
	tableData = append(tableData, []string{pterm.LightCyan("Type"), string(c.store.Payment.Type)})
	tableData = append(tableData, []string{pterm.LightCyan("Scheme"), string(c.store.Payment.Scheme)})
	tableData = append(tableData, []string{pterm.LightCyan("Status"), string(c.store.Payment.Status)})
	tableData = append(tableData, []string{pterm.LightCyan("DestinationAccountID"), c.store.Payment.DestinationAccountID})
	tableData = append(tableData, []string{pterm.LightCyan("SourceAccountID"), c.store.Payment.SourceAccountID})
	tableData = append(tableData, []string{pterm.LightCyan("Provider"), func() string {
		if c.store.Payment.Provider == nil {
			return ""
		}
		return string(*c.store.Payment.Provider)
	}()})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	tableData = fctl.Map(c.store.Payment.Adjustments, func(pa shared.PaymentAdjustment) []string {
		return []string{
			pa.Reference,
			string(pa.Status),
			pa.CreatedAt.Format(time.RFC3339),
			pa.Amount.String(),
		}
	})

	tableData = fctl.Prepend(tableData, []string{"Reference", "Status", "CreatedAt", "Amount"})
	if err := pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return fctl.PrintMetadata(cmd.OutOrStdout(), c.store.Payment.Metadata)
}
