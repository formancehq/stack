package configurations

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/api/stack.formance.com/v1beta3"
	. "github.com/formancehq/operator/internal/core"
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
		_, _, err := CreateOrUpdate[*v1beta1.DatabaseConfiguration](ctx, types.NamespacedName{
			Name: fmt.Sprintf("%s-%s", configuration.Name, cfg.name),
		}, func(t *v1beta1.DatabaseConfiguration) {
			if t.Labels == nil {
				t.Labels = map[string]string{}
			}
			t.Spec = v1beta1.DatabaseConfigurationSpec{
				Port:                  cfg.config.Port,
				Host:                  cfg.config.Host,
				Username:              cfg.config.Username,
				Password:              cfg.config.Password,
				CredentialsFromSecret: cfg.config.CredentialsFromSecret,
				DisableSSLMode:        cfg.config.DisableSSLMode,
				ConfigurationProperties: v1beta1.ConfigurationProperties{
					Stacks: Map(stacks.Items, func(from v1beta3.Stack) string {
						return from.GetName()
					}),
				},
				Service: cfg.name,
			}
		})
		if err != nil {
			return errors.Wrapf(err, "creating database configuration for service %s", cfg.name)
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

	_, _, err := CreateOrUpdate[*v1beta1.BrokerConfiguration](ctx, types.NamespacedName{
		Name: configuration.Name,
	}, func(t *v1beta1.BrokerConfiguration) {
		t.Spec = v1beta1.BrokerConfigurationSpec{
			ConfigurationProperties: v1beta1.ConfigurationProperties{
				Stacks: Map(stacks.Items, func(from v1beta3.Stack) string {
					return from.GetName()
				}),
			},
			Kafka: func() *v1beta1.BrokerKafkaConfig {
				if configuration.Spec.Broker.Kafka == nil {
					return nil
				}
				return &v1beta1.BrokerKafkaConfig{
					Brokers: configuration.Spec.Broker.Kafka.Brokers,
					TLS:     configuration.Spec.Broker.Kafka.TLS,
					SASL: func() *v1beta1.BrokerKafkaSASLConfig {
						if configuration.Spec.Broker.Kafka.SASL == nil {
							return nil
						}
						return &v1beta1.BrokerKafkaSASLConfig{
							Username:     configuration.Spec.Broker.Kafka.SASL.Username,
							Password:     configuration.Spec.Broker.Kafka.SASL.Password,
							Mechanism:    configuration.Spec.Broker.Kafka.SASL.Mechanism,
							ScramSHASize: configuration.Spec.Broker.Kafka.SASL.ScramSHASize,
						}
					}(),
				}
			}(),
			Nats: func() *v1beta1.BrokerNatsConfig {
				if configuration.Spec.Broker.Nats == nil {
					return nil
				}
				return &v1beta1.BrokerNatsConfig{
					URL:      configuration.Spec.Broker.Nats.URL,
					Replicas: configuration.Spec.Broker.Nats.Replicas,
				}
			}(),
		}
	})
	if err != nil {
		return errors.Wrap(err, "creating broker configuration for service")
	}

	_, _, err = CreateOrUpdate[*v1beta1.TemporalConfiguration](ctx, types.NamespacedName{
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
