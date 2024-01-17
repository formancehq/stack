package legacystacks

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/api/stack.formance.com/v1beta3"
	. "github.com/formancehq/operator/internal/core"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ctrl "sigs.k8s.io/controller-runtime"
)

// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks/finalizers,verbs=update
// +kubebuilder:rbac:groups=stack.formance.com,resources=configurations,verbs=get;list;watch
// +kubebuilder:rbac:groups=stack.formance.com,resources=versions,verbs=get;list;watch

func Reconcile(ctx Context, stack *v1beta3.Stack) error {
	if stack.Spec.Versions == "" {
		patch := client.MergeFrom(stack.DeepCopy())
		stack.Spec.Versions = "default"
		if err := ctx.GetClient().Patch(ctx, stack, patch); err != nil {
			return err
		}
		return ErrPending
	}

	configuration := &v1beta3.Configuration{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: stack.Spec.Seed,
	}, configuration); err != nil {
		return err
	}

	versions := &v1beta3.Versions{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: func() string {
			if stack.Spec.Versions == "" {
				return "default"
			}
			return stack.Spec.Versions
		}(),
	}, versions); err != nil {
		return err
	}

	for _, object := range []client.Object{
		&corev1.ConfigMap{},
		&corev1.Secret{},
		&corev1.Service{},
		&appsv1.Deployment{},
		&networkingv1.Ingress{},
		&batchv1.Job{},
		&batchv1.CronJob{},
	} {
		kinds, _, err := ctx.GetScheme().ObjectKinds(object)
		if err != nil {
			return err
		}

		list := &unstructured.UnstructuredList{}
		list.SetGroupVersionKind(kinds[0])
		if err := ctx.GetClient().List(ctx, list, client.InNamespace(stack.Name)); err != nil {
			return err
		}

	l:
		for _, item := range list.Items {
			ownerReferences := item.GetOwnerReferences()
			for i, reference := range ownerReferences {
				if reference.APIVersion != "stack.formance.com/v1beta3" {
					continue
				}

				patch := client.MergeFrom(item.DeepCopy())
				if i < len(ownerReferences)-1 {
					ownerReferences = append(ownerReferences[:i], ownerReferences[i+1:]...)
				} else {
					ownerReferences = ownerReferences[:i]
				}
				item.SetOwnerReferences(ownerReferences)

				labels := item.GetLabels()
				labels["formance.com/migrate"] = "true"
				item.SetLabels(labels)

				if err := ctx.GetClient().Patch(ctx, &item, patch); err != nil {
					return err
				}

				continue l
			}
		}
	}

	_, _, err := CreateOrUpdate[*v1beta1.Stack](ctx, types.NamespacedName{
		Name: stack.Name,
	}, func(t *v1beta1.Stack) {
		t.Spec.Dev = stack.Spec.Dev
		t.Spec.Debug = stack.Spec.Debug
		if configuration.Spec.Services.Gateway.EnableAuditPlugin != nil {
			t.Spec.EnableAudit = *configuration.Spec.Services.Gateway.EnableAuditPlugin
		}
		t.Spec.Disabled = stack.Spec.Disabled
	})
	if err != nil {
		return errors.Wrap(err, "creating stack")
	}

	if stack.Spec.Disabled {
		return nil
	}

	ns := &corev1.Namespace{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: stack.Name,
	}, ns); client.IgnoreNotFound(err) != nil {
		return err
	} else if err == nil {
		ownerReferences := ns.GetOwnerReferences()
		for i, reference := range ownerReferences {
			if reference.APIVersion != "stack.formance.com/v1beta3" {
				continue
			}

			patch := client.MergeFrom(ns.DeepCopy())
			if i < len(ownerReferences)-1 {
				ownerReferences = append(ownerReferences[:i], ownerReferences[i+1:]...)
			} else {
				ownerReferences = ownerReferences[:i]
			}
			ns.SetOwnerReferences(ownerReferences)

			if err := ctx.GetClient().Patch(ctx, ns, patch); err != nil {
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
	} {
		_, _, err := CreateOrUpdate[*v1beta1.DatabaseConfiguration](ctx, types.NamespacedName{
			Name: fmt.Sprintf("%s-%s", stack.Name, cfg.name),
		}, func(t *v1beta1.DatabaseConfiguration) {
			t.Spec = v1beta1.DatabaseConfigurationSpec{
				Port:                  cfg.config.Port,
				Host:                  cfg.config.Host,
				Username:              cfg.config.Username,
				Password:              cfg.config.Password,
				CredentialsFromSecret: cfg.config.CredentialsFromSecret,
				DisableSSLMode:        cfg.config.DisableSSLMode,
			}
			t.Labels = map[string]string{
				StackLabel:   stack.Name,
				ServiceLabel: cfg.name,
			}
		})
		if err != nil {
			return errors.Wrapf(err, "creating database configuration for service %s", cfg.name)
		}
	}

	if configuration.Spec.Monitoring != nil {
		_, _, err := CreateOrUpdate[*v1beta1.OpenTelemetryConfiguration](ctx, types.NamespacedName{
			Name: stack.Name,
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
			}
			t.Labels = map[string]string{
				StackLabel: stack.Name,
			}
		})
		if err != nil {
			return errors.Wrap(err, "creating opentelemetry configuration for service")
		}
	}

	_, _, err = CreateOrUpdate[*v1beta1.BrokerConfiguration](ctx, types.NamespacedName{
		Name: stack.Name,
	}, func(t *v1beta1.BrokerConfiguration) {
		t.Spec = v1beta1.BrokerConfigurationSpec{
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
		t.Labels = map[string]string{
			StackLabel: stack.Name,
		}
	})
	if err != nil {
		return errors.Wrap(err, "creating broker configuration for service")
	}

	_, _, err = CreateOrUpdate[*v1beta1.TemporalConfiguration](ctx, types.NamespacedName{
		Name: stack.Name,
	}, func(t *v1beta1.TemporalConfiguration) {
		t.Spec = v1beta1.TemporalConfigurationSpec{
			Address:   configuration.Spec.Temporal.Address,
			Namespace: configuration.Spec.Temporal.Namespace,
			TLS: v1beta1.TemporalTLSConfig{
				CRT:        configuration.Spec.Temporal.TLS.CRT,
				Key:        configuration.Spec.Temporal.TLS.Key,
				SecretName: configuration.Spec.Temporal.TLS.SecretName,
			},
		}
		t.Labels = map[string]string{
			StackLabel: stack.Name,
		}
	})
	if err != nil {
		return errors.Wrap(err, "creating temporal configuration for service")
	}

	_, _, err = CreateOrUpdate[*v1beta1.ElasticSearchConfiguration](ctx, types.NamespacedName{
		Name: stack.Name,
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
		}
		t.Labels = map[string]string{
			StackLabel: stack.Name,
		}
	})
	if err != nil {
		return errors.Wrap(err, "creating elasticsearch configuration for service")
	}

	if len(configuration.Spec.Registries) > 0 {
		_, _, err = CreateOrUpdate[*v1beta1.RegistriesConfiguration](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.RegistriesConfiguration) {
			registries := make(map[string]v1beta1.RegistryConfigurationSpec)
			for k, config := range configuration.Spec.Registries {
				registries[k] = v1beta1.RegistryConfigurationSpec{
					Endpoint: config.Endpoint,
				}
			}
			t.Spec = v1beta1.RegistriesConfigurationSpec{
				Registries: registries,
			}
			t.Labels = map[string]string{
				StackLabel: stack.Name,
			}
		})
		if err != nil {
			return errors.Wrap(err, "creating registries configuration")
		}
	}

	ready := true

	if !isDisabled(stack, configuration, false, "ledger") {
		ledger, _, err := CreateOrUpdate[*v1beta1.Ledger](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Ledger) {
			t.Spec.Stack = stack.Name
			t.Spec.Version = versions.Spec.Ledger
			t.Spec.DeploymentStrategy = v1beta1.DeploymentStrategy(configuration.Spec.Services.Ledger.DeploymentStrategy)
			t.Spec.ResourceRequirements = resourceRequirements(configuration.Spec.Services.Ledger.ResourceProperties)
			if annotations := configuration.Spec.Services.Ledger.Annotations.Service; annotations != nil {
				t.Spec.Service = &v1beta1.ServiceConfiguration{
					Annotations: annotations,
				}
			}
			t.Spec.Locking = v1beta1.LockingStrategy{
				Strategy: configuration.Spec.Services.Ledger.Locking.Strategy,
				Redis: func() *v1beta1.LockingStrategyRedisConfig {
					if configuration.Spec.Services.Ledger.Locking.Strategy != "redis" {
						return nil
					}

					return &v1beta1.LockingStrategyRedisConfig{
						Uri:         configuration.Spec.Services.Ledger.Locking.Redis.Uri,
						TLS:         configuration.Spec.Services.Ledger.Locking.Redis.TLS,
						InsecureTLS: configuration.Spec.Services.Ledger.Locking.Redis.InsecureTLS,
						Duration:    configuration.Spec.Services.Ledger.Locking.Redis.Duration,
						Retry:       configuration.Spec.Services.Ledger.Locking.Redis.Retry,
					}
				}(),
			}
		})
		if err != nil {
			return errors.Wrap(err, "creating ledger service")
		}
		ready = ready && ledger.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "payments") {
		payments, _, err := CreateOrUpdate[*v1beta1.Payments](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Payments) {
			t.Spec.Stack = stack.Name
			t.Spec.Version = versions.Spec.Payments
			t.Spec.ResourceRequirements = resourceRequirements(configuration.Spec.Services.Payments.ResourceProperties)
			t.Spec.EncryptionKey = configuration.Spec.Services.Payments.EncryptionKey
			if annotations := configuration.Spec.Services.Payments.Annotations.Service; annotations != nil {
				t.Spec.Service = &v1beta1.ServiceConfiguration{
					Annotations: annotations,
				}
			}
		})
		if err != nil {
			return errors.Wrap(err, "creating payments service")
		}
		ready = ready && payments.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "wallets") {
		wallets, _, err := CreateOrUpdate[*v1beta1.Wallets](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Wallets) {
			t.Spec.Stack = stack.Name
			t.Spec.Version = versions.Spec.Wallets
			t.Spec.ResourceRequirements = resourceRequirements(configuration.Spec.Services.Wallets.ResourceProperties)
			if annotations := configuration.Spec.Services.Wallets.Annotations.Service; annotations != nil {
				t.Spec.Service = &v1beta1.ServiceConfiguration{
					Annotations: annotations,
				}
			}
		})
		if err != nil {
			return errors.Wrap(err, "creating wallets service")
		}
		ready = ready && wallets.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "orchestration") {
		orchestration, _, err := CreateOrUpdate[*v1beta1.Orchestration](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Orchestration) {
			t.Spec.Stack = stack.Name
			t.Spec.Version = versions.Spec.Orchestration
			t.Spec.ResourceRequirements = resourceRequirements(configuration.Spec.Services.Orchestration.ResourceProperties)
			if annotations := configuration.Spec.Services.Orchestration.Annotations.Service; annotations != nil {
				t.Spec.Service = &v1beta1.ServiceConfiguration{
					Annotations: annotations,
				}
			}
		})
		if err != nil {
			return errors.Wrap(err, "creating orchestration service")
		}
		ready = ready && orchestration.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "webhooks") {
		webhooks, _, err := CreateOrUpdate[*v1beta1.Webhooks](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Webhooks) {
			t.Spec.Stack = stack.Name
			t.Spec.Version = versions.Spec.Webhooks
			t.Spec.ResourceRequirements = resourceRequirements(configuration.Spec.Services.Webhooks.ResourceProperties)
			if annotations := configuration.Spec.Services.Webhooks.Annotations.Service; annotations != nil {
				t.Spec.Service = &v1beta1.ServiceConfiguration{
					Annotations: annotations,
				}
			}
		})
		if err != nil {
			return errors.Wrap(err, "creating webhooks service")
		}
		ready = ready && webhooks.Status.Ready
	}

	// note(gfyrag): reconciliation declared as EE.
	// We should also declare some other services EE but to keep compatibility, today, we just configuration
	// reconciliation as EE
	if !isDisabled(stack, configuration, true, "reconciliation") {
		reconciliation, _, err := CreateOrUpdate[*v1beta1.Reconciliation](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Reconciliation) {
			t.Spec.Stack = stack.Name
			t.Spec.Version = versions.Spec.Reconciliation
			t.Spec.ResourceRequirements = resourceRequirements(configuration.Spec.Services.Reconciliation.ResourceProperties)
			if annotations := configuration.Spec.Services.Reconciliation.Annotations.Service; annotations != nil {
				t.Spec.Service = &v1beta1.ServiceConfiguration{
					Annotations: annotations,
				}
			}
		})
		if err != nil {
			return errors.Wrap(err, "creating reconciliation service")
		}
		ready = ready && reconciliation.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "search") {
		search, _, err := CreateOrUpdate[*v1beta1.Search](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Search) {
			t.Spec.Stack = stack.Name
			t.Spec.Version = versions.Spec.Search
			t.Spec.Batching = &v1beta1.Batching{
				Count:  configuration.Spec.Services.Search.Batching.Count,
				Period: configuration.Spec.Services.Search.Batching.Period,
			}
			if resourceProperties := configuration.Spec.Services.Search.BenthosResourceProperties; resourceProperties != nil {
				t.Spec.StreamProcessor = &v1beta1.SearchStreamProcessorSpec{
					ResourceRequirements: resourceRequirements(resourceProperties),
				}
			}
			if annotations := configuration.Spec.Services.Search.Annotations.Service; annotations != nil {
				t.Spec.Service = &v1beta1.ServiceConfiguration{
					Annotations: annotations,
				}
			}
		})
		if err != nil {
			return errors.Wrap(err, "creating search service")
		}
		ready = ready && search.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "auth") {
		auth, _, err := CreateOrUpdate[*v1beta1.Auth](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Auth) {
			t.Spec.Stack = stack.Name
			t.Spec.Version = versions.Spec.Auth
			t.Spec.ResourceRequirements = resourceRequirements(configuration.Spec.Services.Auth.ResourceProperties)
			if annotations := configuration.Spec.Services.Auth.Annotations.Service; annotations != nil {
				t.Spec.Service = &v1beta1.ServiceConfiguration{
					Annotations: annotations,
				}
			}
			t.Spec.DelegatedOIDCServer = &v1beta1.DelegatedOIDCServerConfiguration{
				Issuer:       stack.Spec.Auth.DelegatedOIDCServer.Issuer,
				ClientID:     stack.Spec.Auth.DelegatedOIDCServer.ClientID,
				ClientSecret: stack.Spec.Auth.DelegatedOIDCServer.ClientSecret,
			}
		})
		if err != nil {
			return errors.Wrap(err, "creating auth service")
		}
		ready = ready && auth.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "gateway") {
		gateway, _, err := CreateOrUpdate[*v1beta1.Gateway](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Gateway) {
			t.Spec.Stack = stack.Name
			t.Spec.Version = versions.Spec.Gateway
			t.Spec.Ingress = &v1beta1.GatewayIngress{
				Host:   stack.Spec.Host,
				Scheme: stack.Spec.Scheme,
				TLS: func() *v1beta1.GatewayIngressTLS {
					if configuration.Spec.Ingress.TLS == nil {
						return nil
					}
					return &v1beta1.GatewayIngressTLS{
						SecretName: configuration.Spec.Ingress.TLS.SecretName,
					}
				}(),
				Annotations: configuration.Spec.Ingress.Annotations,
			}
			t.Spec.ResourceRequirements = resourceRequirements(configuration.Spec.Services.Gateway.ResourceProperties)
			if annotations := configuration.Spec.Services.Gateway.Annotations.Service; annotations != nil {
				t.Spec.Service = &v1beta1.ServiceConfiguration{
					Annotations: annotations,
				}
			}
		})
		if err != nil {
			return errors.Wrap(err, "creating gateway service")
		}
		ready = ready && gateway.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "stargate") && stack.Spec.Stargate != nil {
		parts := strings.Split(stack.Name, "-")
		if len(parts) == 2 {
			organizationID := parts[0]
			stackID := parts[1]

			stargate, _, err := CreateOrUpdate[*v1beta1.Stargate](ctx, types.NamespacedName{
				Name: stack.Name,
			}, func(t *v1beta1.Stargate) {
				t.Spec.Stack = stack.Name
				t.Spec.Version = versions.Spec.Stargate
				t.Spec.ResourceRequirements = resourceRequirements(configuration.Spec.Services.Stargate.ResourceProperties)
				t.Spec.ServerURL = stack.Spec.Stargate.StargateServerURL
				t.Spec.OrganizationID = organizationID
				t.Spec.StackID = stackID
				t.Spec.Auth.Issuer = stack.Spec.Auth.DelegatedOIDCServer.Issuer
				t.Spec.Auth.ClientID = stack.Spec.Auth.DelegatedOIDCServer.ClientID
				t.Spec.Auth.ClientSecret = stack.Spec.Auth.DelegatedOIDCServer.ClientSecret
			})
			if err != nil {
				return errors.Wrap(err, "creating stargate service")
			}
			ready = ready && stargate.Status.Ready
		}
	}

	for _, client := range stack.Spec.Auth.StaticClients {
		_, _, err = CreateOrUpdate[*v1beta1.AuthClient](ctx, types.NamespacedName{
			Name: fmt.Sprintf("%s-%s", stack.Name, client.ID),
		}, func(t *v1beta1.AuthClient) {
			t.Spec = v1beta1.AuthClientSpec{
				ID:     client.ID,
				Public: client.Public,
				Description: func() string {
					if client.Description == nil {
						return ""
					}
					return *client.Description
				}(),
				RedirectUris:           client.RedirectUris,
				PostLogoutRedirectUris: client.PostLogoutRedirectUris,
				Scopes:                 client.Scopes,
				Secret: func() string {
					if len(client.Secrets) > 0 {
						return client.Secrets[0]
					}
					return ""
				}(),
			}
			t.Spec.Stack = stack.Name
		})
		if err != nil {
			return errors.Wrap(err, "creating auth client service")
		}
	}

	if ready {
		for _, object := range []client.Object{
			&corev1.ConfigMap{},
			&corev1.Secret{},
			&corev1.Service{},
			&appsv1.Deployment{},
			&networkingv1.Ingress{},
			&batchv1.Job{},
			&batchv1.CronJob{},
		} {
			kinds, _, err := ctx.GetScheme().ObjectKinds(object)
			if err != nil {
				return err
			}

			list := &unstructured.UnstructuredList{}
			list.SetGroupVersionKind(kinds[0])
			if err := ctx.GetClient().List(ctx, list, client.InNamespace(stack.Name)); err != nil {
				return err
			}

		l2:
			for _, item := range list.Items {
				ownerReferences := item.GetOwnerReferences()
				for _, reference := range ownerReferences {
					if reference.APIVersion == "formance.com/v1beta1" {
						continue l2
					}
				}
				if item.GetLabels()["formance.com/migrate"] == "true" {
					if err := ctx.GetClient().Delete(ctx, &item); err != nil {
						return err
					}
				}
			}
		}

		list := &unstructured.UnstructuredList{}
		list.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "stack.formance.com",
			Version: "v1beta3",
			Kind:    "Migration",
		})
		if err := ctx.GetClient().List(ctx, list, client.InNamespace(stack.Name)); err != nil {
			return err
		}

		for _, item := range list.Items {
			if err := ctx.GetClient().Delete(ctx, &item); err != nil {
				return err
			}
		}
	}

	return nil
}

func isDisabled(stack *v1beta3.Stack, configuration *v1beta3.Configuration, isEE bool, name string) bool {
	if isEE {
		return !stack.Spec.Services.IsExplicitlyEnabled(name) && !configuration.Spec.Services.IsExplicitlyEnabled(name)
	} else {
		return stack.Spec.Services.IsExplicitlyDisabled(name) || configuration.Spec.Services.IsExplicitlyDisabled(name)
	}
}

func resourceRequirements(resourceProperties *v1beta3.ResourceProperties) *corev1.ResourceRequirements {
	if resourceProperties == nil {
		return nil
	}
	resources := &corev1.ResourceRequirements{}
	if resourceProperties.Request != nil {
		if resources.Requests == nil {
			resources.Requests = make(corev1.ResourceList)
		}

		if resourceProperties.Request.Cpu != "" {
			resources.Requests[corev1.ResourceCPU] = resource.MustParse(resourceProperties.Request.Cpu)
		}

		if resourceProperties.Request.Memory != "" {
			resources.Requests[corev1.ResourceMemory] = resource.MustParse(resourceProperties.Request.Memory)
		}
	}

	if resourceProperties.Limits != nil {
		if resources.Limits == nil {
			resources.Limits = make(corev1.ResourceList)
		}

		if resourceProperties.Limits.Cpu != "" {
			resources.Limits[corev1.ResourceCPU] = resource.MustParse(resourceProperties.Limits.Cpu)
		}

		if resourceProperties.Limits.Memory != "" {
			resources.Limits[corev1.ResourceMemory] = resource.MustParse(resourceProperties.Limits.Memory)
		}
	}

	return resources
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

func listStacksAndReconcile(mgr ctrl.Manager, opts ...client.ListOption) []reconcile.Request {
	stacks := &v1beta3.StackList{}
	err := mgr.GetClient().List(context.TODO(), stacks, opts...)
	if err != nil {
		panic(err)
	}

	return Map(stacks.Items, func(s v1beta3.Stack) reconcile.Request {
		return reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      s.GetName(),
				Namespace: s.GetNamespace(),
			},
		}
	})
}

func watch[T client.Object](field string) func(ctx Context, object T) []reconcile.Request {
	return func(ctx Context, object T) []reconcile.Request {
		return listStacksAndReconcile(ctx, &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(field, object.GetName()),
			Namespace:     object.GetNamespace(),
		})
	}
}

func init() {
	Init(
		WithReconciler(Reconcile,
			WithOwn(&v1beta1.DatabaseConfiguration{}),
			WithOwn(&v1beta1.BrokerConfiguration{}),
			WithOwn(&v1beta1.ElasticSearchConfiguration{}),
			WithOwn(&v1beta1.OpenTelemetryConfiguration{}),
			WithOwn(&v1beta1.Ledger{}),
			WithOwn(&v1beta1.Payments{}),
			WithOwn(&v1beta1.Orchestration{}),
			WithOwn(&v1beta1.Wallets{}),
			WithOwn(&v1beta1.Webhooks{}),
			WithOwn(&v1beta1.Auth{}),
			WithOwn(&v1beta1.Gateway{}),
			WithOwn(&v1beta1.Stargate{}),
			WithOwn(&v1beta1.Stack{}),
			WithWatch(watch[*v1beta3.Configuration](".spec.seed")),
			WithWatch(watch[*v1beta3.Versions](".spec.seed")),
		),
		WithIndex[*v1beta3.Stack](".spec.seed", func(t *v1beta3.Stack) string {
			return t.Spec.Seed
		}),
		WithIndex[*v1beta3.Stack](".spec.versions", func(t *v1beta3.Stack) string {
			return t.Spec.Versions
		}),
	)
}
