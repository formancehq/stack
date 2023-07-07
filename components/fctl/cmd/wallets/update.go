package wallets

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type UpdateStore struct {
	Success bool `json:"success"`
}
type UpdateController struct {
	store        *UpdateStore
	metadataFlag string
}

var _ fctl.Controller[*UpdateStore] = (*UpdateController)(nil)

func NewDefaultUpdateStore() *UpdateStore {
	return &UpdateStore{}
}

func NewUpdateController() *UpdateController {
	return &UpdateController{
		store:        NewDefaultUpdateStore(),
		metadataFlag: "metadata",
	}
}

func NewUpdateCommand() *cobra.Command {
	c := NewUpdateController()
	return fctl.NewCommand("update <wallet-id>",
		fctl.WithShortDescription("Update a wallets"),
		fctl.WithAliases("up"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithStringSliceFlag(c.metadataFlag, []string{""}, "Metadata to use"),
		fctl.WithController[*UpdateStore](c),
	)
}

func (c *UpdateController) GetStore() *UpdateStore {
	return c.store
}

func (c *UpdateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

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

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to update a wallets") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, c.metadataFlag))
	if err != nil {
		return nil, err
	}

	response, err := client.Wallets.UpdateWallet(cmd.Context(), operations.UpdateWalletRequest{
		RequestBody: &operations.UpdateWalletRequestBody{
			Metadata: metadata,
		},
		ID: args[0],
	})
	if err != nil {
		return nil, errors.Wrap(err, "updating wallet")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	return c, nil
}

func (c *UpdateController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Wallet updated successfully!")
	return nil
}
