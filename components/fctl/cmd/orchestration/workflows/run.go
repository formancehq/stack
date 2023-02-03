package workflows

import (
	"strings"

	"github.com/formancehq/fctl/cmd/orchestration/internal"
	fctl "github.com/formancehq/fctl/pkg"
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

			res, _, err := client.OrchestrationApi.
				RunWorkflow(cmd.Context(), args[0]).
				RequestBody(variables).
				Wait(wait).
				Execute()
			if err != nil {
				return errors.Wrap(err, "running workflow")
			}

			pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Workflow occurrence created with ID: %s", res.Data.Id)
			if wait {
				w, _, err := client.OrchestrationApi.GetWorkflow(cmd.Context(), args[0]).Execute()
				if err != nil {
					panic(err)
				}

				return internal.PrintWorkflowInstance(cmd.OutOrStdout(), w.Data, res.Data)
			}

			return nil
		}),
	)
}
