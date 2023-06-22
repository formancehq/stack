package instances

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type InstancesStopStore struct {
	Success    bool   `json:"success"`
	InstanceID string `json:"instance_id"`
}
type InstancesStopController struct {
	store *InstancesStopStore
}

var _ fctl.Controller[*InstancesStopStore] = (*InstancesStopController)(nil)

func NewDefaultInstancesStopStore() *InstancesStopStore {
	return &InstancesStopStore{}
}

func NewInstancesStopController() *InstancesStopController {
	return &InstancesStopController{
		store: NewDefaultInstancesStopStore(),
	}
}

func NewStopCommand() *cobra.Command {
	return fctl.NewCommand("stop <instance-id>",
		fctl.WithShortDescription("Stop a specific workflow instance"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*InstancesStopStore](NewInstancesStopController()),
	)
}

func (c *InstancesStopController) GetStore() *InstancesStopStore {
	return c.store
}

func (c *InstancesStopController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	soc, err := fctl.GetStackOrganizationConfig(cmd)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	response, err := client.Orchestration.CancelEvent(cmd.Context(), operations.CancelEventRequest{
		InstanceID: args[0],
	})
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("%s: %s", response.Error.ErrorCode, response.Error.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.InstanceID = args[0]

	return c, nil
}

func (c *InstancesStopController) Render(cmd *cobra.Command, args []string) error {
	// Print the instance information
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Workflow Instance with ID: %s successfully canceled ", c.store.InstanceID)

	return nil
}
