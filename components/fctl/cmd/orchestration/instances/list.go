package instances

import (
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	const (
		workflowFlag = "workflow"
		runningFlag  = "running"
	)
	return fctl.NewCommand("list",
		fctl.WithShortDescription("List all workflows instances"),
		fctl.WithAliases("ls", "l"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithStringFlag(workflowFlag, "", "Filter on workflow id"),
		fctl.WithBoolFlag(runningFlag, false, "Filter on running instances"),
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

			response, err := client.Orchestration.ListInstances(cmd.Context(), operations.ListInstancesRequest{
				Running:    fctl.Ptr(fctl.GetBool(cmd, runningFlag)),
				WorkflowID: fctl.Ptr(fctl.GetString(cmd, workflowFlag)),
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

			if len(response.ListRunsResponse.Data) == 0 {
				fctl.Println("No workflows found.")
				return nil
			}

			if err := pterm.DefaultTable.
				WithHasHeader(true).
				WithWriter(cmd.OutOrStdout()).
				WithData(
					fctl.Prepend(
						fctl.Map(response.ListRunsResponse.Data,
							func(src shared.WorkflowInstance) []string {
								return []string{
									src.ID,
									src.WorkflowID,
									src.CreatedAt.Format(time.RFC3339),
									src.UpdatedAt.Format(time.RFC3339),
									func() string {
										if src.Terminated {
											return src.TerminatedAt.Format(time.RFC3339)
										}
										return ""
									}(),
								}
							}),
						[]string{"ID", "Workflow ID", "Created at", "Updated at", "Terminated at"},
					),
				).Render(); err != nil {
				return errors.Wrap(err, "rendering table")
			}

			return nil
		}),
	)
}
