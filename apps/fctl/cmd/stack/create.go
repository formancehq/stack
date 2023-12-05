package stack

import (
	"fmt"
	"net/http"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	unprotectFlag = "unprotect"
	regionFlag    = "region"
	nowaitFlag    = "no-wait"
)

type StackCreateStore struct {
	Stack    *membershipclient.Stack
	Versions *shared.GetVersionsResponse
}

type StackCreateController struct {
	store   *StackCreateStore
	profile *fctl.Profile
}

var _ fctl.Controller[*StackCreateStore] = (*StackCreateController)(nil)

func NewDefaultStackCreateStore() *StackCreateStore {
	return &StackCreateStore{
		Stack:    &membershipclient.Stack{},
		Versions: &shared.GetVersionsResponse{},
	}
}
func NewStackCreateController() *StackCreateController {
	return &StackCreateController{
		store: NewDefaultStackCreateStore(),
	}
}

func NewCreateCommand() *cobra.Command {
	return fctl.NewMembershipCommand("create [name]",
		fctl.WithShortDescription("Create a new stack"),
		fctl.WithAliases("c", "cr"),
		fctl.WithArgs(cobra.RangeArgs(0, 1)),
		fctl.WithBoolFlag(unprotectFlag, false, "Unprotect stacks (no confirmation on write commands)"),
		fctl.WithStringFlag(regionFlag, "", "Region on which deploy the stack"),
		fctl.WithBoolFlag(nowaitFlag, false, "Not wait stack availability"),
		fctl.WithController[*StackCreateStore](NewStackCreateController()),
	)
}
func (c *StackCreateController) GetStore() *StackCreateStore {
	return c.store
}

func (c *StackCreateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

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

	protected := !fctl.GetBool(cmd, unprotectFlag)
	metadata := map[string]string{
		fctl.ProtectedStackMetadata: fctl.BoolPointerToString(&protected),
	}

	name := ""
	if len(args) > 0 {
		name = args[0]
	} else {
		name, err = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter a name")
		if err != nil {
			return nil, err
		}
	}

	region := fctl.GetString(cmd, regionFlag)
	if region == "" {
		regions, _, err := apiClient.DefaultApi.ListRegions(cmd.Context(), organization).Execute()
		if err != nil {
			return nil, errors.Wrap(err, "listing regions")
		}

		var options []string
		for _, region := range regions.Data {
			privacy := "Private"
			if region.Public {
				privacy = "Public "
			}
			name := "<noname>"
			if region.Name != "" {
				name = region.Name
			}
			options = append(options, fmt.Sprintf("%s | %s | %s", region.Id, privacy, name))
		}

		printer := pterm.DefaultInteractiveSelect.WithOptions(options)
		selectedOption, err := printer.Show("Please select a region")
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(options); i++ {
			if selectedOption == options[i] {
				region = regions.Data[i].Id
				break
			}
		}
	}

	stackResponse, _, err := apiClient.DefaultApi.CreateStack(cmd.Context(), organization).CreateStackRequest(membershipclient.CreateStackRequest{
		Name:     name,
		Metadata: metadata,
		RegionID: region,
	}).Execute()
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

	fctl.BasicTextCyan.WithWriter(cmd.OutOrStdout()).Printfln("Your dashboard will be reachable on: %s",
		profile.ServicesBaseUrl(stackResponse.Data).String())

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

func (c *StackCreateController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintStackInformation(cmd.OutOrStdout(), c.profile, c.store.Stack, c.store.Versions)
}
