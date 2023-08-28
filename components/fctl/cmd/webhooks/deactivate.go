package webhooks

import (
	"flag"
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDesactivate   = "deactivate <config-id>"
	shortDesactivate = "Deactivate one config"
)

type DesactivateStore struct {
	Success bool `json:"success"`
}

func NewDesactivateStore() *DesactivateStore {
	return &DesactivateStore{
		Success: true,
	}
}

func NewDesactivateConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useDesactivate, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)

	return fctl.NewControllerConfig(
		useDesactivate,
		shortDesactivate,
		shortDesactivate,
		[]string{
			"deactivate",
			"des",
		},
		flags,
	)
}

var _ fctl.Controller[*DesactivateStore] = (*DesactivateController)(nil)

type DesactivateController struct {
	store  *DesactivateStore
	config *fctl.ControllerConfig
}

func NewDesactivateController(config *fctl.ControllerConfig) *DesactivateController {
	return &DesactivateController{
		store:  NewDesactivateStore(),
		config: config,
	}
}

func (c *DesactivateController) GetStore() *DesactivateStore {
	return c.store
}

func (c *DesactivateController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *DesactivateController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	out := c.config.GetOut()
	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "fctl.GetConfig")
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are about to deactivate a webhook") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	if len(c.config.GetArgs()) < 1 {
		return nil, fmt.Errorf("missing config id")
	}

	request := operations.DeactivateConfigRequest{
		ID: c.config.GetArgs()[0],
	}
	response, err := client.Webhooks.DeactivateConfig(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "deactivating config")
	}

	c.store.Success = !response.ConfigResponse.Data.Active

	// Check if there is an error
	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	// Check if the status code is >= 300
	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return c, nil
}

func (c *DesactivateController) Render() error {

	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Config deactivated successfully")

	return nil
}

func NewDeactivateCommand() *cobra.Command {

	config := NewDesactivateConfig()

	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*DesactivateStore](NewDesactivateController(config)),
	)
}
