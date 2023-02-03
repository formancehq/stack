package instances

import (
	"time"

	"github.com/formancehq/fctl/cmd/orchestration/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show <instance-id>",
		fctl.WithShortDescription("Show a specific workflow instance"),
		fctl.WithArgs(cobra.ExactArgs(1)),
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

			res, _, err := client.OrchestrationApi.GetInstance(cmd.Context(), args[0]).Execute()
			if err != nil {
				return errors.Wrap(err, "reading instance")
			}

			fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
			tableData := pterm.TableData{}
			tableData = append(tableData, []string{pterm.LightCyan("ID"), res.Data.Id})
			tableData = append(tableData, []string{pterm.LightCyan("Created at"), res.Data.CreatedAt.Format(time.RFC3339)})
			tableData = append(tableData, []string{pterm.LightCyan("Updated at"), res.Data.UpdatedAt.Format(time.RFC3339)})
			if res.Data.Terminated {
				tableData = append(tableData, []string{pterm.LightCyan("Terminated at"), res.Data.TerminatedAt.Format(time.RFC3339)})
			}

			if err := pterm.DefaultTable.
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render(); err != nil {
				return err
			}

			w, _, err := client.OrchestrationApi.GetWorkflow(cmd.Context(), res.Data.WorkflowID).Execute()
			if err != nil {
				return errors.Wrap(err, "reading workflow")
			}

			if err := internal.PrintWorkflowInstance(cmd.OutOrStdout(), w.Data, res.Data); err != nil {
				return err
			}

			return nil
		}),
	)
}
