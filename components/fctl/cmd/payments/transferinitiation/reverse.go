package transferinitiation

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ReverseStore struct {
	TransferID string `json:"transferId"`
	Success    bool   `json:"success"`
}

type ReverseController struct {
	PaymentsVersion versions.Version

	store *ReverseStore
}

func (c *ReverseController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
}

var _ fctl.Controller[*ReverseStore] = (*ReverseController)(nil)

func NewReverseStore() *ReverseStore {
	return &ReverseStore{}
}

func NewReverseController() *ReverseController {
	return &ReverseController{
		store: NewReverseStore(),
	}
}
func NewReverseCommand() *cobra.Command {
	c := NewReverseController()
	return fctl.NewCommand("reverse <transferID> <file>|-",
		fctl.WithConfirmFlag(),
		fctl.WithAliases("re", "r"),
		fctl.WithShortDescription("Reverse a transfer Initiation"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*ReverseStore](c),
	)
}

func (c *ReverseController) GetStore() *ReverseStore {
	return c.store
}

func (c *ReverseController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	if err := versions.GetPaymentsVersion(cmd, args, c); err != nil {
		return nil, err
	}

	if c.PaymentsVersion < versions.V1 {
		return nil, fmt.Errorf("transfer initiation are only supported in >= v1.0.0")
	}

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

	script, err := fctl.ReadFile(cmd, stack, args[1])
	if err != nil {
		return nil, err
	}

	request := shared.ReverseTransferInitiationRequest{}
	if err := json.Unmarshal([]byte(script), &request); err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to delete '%s'", args[0]) {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	response, err := client.Payments.ReverseTransferInitiation(
		cmd.Context(),
		operations.ReverseTransferInitiationRequest{
			TransferID:                       args[0],
			ReverseTransferInitiationRequest: request,
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

func (c *ReverseController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Transfer Initiation %s reversed!", c.store.TransferID)
	return nil
}
