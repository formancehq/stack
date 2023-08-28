package organizations

import (
	"flag"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useList   = "list"
	shortList = "List all organizations"
)

type OrgRow struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	OwnerID    string `json:"ownerId"`
	OwnerEmail string `json:"ownerEmail"`
	IsMine     string `json:"isMine"`
}

type ListStore struct {
	Organizations []*OrgRow `json:"organizations"`
}

func NewListStore() *ListStore {
	return &ListStore{}
}
func NewListConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	return fctl.NewControllerConfig(
		useList,
		shortList,
		shortList,
		[]string{
			"ls", "l",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

type ListController struct {
	store  *ListStore
	config *fctl.ControllerConfig
}

func NewListController(config *fctl.ControllerConfig) *ListController {
	return &ListController{
		store:  NewListStore(),
		config: config,
	}
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *ListController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	organizations, _, err := apiClient.DefaultApi.ListOrganizationsExpanded(ctx).Execute()
	if err != nil {
		return nil, err
	}

	currentProfile := fctl.GetCurrentProfile(flags, cfg)
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

func (c *ListController) Render() error {
	OrgMap := fctl.Map(c.store.Organizations, func(o *OrgRow) []string {
		return []string{o.ID, o.Name, o.OwnerID, o.OwnerEmail, o.IsMine}
	})

	tableData := fctl.Prepend(OrgMap, []string{"ID", "Name", "Owner ID", "Owner email", "Is mine?"})

	return pterm.DefaultTable.WithHasHeader().WithWriter(c.config.GetOut()).WithData(tableData).Render()
}

func NewListCommand() *cobra.Command {
	config := NewListConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController(config)),
	)
}
