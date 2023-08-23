package transferinitiation

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type TransferInitiationDeleteStore struct {
	TransferID string `json:"transferId"`
	Success    bool   `json:"success"`
}

type TransferInitiationDeleteController struct {
	store *TransferInitiationDeleteStore
}

var _ fctl.Controller[*TransferInitiationDeleteStore] = (*TransferInitiationDeleteController)(nil)

func NewDefaultTransferInitiationDeleteStore() *TransferInitiationDeleteStore {
	return &TransferInitiationDeleteStore{}
}

func NewTransferInitiationDeleteController() *TransferInitiationDeleteController {
	return &TransferInitiationDeleteController{
		store: NewDefaultTransferInitiationDeleteStore(),
	}
}
func NewDeleteCommand() *cobra.Command {
	return fctl.NewCommand("delete <transferID>",
		fctl.WithAliases("del", "d"),
		fctl.WithShortDescription("Delete a transfer Initiation"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*TransferInitiationDeleteStore](NewTransferInitiationDeleteController()),
	)
}

func (c *TransferInitiationDeleteController) GetStore() *TransferInitiationDeleteStore {
	return c.store
}

func (c *TransferInitiationDeleteController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
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
		return nil, errors.Wrap(err, "creating stack client")
	}

	response, err := client.Payments.DeleteTransferInitiation(
		cmd.Context(),
		operations.DeleteTransferInitiationRequest{
			TransferID: args[0],
		},
	)

	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.TransferID = args[0]
	c.store.Success = true

	return c, nil
}

func (c *TransferInitiationDeleteController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithShowLineNumber().Printfln("Transfer Initiation %s Deleted!", c.store.TransferID)
	return nil
}
