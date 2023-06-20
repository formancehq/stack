package stack

import (
	"fmt"
	"net/http"

	"github.com/formancehq/fctl/cmd/stack/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewRestoreStackCommand() *cobra.Command {
	const stackNameFlag = "name"

	return fctl.NewMembershipCommand("restore <stack-id>",
		fctl.WithShortDescription("Restore a deleted stack"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithStringFlag(stackNameFlag, "", ""),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}

			organization, err := fctl.ResolveOrganizationID(cmd, cfg)
			if err != nil {
				return errors.Wrap(err, "searching default organization")
			}

			apiClient, err := fctl.NewMembershipClient(cmd, cfg)
			if err != nil {
				return err
			}

			response, _, err := apiClient.DefaultApi.
				RestoreStack(cmd.Context(), organization, args[0]).
				Execute()
			if err != nil {
				return err
			}

			profile := fctl.GetCurrentProfile(cmd, cfg)

			if err := waitStackReady(cmd, profile, response.Data); err != nil {
				return err
			}

			stackClient, err := fctl.NewStackClient(cmd, cfg, response.Data)
			if err != nil {
				return err
			}

			versions, err := stackClient.GetVersions(cmd.Context())
			if err != nil {
				return err
			}

			if versions.StatusCode != http.StatusOK {
				return fmt.Errorf("unexpected status code %d when reading versions", versions.StatusCode)
			}

			return internal.PrintStackInformation(cmd.OutOrStdout(), fctl.GetCurrentProfile(cmd, cfg), response.Data, versions.GetVersionsResponse)
		}),
	)
}
