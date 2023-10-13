package transferinitiation

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type UpdateStatusStore struct {
	TransferID string `json:"transferId"`
	Status     string `json:"status"`
	Success    bool   `json:"success"`
}
type UpdateStatusController struct {
	store *UpdateStatusStore
}

var _ fctl.Controller[*UpdateStatusStore] = (*UpdateStatusController)(nil)

func NewUpdateStatusStore() *UpdateStatusStore {
	return &UpdateStatusStore{}
}

func NewUpdateStatusController() *UpdateStatusController {
	return &UpdateStatusController{
		store: NewUpdateStatusStore(),
	}
}

func NewUpdateStatusCommand() *cobra.Command {
	return fctl.NewCommand("update_status <transferID> <status>",
		fctl.WithShortDescription("Update the status of a transfer initiation"),
		fctl.WithAliases("u"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*UpdateStatusStore](NewUpdateStatusController()),
	)
}

func (c *UpdateStatusController) GetStore() *UpdateStatusStore {
	return c.store
}

func (c *UpdateStatusController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	soc, err := fctl.GetStackOrganizationConfig(cmd)
	if err != nil {
		return nil, err
	}
	client, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	//nolint:gosimple
	response, err := client.Payments.UdpateTransferInitiationStatus(cmd.Context(), operations.UdpateTransferInitiationStatusRequest{
		UpdateTransferInitiationStatusRequest: shared.UpdateTransferInitiationStatusRequest{
			Status: shared.UpdateTransferInitiationStatusRequestStatus(args[1]),
		},
		TransferID: args[0],
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.TransferID = args[0]
	c.store.Status = args[1]
	c.store.Success = true

	return c, nil
}

func (c *UpdateStatusController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Update Transfer Initiation status with ID: %s and status %s", c.store.TransferID, c.store.Status)

	return nil
}
