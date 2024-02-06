package stack

import (
	"fmt"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/pointer"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type StackUpdateStore struct {
	Stack    *membershipclient.Stack
	Versions *shared.GetVersionsResponse
}

type StackUpdateController struct {
	store   *StackUpdateStore
	profile *fctl.Profile
}

var _ fctl.Controller[*StackUpdateStore] = (*StackUpdateController)(nil)

func NewDefaultStackUpdateStore() *StackUpdateStore {
	return &StackUpdateStore{
		Stack:    &membershipclient.Stack{},
		Versions: &shared.GetVersionsResponse{},
	}
}
func NewStackUpdateController() *StackUpdateController {
	return &StackUpdateController{
		store: NewDefaultStackUpdateStore(),
	}
}

func NewUpdateCommand() *cobra.Command {
	return fctl.NewMembershipCommand("update <stack-id> [name]",
		fctl.WithShortDescription("Update a created stack"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithStringFlag(versionFlag, "", "Version of the stack"),
		fctl.WithController[*StackUpdateStore](NewStackUpdateController()),
	)
}
func (c *StackUpdateController) GetStore() *StackUpdateStore {
	return c.store
}

func (c *StackUpdateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	organization, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	if len(args) < 2 {
		return nil, errors.New("missing name or stack-id argument")
	}

	name := ""
	if len(args) > 1 {
		name = args[1]
	} else {
		name, err = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter a name")
		if err != nil {
			return nil, err
		}
	}

	req := membershipclient.UpdateStackRequest{
		Name: name,
	}

	//Retrieve region from the stack
	stack, res, err := apiClient.DefaultApi.GetStack(cmd.Context(), organization, args[0]).Execute()
	if err != nil {
		return nil, errors.Wrap(err, "retrieving stack")
	}

	if res.StatusCode > 300 {
		return nil, errors.New("stack not found")
	}

	availableVersions, httpResponse, err := apiClient.DefaultApi.GetRegionVersions(cmd.Context(), organization, stack.Data.RegionID).Execute()
	if httpResponse == nil {
		return nil, err
	}

	specifiedVersion := fctl.GetString(cmd, versionFlag)
	switch {
	case httpResponse.StatusCode == http.StatusOK && specifiedVersion == "":
		var options []string
		for _, version := range availableVersions.Data {
			options = append(options, version.Name)
		}

		printer := pterm.DefaultInteractiveSelect.WithOptions(options)
		selectedOption, err := printer.Show("Please select a version")
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(options); i++ {
			if selectedOption == options[i] {
				specifiedVersion = availableVersions.Data[i].Name
				break
			}
		}
	case httpResponse.StatusCode != http.StatusOK && specifiedVersion == "":
		// nothing to do, we cannot set a specific version for this membership version
	case httpResponse.StatusCode == http.StatusOK && specifiedVersion != "":
		// nothing to do, let membership handle the case
	case httpResponse.StatusCode != http.StatusOK && specifiedVersion != "":
		return nil, errors.New("--version flag can not be used with the actual membership api")
	}

	if specifiedVersion != "" {
		req.Version = pointer.For(specifiedVersion)
	}

	stackResponse, _, err := apiClient.DefaultApi.
		UpdateStack(cmd.Context(), organization, args[0]).
		UpdateStackRequest(req).
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, "creating stack")
	}

	profile := fctl.GetCurrentProfile(cmd, cfg)

	if !fctl.GetBool(cmd, nowaitFlag) {
		spinner, err := pterm.DefaultSpinner.Start("Waiting services availability")
		if err != nil {
			return nil, err
		}

		if err := waitStackReady(cmd, profile, stackResponse.Data); err != nil {
			return nil, err
		}

		if err := spinner.Stop(); err != nil {
			return nil, err
		}
	}

	fctl.BasicTextCyan.WithWriter(cmd.OutOrStdout()).Println("Your dashboard will be reachable on: https://console.formance.cloud")

	stackClient, err := fctl.NewStackClient(cmd, cfg, stackResponse.Data)
	if err != nil {
		return nil, err
	}

	versions, err := stackClient.GetVersions(cmd.Context())
	if err != nil {
		return nil, err
	}
	if versions.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when reading versions", versions.StatusCode)
	}

	c.store.Stack = stackResponse.Data
	c.store.Versions = versions.GetVersionsResponse
	c.profile = profile

	return c, nil
}

func (c *StackUpdateController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintStackInformation(cmd.OutOrStdout(), c.profile, c.store.Stack, c.store.Versions)
}
