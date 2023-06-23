package organizations

import (
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type OrgRow struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	OwnerID    string `json:"owner_id"`
	OwnerEmail string `json:"owner_email"`
	IsMine     string `json:"is_mine"`
}

type ListStore struct {
	Organizations []*OrgRow `json:"organizations"`
}
type ListController struct {
	store *ListStore
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewDefaultListStore() *ListStore {
	return &ListStore{}
}

func NewListController() *ListController {
	return &ListController{
		store: NewDefaultListStore(),
	}
}

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List organizations"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController()),
	)
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	organizations, _, err := apiClient.DefaultApi.ListOrganizationsExpanded(cmd.Context()).Execute()
	if err != nil {
		return nil, err
	}

	currentProfile := fctl.GetCurrentProfile(cmd, cfg)
	claims, err := currentProfile.GetClaims()
	if err != nil {
		return nil, err
	}

	c.store.Organizations = fctl.Map(organizations.Data, func(o membershipclient.ListOrganizationExpandedResponseDataInner) *OrgRow {
		isMine := fctl.BoolToString(o.OwnerId == claims["sub"].(string))
		return &OrgRow{
			ID:         o.Id,
			Name:       o.Name,
			OwnerID:    o.OwnerId,
			OwnerEmail: o.Owner.Email,
			IsMine:     isMine,
		}
	})

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	OrgMap := fctl.Map(c.store.Organizations, func(o *OrgRow) []string {
		return []string{o.ID, o.Name, o.OwnerID, o.OwnerEmail, o.IsMine}
	})

	tableData := fctl.Prepend(OrgMap, []string{"ID", "Name", "Owner ID", "Owner email", "Is mine?"})

	return pterm.DefaultTable.WithHasHeader().WithWriter(cmd.OutOrStdout()).WithData(tableData).Render()
}
