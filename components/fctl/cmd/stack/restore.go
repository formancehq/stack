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
	"github.com/spf13/cobra"
)

const (
	useRestore         = "restore <stack-id>"
	descriptionRestore = "Restore a stack"
)

type StackRestoreStore struct {
	Stack    *membershipclient.Stack     `json:"stack"`
	Versions *shared.GetVersionsResponse `json:"versions"`
}

func NewDefaultVersionStore() *StackRestoreStore {
	return &StackRestoreStore{
		Stack:    &membershipclient.Stack{},
		Versions: &shared.GetVersionsResponse{},
	}
}

type StackRestoreControllerConfig struct {
	context     context.Context
	use         string
	description string
	aliases     []string
	out         io.Writer
	flags       *flag.FlagSet
	args        []string
	fctlConfig  *fctl.Config
}

func NewStackRestoreControllerConfig() *StackRestoreControllerConfig {
	flags := flag.NewFlagSet(useRestore, flag.ExitOnError)
	flags.String(internal.StackNameFlag, "", "Stack name")

	fctl.WithGlobalFlags(flags)

	return &StackRestoreControllerConfig{
		context:     nil,
		use:         useRestore,
		description: descriptionRestore,
		aliases: []string{
			"res",
			"r",
		},
		out:   os.Stdout,
		flags: flags,
		args:  []string{},
	}
}

var _ fctl.Controller[*StackRestoreStore] = (*StackRestoreController)(nil)

type StackRestoreController struct {
	store  *StackRestoreStore
	config StackRestoreControllerConfig
}

func NewStackRestoreController(config StackRestoreControllerConfig) *StackRestoreController {
	return &StackRestoreController{
		store:  NewDefaultVersionStore(),
		config: config,
	}
}

func (c *StackRestoreController) GetFlags() *flag.FlagSet {
	return c.config.flags
}

func (c *StackRestoreController) GetContext() context.Context {
	return c.config.context
}

func (c *StackRestoreController) SetContext(ctx context.Context) {
	c.config.context = ctx
}

func (c *StackRestoreController) GetStore() *StackRestoreStore {
	return c.store
}

func (c *StackRestoreController) SetArgs(args []string) {
	c.config.args = append([]string{}, args...)
}

func (c *StackRestoreController) Run() (fctl.Renderable, error) {
	flags := c.config.flags
	ctx := c.config.context

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organization, err := fctl.ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	if len(c.config.args) == 0 {
		return nil, fmt.Errorf("stack id is required")
	}

	response, _, err := apiClient.DefaultApi.
		RestoreStack(ctx, organization, c.config.args[0]).
		Execute()
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(flags, cfg)

	if err := waitStackReady(ctx, flags, profile, response.Data); err != nil {
		return nil, err
	}

	stackClient, err := fctl.NewStackClient(flags, ctx, cfg, response.Data)
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

	c.store.Stack = response.Data
	c.store.Versions = versions.GetVersionsResponse
	c.config.fctlConfig = cfg

	return c, nil
}

func (c *StackRestoreController) Render() error {
	return internal.PrintStackInformation(c.config.out, fctl.GetCurrentProfile(c.config.flags, c.config.fctlConfig), c.store.Stack, c.store.Versions)
}

func NewRestoreStackCommand() *cobra.Command {
	config := NewStackRestoreControllerConfig()
	return fctl.NewMembershipCommand(config.use,
		fctl.WithShortDescription(config.description),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithAliases(config.aliases...),
		fctl.WithGoFlagSet(config.flags),
		fctl.WithController[*StackRestoreStore](NewStackRestoreController(*config)),
	)
}
