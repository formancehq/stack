package stack_formance_com

import (
	"context"
	"fmt"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/api/stack.formance.com/v1beta3"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	controllererrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ctrl "sigs.k8s.io/controller-runtime"
)

// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks/finalizers,verbs=update
// +kubebuilder:rbac:groups=stack.formance.com,resources=configurations,verbs=get;list;watch
// +kubebuilder:rbac:groups=stack.formance.com,resources=versions,verbs=get;list;watch

// Reconciler reconciles a Stack object
type StackController struct{}

func (r *StackController) Reconcile(ctx Context, stack *v1beta3.Stack) error {

	if stack.Spec.Disabled {
		return client.IgnoreNotFound(ctx.GetClient().Delete(ctx, &corev1.Namespace{
			ObjectMeta: v1.ObjectMeta{
				Name: stack.Name,
			},
		}))
	}

	ns := &corev1.Namespace{}
	err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: stack.Name,
	}, ns)
	switch {
	case controllererrors.IsNotFound(err):
	case err == nil: // Namespace found
		if ns.Annotations[OperatorVersionKey] != OperatorVersion {
			return ctx.GetClient().Delete(ctx, &corev1.Namespace{
				ObjectMeta: v1.ObjectMeta{
					Name: stack.Name,
				},
			})
		}
	default:
		return err
	}

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

	_, _, err = CreateOrUpdate[*v1beta1.Stack](ctx, types.NamespacedName{
		Name: stack.Name,
	}, func(t *v1beta1.Stack) {
	}, WithController[*v1beta1.Stack](ctx.GetScheme(), stack), func(t *v1beta1.Stack) {
		t.Spec.Dev = stack.Spec.Dev
		t.Spec.Debug = stack.Spec.Debug
		if configuration.Spec.Services.Gateway.EnableAuditPlugin != nil {
			t.Spec.EnableAudit = *configuration.Spec.Services.Gateway.EnableAuditPlugin
		}
	})
	if err != nil {
		return errors.Wrap(err, "creating stack")
	}

	type databaseDescriptor struct {
		config v1beta1.DatabaseConfigurationSpec
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
		}, WithController[*v1beta1.DatabaseConfiguration](ctx.GetScheme(), stack))
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
		}, WithController[*v1beta1.OpenTelemetryConfiguration](ctx.GetScheme(), stack))
		if err != nil {
			return errors.Wrap(err, "creating opentelemetry configuration for service")
		}
	}

	_, _, err = CreateOrUpdate[*v1beta1.BrokerConfiguration](ctx, types.NamespacedName{
		Name: stack.Name,
	}, func(t *v1beta1.BrokerConfiguration) {
		t.Spec = configuration.Spec.Broker
		t.Labels = map[string]string{
			StackLabel: stack.Name,
		}
	}, WithController[*v1beta1.BrokerConfiguration](ctx.GetScheme(), stack))
	if err != nil {
		return errors.Wrap(err, "creating broker configuration for service")
	}

	_, _, err = CreateOrUpdate[*v1beta1.ElasticSearchConfiguration](ctx, types.NamespacedName{
		Name: stack.Name,
	}, func(t *v1beta1.ElasticSearchConfiguration) {
		t.Spec = configuration.Spec.Services.Search.ElasticSearchConfig
		t.Labels = map[string]string{
			StackLabel: stack.Name,
		}
	}, WithController[*v1beta1.ElasticSearchConfiguration](ctx.GetScheme(), stack))
	if err != nil {
		return errors.Wrap(err, "creating elasticsearch configuration for service")
	}

	if len(configuration.Spec.Registries) > 0 {
		_, _, err = CreateOrUpdate[*v1beta1.RegistriesConfiguration](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.RegistriesConfiguration) {
			t.Spec.Registries = configuration.Spec.Registries
			t.Labels = map[string]string{
				StackLabel: stack.Name,
			}
		}, WithController[*v1beta1.RegistriesConfiguration](ctx.GetScheme(), stack))
		if err != nil {
			return errors.Wrap(err, "creating registries configuration")
		}
	}

	if !isDisabled(stack, configuration, false, "ledger") {
		_, _, err = CreateOrUpdate[*v1beta1.Ledger](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Ledger) {
		}, WithController[*v1beta1.Ledger](ctx.GetScheme(), stack), func(t *v1beta1.Ledger) {
			t.Spec.Stack = stack.Name
			t.Spec.Version = versions.Spec.Ledger
			t.Spec.DeploymentStrategy = configuration.Spec.Services.Ledger.DeploymentStrategy
			t.Spec.ResourceRequirements = resourceRequirements(configuration.Spec.Services.Ledger.ResourceProperties)
			if annotations := configuration.Spec.Services.Ledger.Annotations.Service; annotations != nil {
				t.Spec.Service = &v1beta1.ServiceConfiguration{
					Annotations: annotations,
				}
			}
		})
		if err != nil {
			return errors.Wrap(err, "creating ledger service")
		}
	}

	if !isDisabled(stack, configuration, false, "payments") {
		_, _, err = CreateOrUpdate[*v1beta1.Payments](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Payments) {
		}, WithController[*v1beta1.Payments](ctx.GetScheme(), stack), func(t *v1beta1.Payments) {
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
	}

	if !isDisabled(stack, configuration, false, "wallets") {
		_, _, err = CreateOrUpdate[*v1beta1.Wallets](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Wallets) {
		}, WithController[*v1beta1.Wallets](ctx.GetScheme(), stack), func(t *v1beta1.Wallets) {
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
	}

	if !isDisabled(stack, configuration, false, "orchestration") {
		_, _, err = CreateOrUpdate[*v1beta1.Orchestration](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Orchestration) {
		}, WithController[*v1beta1.Orchestration](ctx.GetScheme(), stack), func(t *v1beta1.Orchestration) {
			t.Spec.Stack = stack.Name
			t.Spec.Version = versions.Spec.Orchestration
			t.Spec.ResourceRequirements = resourceRequirements(configuration.Spec.Services.Orchestration.ResourceProperties)
			if annotations := configuration.Spec.Services.Orchestration.Annotations.Service; annotations != nil {
				t.Spec.Service = &v1beta1.ServiceConfiguration{
					Annotations: annotations,
				}
			}
			t.Spec.Temporal = configuration.Spec.Temporal
		})
		if err != nil {
			return errors.Wrap(err, "creating orchestration service")
		}
	}

	if !isDisabled(stack, configuration, false, "webhooks") {
		_, _, err = CreateOrUpdate[*v1beta1.Webhooks](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Webhooks) {
		}, WithController[*v1beta1.Webhooks](ctx.GetScheme(), stack), func(t *v1beta1.Webhooks) {
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
	}

	// note(gfyrag): reconciliation declared as EE.
	// We should also declare some other services EE but to keep compatibility, today, we just configuration
	// reconciliation as EE
	if !isDisabled(stack, configuration, true, "reconciliation") {
		_, _, err = CreateOrUpdate[*v1beta1.Reconciliation](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Reconciliation) {
		}, WithController[*v1beta1.Reconciliation](ctx.GetScheme(), stack), func(t *v1beta1.Reconciliation) {
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
	}

	if !isDisabled(stack, configuration, false, "search") {
		_, _, err = CreateOrUpdate[*v1beta1.Search](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Search) {
		}, WithController[*v1beta1.Search](ctx.GetScheme(), stack), func(t *v1beta1.Search) {
			t.Spec.Stack = stack.Name
			t.Spec.Version = versions.Spec.Search
			t.Spec.Batching = &configuration.Spec.Services.Search.Batching
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
	}

	if !isDisabled(stack, configuration, false, "auth") {
		_, _, err = CreateOrUpdate[*v1beta1.Auth](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Auth) {
		}, WithController[*v1beta1.Auth](ctx.GetScheme(), stack), func(t *v1beta1.Auth) {
			t.Spec.Stack = stack.Name
			t.Spec.Version = versions.Spec.Auth
			t.Spec.ResourceRequirements = resourceRequirements(configuration.Spec.Services.Auth.ResourceProperties)
			if annotations := configuration.Spec.Services.Auth.Annotations.Service; annotations != nil {
				t.Spec.Service = &v1beta1.ServiceConfiguration{
					Annotations: annotations,
				}
			}
			t.Spec.DelegatedOIDCServer = &stack.Spec.Auth.DelegatedOIDCServer
		})
		if err != nil {
			return errors.Wrap(err, "creating auth service")
		}
	}

	if !isDisabled(stack, configuration, false, "gateway") {
		_, _, err = CreateOrUpdate[*v1beta1.Gateway](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Gateway) {
		}, WithController[*v1beta1.Gateway](ctx.GetScheme(), stack), func(t *v1beta1.Gateway) {
			t.Spec.Stack = stack.Name
			t.Spec.Version = versions.Spec.Gateway
			t.Spec.Ingress = &v1beta1.GatewayIngress{
				Host:        stack.Spec.Host,
				Scheme:      stack.Spec.Scheme,
				TLS:         configuration.Spec.Ingress.TLS,
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
	}

	if !isDisabled(stack, configuration, false, "stargate") && stack.Spec.Stargate != nil {
		parts := strings.Split(stack.Name, "-")
		if len(parts) == 2 {
			organizationID := parts[0]
			stackID := parts[1]

			_, _, err = CreateOrUpdate[*v1beta1.Stargate](ctx, types.NamespacedName{
				Name: stack.Name,
			}, func(t *v1beta1.Stargate) {
			}, WithController[*v1beta1.Stargate](ctx.GetScheme(), stack), func(t *v1beta1.Stargate) {
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
		}
	}

	for _, client := range stack.Spec.Auth.StaticClients {
		_, _, err = CreateOrUpdate[*v1beta1.AuthClient](ctx, types.NamespacedName{
			Name: fmt.Sprintf("%s-%s", stack.Name, client.ID),
		}, func(t *v1beta1.AuthClient) {
		}, WithController[*v1beta1.AuthClient](ctx.GetScheme(), stack), func(t *v1beta1.AuthClient) {
			t.Spec = *client
			t.Spec.Stack = stack.Name
		})
		if err != nil {
			return errors.Wrap(err, "creating auth client service")
		}
	}

	stack.Status.Ready = true

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StackController) SetupWithManager(mgr Manager) (*builder.Builder, error) {
	if err := mgr.GetFieldIndexer().
		IndexField(context.Background(), &v1beta3.Stack{}, ".spec.seed", func(rawObj client.Object) []string {
			return []string{rawObj.(*v1beta3.Stack).Spec.Seed}
		}); err != nil {
		return nil, err
	}
	if err := mgr.GetFieldIndexer().
		IndexField(context.Background(), &v1beta3.Stack{}, ".spec.versions", func(rawObj client.Object) []string {
			return []string{rawObj.(*v1beta3.Stack).Spec.Versions}
		}); err != nil {
		return nil, err
	}

	return ctrl.NewControllerManagedBy(mgr).
		Owns(&v1beta1.DatabaseConfiguration{}).
		Owns(&v1beta1.BrokerConfiguration{}).
		Owns(&v1beta1.ElasticSearchConfiguration{}).
		Owns(&v1beta1.OpenTelemetryConfiguration{}).
		Owns(&v1beta1.Ledger{}).
		Owns(&v1beta1.Payments{}).
		Owns(&v1beta1.Orchestration{}).
		Owns(&v1beta1.Wallets{}).
		Owns(&v1beta1.Webhooks{}).
		Owns(&v1beta1.Auth{}).
		Owns(&v1beta1.Gateway{}).
		Owns(&v1beta1.Stargate{}).
		Owns(&v1beta1.Stack{}).
		Watches(
			&v1beta3.Configuration{},
			watch(mgr, ".spec.seed"),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&v1beta3.Versions{},
			watch(mgr, ".spec.versions"),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		For(&v1beta3.Stack{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForStack() *StackController {
	return &StackController{}
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

	return collectionutils.Map(stacks.Items, func(s v1beta3.Stack) reconcile.Request {
		return reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      s.GetName(),
				Namespace: s.GetNamespace(),
			},
		}
	})
}

func watch(mgr ctrl.Manager, field string) handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
		return listStacksAndReconcile(mgr, &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(field, object.GetName()),
			Namespace:     object.GetNamespace(),
		})
	})
}
