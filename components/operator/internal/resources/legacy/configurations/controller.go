package configurations

import (
	"fmt"
	"github.com/formancehq/operator/api/stack.formance.com/v1beta3"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
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

	if monitoring := configuration.Spec.Monitoring; monitoring != nil {

		createSettings := func(discr string) error {
			_, err := settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-otel-traces-enabled", configuration.Name),
				fmt.Sprintf("opentelemetry.%s.enabled", discr), "true", stackNames...)
			if err != nil {
				return err
			}
			_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-otel-traces-endpoint", configuration.Name),
				fmt.Sprintf("opentelemetry.%s.endpoint", discr), monitoring.Traces.Otlp.Endpoint, stackNames...)
			if err != nil {
				return err
			}
			if monitoring.Traces.Otlp.Insecure {
				_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-otel-traces-insecure", configuration.Name),
					fmt.Sprintf("opentelemetry.%s.insecure", discr), "true", stackNames...)
				if err != nil {
					return err
				}
			}
			if monitoring.Traces.Otlp.Port != 4317 {
				_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-otel-traces-port", configuration.Name),
					fmt.Sprintf("opentelemetry.%s.port", discr), fmt.Sprint(monitoring.Traces.Otlp.Port), stackNames...)
				if err != nil {
					return err
				}
			}
			if monitoring.Traces.Otlp.Mode != "grpc" {
				_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-otel-traces-mode", configuration.Name),
					fmt.Sprintf("opentelemetry.%s.insecure", discr), monitoring.Traces.Otlp.Mode, stackNames...)
				if err != nil {
					return err
				}
			}
			if monitoring.Traces.Otlp.ResourceAttributes != "" {
				_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-otel-traces-mode", configuration.Name),
					fmt.Sprintf("opentelemetry.%s.resource-attributes", discr), monitoring.Traces.Otlp.ResourceAttributes, stackNames...)
				if err != nil {
					return err
				}
			}

			return nil
		}

		if monitoring.Traces != nil && monitoring.Traces.Otlp != nil {
			if err := createSettings("traces"); err != nil {
				return err
			}
		}

		if monitoring.Metrics != nil && monitoring.Metrics.Otlp != nil {
			if err := createSettings("metrics"); err != nil {
				return err
			}
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

	_, err := settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-temporal-address", configuration.Name),
		"temporal-address", configuration.Spec.Temporal.Address, stackNames...)
	if err != nil {
		return err
	}

	_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-temporal-namespace", configuration.Name),
		"temporal-namespace", configuration.Spec.Temporal.Namespace, stackNames...)
	if err != nil {
		return err
	}

	if configuration.Spec.Temporal.TLS.CRT != "" {
		_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-temporal-tls-crt", configuration.Name),
			"temporal.tls.crt", configuration.Spec.Temporal.TLS.CRT, stackNames...)
		if err != nil {
			return err
		}
		_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-temporal-tls-key", configuration.Name),
			"temporal.tls.key", configuration.Spec.Temporal.TLS.Key, stackNames...)
		if err != nil {
			return err
		}
	} else if configuration.Spec.Temporal.TLS.SecretName != "" {
		_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-temporal-tls-secret", configuration.Name),
			"temporal.tls.secret", configuration.Spec.Temporal.TLS.SecretName, stackNames...)
		if err != nil {
			return err
		}
	}

	_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-elasticsearch-host", configuration.Name),
		"elasticsearch.host", configuration.Spec.Services.Search.ElasticSearchConfig.Host, stackNames...)
	if err != nil {
		return err
	}

	if configuration.Spec.Services.Search.ElasticSearchConfig.Scheme != "" {
		_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-elasticsearch-scheme", configuration.Name),
			"elasticsearch.scheme", configuration.Spec.Services.Search.ElasticSearchConfig.Scheme, stackNames...)
		if err != nil {
			return err
		}
	}

	if configuration.Spec.Services.Search.ElasticSearchConfig.Port != 9200 && configuration.Spec.Services.Search.ElasticSearchConfig.Port != 0 {
		_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-elasticsearch-port", configuration.Name),
			"elasticsearch.port", configuration.Spec.Services.Search.ElasticSearchConfig.Scheme, stackNames...)
		if err != nil {
			return err
		}
	}

	if configuration.Spec.Services.Search.ElasticSearchConfig.TLS.Enabled {
		_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-elasticsearch-tls-enabled", configuration.Name),
			"elasticsearch.tls.enabled", "true", stackNames...)
		if err != nil {
			return err
		}
		if configuration.Spec.Services.Search.ElasticSearchConfig.TLS.SkipCertVerify {
			_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-elasticsearch-tls-skip-cert-verify", configuration.Name),
				"elasticsearch.tls.skip-cert-verify", "true", stackNames...)
			if err != nil {
				return err
			}
		}
	}

	if basicAuth := configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth; basicAuth != nil {
		if basicAuth.Username != "" {
			_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-elasticsearch-basic-auth-username", configuration.Name),
				"elasticsearch.basic-auth.username", basicAuth.Username, stackNames...)
			if err != nil {
				return err
			}
			_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-elasticsearch-basic-auth-password", configuration.Name),
				"elasticsearch.basic-auth.password", basicAuth.Password, stackNames...)
			if err != nil {
				return err
			}
		}
		if basicAuth.SecretName != "" {
			_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-elasticsearch-basic-auth-secret", configuration.Name),
				"elasticsearch.basic-auth.secret", basicAuth.SecretName, stackNames...)
			if err != nil {
				return err
			}
		}
	}

	if configuration.Spec.Services.Search.Batching.Count != 0 {
		_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-search-batching-count", configuration.Name),
			"search.batching.count", fmt.Sprint(configuration.Spec.Services.Search.Batching.Count), stackNames...)
		if err != nil {
			return err
		}
	}

	if configuration.Spec.Services.Search.Batching.Period != "" {
		_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-search-batching-count", configuration.Name),
			"search.batching.period", configuration.Spec.Services.Search.Batching.Period, stackNames...)
		if err != nil {
			return err
		}
	}

	if len(configuration.Spec.Registries) > 0 {
		for name, config := range configuration.Spec.Registries {
			_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-registries-%s", configuration.Name, name),
				fmt.Sprintf("registries.%s.endpoint", name), config.Endpoint, stackNames...)
			if err != nil {
				return err
			}
		}
	}

	if configuration.Spec.Services.Payments.EncryptionKey != "" {
		_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-payments-encryption-key", configuration.Name),
			"payments.encryption-key", configuration.Spec.Services.Payments.EncryptionKey, stackNames...)
		if err != nil {
			return err
		}
	}

	return nil
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
