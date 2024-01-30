package versions

import (
	"fmt"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type ListStore struct {
	Versions []membershipclient.Version `json:"versions"`
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
	return fctl.NewCommand("list <region-id>",
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithAliases("ls", "l"),
		fctl.WithPreRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}

			apiClient, err := fctl.NewMembershipClient(cmd, cfg)
			if err != nil {
				return err
			}

			version := fctl.MembershipServerInfo(cmd.Context(), apiClient)
			if !semver.IsValid(version) {
				return nil
			}

			if semver.Compare(version, "v0.26.1") >= 0 {
				return nil
			}

			return fmt.Errorf("unsupported membership server version: %s", version)
		}),
		fctl.WithShortDescription("List all versions installed on a region"),
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

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	if len(args) != 1 {
		return nil, fmt.Errorf("missing region id")
	}

	regionId := args[0]

	regionalVersion, _, err := apiClient.DefaultApi.GetRegionVersions(cmd.Context(), organizationID, regionId).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Versions = regionalVersion.Data

	return c, nil
}

func printVersion(version *membershipclient.Version) pterm.TreeNode {
	node := pterm.TreeNode{
		Text:     pterm.LightCyan(version.Name),
		Children: []pterm.TreeNode{},
	}

	for name, v := range version.Versions {
		node.Children = append(node.Children, pterm.TreeNode{
			Text: fmt.Sprintf("%s: %s", pterm.Cyan(name), v),
		})
	}
	return node

}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	tree := pterm.TreeNode{
		Text:     pterm.Cyan("Versions"),
		Children: []pterm.TreeNode{},
	}
	for _, version := range c.store.Versions {
		tree.Children = append(tree.Children, printVersion(&version))
	}

	return pterm.DefaultTree.WithRoot(tree).Render()
}
