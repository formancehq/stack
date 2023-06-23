package secrets

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type CreateStore struct {
	SecretId string `json:"secret_id"`
	Name     string `json:"name"`
	Clear    string `json:"clear"`
}
type CreateController struct {
	store *CreateStore
}

var _ fctl.Controller[*CreateStore] = (*CreateController)(nil)

func NewDefaultCreateStore() *CreateStore {
	return &CreateStore{}
}

func NewCreateController() *CreateController {
	return &CreateController{
		store: NewDefaultCreateStore(),
	}
}

func NewCreateCommand() *cobra.Command {
	return fctl.NewCommand("create <client-id> <secret-name>",
		fctl.WithAliases("c"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithShortDescription("Create secret"),
		fctl.WithConfirmFlag(),
		fctl.WithController[*CreateStore](NewCreateController()),
	)
}

func (c *CreateController) GetStore() *CreateStore {
	return c.store
}

func (c *CreateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to create a new client secret") {
		return nil, fctl.ErrMissingApproval
	}

	authClient, err := fctl.NewStackClient(cmd, cfg, stack)
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
	response, err := authClient.Auth.CreateSecret(cmd.Context(), request)
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

func (c *CreateController) Render(cmd *cobra.Command, args []string) error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.SecretId})
	tableData = append(tableData, []string{pterm.LightCyan("Name"), c.store.Name})
	tableData = append(tableData, []string{pterm.LightCyan("Clear"), c.store.Clear})
	return pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()

}
