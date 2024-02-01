package configurations

import (
	"fmt"
	"net/url"

	"github.com/formancehq/operator/api/stack.formance.com/v1beta3"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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

	type resourceRequirementDescriptor struct {
		requirements *v1beta3.ResourceProperties
		deployment   string
	}
	for _, cfg := range []resourceRequirementDescriptor{
		{
			requirements: configuration.Spec.Services.Ledger.ResourceProperties,
			deployment:   "ledger",
		},
		{
			requirements: configuration.Spec.Services.Payments.ResourceProperties,
			deployment:   "payments",
		},
		{
			requirements: configuration.Spec.Services.Orchestration.ResourceProperties,
			deployment:   "orchestration",
		},
		{
			requirements: configuration.Spec.Services.Auth.ResourceProperties,
			deployment:   "auth",
		},
		{
			requirements: configuration.Spec.Services.Webhooks.ResourceProperties,
			deployment:   "webhooks",
		},
		{
			requirements: configuration.Spec.Services.Reconciliation.ResourceProperties,
			deployment:   "reconciliation",
		},
		{
			requirements: configuration.Spec.Services.Gateway.ResourceProperties,
			deployment:   "gateway",
		},
		{
			requirements: configuration.Spec.Services.Wallets.ResourceProperties,
			deployment:   "wallets",
		},
		{
			requirements: configuration.Spec.Services.Stargate.ResourceProperties,
			deployment:   "stargate",
		},
		{
			requirements: configuration.Spec.Services.Search.SearchResourceProperties,
			deployment:   "search",
		},
		{
			requirements: configuration.Spec.Services.Search.BenthosResourceProperties,
			deployment:   "benthos",
		},
	} {
		if cfg.requirements == nil {
			continue
		}
		var computeResourceList = func(resource *v1beta3.Resource) string {
			if resource == nil {
				return ""
			}
			limits := ""
			if resource.Cpu != "" {
				limits = limits + "cpu=" + resource.Cpu
			}
			if resource.Memory != "" {
				if limits != "" {
					limits = limits + ","
				}
				limits = limits + "memory=" + resource.Memory
			}
			return limits
		}

		if limits := computeResourceList(cfg.requirements.Limits); limits != "" {
			settingName := fmt.Sprintf("%s-%s-resource-limits", configuration.Name, cfg.deployment)
			settingKey := fmt.Sprintf("deployments.%s.containers.*.resource-requirements.limits", cfg.deployment)

			_, err := settings.CreateOrUpdate(ctx, settingName, settingKey, limits, stackNames...)
			if err != nil {
				return err
			}
		}

		if requests := computeResourceList(cfg.requirements.Request); requests != "" {
			settingName := fmt.Sprintf("%s-%s-resource-requests", configuration.Name, cfg.deployment)
			settingKey := fmt.Sprintf("deployments.%s.containers.*.resource-requirements.requests", cfg.deployment)

			_, err := settings.CreateOrUpdate(ctx, settingName, settingKey, requests, stackNames...)
			if err != nil {
				return err
			}
		}
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
		basicAuth := ""
		if cfg.config.Username != "" {
			basicAuth = fmt.Sprintf("%s:%s@", cfg.config.Username, cfg.config.Password)
		}
		endpoint := fmt.Sprintf("%s:%d", cfg.config.Host, cfg.config.Port)

		options := url.Values{}
		if cfg.config.DisableSSLMode {
			options.Set("disableSSLMode", "true")
		}
		if cfg.config.CredentialsFromSecret != "" {
			options.Set("secret", cfg.config.CredentialsFromSecret)
		}

		_, err := settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-%s-postgres-uri", configuration.Name, cfg.name),
			fmt.Sprintf("postgres.%s.uri", cfg.name),
			fmt.Sprintf("postgresql://%s%s?%s", basicAuth, endpoint, options.Encode()),
			stackNames...)
		if err != nil {
			return err
		}
	}

	if monitoring := configuration.Spec.Monitoring; monitoring != nil {

		var createSettings = func(discr string, spec *v1beta3.OtlpSpec) error {
			options := url.Values{}
			if spec.Insecure {
				options.Set("insecure", "true")
			}
			dsn := fmt.Sprintf(
				"%s://%s:%d?%s",
				spec.Mode,
				spec.Endpoint,
				spec.Port,
				options.Encode(),
			)

			_, err := settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-otel-traces", configuration.Name),
				fmt.Sprintf("opentelemetry.%s.dsn", discr), dsn, stackNames...)
			if err != nil {
				return err
			}

			return nil
		}

		if monitoring.Traces != nil && monitoring.Traces.Otlp != nil {
			if err := createSettings("traces", monitoring.Traces.Otlp); err != nil {
				return err
			}
		}

		if monitoring.Metrics != nil && monitoring.Metrics.Otlp != nil {
			if err := createSettings("metrics", monitoring.Metrics.Otlp); err != nil {
				return err
			}
		}
	}

	if configuration.Spec.Broker.Kafka != nil {
		options := url.Values{}
		if configuration.Spec.Broker.Kafka.TLS {
			options.Set("tls", "true")
		}
		if sasl := configuration.Spec.Broker.Kafka.SASL; sasl != nil {
			options.Set("saslEnabled", "true")
			if sasl.Username != "" {
				options.Set("saslUsername", sasl.Username)
			}
			if sasl.Password != "" {
				options.Set("saslPassword", sasl.Password)
			}
			if sasl.Mechanism != "" {
				options.Set("saslMechanism", sasl.Mechanism)
			}
			if sasl.ScramSHASize != "" {
				options.Set("saslSCRAMSHASize", sasl.ScramSHASize)
			}
		}
		dsn := fmt.Sprintf("kafka://%s?%s", configuration.Spec.Broker.Kafka.Brokers[0], options.Encode())

		_, err := settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-broker", configuration.Name),
			"broker.dsn", dsn, stackNames...)
		if err != nil {
			return err
		}
	} else if configuration.Spec.Broker.Nats != nil {

		options := url.Values{}
		if configuration.Spec.Broker.Nats.Replicas != 0 {
			options.Set("replicas", fmt.Sprint(configuration.Spec.Broker.Nats.Replicas))
		}
		dsn := fmt.Sprintf("nats://%s?%s", configuration.Spec.Broker.Nats.URL, options.Encode())

		_, err := settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-broker", configuration.Name),
			"broker.dsn", dsn, stackNames...)
		if err != nil {
			return err
		}
	}

	options := url.Values{}
	if configuration.Spec.Temporal.TLS.SecretName != "" {
		options.Set("secret", configuration.Spec.Temporal.TLS.SecretName)
	}
	_, err := settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-temporal-dsn", configuration.Name),
		"temporal.dsn",
		fmt.Sprintf("temporal://%s/%s?%s", configuration.Spec.Temporal.Address, configuration.Spec.Temporal.Namespace, options.Encode()),
		stackNames...,
	)
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
	}

	basicAuth := ""
	options = url.Values{}
	if basicAuthConf := configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth; basicAuthConf != nil {
		if basicAuthConf.Username != "" {
			basicAuth = fmt.Sprintf("%s:%s@", basicAuthConf.Username, basicAuthConf.Password)
		}
		if basicAuthConf.SecretName != "" {
			options.Set("secret", basicAuthConf.SecretName)
		}
	}
	if configuration.Spec.Services.Search.ElasticSearchConfig.TLS.Enabled {
		options.Set("tls", "true")
		if configuration.Spec.Services.Search.ElasticSearchConfig.TLS.SkipCertVerify {
			options.Set("skipCertVerify", "true")
		}
	}

	elasticSearchDSN := fmt.Sprintf(
		"%s://%s%s:%d?%s",
		configuration.Spec.Services.Search.ElasticSearchConfig.Scheme,
		basicAuth,
		configuration.Spec.Services.Search.ElasticSearchConfig.Host,
		configuration.Spec.Services.Search.ElasticSearchConfig.Port,
		options.Encode(),
	)

	_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-elasticsearch-dsn", configuration.Name),
		"elasticsearch.dsn", elasticSearchDSN, stackNames...)
	if err != nil {
		return err
	}

	batchingConfig := ""

	if configuration.Spec.Services.Search.Batching.Count != 0 {
		batchingConfig = fmt.Sprintf("count=%d", configuration.Spec.Services.Search.Batching.Count)
	}

	if configuration.Spec.Services.Search.Batching.Period != "" {
		if len(batchingConfig) > 0 {
			batchingConfig = batchingConfig + ","
		}
		batchingConfig = batchingConfig + "period=" + configuration.Spec.Services.Search.Batching.Period
	}
	if batchingConfig != "" {
		_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-search-batching", configuration.Name),
			"search.batching", batchingConfig, stackNames...)
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

	type serviceAnnotationsDescriptor struct {
		annotations map[string]string
		name        string
	}
	for _, cfg := range []serviceAnnotationsDescriptor{
		{
			annotations: configuration.Spec.Services.Ledger.Annotations.Service,
			name:        "ledger",
		},
		{
			annotations: configuration.Spec.Services.Payments.Annotations.Service,
			name:        "payments",
		},
		{
			annotations: configuration.Spec.Services.Orchestration.Annotations.Service,
			name:        "orchestration",
		},
		{
			annotations: configuration.Spec.Services.Auth.Annotations.Service,
			name:        "auth",
		},
		{
			annotations: configuration.Spec.Services.Webhooks.Annotations.Service,
			name:        "webhooks",
		},
		{
			annotations: configuration.Spec.Services.Reconciliation.Annotations.Service,
			name:        "reconciliation",
		},
		{
			annotations: configuration.Spec.Services.Gateway.Annotations.Service,
			name:        "gateway",
		},
		{
			annotations: configuration.Spec.Services.Wallets.Annotations.Service,
			name:        "wallets",
		},
		{
			annotations: configuration.Spec.Services.Stargate.Annotations.Service,
			name:        "stargate",
		},
		{
			annotations: configuration.Spec.Services.Search.Annotations.Service,
			name:        "search",
		},
	} {
		if cfg.annotations == nil || len(cfg.annotations) == 0 {
			continue
		}

		computed := ""
		for key, value := range cfg.annotations {
			computed = fmt.Sprintf("%s=%s,%s", key, value, computed)
		}
		computed = computed[:len(computed)-1]

		_, err = settings.CreateOrUpdate(ctx, fmt.Sprintf("%s-%s-service-annotations", configuration.Name, cfg.name),
			fmt.Sprintf("services.%s.annotations", cfg.name),
			computed, stackNames...)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	Init(
		WithReconciler[*v1beta3.Configuration](Reconcile,
			WithWatch[*v1beta3.Configuration, *v1beta3.Stack](func(ctx Context, object *v1beta3.Stack) []reconcile.Request {
				return []reconcile.Request{{
					NamespacedName: types.NamespacedName{
						Name: object.Spec.Seed,
					},
				}}
			}),
		),
	)
}
