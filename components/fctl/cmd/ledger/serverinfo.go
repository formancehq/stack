package ledger

import (
	"flag"
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useServerInfo   = "server-infos"
	shortServerInfo = "Read server info"
)

type ServerInfoStore struct {
	Server        string   `json:"server"`
	Version       string   `json:"version"`
	StorageDriver string   `json:"storageDriver"`
	Ledgers       []string `json:"ledgers"`
}

func NewServerInfoStore() *ServerInfoStore {
	return &ServerInfoStore{
		Server:        "unknown",
		Version:       "unknown",
		StorageDriver: "unknown",
		Ledgers:       []string{},
	}
}

func NewServerInfoConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useServerInfo, flag.ExitOnError)

	return fctl.NewControllerConfig(
		useServerInfo,
		shortServerInfo,
		shortServerInfo,
		[]string{
			"si",
		},
		flags,
		fctl.Organization, fctl.Stack, fctl.Ledger,
	)
}

var _ fctl.Controller[*ServerInfoStore] = (*ServerInfoController)(nil)

type ServerInfoController struct {
	store  *ServerInfoStore
	config *fctl.ControllerConfig
}

func NewServerInfoController(config *fctl.ControllerConfig) *ServerInfoController {
	return &ServerInfoController{
		store:  NewServerInfoStore(),
		config: config,
	}
}

func (c *ServerInfoController) GetStore() *ServerInfoStore {
	return c.store
}

func (c *ServerInfoController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *ServerInfoController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
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

	ledgerClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	response, err := ledgerClient.Ledger.GetInfo(ctx)
	if err != nil {
		return nil, err
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Server = response.ConfigInfoResponse.Data.Server
	c.store.Version = response.ConfigInfoResponse.Data.Version
	c.store.StorageDriver = response.ConfigInfoResponse.Data.Config.Storage.Driver
	c.store.Ledgers = response.ConfigInfoResponse.Data.Config.Storage.Ledgers

	return c, nil
}

func (c *ServerInfoController) Render() error {
	out := c.config.GetOut()

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Server"), fmt.Sprint(c.store.Server)})
	tableData = append(tableData, []string{pterm.LightCyan("Version"), fmt.Sprint(c.store.Version)})
	tableData = append(tableData, []string{pterm.LightCyan("Storage driver"), fmt.Sprint(c.store.StorageDriver)})

	if err := pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fctl.BasicTextCyan.WithWriter(out).Printfln("Ledgers :")
	if err := pterm.DefaultBulletList.
		WithWriter(out).
		WithItems(fctl.Map(c.store.Ledgers, func(ledger string) pterm.BulletListItem {
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

func NewServerInfoCommand() *cobra.Command {

	config := NewServerInfoConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ServerInfoStore](NewServerInfoController(config)),
	)
}
