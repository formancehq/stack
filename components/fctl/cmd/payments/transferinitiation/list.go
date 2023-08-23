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

type TransferInitiationListStore struct {
	Cursor *shared.TransferInitiationsCursorCursor `json:"cursor"`
}

type TransferInitiationListController struct {
	store *TransferInitiationListStore
}

var _ fctl.Controller[*TransferInitiationListStore] = (*TransferInitiationListController)(nil)

func NewDefaultTransferInitiationListStore() *TransferInitiationListStore {
	return &TransferInitiationListStore{
		Cursor: &shared.TransferInitiationsCursorCursor{},
	}
}

func NewTransferInitiationListController() *TransferInitiationListController {
	return &TransferInitiationListController{
		store: NewDefaultTransferInitiationListStore(),
	}
}

func (c *TransferInitiationListController) GetStore() *TransferInitiationListStore {
	return c.store
}

func (c *TransferInitiationListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

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

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	response, err := client.Payments.ListTransferInitiations(
		cmd.Context(),
		operations.ListTransferInitiationsRequest{},
	)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Cursor = &response.TransferInitiationsCursor.Cursor

	return c, nil
}

func (c *TransferInitiationListController) Render(cmd *cobra.Command, args []string) error {
	tableData := fctl.Map(c.store.Cursor.Data, func(tf shared.TransferInitiation) []string {
		return []string{
			tf.ID,
			tf.CreatedAt.Format(time.RFC3339),
			tf.UpdatedAt.Format(time.RFC3339),
			tf.Description,
			tf.SourceAccountID,
			tf.DestinationAccountID,
			string(tf.Provider),
			string(tf.Type),
			fmt.Sprint(tf.Amount),
			tf.Asset,
			string(tf.Status),
			tf.Error,
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "CreatedAt", "UpdatedAt", "Description", "Source Account ID",
		"Destination Account ID", "Provider", "Type", "Amount", "Asset", "Status", "Error"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}

func NewListTransferInitiationCommand() *cobra.Command {
	return fctl.NewCommand("list_transfer_initiations",
		fctl.WithAliases("ls"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithShortDescription("List transfer initiation"),
		fctl.WithController[*TransferInitiationListStore](NewTransferInitiationListController()),
	)
}
