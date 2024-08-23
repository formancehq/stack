package cmd

import (
	"io"
	"text/template"

	"github.com/formancehq/stack/ee/ingester/internal/drivers"
	"github.com/formancehq/stack/ee/ingester/internal/modules"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"

	"github.com/formancehq/stack/ee/ingester/internal/controller"
	"github.com/formancehq/stack/ee/ingester/internal/drivers/all"
	"github.com/formancehq/stack/ee/ingester/internal/runner"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"

	"github.com/formancehq/stack/ee/ingester/internal/httpclient"

	"github.com/formancehq/stack/ee/ingester/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"

	"github.com/formancehq/stack/ee/ingester/internal/api"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

const (
	StackFlag              = "stack"
	stackIssuerFlag        = "stack-issuer"
	stackClientIDFlag      = "stack-client-id"
	stackClientSecretFlag  = "stack-client-secret"
	BindFlag               = "bind"
	ModuleURLTplFlag       = "module-url-tpl"
	ModulePullPageSizeFlag = "module-pull-page-size"

	ServiceName = "ingester"

	Version = "develop"
)

func NewRootCommand() *cobra.Command {
	ret := &cobra.Command{
		Use:          "ingester",
		Short:        "General ingestion system",
		SilenceUsage: true,
	}
	ret.AddCommand(NewServeCommand())

	return ret
}

func NewServeCommand() *cobra.Command {
	ret := &cobra.Command{
		RunE: serve,
		Use:  "serve",
		Args: cobra.ExactArgs(0),
	}
	ret.Flags().String(stackClientIDFlag, "", "Stack client ID")
	ret.Flags().String(stackClientSecretFlag, "", "Stack client secret")
	ret.Flags().String(stackIssuerFlag, "", "Stack issuer flag")
	ret.Flags().String(StackFlag, "", "Stack")
	ret.Flags().String(BindFlag, ":8080", "HTTP server listen address")
	ret.Flags().String(ModuleURLTplFlag, "", "Module url template")
	ret.Flags().Int(ModulePullPageSizeFlag, 100, "Page size to pull modules")
	ret.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Register generic service flags
	service.AddFlags(ret.Flags())
	publish.AddFlags(ServiceName, ret.Flags())
	bunconnect.AddFlags(ret.Flags())
	otlptraces.AddFlags(ret.Flags())
	auth.AddFlags(ret.Flags())

	return ret
}

type runConfiguration struct {
	stack             string
	debug             bool
	pullConfiguration modules.PullConfiguration
	authConfig        httpclient.OAuth2Config
	bind              string
	connectionOptions bunconnect.ConnectionOptions
}

func evalServeConfiguration(cmd *cobra.Command) (ret *runConfiguration, err error) {
	ret = &runConfiguration{}

	ret.stack, _ = cmd.Flags().GetString(StackFlag)
	if ret.stack == "" {
		return nil, errors.New("no stack specified")
	}
	ret.debug = service.IsDebug(cmd)
	moduleURLTemplate, _ := cmd.Flags().GetString(ModuleURLTplFlag)
	if moduleURLTemplate == "" {
		return nil, errors.New("no module url template specified")
	}

	ret.pullConfiguration.ModuleURLTpl, err = template.New("").Parse(moduleURLTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse module url template")
	}
	if err := ret.pullConfiguration.ModuleURLTpl.Execute(io.Discard, map[string]interface{}{
		"module": "test",
	}); err != nil {
		return nil, err
	}
	ret.pullConfiguration.PullPageSize, _ = cmd.Flags().GetInt(ModulePullPageSizeFlag)

	ret.authConfig.ClientID, _ = cmd.Flags().GetString(stackClientIDFlag)
	if ret.authConfig.ClientID != "" {
		ret.authConfig.ClientSecret, _ = cmd.Flags().GetString(stackClientSecretFlag)
		if ret.authConfig.ClientSecret == "" {
			return nil, errors.New("no stack client secret specified")
		}

		ret.authConfig.Issuer, _ = cmd.Flags().GetString(stackIssuerFlag)
		if ret.authConfig.Issuer == "" {
			return nil, errors.New("no stack issuer specified")
		}
	}

	ret.bind, _ = cmd.Flags().GetString(BindFlag)

	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "parsing postgres connection flags")
	}
	ret.connectionOptions = *connectionOptions

	return ret, nil
}

func serve(cmd *cobra.Command, _ []string) error {

	serveConfiguration, err := evalServeConfiguration(cmd)
	if err != nil {
		return err
	}

	return service.New(cmd.OutOrStdout(),
		publish.FXModuleFromFlags(cmd, service.IsDebug(cmd)),
		otlptraces.FXModuleFromFlags(cmd),
		auth.FXModuleFromFlags(cmd),
		fx.Supply(drivers.ServiceConfig{
			Stack: serveConfiguration.stack,
			Debug: serveConfiguration.debug,
		}),
		fx.Supply(serveConfiguration.authConfig),
		fx.Invoke(all.Register),
		storage.Module(serveConfiguration.debug, serveConfiguration.connectionOptions),
		fx.Provide(httpclient.NewClient),
		fx.Provide(httpclient.NewStackAuthenticatedClient),
		controller.NewModule(),
		runner.NewModule(serveConfiguration.stack, serveConfiguration.pullConfiguration),
		fx.Supply(sharedapi.ServiceInfo{
			Version: Version,
			Debug:   serveConfiguration.debug,
		}),
		api.NewModule(serveConfiguration.bind),
	).Run(cmd)
}

func Execute() {
	service.Execute(NewRootCommand())
}
