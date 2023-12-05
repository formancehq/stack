package invitations

import (
	"time"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type Invitations struct {
	Id           string    `json:"id"`
	UserEmail    string    `json:"userEmail"`
	Status       string    `json:"status"`
	CreationDate time.Time `json:"creationDate"`
}

type ListStore struct {
	Invitations []Invitations `json:"invitations"`
}
type ListController struct {
	store      *ListStore
	statusFlag string
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewDefaultListStore() *ListStore {
	return &ListStore{}
}

func NewListController() *ListController {
	return &ListController{
		store:      NewDefaultListStore(),
		statusFlag: "status",
	}
}

func NewListCommand() *cobra.Command {
	c := NewListController()
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithShortDescription("List invitations"),
		fctl.WithAliases("s"),
		fctl.WithStringFlag(c.statusFlag, "", "Filter invitations by status"),
		fctl.WithController[*ListStore](c),
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

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	listInvitationsResponse, _, err := apiClient.DefaultApi.
		ListOrganizationInvitations(cmd.Context(), organizationID).
		Status(fctl.GetString(cmd, c.statusFlag)).
		Execute()
	if err != nil {
		return nil, err
	}

	c.store.Invitations = fctl.Map(listInvitationsResponse.Data, func(i membershipclient.Invitation) Invitations {
		return Invitations{
			Id:           i.Id,
			UserEmail:    i.UserEmail,
			Status:       i.Status,
			CreationDate: i.CreationDate,
		}
	})

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	tableData := fctl.Map(c.store.Invitations, func(i Invitations) []string {
		return []string{
			i.Id,
			i.UserEmail,
			i.Status,
			i.CreationDate.Format(time.RFC3339),
		}
	})

	tableData = fctl.Prepend(tableData, []string{"ID", "Email", "Status", "Creation date"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()

}
