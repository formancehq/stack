package payments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
	"golang.org/x/mod/semver"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
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
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				Command:  []string{"migrate", "up"},
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
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
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				Command:  []string{"migrate", "up"},
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
			PostUpgrade: PostUpgradePreV1,
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.6.7": {
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				Command:  []string{"migrate", "up"},
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.6.8": {
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				Command:  []string{"migrate", "up"},
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.7.0": {
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
			PostUpgrade: PostUpgradePreV1,
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.8.0": {
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
			PostUpgrade: PostUpgradePreV1,
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.8.1": {
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
			PostUpgrade: PostUpgradePreV1,
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.9.0": {
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.9.1": {
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.9.4": {
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v0.10.0": {
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v1.0.0-alpha.1": {
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v1.0.0-alpha.3": {
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServices(paymentsEnvVars)
			},
		},
		"v1.0.0-alpha.6": {
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServicesSplitted(paymentsEnvVars)
			},
		},
		"v1.0.0-beta.3": {
			PreUpgrade: func(ctx context.Context, jobRunner modules.JobRunner, config modules.MigrationConfig) (bool, error) {
				ok, err := jobRunner.RunJob(ctx, "delete-search-data", nil, deleteSearchDataJob(config.ReconciliationConfig, config.ReconciliationConfig.Stack.Name))
				if err != nil {
					return false, err
				}
				if !ok {
					return false, nil
				}

				return true, nil
			},
			DatabaseMigration: &modules.DatabaseMigration{
				Shutdown: true,
				AdditionalEnv: func(ctx modules.ReconciliationConfig) []modules.EnvVar {
					return []modules.EnvVar{
						modules.Env("CONFIG_ENCRYPTION_KEY", ctx.Configuration.Spec.Services.Payments.EncryptionKey),
					}
				},
			},
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return paymentsServicesSplitted(paymentsEnvVars)
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

func paymentsServicesSplitted(
	env func(resolveContext modules.ContainerResolutionConfiguration) modules.ContainerEnv,
) modules.Services {
	return modules.Services{
		{
			Name:                    "api",
			InjectPostgresVariables: true,
			HasVersionEndpoint:      true,
			ListenEnvVar:            "LISTEN",
			ExposeHTTP:              modules.DefaultExposeHTTP,
			Paths: []modules.Path{
				{
					Path:    "/payments",
					Methods: []string{"GET", "PATCH"},
					Name:    "payments",
				},
				{
					Path:    "/accounts",
					Methods: []string{"GET"},
					Name:    "accounts",
				},
				{
					Path:    "/bank-accounts",
					Methods: []string{"GET"},
					Name:    "bank-accounts-read",
				},
				{
					Path:    "/transfer-initiations",
					Methods: []string{"GET"},
					Name:    "transfer-initiations-read",
				},
			},
			NeedTopic: true,
			Liveness:  modules.LivenessLegacy,
			Container: func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
				return modules.Container{
					Args:  []string{"api", "serve"},
					Env:   env(resolveContext),
					Image: modules.GetImage("payments", resolveContext.Versions.Spec.Payments),
					Resources: modules.GetResourcesWithDefault(
						resolveContext.Configuration.Spec.Services.Payments.ResourceProperties,
						modules.ResourceSizeSmall(),
					),
				}
			},
		},
		{
			Name:                    "connectors",
			InjectPostgresVariables: true,
			HasVersionEndpoint:      true,
			ListenEnvVar:            "LISTEN",
			ExposeHTTP:              modules.DefaultExposeHTTP,
			Paths: []modules.Path{
				{
					Path: "/connectors",
					Name: "connectors",
				},
				{
					Path: "/configs",
					Name: "configs",
				},
				{
					Path:    "/transfer-initiations",
					Methods: []string{"POST", "DELETE"},
					Name:    "transfer-initiations-write",
				},
				{
					Path:    "/bank-accounts",
					Methods: []string{"POST"},
					Name:    "bank-accounts-write",
				},
			},
			NeedTopic: true,
			Liveness:  modules.LivenessLegacy,
			Container: func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
				return modules.Container{
					Args:  []string{"connectors", "serve"},
					Env:   env(resolveContext),
					Image: modules.GetImage("payments", resolveContext.Versions.Spec.Payments),
					Resources: modules.GetResourcesWithDefault(
						resolveContext.Configuration.Spec.Services.Payments.ResourceProperties,
						modules.ResourceSizeSmall(),
					),
				}
			},
		},
	}
}

func PostUpgradePreV1(ctx context.Context, upgrader modules.JobRunner, config modules.MigrationConfig) (bool, error) {
	if config.Version == "latest" {
		return true, nil
	}

	// That mean the version is a commit, and we assume that a commit is always
	// the highest version.
	// cf https://github.com/formancehq/stack/blob/main/components/operator/apis/stack/v1beta3/versions_types.go#L89
	if !semver.IsValid(config.Version) {
		return true, nil
	}

	switch semver.Compare(config.Version, "v1.0.0-alpha.6") {
	case 0, 1:
		// Services are splitted, nothing to do
		return true, nil
	case -1:
		return resetConnectors(ctx, config)
	}

	return true, nil
}

func resetConnectors(ctx context.Context, config modules.MigrationConfig) (bool, error) {
	if err := resetConnector(ctx, config, "stripe"); err != nil {
		return false, err
	}
	if err := resetConnector(ctx, config, "wise"); err != nil {
		return false, err
	}
	if err := resetConnector(ctx, config, "modulr"); err != nil {
		return false, err
	}
	if err := resetConnector(ctx, config, "banking-circle"); err != nil {
		return false, err
	}
	if err := resetConnector(ctx, config, "currency-cloud"); err != nil {
		return false, err
	}
	if err := resetConnector(ctx, config, "dummy-pay"); err != nil {
		return false, err
	}
	if err := resetConnector(ctx, config, "mangopay"); err != nil {
		return false, err
	}
	if err := resetConnector(ctx, config, "moneycorp"); err != nil {
		return false, err
	}

	return true, nil
}

func resetConnector(ctx context.Context, config modules.MigrationConfig, connector string) error {
	endpoint := fmt.Sprintf(
		"http://payments.%s.svc:%d/connectors/%s/reset",
		config.ReconciliationConfig.Stack.Name,
		config.ReconciliationConfig.Stack.Status.Ports["payments"]["payments"],
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

func deleteSearchDataJob(config modules.ReconciliationConfig, name string) func(t *batchv1.Job) {
	return func(t *batchv1.Job) {
		t.Spec = batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{{
						Name:  "delete-index",
						Image: "curlimages/curl",
						Command: modules.ShellCommand(`
							curl -X POST -H 'Content-Type: application/json' ${OPEN_SEARCH_SCHEME}://${OPEN_SEARCH_SERVICE}/stacks/_delete_by_query -d '{"query": {"bool": {"must": [{"bool": {"should": [{"match": {"kind": "PAYMENT"}},{"match": {"kind": "PAYMENT_ACCOUNT"}},{"match": {"kind": "PAYMENT_BALANCE"}},{"match": {"kind": "PAYMENT_BANK_ACCOUNT"}},{"match": {"kind": "PAYMENT_TRANSFER_INITIATION"}}]}},{"match": {"stack": "%s"}}]}}}'`, name),
						Env: modules.SearchEnvVars(config).ToCoreEnv(),
					}},
				},
			},
		}
	}
}
