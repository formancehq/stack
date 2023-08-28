package secrets

import (
	"flag"
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDelete   = "delete <client-id> <secret-id>"
	shortDelete = "Delete secret specific secret for a client"
)

type DeleteStore struct {
	Success  bool   `json:"success"`
	SecretId string `json:"secretId"`
}

func NewDeleteStore() *DeleteStore {
	return &DeleteStore{}
}

func NewDeleteConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useDelete, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)
	return fctl.NewControllerConfig(
		useDelete,
		shortDelete,
		shortDelete,
		[]string{
			"d",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

var _ fctl.Controller[*DeleteStore] = (*DeleteController)(nil)

type DeleteController struct {
	store  *DeleteStore
	config *fctl.ControllerConfig
}

func NewDeleteController(config *fctl.ControllerConfig) *DeleteController {
	return &DeleteController{
		store:  NewDeleteStore(),
		config: config,
	}
}

func (c *DeleteController) GetStore() *DeleteStore {
	return c.store
}

func (c *DeleteController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *DeleteController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()
	out := c.config.GetOut()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are about to delete a client secret") {
		return nil, fctl.ErrMissingApproval
	}

	authClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	request := operations.DeleteSecretRequest{
		ClientID: args[0],
		SecretID: args[1],
	}
	response, err := authClient.Auth.DeleteSecret(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.SecretId = args[1]
	c.store.Success = true

	return c, nil
}

func (c *DeleteController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Secret %s successfully deleted!", c.store.SecretId)

	return nil

}

func NewDeleteCommand() *cobra.Command {
	config := NewDeleteConfig()
	return fctl.NewCommand("delete <client-id> <secret-id>",
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*DeleteStore](NewDeleteController(config)),
	)
}
