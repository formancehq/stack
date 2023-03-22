package workflows

import (
	"fmt"
	"strings"

	"github.com/formancehq/fctl/cmd/orchestration/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewRunCommand() *cobra.Command {
	const (
		variableFlag = "variable"
		waitFlag     = "wait"
	)
	return fctl.NewCommand("run <id>",
		fctl.WithShortDescription("Run a workflow"),
		fctl.WithAliases("r"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithBoolFlag(waitFlag, false, "Wait end of the run"),
		fctl.WithStringSliceFlag(variableFlag, []string{}, "Variable to pass to the workflow"),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return errors.Wrap(err, "retrieving config")
			}

			organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
			if err != nil {
				return err
			}

			stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
			if err != nil {
				return err
			}

			client, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return errors.Wrap(err, "creating stack client")
			}

			wait := fctl.GetBool(cmd, waitFlag)
			variables := make(map[string]string)
			for _, variable := range fctl.GetStringSlice(cmd, variableFlag) {
				parts := strings.SplitN(variable, "=", 2)
				if len(parts) != 2 {
					return errors.New("malformed flag: " + variable)
				}
				variables[parts[0]] = parts[1]
			}

			response, err := client.Orchestration.
				RunWorkflow(cmd.Context(), operations.RunWorkflowRequest{
					RequestBody: variables,
					Wait:        &wait,
					WorkflowID:  args[0],
				})
			if err != nil {
				return err
			}

			if response.Error != nil {
				return fmt.Errorf("%s: %s", response.Error.ErrorCode, response.Error.ErrorMessage)
			}

			if response.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Workflow instance created with ID: %s", response.RunWorkflowResponse.Data.ID)
			if wait {
				w, err := client.Orchestration.GetWorkflow(cmd.Context(), operations.GetWorkflowRequest{
					FlowID: args[0],
				})
				if err != nil {
					panic(err)
				}

				return internal.PrintWorkflowInstance(cmd.OutOrStdout(), w.GetWorkflowResponse.Data, response.RunWorkflowResponse.Data)
			}

			return nil
		}),
	)
}
