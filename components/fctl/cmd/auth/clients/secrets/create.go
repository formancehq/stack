package secrets

import (
	"flag"
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useCreate         = "create <client-id> <secret-name>"
	descriptionCreate = "Create a new secret for a client. You can list all clients with `fctl auth clients list`"
	shortCreate       = "Create a new secret for a client"
)

type CreateStore struct {
	SecretId string `json:"secretId"`
	Name     string `json:"name"`
	Clear    string `json:"clear"`
}

func NewCreateStore() *CreateStore {
	return &CreateStore{}
}
func NewSetupConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useCreate, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)

	return fctl.NewControllerConfig(
		useCreate,
		descriptionCreate,
		shortCreate,
		[]string{
			"c",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

var _ fctl.Controller[*CreateStore] = (*CreateController)(nil)

type CreateController struct {
	store  *CreateStore
	config *fctl.ControllerConfig
}

func NewCreateController(config *fctl.ControllerConfig) *CreateController {
	return &CreateController{
		store:  NewCreateStore(),
		config: config,
	}
}

func (c *CreateController) GetStore() *CreateStore {
	return c.store
}

func (c *CreateController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *CreateController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()
	out := c.config.GetOut()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are about to create a new client secret") {
		return nil, fctl.ErrMissingApproval
	}

	authClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	request := operations.CreateSecretRequest{
		ClientID: args[0],
		CreateSecretRequest: &shared.CreateSecretRequest{
			Name:     args[1],
			Metadata: nil,
		},
	}
	response, err := authClient.Auth.CreateSecret(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.SecretId = response.CreateSecretResponse.Data.ID
	c.store.Name = response.CreateSecretResponse.Data.Name
	c.store.Clear = response.CreateSecretResponse.Data.Clear

	return c, nil
}

func (c *CreateController) Render() error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.SecretId})
	tableData = append(tableData, []string{pterm.LightCyan("Name"), c.store.Name})
	tableData = append(tableData, []string{pterm.LightCyan("Clear"), c.store.Clear})
	return pterm.DefaultTable.
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()

}

func NewCreateCommand() *cobra.Command {
	config := NewSetupConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*CreateStore](NewCreateController(config)),
	)
}
