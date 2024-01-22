package configurations

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/api/stack.formance.com/v1beta3"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

// +kubebuilder:rbac:groups=stack.formance.com,resources=configurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stack.formance.com,resources=configurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stack.formance.com,resources=configurations/finalizers,verbs=update

func Reconcile(ctx Context, configuration *v1beta3.Configuration) error {

	stacks := &v1beta3.StackList{}
	if err := ctx.GetClient().List(ctx, stacks, client.MatchingFields{
		".spec.seed": configuration.Name,
	}); err != nil {
		return err
	}

	stackNames := Map(stacks.Items, func(from v1beta3.Stack) string {
		return from.GetName()
	})

	type databaseDescriptor struct {
		config v1beta3.PostgresConfig
		name   string
	}

	for _, cfg := range []databaseDescriptor{
		{
			config: configuration.Spec.Services.Ledger.Postgres,
			name:   "ledger",
		},
		{
			config: configuration.Spec.Services.Payments.Postgres,
			name:   "payments",
		},
		{
			config: configuration.Spec.Services.Orchestration.Postgres,
			name:   "orchestration",
		},
		{
			config: configuration.Spec.Services.Auth.Postgres,
			name:   "auth",
		},
		{
			config: configuration.Spec.Services.Webhooks.Postgres,
			name:   "webhooks",
		},
		{
			config: configuration.Spec.Services.Reconciliation.Postgres,
			name:   "reconciliation",
		},
	} {
		_, err := settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-%s-database-host", configuration.Name, cfg.name),
			fmt.Sprintf("databases.%s.host", cfg.name), cfg.config.Host, stackNames...)
		if err != nil {
			return err
		}

		if cfg.config.Port != 0 && cfg.config.Port != 5432 {
			_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-%s-database-port", configuration.Name, cfg.name),
				fmt.Sprintf("databases.%s.port", cfg.name), cfg.config.Port, stackNames...)
			if err != nil {
				return err
			}
		}

		if cfg.config.Username != "" {
			_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-%s-database-username", configuration.Name, cfg.name),
				fmt.Sprintf("databases.%s.username", cfg.name), cfg.config.Username, stackNames...)
			if err != nil {
				return err
			}
		}
		if cfg.config.Password != "" {
			_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-%s-database-password", configuration.Name, cfg.name),
				fmt.Sprintf("databases.%s.password", cfg.name), cfg.config.Password, stackNames...)
			if err != nil {
				return err
			}
		}
		if cfg.config.CredentialsFromSecret != "" {
			_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-%s-database-secret", configuration.Name, cfg.name),
				fmt.Sprintf("databases.%s.secret", cfg.name), cfg.config.CredentialsFromSecret, stackNames...)
			if err != nil {
				return err
			}
		}
		if cfg.config.DisableSSLMode {
			_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-%s-database-ssl", configuration.Name, cfg.name),
				fmt.Sprintf("databases.%s.ssl.disable", cfg.name), cfg.config.DisableSSLMode, stackNames...)
			if err != nil {
				return err
			}
		}
	}

	if configuration.Spec.Monitoring != nil {
		_, _, err := CreateOrUpdate[*v1beta1.OpenTelemetryConfiguration](ctx, types.NamespacedName{
			Name: configuration.Name,
		}, func(t *v1beta1.OpenTelemetryConfiguration) {
			t.Spec = v1beta1.OpenTelemetryConfigurationSpec{
				Traces: func() *v1beta1.TracesSpec {
					traces := configuration.Spec.Monitoring.Traces
					if traces == nil {
						return nil
					}
					return &v1beta1.TracesSpec{
						Otlp: convertOtlpSpec(traces.Otlp),
					}
				}(),
				Metrics: func() *v1beta1.MetricsSpec {
					metrics := configuration.Spec.Monitoring.Metrics
					if metrics == nil {
						return nil
					}
					return &v1beta1.MetricsSpec{
						Otlp: convertOtlpSpec(metrics.Otlp),
					}
				}(),
				ConfigurationProperties: v1beta1.ConfigurationProperties{
					Stacks: Map(stacks.Items, func(from v1beta3.Stack) string {
						return from.GetName()
					}),
				},
			}
		})
		if err != nil {
			return errors.Wrap(err, "creating opentelemetry configuration for service")
		}
	}

	if configuration.Spec.Broker.Kafka != nil {
		_, err := settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-broker-kind", configuration.Name),
			"broker.kind", "kafka", stackNames...)
		if err != nil {
			return err
		}

		if len(configuration.Spec.Broker.Kafka.Brokers) > 0 {
			_, err := settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-broker-kind", configuration.Name),
				"broker.kafka.endpoints", strings.Join(configuration.Spec.Broker.Kafka.Brokers, ","), stackNames...)
			if err != nil {
				return err
			}
		}

		if configuration.Spec.Broker.Kafka.TLS {
			_, err := settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-broker-ssl", configuration.Name),
				"broker.kafka.ssl.enabled", "true", stackNames...)
			if err != nil {
				return err
			}
		}

		if sasl := configuration.Spec.Broker.Kafka.SASL; sasl != nil {
			_, err := settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-broker-kafka-username", configuration.Name),
				"broker.kafka.sasl.username", sasl.Username, stackNames...)
			if err != nil {
				return err
			}
			_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-broker-kafka-password", configuration.Name),
				"broker.kafka.sasl.password", sasl.Password, stackNames...)
			if err != nil {
				return err
			}
			_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-broker-kafka-mechanism", configuration.Name),
				"broker.kafka.sasl.mechanism", sasl.Mechanism, stackNames...)
			if err != nil {
				return err
			}
			_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-broker-kafka-scram-sha-size", configuration.Name),
				"broker.kafka.sasl.scram-sha-size", sasl.ScramSHASize, stackNames...)
			if err != nil {
				return err
			}
		}
	} else if configuration.Spec.Broker.Nats != nil {
		_, err := settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-broker-kind", configuration.Name),
			"broker.kind", "nats", stackNames...)
		if err != nil {
			return err
		}

		_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-broker-nats-url", configuration.Name),
			"broker.nats.endpoint", configuration.Spec.Broker.Nats.URL, stackNames...)
		if err != nil {
			return err
		}

		_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-broker-nats-replicas", configuration.Name),
			"broker.nats.replicas", configuration.Spec.Broker.Nats.Replicas, stackNames...)
		if err != nil {
			return err
		}
	}

	_, _, err := CreateOrUpdate[*v1beta1.TemporalConfiguration](ctx, types.NamespacedName{
		Name: configuration.Name,
	}, func(t *v1beta1.TemporalConfiguration) {
		t.Spec = v1beta1.TemporalConfigurationSpec{
			ConfigurationProperties: v1beta1.ConfigurationProperties{
				Stacks: Map(stacks.Items, func(from v1beta3.Stack) string {
					return from.GetName()
				}),
			},
			Address:   configuration.Spec.Temporal.Address,
			Namespace: configuration.Spec.Temporal.Namespace,
			TLS: v1beta1.TemporalTLSConfig{
				CRT:        configuration.Spec.Temporal.TLS.CRT,
				Key:        configuration.Spec.Temporal.TLS.Key,
				SecretName: configuration.Spec.Temporal.TLS.SecretName,
			},
		}
	})
	if err != nil {
		return errors.Wrap(err, "creating temporal configuration for service")
	}

	_, _, err = CreateOrUpdate[*v1beta1.ElasticSearchConfiguration](ctx, types.NamespacedName{
		Name: configuration.Name,
	}, func(t *v1beta1.ElasticSearchConfiguration) {
		t.Spec = v1beta1.ElasticSearchConfigurationSpec{
			Scheme: configuration.Spec.Services.Search.ElasticSearchConfig.Scheme,
			Host:   configuration.Spec.Services.Search.ElasticSearchConfig.Host,
			Port:   configuration.Spec.Services.Search.ElasticSearchConfig.Port,
			TLS: v1beta1.ElasticSearchTLSConfig{
				Enabled:        configuration.Spec.Services.Search.ElasticSearchConfig.TLS.Enabled,
				SkipCertVerify: configuration.Spec.Services.Search.ElasticSearchConfig.TLS.SkipCertVerify,
			},
			BasicAuth: func() *v1beta1.ElasticSearchBasicAuthConfig {
				if configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth == nil {
					return nil
				}
				return &v1beta1.ElasticSearchBasicAuthConfig{
					Username:   configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Username,
					Password:   configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Password,
					SecretName: configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.SecretName,
				}
			}(),
			ConfigurationProperties: v1beta1.ConfigurationProperties{
				Stacks: Map(stacks.Items, func(from v1beta3.Stack) string {
					return from.GetName()
				}),
			},
		}
	})
	if err != nil {
		return errors.Wrap(err, "creating elasticsearch configuration for service")
	}

	_, _, err = CreateOrUpdate[*v1beta1.SearchBatchingConfiguration](ctx, types.NamespacedName{
		Name: configuration.Name,
	}, func(t *v1beta1.SearchBatchingConfiguration) {
		t.Spec = v1beta1.SearchBatchingConfigurationSpec{
			ConfigurationProperties: v1beta1.ConfigurationProperties{
				Stacks: Map(stacks.Items, func(from v1beta3.Stack) string {
					return from.GetName()
				}),
			},
			Batching: v1beta1.Batching{
				Count:  configuration.Spec.Services.Search.Batching.Count,
				Period: configuration.Spec.Services.Search.Batching.Period,
			},
		}
	})
	if err != nil {
		return errors.Wrap(err, "creating registries configuration")
	}

	if len(configuration.Spec.Registries) > 0 {
		_, _, err = CreateOrUpdate[*v1beta1.RegistriesConfiguration](ctx, types.NamespacedName{
			Name: configuration.Name,
		}, func(t *v1beta1.RegistriesConfiguration) {
			registries := make(map[string]v1beta1.RegistryConfigurationSpec)
			for k, config := range configuration.Spec.Registries {
				registries[k] = v1beta1.RegistryConfigurationSpec{
					Endpoint: config.Endpoint,
				}
			}
			t.Spec = v1beta1.RegistriesConfigurationSpec{
				ConfigurationProperties: v1beta1.ConfigurationProperties{
					Stacks: Map(stacks.Items, func(from v1beta3.Stack) string {
						return from.GetName()
					}),
				},
				Registries: registries,
			}
		})
		if err != nil {
			return errors.Wrap(err, "creating registries configuration")
		}
	}

	return nil
}

func convertOtlpSpec(otlp *v1beta3.OtlpSpec) *v1beta1.OtlpSpec {
	if otlp == nil {
		return nil
	}

	return &v1beta1.OtlpSpec{
		Endpoint: otlp.Endpoint,
		Port:     otlp.Port,
		Insecure: otlp.Insecure,
		Mode:     otlp.Mode,
		ResourceAttributes: func() map[string]string {
			if otlp.ResourceAttributes == "" {
				return nil
			}
			parts := strings.Split(otlp.ResourceAttributes, " ")
			ret := make(map[string]string)
			for _, part := range parts {
				parts := strings.Split(part, "=")
				ret[parts[0]] = parts[1]
			}

			return ret
		}(),
	}
}

func init() {
	Init(
		WithReconciler[*v1beta3.Configuration](Reconcile,
			WithWatch[*v1beta3.Stack](func(ctx Context, object *v1beta3.Stack) []reconcile.Request {
				return []reconcile.Request{{
					NamespacedName: types.NamespacedName{
						Name: object.Spec.Seed,
					},
				}}
			}),
		),
	)
}
