package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
	"github.com/formancehq/payments/cmd"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func init() {
	env := func(resolveContext modules.ContainerResolutionContext) modules.ContainerEnv {
		return modules.BrokerEnvVars(resolveContext.Configuration.Spec.Broker, "payments").
			Append(
				modules.Env("POSTGRES_DATABASE_NAME", "$(POSTGRES_DATABASE)"),
				modules.Env("CONFIG_ENCRYPTION_KEY", resolveContext.Configuration.Spec.Services.Payments.EncryptionKey),
				modules.Env("PUBLISHER_TOPIC_MAPPING", "*:"+resolveContext.Stack.GetServiceName("payments")),
			)
	}
	modules.Register("payments", modules.Module{
		Postgres: func(ctx modules.Context) v1beta3.PostgresConfig {
			return ctx.Configuration.Spec.Services.Payments.Postgres
		},
		Versions: map[string]modules.Version{
			"v0.0.0": {
				Services: func(ctx modules.ModuleContext) modules.Services {
					migrateCommand := []string{"payments", "migrate"}
					if ctx.HasVersionLower("v0.7.0") {
						migrateCommand = append(migrateCommand, "up")
					}
					return modules.Services{{
						InjectPostgresVariables: true,
						HasVersionEndpoint:      true,
						ListenEnvVar:            "LISTEN",
						ExposeHTTP:              true,
						NeedTopic:               true,
						Liveness:                modules.LivenessLegacy,
						Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
							return modules.Container{
								Env:       env(resolveContext),
								Image:     modules.GetImage("payments", resolveContext.Versions.Spec.Payments),
								Resources: modules.ResourceSizeSmall(),
							}
						},
						InitContainer: func(resolveContext modules.ContainerResolutionContext) []modules.Container {
							return []modules.Container{{
								Name:    "migrate",
								Image:   modules.GetImage("payments", resolveContext.Versions.Spec.Payments),
								Env:     env(resolveContext),
								Command: migrateCommand,
							}}
						},
					}}
				},
			},
			"v0.6.5": {
				PreUpgrade: func(ctx modules.Context) error {
					return paymentsPreUpgradeMigration(ctx)
				},
				PostUpgrade: func(ctx modules.PostInstallContext) error {
					return resetConnectors(ctx)
				},
				Services: func(ctx modules.ModuleContext) modules.Services {
					return paymentsServices(ctx, env)
				},
			},
			"v0.6.7": {
				PreUpgrade: func(ctx modules.Context) error {
					return paymentsPreUpgradeMigration(ctx)
				},
				Services: func(ctx modules.ModuleContext) modules.Services {
					return paymentsServices(ctx, env)
				},
			},
			"v0.6.8": {
				PreUpgrade: func(ctx modules.Context) error {
					return paymentsPreUpgradeMigration(ctx)
				},
				Services: func(ctx modules.ModuleContext) modules.Services {
					return paymentsServices(ctx, env)
				},
			},
			"v0.7.0": {
				PreUpgrade: func(ctx modules.Context) error {
					return paymentsPreUpgradeMigration(ctx)
				},
				PostUpgrade: func(ctx modules.PostInstallContext) error {
					return resetConnectors(ctx)
				},
				Services: func(ctx modules.ModuleContext) modules.Services {
					return paymentsServices(ctx, env)
				},
			},
			"v0.8.0": {
				PreUpgrade: func(ctx modules.Context) error {
					// Add payment accounts
					return paymentsPreUpgradeMigration(ctx)
				},
				PostUpgrade: func(ctx modules.PostInstallContext) error {
					return resetConnectors(ctx)
				},
				Services: func(ctx modules.ModuleContext) modules.Services {
					return paymentsServices(ctx, env)
				},
			},
			"v0.8.1": {
				PostUpgrade: func(ctx modules.PostInstallContext) error {
					return resetConnectors(ctx)
				},
				Services: func(ctx modules.ModuleContext) modules.Services {
					return paymentsServices(ctx, env)
				},
			},
			"v0.9.0": {
				PreUpgrade: func(ctx modules.Context) error {
					// Add payment accounts
					return paymentsPreUpgradeMigration(ctx)
				},
				Services: func(ctx modules.ModuleContext) modules.Services {
					return paymentsServices(ctx, env)
				},
			},
			"v0.9.1": {
				PreUpgrade: func(ctx modules.Context) error {
					// Add payment accounts
					return paymentsPreUpgradeMigration(ctx)
				},
				Services: func(ctx modules.ModuleContext) modules.Services {
					return paymentsServices(ctx, env)
				},
			},
			"v0.9.4": {
				PreUpgrade: func(ctx modules.Context) error {
					// Add payment accounts
					return paymentsPreUpgradeMigration(ctx)
				},
				Services: func(ctx modules.ModuleContext) modules.Services {
					return paymentsServices(ctx, env)
				},
			},
		},
	})
}

func paymentsPreUpgradeMigration(ctx modules.Context) error {
	postgresUri := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		ctx.Configuration.Spec.Services.Payments.Postgres.Username,
		ctx.Configuration.Spec.Services.Payments.Postgres.Password,
		ctx.Configuration.Spec.Services.Payments.Postgres.Host,
		ctx.Configuration.Spec.Services.Payments.Postgres.Port,
		ctx.Stack.GetServiceName("payments"),
	)

	db, err := sql.Open("postgres", postgresUri)
	if err != nil {
		return err
	}

	bunDB := bun.NewDB(db, pgdialect.New())
	defer bunDB.Close()

	return cmd.Migrate(ctx.Context, bunDB)
}

func paymentsServices(
	ctx modules.ModuleContext,
	env func(resolveContext modules.ContainerResolutionContext) modules.ContainerEnv,
) modules.Services {
	return modules.Services{{
		InjectPostgresVariables: true,
		HasVersionEndpoint:      true,
		ListenEnvVar:            "LISTEN",
		ExposeHTTP:              true,
		NeedTopic:               true,
		Liveness:                modules.LivenessLegacy,
		Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
			return modules.Container{
				Env:       env(resolveContext),
				Image:     modules.GetImage("payments", resolveContext.Versions.Spec.Payments),
				Resources: modules.ResourceSizeSmall(),
			}
		},
	}}
}

func resetConnectors(ctx modules.PostInstallContext) error {
	if err := resetConnector(ctx, "stripe"); err != nil {
		return err
	}
	if err := resetConnector(ctx, "wise"); err != nil {
		return err
	}
	if err := resetConnector(ctx, "modulr"); err != nil {
		return err
	}
	if err := resetConnector(ctx, "banking-circle"); err != nil {
		return err
	}
	if err := resetConnector(ctx, "currency-cloud"); err != nil {
		return err
	}
	if err := resetConnector(ctx, "dummy-pay"); err != nil {
		return err
	}
	if err := resetConnector(ctx, "mangopay"); err != nil {
		return err
	}
	if err := resetConnector(ctx, "moneycorp"); err != nil {
		return err
	}

	return nil
}

func resetConnector(ctx modules.PostInstallContext, connector string) error {
	endpoint := fmt.Sprintf(
		"http://payments.%s.svc:%d/connectors/%s/reset",
		ctx.Stack.Name,
		ctx.Stack.Status.Ports[ctx.ModuleName]["payments"],
		connector,
	)
	res, err := http.Post(endpoint, "", nil)
	if err != nil {
		logger := log.FromContext(ctx)
		logger.WithValues("endpoint", endpoint).Error(err, "failed to reset connector")
		return err
	}

	switch res.StatusCode {
	case http.StatusOK, http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		// Connector is not installed, we can directly return nil, nothing to do
		return nil
	default:
		// Return an error to retry the migration. It can be the case when the
		// pod is up, but not the http server.
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
}
