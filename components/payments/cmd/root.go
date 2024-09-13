package cmd

import (
	"errors"
	"os"

	_ "github.com/bombsimon/logrusr/v3"
	sharedapi "github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/auth"
	"github.com/formancehq/go-libs/bun/bunconnect"
	"github.com/formancehq/go-libs/bun/bunmigrate"
	"github.com/formancehq/go-libs/health"
	"github.com/formancehq/go-libs/licence"
	"github.com/formancehq/go-libs/otlp/otlptraces"
	"github.com/formancehq/go-libs/publish"
	"github.com/formancehq/go-libs/service"
	"github.com/formancehq/go-libs/temporal"
	"github.com/formancehq/payments/internal/api"
	v2 "github.com/formancehq/payments/internal/api/v2"
	v3 "github.com/formancehq/payments/internal/api/v3"
	"github.com/formancehq/payments/internal/connectors/engine"
	"github.com/formancehq/payments/internal/storage"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// TODO(polo/crimson): add profiling

var (
	ServiceName = "payments"
	Version     = "develop"
	BuildDate   = "-"
	Commit      = "-"
)

const (
	pluginsDirectoryPathFlag = "plugin-directory-path"
	configEncryptionKeyFlag  = "config-encryption-key"
	listenFlag               = "listen"
	stackFlag                = "stack"
	stackPublicURLFlag       = "stack-public-url"
)

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:               "payments",
		Short:             "payments",
		DisableAutoGenTag: true,
		Version:           Version,
	}

	root.PersistentFlags().String(configEncryptionKeyFlag, "", "Config encryption key")

	version := newVersion()
	root.AddCommand(version)

	migrate := newMigrate()
	root.AddCommand(migrate)

	server := newServer()
	addAutoMigrateCommand(server)
	server.Flags().String(listenFlag, ":8080", "Listen address")
	server.Flags().String(pluginsDirectoryPathFlag, "", "Plugin directory path")
	server.Flags().String(stackFlag, "", "Stack name")
	root.AddCommand(server)

	return root
}

func Execute() {
	service.Execute(NewRootCommand())
}

func addAutoMigrateCommand(cmd *cobra.Command) {
	cmd.Flags().Bool(autoMigrateFlag, false, "Auto migrate database")
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		autoMigrate, _ := cmd.Flags().GetBool(autoMigrateFlag)
		if autoMigrate {
			return bunmigrate.Run(cmd, args, Migrate)
		}
		return nil
	}
}

func commonOptions(cmd *cobra.Command) (fx.Option, error) {
	configEncryptionKey, _ := cmd.Flags().GetString(configEncryptionKeyFlag)
	if configEncryptionKey == "" {
		return nil, errors.New("missing config encryption key")
	}

	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd)
	if err != nil {
		return nil, err
	}

	path, _ := cmd.Flags().GetString(pluginsDirectoryPathFlag)
	pluginPaths, err := getPluginsMap(path)
	if err != nil {
		return nil, err
	}

	listen, _ := cmd.Flags().GetString(listenFlag)
	stack, _ := cmd.Flags().GetString(stackFlag)
	stackPublicURL, _ := cmd.Flags().GetString(stackPublicURLFlag)

	return fx.Options(
		fx.Provide(func() *bunconnect.ConnectionOptions {
			return connectionOptions
		}),
		fx.Provide(func() sharedapi.ServiceInfo {
			return sharedapi.ServiceInfo{
				Version: Version,
			}
		}),
		otlptraces.FXModuleFromFlags(cmd),
		temporal.FXModuleFromFlags(
			cmd,
			engine.Tracer,
			temporal.SearchAttributes{
				SearchAttributes: engine.SearchAttributes,
			},
		),
		auth.FXModuleFromFlags(cmd),
		health.Module(),
		publish.FXModuleFromFlags(cmd, service.IsDebug(cmd)),
		licence.FXModuleFromFlags(cmd, ServiceName),
		storage.Module(cmd, *connectionOptions, configEncryptionKey),
		api.NewModule(listen, service.IsDebug(cmd)),
		engine.Module(pluginPaths, stack, stackPublicURL),
		v2.NewModule(),
		v3.NewModule(),
	), nil
}

func getPluginsMap(pluginsDirectoryPath string) (map[string]string, error) {
	if pluginsDirectoryPath == "" {
		return nil, errors.New("missing plugin directory path")
	}

	files, err := os.ReadDir(pluginsDirectoryPath)
	if err != nil {
		return nil, errors.New("failed to read plugins directory")
	}

	plugins := make(map[string]string)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		plugins[file.Name()] = pluginsDirectoryPath + "/" + file.Name()
	}

	return plugins, nil
}
