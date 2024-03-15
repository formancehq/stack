package reconciliation

import (
	"github.com/formancehq/fctl/cmd/reconciliation/policies"
	"github.com/formancehq/fctl/cmd/reconciliation/store"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewStackCommand("reconciliation",
		fctl.WithShortDescription("Reconciliation management"),
		fctl.WithChildCommands(
			policies.NewPoliciesCommand(),
			NewListCommand(),
			NewShowCommand(),
		),
		fctl.WithPersistentPreRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}
			apiClient, err := fctl.NewMembershipClient(cmd, cfg)
			if err != nil {
				return err
			}
			organizationID, err := fctl.ResolveOrganizationID(cmd, cfg, apiClient.DefaultApi)
			if err != nil {
				return err
			}

			stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
			if err != nil {
				return err
			}

			stackClient, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return err
			}
			cmd.SetContext(store.ContextWithStore(cmd.Context(), store.ReconciliationNode(cfg, stack, organizationID, stackClient)))
			return nil
		}),
	)
}
