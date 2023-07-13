package stack

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	unprotectFlag     = "unprotect"
	regionFlag        = "region"
	nowaitFlag        = "no-wait"
	useCreate         = "create [name]"
	descriptionCreate = "Create a new stack"
)

type StackCreateStore struct {
	Stack    *membershipclient.Stack
	Versions *shared.GetVersionsResponse
}

func NewDefaultStackCreateStore() *StackCreateStore {
	return &StackCreateStore{
		Stack:    &membershipclient.Stack{},
		Versions: &shared.GetVersionsResponse{},
	}
}

type StackCreateControllerConfig struct {
	context     context.Context
	use         string
	description string
	aliases     []string
	out         io.Writer
	flags       *flag.FlagSet
	args        []string
}

func NewStackCreateControllerConfig(ctx *context.Context) *StackCreateControllerConfig {
	flags := flag.NewFlagSet(useCreate, flag.ExitOnError)
	flags.Bool(unprotectFlag, false, "Unprotect stacks (no confirmation on write commands)")
	flags.String(regionFlag, "", "Region on which deploy the stack")
	flags.Bool(nowaitFlag, false, "Not wait stack availability")

	fctl.WithGlobalFlags(flags)
	return &StackCreateControllerConfig{
		context:     *ctx,
		use:         useCreate,
		description: descriptionCreate,
		aliases: []string{
			"cr",
			"c",
		},
		out:   os.Stdout,
		flags: flags,
		args:  []string{},
	}
}

var _ fctl.Controller[*StackCreateStore] = (*StackCreateController)(nil)

type StackCreateController struct {
	store   *StackCreateStore
	profile *fctl.Profile
	config  StackCreateControllerConfig
}

func NewStackCreateController(config StackCreateControllerConfig) *StackCreateController {
	return &StackCreateController{
		store:  NewDefaultStackCreateStore(),
		config: config,
	}
}

func (c *StackCreateController) GetFlags() *flag.FlagSet {
	return c.config.flags
}

func (c *StackCreateController) GetContext() context.Context {
	return c.config.context
}

func (c *StackCreateController) SetContext(ctx context.Context) {
	c.config.context = ctx
}

func (c *StackCreateController) GetStore() *StackCreateStore {
	return c.store
}

func (c *StackCreateController) SetArgs(args []string) {
	c.config.args = append([]string{}, args...)
}

func (c *StackCreateController) Run() (fctl.Renderable, error) {
	flags := c.config.flags
	ctx := c.config.context

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organization, err := fctl.ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	protected := !fctl.GetBool(flags, unprotectFlag)
	metadata := map[string]string{
		fctl.ProtectedStackMetadata: fctl.BoolPointerToString(&protected),
	}

	name := ""
	if len(c.config.args) > 0 {
		name = c.config.args[0]
	} else {
		name, err = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter a name")
		if err != nil {
			return nil, err
		}
	}

	region := fctl.GetString(flags, regionFlag)
	if region == "" {
		regions, _, err := apiClient.DefaultApi.ListRegions(ctx, organization).Execute()
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

	stackResponse, _, err := apiClient.DefaultApi.CreateStack(ctx, organization).CreateStackRequest(membershipclient.CreateStackRequest{
		Name:     name,
		Metadata: metadata,
		RegionID: region,
	}).Execute()
	if err != nil {
		return nil, errors.Wrap(err, "creating stack")
	}

	profile := fctl.GetCurrentProfile(flags, cfg)

	if !fctl.GetBool(flags, nowaitFlag) {
		spinner, err := pterm.DefaultSpinner.Start("Waiting services availability")
		if err != nil {
			return nil, err
		}

		if err := waitStackReady(ctx, flags, profile, stackResponse.Data); err != nil {
			return nil, err
		}

		if err := spinner.Stop(); err != nil {
			return nil, err
		}
	}

	fctl.BasicTextCyan.WithWriter(c.config.out).Printfln("Your dashboard will be reachable on: %s",
		profile.ServicesBaseUrl(stackResponse.Data).String())

	stackClient, err := fctl.NewStackClient(flags, ctx, cfg, stackResponse.Data)
	if err != nil {
		return nil, err
	}

	versions, err := stackClient.GetVersions(ctx)
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

func (c *StackCreateController) Render() error {
	return internal.PrintStackInformation(c.config.out, c.profile, c.store.Stack, c.store.Versions)
}

func NewCreateCommand() *cobra.Command {

	config := NewStackCreateControllerConfig(nil)

	return fctl.NewMembershipCommand(config.use,
		fctl.WithShortDescription(config.description),
		fctl.WithAliases(config.aliases...),
		fctl.WithArgs(cobra.RangeArgs(0, 1)),
		fctl.WithGoFlagSet(config.flags),
		fctl.WithController[*StackCreateStore](NewStackCreateController(*config)),
	)
}
