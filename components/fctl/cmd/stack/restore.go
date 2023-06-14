package stack

import (
	"fmt"
	"net/http"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewRestoreStackCommand() *cobra.Command {
	const stackNameFlag = "name"

	return fctl.NewMembershipCommand("restore <stack-id>",
		fctl.WithShortDescription("Restore a deleted stack"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithStringFlag(stackNameFlag, "", ""),
		fctl.WithRunE(restoreCommand),
		fctl.WrapOutputPostRunE(viewRestore),
	)
}

func restoreCommand(cmd *cobra.Command, args []string) error {
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

	fctl.SetSharedData(&RestoreStackInformation{
		Stack:    response.Data,
		Versions: versions.GetVersionsResponse,
	}, profile, cfg)

	return nil
}

type RestoreStackInformation struct {
	Stack    *membershipclient.Stack     `json:"stack"`
	Versions *shared.GetVersionsResponse `json:"versions"`
}

func viewRestore(cmd *cobra.Command, args []string) error {
	data := fctl.GetSharedData().(*RestoreStackInformation)

	return internal.PrintStackInformation(cmd.OutOrStdout(), fctl.GetCurrentProfile(cmd, fctl.GetSharedConfig()), data.Stack, data.Versions)
}
