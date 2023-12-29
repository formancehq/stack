/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	_ "embed"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/brokerconfigurations"
	common "github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/databases"
	"github.com/formancehq/operator/v2/internal/deployments"
	"github.com/formancehq/operator/v2/internal/httpapis"
	"github.com/formancehq/operator/v2/internal/ledgers"
	"github.com/formancehq/operator/v2/internal/services"
	"github.com/formancehq/operator/v2/internal/stacks"
	"github.com/formancehq/operator/v2/internal/streams"
	"github.com/formancehq/operator/v2/internal/topics"
	"github.com/formancehq/search/benthos"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	pkgError "github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// LedgerController reconciles a Ledger object
type LedgerController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=ledgers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=ledgers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=ledgers/finalizers,verbs=update

func (r *LedgerController) Reconcile(ctx common.Context, ledger *v1beta1.Ledger) error {

	stack, err := stacks.GetStack(ctx, ledger.Spec)
	if err != nil {
		return err
	}

	database, err := databases.Create(ctx, stack, "ledger")
	if err != nil {
		if pkgError.Is(err, common.ErrPending) {
			return nil
		}
		return err
	}

	err = r.installLedger(ctx, stack, ledger, database, common.GetVersion(stack, ledger.Spec.Version))
	if err != nil {
		return err
	}

	if err := httpapis.Create(ctx, stack, ledger, "ledger"); err != nil {
		return err
	}

	if err := streams.LoadFromFileSystem(ctx, benthos.Streams, ledger.Spec.Stack, "streams/ledger/v2.0.0"); err != nil {
		return err
	}

	return nil
}

func (r *LedgerController) installLedger(ctx common.Context, stack *v1beta1.Stack,
	ledger *v1beta1.Ledger, database *v1beta1.Database, version string) error {

	switch ledger.Spec.DeploymentStrategy {
	case v1beta1.DeploymentStrategyMonoWriterMultipleReader:
		if err := common.DeleteIfExists[*appsv1.Deployment](ctx, common.GetNamespacedResourceName(stack.Name, "ledger")); err != nil {
			return err
		}
		return r.installLedgerV2MonoWriterMultipleReader(ctx, stack, ledger, database, version)
	default:
		if err := r.uninstallLedgerV2MonoWriterMultipleReader(ctx, stack); err != nil {
			return err
		}
		return r.installLedgerV2SingleInstance(ctx, stack, ledger, database, version)
	}
}

func (r *LedgerController) installLedgerV2SingleInstance(ctx common.Context, stack *v1beta1.Stack,
	ledger *v1beta1.Ledger, database *v1beta1.Database, version string) error {
	container, err := r.createLedgerContainerV2Full(ctx, stack)
	if err != nil {
		return err
	}

	err = r.setCommonContainerConfiguration(ctx, stack, ledger, version, database, container)
	if err != nil {
		return err
	}

	if err := r.createDeployment(ctx, stack, ledger, "ledger", *container,
		deployments.WithReplicas(1),
		r.setInitContainer(database, version),
	); err != nil {
		return err
	}

	return nil
}

func (r *LedgerController) setInitContainer(database *v1beta1.Database, version string) func(t *appsv1.Deployment) {
	return func(t *appsv1.Deployment) {
		t.Spec.Template.Spec.InitContainers = []corev1.Container{
			databases.MigrateDatabaseContainer(
				"ledger",
				database.Status.Configuration.DatabaseConfigurationSpec,
				database.Status.Configuration.Database,
				version,
				func(m *databases.MigrationConfiguration) {
					m.Command = []string{"buckets", "upgrade-all"}
					m.AdditionalEnv = []corev1.EnvVar{
						common.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"),
					}
				},
			),
		}
	}
}

func (r *LedgerController) installLedgerV2MonoWriterMultipleReader(ctx common.Context, stack *v1beta1.Stack,
	ledger *v1beta1.Ledger, database *v1beta1.Database, version string) error {

	createDeployment := func(name string, container corev1.Container, mutators ...common.ObjectMutator[*appsv1.Deployment]) error {
		err := r.setCommonContainerConfiguration(ctx, stack, ledger, version, database, &container)
		if err != nil {
			return err
		}

		if err := r.createDeployment(ctx, stack, ledger, name, container, mutators...); err != nil {
			return err
		}

		if err := r.createK8SService(ctx, stack, ledger, name); err != nil {
			return err
		}

		return nil
	}

	container, err := r.createLedgerContainerV2WriteOnly(ctx, stack)
	if err != nil {
		return err
	}
	if err := createDeployment("ledger-write", *container,
		deployments.WithReplicas(1),
		r.setInitContainer(database, version),
	); err != nil {
		return err
	}

	container = r.createLedgerContainerV2ReadOnly()
	if err := createDeployment("ledger-read", *container); err != nil {
		return err
	}

	if err := r.createGatewayDeployment(ctx, stack, ledger); err != nil {
		return err
	}

	return nil
}

func (r *LedgerController) uninstallLedgerV2MonoWriterMultipleReader(ctx common.Context, stack *v1beta1.Stack) error {

	remove := func(name string) error {
		if err := common.DeleteIfExists[*appsv1.Deployment](ctx, common.GetNamespacedResourceName(stack.Name, name)); err != nil {
			return err
		}
		if err := common.DeleteIfExists[*corev1.Service](ctx, common.GetNamespacedResourceName(stack.Name, name)); err != nil {
			return err
		}

		return nil
	}

	if err := remove("ledger-write"); err != nil {
		return err
	}

	if err := remove("ledger-read"); err != nil {
		return err
	}

	if err := common.DeleteIfExists[*appsv1.Deployment](ctx, common.GetNamespacedResourceName(stack.Name, "ledger-gateway")); err != nil {
		return err
	}

	return nil
}

func (r *LedgerController) createK8SService(ctx common.Context, stack *v1beta1.Stack, owner *v1beta1.Ledger, name string) error {
	_, _, err := common.CreateOrUpdate[*corev1.Service](ctx, types.NamespacedName{
		Name:      name,
		Namespace: stack.Name,
	},
		services.ConfigureK8SService(name),
		common.WithController[*corev1.Service](ctx.GetScheme(), owner),
	)
	return err
}

func (r *LedgerController) createDeployment(ctx common.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger,
	name string, container corev1.Container, mutators ...common.ObjectMutator[*appsv1.Deployment]) error {
	mutators = append([]common.ObjectMutator[*appsv1.Deployment]{
		deployments.WithContainers(container),
		deployments.WithMatchingLabels(name),
		common.WithController[*appsv1.Deployment](ctx.GetScheme(), ledger),
	}, mutators...)
	_, _, err := common.CreateOrUpdate[*appsv1.Deployment](ctx,
		common.GetNamespacedResourceName(stack.Name, name),
		mutators...,
	)
	return err
}

func (r *LedgerController) setCommonContainerConfiguration(ctx common.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger,
	version string, database *v1beta1.Database, container *corev1.Container) error {

	env, err := GetCommonServicesEnvVars(ctx, stack, "ledger", ledger.Spec)
	if err != nil {
		return err
	}

	container.Resources = common.GetResourcesWithDefault(ledger.Spec.ResourceProperties, common.ResourceSizeSmall())
	container.Image = common.GetImage("ledger", version)
	container.ImagePullPolicy = common.GetPullPolicy(container.Image)
	container.Env = append(container.Env, env...)
	container.Env = append(container.Env, databases.PostgresEnvVars(
		database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database)...)
	container.Env = append(container.Env, common.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"))
	container.Env = append(container.Env, common.Env("STORAGE_DRIVER", "postgres"))
	container.Ports = []corev1.ContainerPort{deployments.StandardHTTPPort()}

	return nil
}

func (r *LedgerController) createBaseLedgerContainerV2() *corev1.Container {
	return &corev1.Container{
		Env: []corev1.EnvVar{
			common.Env("BIND", ":8080"),
		},
		Name: "ledger",
	}
}

func (r *LedgerController) createLedgerContainerV2Full(ctx common.Context, stack *v1beta1.Stack) (*corev1.Container, error) {
	container := r.createBaseLedgerContainerV2()
	needPublisher, err := topics.TopicExists(ctx, stack, "ledger")
	if err != nil {
		return nil, err
	}
	if needPublisher {
		brokerEnvVars, err := brokerconfigurations.GetEnvVars(ctx, stack.Name, "ledger")
		if err != nil {
			return nil, err
		}
		container.Env = append(container.Env, brokerEnvVars...)
		container.Env = append(container.Env, common.Env("PUBLISHER_TOPIC_MAPPING", "*:"+common.GetObjectName(stack.Name, "ledger")))
	}
	return container, nil
}

func (r *LedgerController) createLedgerContainerV2WriteOnly(ctx common.Context, stack *v1beta1.Stack) (*corev1.Container, error) {
	return r.createLedgerContainerV2Full(ctx, stack)
}

func (r *LedgerController) createLedgerContainerV2ReadOnly() *corev1.Container {
	container := r.createBaseLedgerContainerV2()
	container.Env = append(container.Env, common.Env("READ_ONLY", "true"))
	return container
}

func (r *LedgerController) createGatewayDeployment(ctx common.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger) error {

	caddyfileConfigMap, err := CreateCaddyfileConfigMap(ctx, stack, "ledger", ledgers.Caddyfile, map[string]any{
		"Debug": stack.Spec.Debug || ledger.Spec.Debug,
	}, common.WithController[*corev1.ConfigMap](ctx.GetScheme(), ledger))
	if err != nil {
		return err
	}

	env, err := GetCommonServicesEnvVars(ctx, stack, "ledger", ledger.Spec)
	if err != nil {
		return err
	}

	containerEnv := make([]corev1.EnvVar, 0)
	containerEnv = append(containerEnv, env...)

	mutators := common.ConfigureCaddy(caddyfileConfigMap, "caddy:2.7.6-alpine", containerEnv, nil)
	mutators = append(mutators,
		common.WithController[*appsv1.Deployment](ctx.GetScheme(), ledger),
		deployments.WithMatchingLabels("ledger"),
	)

	_, _, err = common.CreateOrUpdate[*appsv1.Deployment](ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "ledger-gateway",
	}, mutators...)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *LedgerController) SetupWithManager(mgr common.Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Ledger{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(
			&v1beta1.Topic{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
				topic := object.(*v1beta1.Topic)
				if topic.Spec.Service != "ledger" {
					return []reconcile.Request{}
				}

				ledgerList := &v1beta1.LedgerList{}
				if err := mgr.GetClient().List(ctx, ledgerList, client.MatchingFields{
					".spec.stack": topic.Spec.Stack,
				}); err != nil {
					return []reconcile.Request{}
				}

				return common.MapObjectToReconcileRequests(
					Map(ledgerList.Items, ToPointer[v1beta1.Ledger])...,
				)
			}),
		).
		Watches(
			&v1beta1.Database{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
				database := object.(*v1beta1.Database)
				if database.Spec.Service != "ledger" {
					return []reconcile.Request{}
				}

				ledgerList := &v1beta1.LedgerList{}
				if err := mgr.GetClient().List(ctx, ledgerList, client.MatchingFields{
					".spec.stack": database.Spec.Stack,
				}); err != nil {
					return []reconcile.Request{}
				}

				return common.MapObjectToReconcileRequests(
					Map(ledgerList.Items, ToPointer[v1beta1.Ledger])...,
				)
			}),
		).
		Owns(&appsv1.Deployment{}).
		Owns(&v1beta1.HTTPAPI{}).
		Owns(&v1beta1.Database{}), nil
}

func ForLedger() *LedgerController {
	return &LedgerController{}
}
