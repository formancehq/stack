package ledger

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ServerInfoStore struct {
	ConfigInfoResponse *shared.ConfigInfo `json:"config_info_response"`
}
type ServerInfoController struct {
	store *ServerInfoStore
}

var _ fctl.Controller[*ServerInfoStore] = (*ServerInfoController)(nil)

func NewDefaultServerInfoStore() *ServerInfoStore {
	return &ServerInfoStore{
		ConfigInfoResponse: &shared.ConfigInfo{},
	}
}

func NewServerInfoController() *ServerInfoController {
	return &ServerInfoController{
		store: NewDefaultServerInfoStore(),
	}
}

func NewServerInfoCommand() *cobra.Command {
	return fctl.NewCommand("server-infos",
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithAliases("si"),
		fctl.WithShortDescription("Read server info"),
		fctl.WithController[*ServerInfoStore](NewServerInfoController()),
	)
}

func (c *ServerInfoController) GetStore() *ServerInfoStore {
	return c.store
}

func (c *ServerInfoController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	ledgerClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	response, err := ledgerClient.Ledger.GetInfo(cmd.Context())
	if err != nil {
		return nil, err
	}

	c.store.ConfigInfoResponse = response.ConfigInfoResponse.Data

	return c, nil
}

func (c *ServerInfoController) Render(cmd *cobra.Command, args []string) error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Server"), fmt.Sprint(c.store.ConfigInfoResponse.Server)})
	tableData = append(tableData, []string{pterm.LightCyan("Version"), fmt.Sprint(c.store.ConfigInfoResponse.Version)})
	tableData = append(tableData, []string{pterm.LightCyan("Storage driver"), fmt.Sprint(c.store.ConfigInfoResponse.Config.Storage.Driver)})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fctl.BasicTextCyan.WithWriter(cmd.OutOrStdout()).Printfln("Ledgers :")
	if err := pterm.DefaultBulletList.
		WithWriter(cmd.OutOrStdout()).
		WithItems(fctl.Map(c.store.ConfigInfoResponse.Config.Storage.Ledgers, func(ledger string) pterm.BulletListItem {
			return pterm.BulletListItem{
				Text:        ledger,
				TextStyle:   pterm.NewStyle(pterm.FgDefault),
				BulletStyle: pterm.NewStyle(pterm.FgLightCyan),
			}
		})).
		Render(); err != nil {
		return err
	}

	return nil
}
