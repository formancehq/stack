package transferinitiation

import (
	"fmt"
	"time"

	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ShowStore struct {
	TransferInitiation *shared.TransferInitiation `json:"transferInitiation"`
}
type ShowController struct {
	PaymentsVersion versions.Version

	store *ShowStore
}

func (c *ShowController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
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
	c := NewShowController()
	return fctl.NewCommand("get <transferID>",
		fctl.WithShortDescription("Get transfer initiation"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithAliases("sh", "s"),
		fctl.WithController[*ShowStore](c),
	)
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())
	if err := versions.GetPaymentsVersion(cmd, args, c); err != nil {
		return nil, err
	}

	if c.PaymentsVersion < versions.V1 {
		return nil, fmt.Errorf("transfer initiation are only supported in >= v1.0.0")
	}

	response, err := store.Client().Payments.GetTransferInitiation(cmd.Context(), operations.GetTransferInitiationRequest{
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
	tableData = append(tableData, []string{pterm.LightCyan("Reference"), c.store.TransferInitiation.Reference})
	tableData = append(tableData, []string{pterm.LightCyan("CreatedAt"), c.store.TransferInitiation.CreatedAt.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("ScheduledAt"), c.store.TransferInitiation.ScheduledAt.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("Description"), c.store.TransferInitiation.Description})
	tableData = append(tableData, []string{pterm.LightCyan("SourceAccountID"), c.store.TransferInitiation.SourceAccountID})
	tableData = append(tableData, []string{pterm.LightCyan("DestinationAccountID"), c.store.TransferInitiation.DestinationAccountID})
	tableData = append(tableData, []string{pterm.LightCyan("ConnectorID"), string(c.store.TransferInitiation.ConnectorID)})
	tableData = append(tableData, []string{pterm.LightCyan("Type"), string(c.store.TransferInitiation.Type)})
	tableData = append(tableData, []string{pterm.LightCyan("Amount"), fmt.Sprint(c.store.TransferInitiation.Amount)})
	tableData = append(tableData, []string{pterm.LightCyan("InitialAmount"), fmt.Sprint(c.store.TransferInitiation.InitialAmount)})
	tableData = append(tableData, []string{pterm.LightCyan("Asset"), c.store.TransferInitiation.Asset})
	tableData = append(tableData, []string{pterm.LightCyan("Status"), string(c.store.TransferInitiation.Status)})
	tableData = append(tableData, []string{pterm.LightCyan("Error"), c.store.TransferInitiation.Error})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	tableData = fctl.Map(c.store.TransferInitiation.RelatedPayments, func(tf shared.TransferInitiationPayments) []string {
		return []string{
			tf.PaymentID,
			tf.CreatedAt.Format(time.RFC3339),
			string(tf.Status),
			tf.Error,
		}
	})
	tableData = fctl.Prepend(tableData, []string{"PaymentID", "CreatedAt", "Status", "Error"})
	if err := pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	tableData = fctl.Map(c.store.TransferInitiation.RelatedAdjustments, func(tf shared.TransferInitiationAdjusments) []string {
		return []string{
			tf.AdjustmentID,
			tf.CreatedAt.Format(time.RFC3339),
			string(tf.Status),
			tf.Error,
		}
	})
	tableData = fctl.Prepend(tableData, []string{"AdjustmentID", "CreatedAt", "Status", "Error"})
	if err := pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return fctl.PrintMetadata(cmd.OutOrStdout(), c.store.TransferInitiation.Metadata)
}
