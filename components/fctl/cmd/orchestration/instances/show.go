package instances

import (
	"fmt"
	"time"

	"github.com/formancehq/fctl/cmd/orchestration/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
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

			res, err := client.Orchestration.GetInstance(cmd.Context(), operations.GetInstanceRequest{
				InstanceID: args[0],
			})
			if err != nil {
				return errors.Wrap(err, "reading instance")
			}

			fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
			tableData := pterm.TableData{}
			tableData = append(tableData, []string{pterm.LightCyan("ID"), res.GetWorkflowInstanceResponse.Data.ID})
			tableData = append(tableData, []string{pterm.LightCyan("Created at"), res.GetWorkflowInstanceResponse.Data.CreatedAt.Format(time.RFC3339)})
			tableData = append(tableData, []string{pterm.LightCyan("Updated at"), res.GetWorkflowInstanceResponse.Data.UpdatedAt.Format(time.RFC3339)})
			if res.GetWorkflowInstanceResponse.Data.Terminated {
				tableData = append(tableData, []string{pterm.LightCyan("Terminated at"), res.GetWorkflowInstanceResponse.Data.TerminatedAt.Format(time.RFC3339)})
			}
			if res.GetWorkflowInstanceResponse.Data.Error != nil && *res.GetWorkflowInstanceResponse.Data.Error != "" {
				tableData = append(tableData, []string{pterm.LightCyan("Error"), *res.GetWorkflowInstanceResponse.Data.Error})
			}

			if err := pterm.DefaultTable.
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render(); err != nil {
				return err
			}

			response, err := client.Orchestration.GetWorkflow(cmd.Context(), operations.GetWorkflowRequest{
				FlowID: res.GetWorkflowInstanceResponse.Data.WorkflowID,
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

			if err := internal.PrintWorkflowInstance(cmd.OutOrStdout(), response.GetWorkflowResponse.Data, res.GetWorkflowInstanceResponse.Data); err != nil {
				return err
			}

			return nil
		}),
	)
}
