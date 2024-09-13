package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/formancehq/payments/internal/connectors/grpc"
	"github.com/formancehq/payments/internal/models"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("plugin not found")
)

type Plugins interface {
	RegisterPlugin(connectorID models.ConnectorID) error
	UnregisterPlugin(connectorID models.ConnectorID) error
	Get(connectorID models.ConnectorID) (models.Plugin, error)
}

// Will start, hold, manage and stop plugins
type plugins struct {
	pluginsPath map[string]string

	plugins map[string]pluginInformation
	rwMutex sync.RWMutex
}

type pluginInformation struct {
	client *plugin.Client
}

func New(pluginsPath map[string]string) *plugins {
	return &plugins{
		pluginsPath: pluginsPath,
		plugins:     make(map[string]pluginInformation),
	}
}

func (p *plugins) RegisterPlugin(connectorID models.ConnectorID) error {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	// Check if plugin is already installed
	_, ok := p.plugins[connectorID.String()]
	if ok {
		return nil
	}

	pluginPath, ok := p.pluginsPath[connectorID.Provider]
	if !ok {
		return errors.Wrap(ErrNotFound, "plugin path not found")
	}

	pc := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  grpc.Handshake,
		Plugins:          grpc.PluginMap,
		Cmd:              exec.Command("sh", "-c", pluginPath),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:   fmt.Sprintf("%s-%s", connectorID.Provider, connectorID.String()),
			Output: os.Stdout,
			Level:  hclog.Debug,
		}),
	})

	p.plugins[connectorID.String()] = pluginInformation{
		client: pc,
	}

	return nil
}

func (p *plugins) UnregisterPlugin(connectorID models.ConnectorID) error {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	pluginInfo, ok := p.plugins[connectorID.String()]
	if !ok {
		// Nothing to do``
		return nil
	}

	// Close the connection
	pluginInfo.client.Kill()

	delete(p.plugins, connectorID.String())

	return nil
}

func (p *plugins) Get(connectorID models.ConnectorID) (models.Plugin, error) {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	pluginInfo, ok := p.plugins[connectorID.String()]
	if !ok {
		return nil, errors.New("plugin not found")
	}

	return getPlugin(pluginInfo.client)
}

func getPlugin(client *plugin.Client) (models.Plugin, error) {
	// Connect via RPC
	conn, err := client.Client()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to plugin")
	}

	raw, err := conn.Dispense("psp")
	if err != nil {
		return nil, errors.Wrap(err, "failed to dispense plugin")
	}

	plugin, ok := raw.(grpc.PSP)
	if !ok {
		return nil, errors.New("failed to cast plugin")
	}

	impl := &impl{
		pluginClient: plugin,
	}

	return impl, nil
}

var _ Plugins = &plugins{}
