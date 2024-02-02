package stacks

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/api/stack.formance.com/v1beta3"
	. "github.com/formancehq/operator/internal/core"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stack.formance.com,resources=stacks/finalizers,verbs=update
// +kubebuilder:rbac:groups=stack.formance.com,resources=configurations,verbs=get;list;watch

func Reconcile(ctx Context, stack *v1beta3.Stack) error {
	if stack.Spec.Versions == "" {
		patch := client.MergeFrom(stack.DeepCopy())
		stack.Spec.Versions = "default"
		if err := ctx.GetClient().Patch(ctx, stack, patch); err != nil {
			return err
		}
		return NewPendingError()
	}

	configuration := &v1beta3.Configuration{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: stack.Spec.Seed,
	}, configuration); err != nil {
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

	newStack, _, err := CreateOrUpdate[*v1beta1.Stack](ctx, types.NamespacedName{
		Name: stack.Name,
	}, func(t *v1beta1.Stack) error {
		if t.Annotations == nil {
			t.Annotations = map[string]string{}
		}
		// Automatically set skip label on creation
		// if the new object stack does not exists actually (we're upgrading)
		// and if the actual stack (old object) is already ready
		// this way, creating with the agent works as expected
		// and old stacks will not be automatically upgraded
		if t.ResourceVersion == "" && stack.Status.Ready {
			t.Annotations[v1beta1.SkipLabel] = "true"
		}
		t.Spec.Dev = stack.Spec.Dev
		t.Spec.Debug = stack.Spec.Debug
		if configuration.Spec.Services.Gateway.EnableAuditPlugin != nil {
			t.Spec.EnableAudit = *configuration.Spec.Services.Gateway.EnableAuditPlugin
		}
		t.Spec.Disabled = stack.Spec.Disabled
		t.Spec.VersionsFromFile = stack.Spec.Versions

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "creating stack")
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

	ready := true

	if !isDisabled(stack, configuration, false, "ledger") {
		ledger, _, err := CreateOrUpdate[*v1beta1.Ledger](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Ledger) error {
			t.Spec.Stack = stack.Name
			t.Spec.DeploymentStrategy = v1beta1.DeploymentStrategy(configuration.Spec.Services.Ledger.DeploymentStrategy)
			t.Spec.Locking = &v1beta1.LockingStrategy{
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

			return nil
		})
		if err != nil {
			return errors.Wrap(err, "creating ledger service")
		}
		ready = ready && ledger.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "payments") {
		payments, _, err := CreateOrUpdate[*v1beta1.Payments](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Payments) error {
			t.Spec.Stack = stack.Name
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "creating payments service")
		}
		ready = ready && payments.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "wallets") {
		wallets, _, err := CreateOrUpdate[*v1beta1.Wallets](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Wallets) error {
			t.Spec.Stack = stack.Name
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "creating wallets service")
		}
		ready = ready && wallets.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "orchestration") {
		orchestration, _, err := CreateOrUpdate[*v1beta1.Orchestration](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Orchestration) error {
			t.Spec.Stack = stack.Name
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "creating orchestration service")
		}
		ready = ready && orchestration.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "webhooks") {
		webhooks, _, err := CreateOrUpdate[*v1beta1.Webhooks](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Webhooks) error {
			t.Spec.Stack = stack.Name
			return nil
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
		}, func(t *v1beta1.Reconciliation) error {
			t.Spec.Stack = stack.Name
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "creating reconciliation service")
		}
		ready = ready && reconciliation.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "search") {
		search, _, err := CreateOrUpdate[*v1beta1.Search](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Search) error {
			t.Spec.Stack = stack.Name
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "creating search service")
		}
		ready = ready && search.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "auth") {
		auth, _, err := CreateOrUpdate[*v1beta1.Auth](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Auth) error {
			t.Spec.Stack = stack.Name
			t.Spec.DelegatedOIDCServer = &v1beta1.DelegatedOIDCServerConfiguration{
				Issuer:       stack.Spec.Auth.DelegatedOIDCServer.Issuer,
				ClientID:     stack.Spec.Auth.DelegatedOIDCServer.ClientID,
				ClientSecret: stack.Spec.Auth.DelegatedOIDCServer.ClientSecret,
			}

			return nil
		})
		if err != nil {
			return errors.Wrap(err, "creating auth service")
		}
		ready = ready && auth.Status.Ready
	}

	if !isDisabled(stack, configuration, false, "gateway") {
		gateway, _, err := CreateOrUpdate[*v1beta1.Gateway](ctx, types.NamespacedName{
			Name: stack.Name,
		}, func(t *v1beta1.Gateway) error {
			t.Spec.Stack = stack.Name
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

			return nil
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
			}, func(t *v1beta1.Stargate) error {
				t.Spec.Stack = stack.Name
				t.Spec.ServerURL = stack.Spec.Stargate.StargateServerURL
				t.Spec.OrganizationID = organizationID
				t.Spec.StackID = stackID
				t.Spec.Auth.Issuer = stack.Spec.Auth.DelegatedOIDCServer.Issuer
				t.Spec.Auth.ClientID = stack.Spec.Auth.DelegatedOIDCServer.ClientID
				t.Spec.Auth.ClientSecret = stack.Spec.Auth.DelegatedOIDCServer.ClientSecret

				return nil
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
		}, func(t *v1beta1.AuthClient) error {
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

			return nil
		})
		if err != nil {
			return errors.Wrap(err, "creating auth client service")
		}
	}

	if ready && newStack.Annotations[v1beta1.SkipLabel] != "true" {
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

		pods := &corev1.PodList{}
		if err := ctx.GetClient().List(ctx, pods, client.InNamespace(stack.Name)); err != nil {
			return err
		}

		for _, item := range pods.Items {
			if item.Status.Phase == "Succeeded" {
				if err := ctx.GetClient().Delete(ctx, &item); err != nil {
					return err
				}
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

func listStacksAndReconcile(ctx Context, opts ...client.ListOption) []reconcile.Request {
	stacks := &v1beta3.StackList{}
	err := ctx.GetClient().List(ctx, stacks, opts...)
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
		WithStdReconciler(Reconcile,
			WithWatch[*v1beta3.Stack, *v1beta3.Configuration](watch[*v1beta3.Configuration](".spec.seed")),
			WithWatch[*v1beta3.Stack, *v1beta3.Versions](watch[*v1beta3.Versions](".spec.seed")),
			WithWatch[*v1beta3.Stack, *v1beta1.Stack](func(ctx Context, object *v1beta1.Stack) []reconcile.Request {
				return []reconcile.Request{{
					NamespacedName: types.NamespacedName{
						Name: object.GetName(),
					},
				}}
			}),
		),
		WithSimpleIndex[*v1beta3.Stack](".spec.seed", func(t *v1beta3.Stack) string {
			return t.Spec.Seed
		}),
		WithSimpleIndex[*v1beta3.Stack](".spec.versions", func(t *v1beta3.Stack) string {
			return t.Spec.Versions
		}),
	)
}
