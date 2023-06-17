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

var errStackNotFound = errors.New("stack not found")

type StackInformation struct {
	Stack    *membershipclient.Stack     `json:"stack"`
	Versions *shared.GetVersionsResponse `json:"versions"`
}

func NewShowCommand() *cobra.Command {
	var stackNameFlag = "name"

	return fctl.NewMembershipCommand("show (<stack-id> | --name=<stack-name>)",
		fctl.WithAliases("s", "sh"),
		fctl.WithShortDescription("Show stack"),
		fctl.WithArgs(cobra.MaximumNArgs(1)),
		fctl.WithStringFlag(stackNameFlag, "", ""),
		fctl.WithRunE(showCommand),
		fctl.WrapOutputPostRunE(viewStackInformation),
	)
}

func showCommand(cmd *cobra.Command, args []string) error {
	var stackNameFlag = "name"

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

	var stack *membershipclient.Stack
	if len(args) == 1 {
		if fctl.GetString(cmd, stackNameFlag) != "" {
			return errors.New("need either an id of a name specified using --name flag")
		}
		stackResponse, httpResponse, err := apiClient.DefaultApi.ReadStack(cmd.Context(), organization, args[0]).Execute()
		if err != nil {
			if httpResponse.StatusCode == http.StatusNotFound {
				return errStackNotFound
			}
			return errors.Wrap(err, "listing stacks")
		}
		stack = stackResponse.Data
	} else {
		if fctl.GetString(cmd, stackNameFlag) == "" {
			return errors.New("need either an id of a name specified using --name flag")
		}
		stacksResponse, _, err := apiClient.DefaultApi.ListStacks(cmd.Context(), organization).Execute()
		if err != nil {
			return errors.Wrap(err, "listing stacks")
		}
		for _, s := range stacksResponse.Data {
			if s.Name == fctl.GetString(cmd, stackNameFlag) {
				stack = &s
				break
			}
		}
	}

	if stack == nil {
		return errStackNotFound
	}

	stackClient, err := fctl.NewStackClient(cmd, cfg, stack)
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

	fctl.SetSharedData(&StackInformation{
		Stack:    stack,
		Versions: versions.GetVersionsResponse,
	}, nil, cfg, nil)

	return nil

}

func viewStackInformation(cmd *cobra.Command, args []string) error {
	data := fctl.GetSharedData().(*StackInformation)
	return internal.PrintStackInformation(cmd.OutOrStdout(), fctl.GetCurrentProfile(cmd, fctl.GetSharedConfig()), data.Stack, data.Versions)
}
