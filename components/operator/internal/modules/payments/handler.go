package payments

import (
	"context"
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

type module struct{}

func (p module) Name() string {
	return "payments"
}

func (p module) Postgres(ctx modules.ReconciliationConfig) v1beta3.PostgresConfig {
	return ctx.Configuration.Spec.Services.Payments.Postgres
}

func (p module) Versions() map[string]modules.Version {
	return map[string]modules.Version{
		"v0.0.0": {
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				migrateCommand := []string{"payments", "migrate"}
				if ctx.Versions.IsHigherOrEqual("payments", "v0.7.0") {
					migrateCommand = append(migrateCommand, "up")
				}
				return modules.Services{{
					InjectPostgresVariables: true,
					HasVersionEndpoint:      true,
					ListenEnvVar:            "LISTEN",
					ExposeHTTP:              modules.DefaultExposeHTTP,
					NeedTopic:               true,
					Liveness:                modules.LivenessLegacy,
					Annotations:             ctx.Configuration.Spec.Services.Payments.Annotations.Service,
					Container: func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
						return modules.Container{
							Env:   paymentsEnvVars(resolveContext),
							Image: modules.GetImage("payments", resolveContext.Versions.Spec.Payments),
							Resources: modules.GetResourcesWithDefault(
								resolveContext.Configuration.Spec.Services.Payments.ResourceProperties,
								modules.ResourceSizeSmall(),
							),
						}
					},
					InitContainer: func(resolveContext modules.ContainerResolutionConfiguration) []modules.Container {
						return []modules.Container{{
							Name:    "migrate",
							Image:   modules.GetImage("payments", resolveContext.Versions.Spec.Payments),
							Env:     paymentsEnvVars(resolveContext),
							Command: migrateCommand,
						}}
					},
				}}
			},
		},
		"v0.6.5": {
			PreUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				return paymentsPreUpgradeMigration(ctx, config)
			},
			PostUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				return resetConnectors(ctx, config)
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.6.7": {
			PreUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				return paymentsPreUpgradeMigration(ctx, config)
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.6.8": {
			PreUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				return paymentsPreUpgradeMigration(ctx, config)
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.7.0": {
			PreUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				return paymentsPreUpgradeMigration(ctx, config)
			},
			PostUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				return resetConnectors(ctx, config)
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.8.0": {
			PreUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				// Add payment accounts
				return paymentsPreUpgradeMigration(ctx, config)
			},
			PostUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				return resetConnectors(ctx, config)
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.8.1": {
			PostUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				return resetConnectors(ctx, config)
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.9.0": {
			PreUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				// Add payment accounts
				return paymentsPreUpgradeMigration(ctx, config)
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.9.1": {
			PreUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				// Add payment accounts
				return paymentsPreUpgradeMigration(ctx, config)
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.9.4": {
			PreUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				// Add payment accounts
				return paymentsPreUpgradeMigration(ctx, config)
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.10.0": {
			PreUpgrade: func(ctx context.Context, config modules.ReconciliationConfig) error {
				// Add payment accounts
				return paymentsPreUpgradeMigration(ctx, config)
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
	}
}

var Module = &module{}

var _ modules.Module = Module
var _ modules.PostgresAwareModule = Module

func init() {
	modules.Register(Module)
}

func paymentsEnvVars(resolveContext modules.ContainerResolutionConfiguration) modules.ContainerEnv {
	return modules.BrokerEnvVars(resolveContext.Configuration.Spec.Broker, "payments").
		Append(
			modules.Env("POSTGRES_DATABASE_NAME", "$(POSTGRES_DATABASE)"),
			modules.Env("CONFIG_ENCRYPTION_KEY", resolveContext.Configuration.Spec.Services.Payments.EncryptionKey),
			modules.Env("PUBLISHER_TOPIC_MAPPING", "*:"+resolveContext.Stack.GetServiceName("payments")),
		)
}

func paymentsPreUpgradeMigration(ctx context.Context, config modules.ReconciliationConfig) error {
	postgresUri := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		config.Configuration.Spec.Services.Payments.Postgres.Username,
		config.Configuration.Spec.Services.Payments.Postgres.Password,
		config.Configuration.Spec.Services.Payments.Postgres.Host,
		config.Configuration.Spec.Services.Payments.Postgres.Port,
		config.Stack.GetServiceName("payments"),
	)

	db, err := sql.Open("postgres", postgresUri)
	if err != nil {
		return err
	}

	bunDB := bun.NewDB(db, pgdialect.New())
	defer bunDB.Close()

	return cmd.Migrate(ctx, bunDB)
}

func paymentsServices(
	env func(resolveContext modules.ContainerResolutionConfiguration) modules.ContainerEnv,
) modules.Services {
	return modules.Services{{
		InjectPostgresVariables: true,
		HasVersionEndpoint:      true,
		ListenEnvVar:            "LISTEN",
		ExposeHTTP:              modules.DefaultExposeHTTP,
		NeedTopic:               true,
		Liveness:                modules.LivenessLegacy,
		Container: func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
			return modules.Container{
				Env:   env(resolveContext),
				Image: modules.GetImage("payments", resolveContext.Versions.Spec.Payments),
				Resources: modules.GetResourcesWithDefault(
					resolveContext.Configuration.Spec.Services.Payments.ResourceProperties,
					modules.ResourceSizeSmall(),
				),
			}
		},
	}}
}

func resetConnectors(ctx context.Context, config modules.ReconciliationConfig) error {
	if err := resetConnector(ctx, config, "stripe"); err != nil {
		return err
	}
	if err := resetConnector(ctx, config, "wise"); err != nil {
		return err
	}
	if err := resetConnector(ctx, config, "modulr"); err != nil {
		return err
	}
	if err := resetConnector(ctx, config, "banking-circle"); err != nil {
		return err
	}
	if err := resetConnector(ctx, config, "currency-cloud"); err != nil {
		return err
	}
	if err := resetConnector(ctx, config, "dummy-pay"); err != nil {
		return err
	}
	if err := resetConnector(ctx, config, "mangopay"); err != nil {
		return err
	}
	if err := resetConnector(ctx, config, "moneycorp"); err != nil {
		return err
	}

	return nil
}

func resetConnector(ctx context.Context, config modules.ReconciliationConfig, connector string) error {
	endpoint := fmt.Sprintf(
		"http://payments.%s.svc:%d/connectors/%s/reset",
		config.Stack.Name,
		config.Stack.Status.Ports["payments"]["payments"],
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
