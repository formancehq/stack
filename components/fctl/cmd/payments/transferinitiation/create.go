package transferinitiation

import (
	"encoding/json"
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type TransferInitiationCreateStore struct {
	TransferInitiationId string `json:"transferInitiationId"`
}
type TransferInitiationCreateController struct {
	store *TransferInitiationCreateStore
}

var _ fctl.Controller[*TransferInitiationCreateStore] = (*TransferInitiationCreateController)(nil)

func NewDefaultTransferInitiationCreateStore() *TransferInitiationCreateStore {
	return &TransferInitiationCreateStore{}
}

func NewTransferInitiationCreateController() *TransferInitiationCreateController {
	return &TransferInitiationCreateController{
		store: NewDefaultTransferInitiationCreateStore(),
	}
}

func NewCreateCommand() *cobra.Command {
	return fctl.NewCommand("create <file>|-",
		fctl.WithShortDescription("Create a transfer initiation"),
		fctl.WithAliases("cr", "c"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*TransferInitiationCreateStore](NewTransferInitiationCreateController()),
	)
}

func (c *TransferInitiationCreateController) GetStore() *TransferInitiationCreateStore {
	return c.store
}

func (c *TransferInitiationCreateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	soc, err := fctl.GetStackOrganizationConfig(cmd)
	if err != nil {
		return nil, err
	}
	client, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	script, err := fctl.ReadFile(cmd, soc.Stack, args[0])
	if err != nil {
		return nil, err
	}

	request := shared.TransferInitiationRequest{}
	if err := json.Unmarshal([]byte(script), &request); err != nil {
		return nil, err
	}

	//nolint:gosimple
	response, err := client.Payments.CreateTransferInitiation(cmd.Context(), request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.TransferInitiationId = response.TransferInitiationResponse.Data.ID

	return c, nil
}

func (c *TransferInitiationCreateController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Transfer Initiation created with ID: %s", c.store.TransferInitiationId)

	return nil
}
