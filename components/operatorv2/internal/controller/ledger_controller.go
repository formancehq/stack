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

package controller

import (
	"context"
	_ "embed"
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/internal"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	pkgError "github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//go:embed templates/Caddyfile.ledger.gotpl
var ledgerCaddyfile string

// LedgerController reconciles a Ledger object
type LedgerController struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=formance.com,resources=ledgers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=ledgers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=ledgers/finalizers,verbs=update

func (r *LedgerController) Reconcile(ctx context.Context, ledger *v1beta1.Ledger) error {

	stack, err := GetStack(ctx, r.Client, ledger.Spec)
	if err != nil {
		return err
	}

	database, err := CreateDatabase(ctx, r.Client, stack, "ledger")
	if err != nil {
		if pkgError.Is(err, ErrPending) {
			return nil
		}
		return err
	}

	err = r.installLedger(ctx, stack, ledger, database, GetVersion(stack, ledger.Spec.Version))
	if err != nil {
		return err
	}

	if err := CreateHTTPAPI(ctx, r.Client, r.Scheme, stack, ledger, "ledger"); err != nil {
		return err
	}

	return nil
}

func (r *LedgerController) installLedger(ctx context.Context, stack *v1beta1.Stack,
	ledger *v1beta1.Ledger, database *v1beta1.Database, version string) error {

	switch ledger.Spec.DeploymentStrategy {
	case v1beta1.DeploymentStrategyMonoWriterMultipleReader:
		if err := DeleteIfExists[*appsv1.Deployment](ctx, r.Client, GetNamespacedResourceName(stack.Name, "ledger")); err != nil {
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

func (r *LedgerController) installLedgerV2SingleInstance(ctx context.Context, stack *v1beta1.Stack,
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
		WithReplicas(1),
		r.setInitContainer(database, version),
	); err != nil {
		return err
	}

	return nil
}

func (r *LedgerController) setInitContainer(database *v1beta1.Database, version string) func(t *appsv1.Deployment) {
	return func(t *appsv1.Deployment) {
		t.Spec.Template.Spec.InitContainers = []corev1.Container{
			MigrateDatabaseContainer(
				"ledger",
				database.Status.Configuration.DatabaseConfigurationSpec,
				database.Status.Configuration.Database,
				version,
				func(m *MigrationConfiguration) {
					m.Command = []string{"buckets", "upgrade-all"}
					m.AdditionalEnv = []corev1.EnvVar{
						Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"),
					}
				},
			),
		}
	}
}

func (r *LedgerController) installLedgerV2MonoWriterMultipleReader(ctx context.Context, stack *v1beta1.Stack,
	ledger *v1beta1.Ledger, database *v1beta1.Database, version string) error {

	createDeployment := func(name string, container corev1.Container, mutators ...ObjectMutator[*appsv1.Deployment]) error {
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
		WithReplicas(1),
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

func (r *LedgerController) uninstallLedgerV2MonoWriterMultipleReader(ctx context.Context, stack *v1beta1.Stack) error {

	remove := func(name string) error {
		if err := DeleteIfExists[*appsv1.Deployment](ctx, r.Client, GetNamespacedResourceName(stack.Name, name)); err != nil {
			return err
		}
		if err := DeleteIfExists[*corev1.Service](ctx, r.Client, GetNamespacedResourceName(stack.Name, name)); err != nil {
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

	if err := DeleteIfExists[*appsv1.Deployment](ctx, r.Client, GetNamespacedResourceName(stack.Name, "ledger-gateway")); err != nil {
		return err
	}

	return nil
}

func (r *LedgerController) createK8SService(ctx context.Context, stack *v1beta1.Stack, owner *v1beta1.Ledger, name string) error {
	_, _, err := CreateOrUpdate[*corev1.Service](ctx, r.Client, types.NamespacedName{
		Name:      name,
		Namespace: stack.Name,
	},
		ConfigureK8SService(name),
		WithController[*corev1.Service](r.Scheme, owner),
	)
	return err
}

func (r *LedgerController) createDeployment(ctx context.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger,
	name string, container corev1.Container, mutators ...ObjectMutator[*appsv1.Deployment]) error {
	mutators = append([]ObjectMutator[*appsv1.Deployment]{
		WithContainers(container),
		WithMatchingLabels(name),
		WithController[*appsv1.Deployment](r.Scheme, ledger),
	}, mutators...)
	_, _, err := CreateOrUpdate[*appsv1.Deployment](ctx, r.Client,
		GetNamespacedResourceName(stack.Name, name),
		mutators...,
	)
	return err
}

func (r *LedgerController) setCommonContainerConfiguration(ctx context.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger,
	version string, database *v1beta1.Database, container *corev1.Container) error {

	env, err := GetCommonServicesEnvVars(ctx, r.Client, stack, "ledger", ledger.Spec)
	if err != nil {
		return err
	}

	container.Resources = GetResourcesWithDefault(ledger.Spec.ResourceProperties, ResourceSizeSmall())
	container.Image = GetImage("ledger", version)
	container.ImagePullPolicy = GetPullPolicy(container.Image)
	container.Env = append(container.Env, env...)
	container.Env = append(container.Env, PostgresEnvVars(
		database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database)...)
	container.Env = append(container.Env, Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"))
	container.Env = append(container.Env, Env("STORAGE_DRIVER", "postgres"))
	container.Ports = []corev1.ContainerPort{StandardHTTPPort()}

	return nil
}

func (r *LedgerController) createBaseLedgerContainerV2() *corev1.Container {
	return &corev1.Container{
		Env: []corev1.EnvVar{
			Env("BIND", ":8080"),
		},
		Name: "ledger",
	}
}

func (r *LedgerController) createLedgerContainerV2Full(ctx context.Context, stack *v1beta1.Stack) (*corev1.Container, error) {
	container := r.createBaseLedgerContainerV2()
	needPublisher, err := TopicExists(ctx, r.Client, stack, "ledger")
	if err != nil {
		return nil, err
	}
	if needPublisher {
		brokerEnvVars, err := GetBrokerEnvVars(ctx, r.Client, stack.Name, "ledger")
		if err != nil {
			return nil, err
		}
		container.Env = append(container.Env, brokerEnvVars...)
		container.Env = append(container.Env, Env("PUBLISHER_TOPIC_MAPPING", "*:"+GetObjectName(stack.Name, "ledger")))
	}
	return container, nil
}

func (r *LedgerController) createLedgerContainerV2WriteOnly(ctx context.Context, stack *v1beta1.Stack) (*corev1.Container, error) {
	return r.createLedgerContainerV2Full(ctx, stack)
}

func (r *LedgerController) createLedgerContainerV2ReadOnly() *corev1.Container {
	container := r.createBaseLedgerContainerV2()
	container.Env = append(container.Env, Env("READ_ONLY", "true"))
	return container
}

func (r *LedgerController) createGatewayDeployment(ctx context.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger) error {

	caddyfileConfigMap, err := CreateCaddyfileConfigMap(ctx, r.Client, stack, "ledger", ledgerCaddyfile, map[string]any{
		"Debug": stack.Spec.Debug || ledger.Spec.Debug,
	}, WithController[*corev1.ConfigMap](r.Scheme, ledger))
	if err != nil {
		return err
	}

	env, err := GetCommonServicesEnvVars(ctx, r.Client, stack, "ledger", ledger.Spec)
	if err != nil {
		return err
	}

	containerEnv := make([]corev1.EnvVar, 0)
	containerEnv = append(containerEnv, env...)

	mutators := ConfigureCaddy(caddyfileConfigMap, "caddy:2.7.6-alpine", containerEnv, nil)
	mutators = append(mutators,
		WithController[*appsv1.Deployment](r.Scheme, ledger),
		WithMatchingLabels("ledger"),
	)

	_, _, err = CreateOrUpdate[*appsv1.Deployment](ctx, r.Client, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "ledger-gateway",
	}, mutators...)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *LedgerController) SetupWithManager(mgr ctrl.Manager) (*builder.Builder, error) {

	indexer := mgr.GetFieldIndexer()
	if err := indexer.IndexField(context.Background(), &v1beta1.Ledger{}, ".spec.stack", func(rawObj client.Object) []string {
		return []string{rawObj.(*v1beta1.Ledger).Spec.Stack}
	}); err != nil {
		return nil, err
	}

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

				return MapObjectToReconcileRequests(
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

				return MapObjectToReconcileRequests(
					Map(ledgerList.Items, ToPointer[v1beta1.Ledger])...,
				)
			}),
		).
		Owns(&appsv1.Deployment{}).
		Owns(&v1beta1.HTTPAPI{}).
		Owns(&v1beta1.Database{}), nil
}

func ForLedger(client client.Client, scheme *runtime.Scheme) *LedgerController {
	return &LedgerController{
		Client: client,
		Scheme: scheme,
	}
}
