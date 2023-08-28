package webhooks

import (
	"flag"
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDelete              = "delete <config-id>"
	shortDescriptionDelete = "Delete a config"
)

type DeleteStore struct {
	ErrorResponse *shared.ErrorResponse `json:"error"`
	Success       bool                  `json:"success"`
}

func NewDeleteStore() *DeleteStore {
	return &DeleteStore{
		Success: true,
	}
}

func NewDeleteConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useDelete, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)

	return fctl.NewControllerConfig(
		useDelete,
		shortDescriptionDelete,
		shortDescriptionDelete,
		[]string{
			"delete",
			"del",
		},
		flags,
	)

}

var _ fctl.Controller[*DeleteStore] = (*DeleteController)(nil)

type DeleteController struct {
	store      *DeleteStore
	config     *fctl.ControllerConfig
	fctlConfig *fctl.Config
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
	out := c.config.GetOut()
	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "fctl.GetConfig")
	}
	c.fctlConfig = cfg

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are about to delete a webhook") {
		return nil, fctl.ErrMissingApproval
	}

	webhookClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	if len(c.config.GetArgs()) == 0 {
		return nil, errors.New("missing config id")
	}

	request := operations.DeleteConfigRequest{
		ID: c.config.GetArgs()[0],
	}
	response, err := webhookClient.Webhooks.DeleteConfig(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "deleting config")
	}

	if response.ErrorResponse != nil {
		if response.ErrorResponse.ErrorCode == "NOT_FOUND" {
			c.store.ErrorResponse = response.ErrorResponse
			c.store.Success = false
			return c, nil
		}
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = response.StatusCode == 200

	return c, nil
}

func (c *DeleteController) Render() error {
	if !c.store.Success {
		pterm.Warning.WithShowLineNumber(false).Printfln("Config %s not found", c.config.GetArgs()[0])
		return nil
	}

	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.ErrorMessage)
		return nil
	}

	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Config deleted successfully")

	return nil
}

func NewDeleteCommand() *cobra.Command {
	config := NewDeleteConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*DeleteStore](NewDeleteController(config)),
	)
}
