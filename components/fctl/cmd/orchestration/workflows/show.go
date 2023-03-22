package workflows

import (
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show <id>",
		fctl.WithShortDescription("Show a workflow"),
		fctl.WithAliases("s"),
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

			response, err := client.Orchestration.
				GetWorkflow(cmd.Context(), operations.GetWorkflowRequest{
					FlowID: args[0],
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

			fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
			tableData := pterm.TableData{}
			tableData = append(tableData, []string{pterm.LightCyan("ID"), response.GetWorkflowResponse.Data.ID})
			tableData = append(tableData, []string{pterm.LightCyan("Name"), func() string {
				if response.GetWorkflowResponse.Data.Config.Name != nil {
					return *response.GetWorkflowResponse.Data.Config.Name
				}
				return ""
			}()})
			tableData = append(tableData, []string{pterm.LightCyan("Created at"), response.GetWorkflowResponse.Data.CreatedAt.Format(time.RFC3339)})
			tableData = append(tableData, []string{pterm.LightCyan("Updated at"), response.GetWorkflowResponse.Data.UpdatedAt.Format(time.RFC3339)})

			if err := pterm.DefaultTable.
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render(); err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout())

			fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Configuration")
			configAsBytes, err := yaml.Marshal(response.GetWorkflowResponse.Data.Config)
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(configAsBytes))

			return nil
		}),
	)
}
