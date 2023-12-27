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
	"bytes"
	"context"
	_ "embed"
	"text/template"

	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/internal"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	pkgError "github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//go:embed templates/Caddyfile.ledger.gotpl
var ledgerCaddyfile string

// LedgerReconciler reconciles a Ledger object
type LedgerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=formance.com,resources=ledgers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=ledgers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=ledgers/finalizers,verbs=update

func (r *LedgerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx, "ledger", req.NamespacedName)
	log.Info("Starting reconciliation")

	ledger := &v1beta1.Ledger{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name: req.Name,
	}, ledger); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	stack := &v1beta1.Stack{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name: ledger.Spec.Stack,
	}, stack); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	database, err := CreateDatabase(ctx, r.Client, stack, "ledger")
	if err != nil {
		if pkgError.Is(err, ErrPending) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	err = r.installLedger(ctx, stack, ledger, database, GetVersion(stack, ledger.Spec.Version))
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := CreateHTTPAPI(ctx, r.Client, r.Scheme, stack, ledger, "ledger"); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Client.Status().Update(ctx, ledger); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *LedgerReconciler) installLedger(ctx context.Context, stack *v1beta1.Stack,
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

func (r *LedgerReconciler) installLedgerV2SingleInstance(ctx context.Context, stack *v1beta1.Stack,
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

func (r *LedgerReconciler) setInitContainer(database *v1beta1.Database, version string) func(t *appsv1.Deployment) {
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

func (r *LedgerReconciler) installLedgerV2MonoWriterMultipleReader(ctx context.Context, stack *v1beta1.Stack,
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

func (r *LedgerReconciler) uninstallLedgerV2MonoWriterMultipleReader(ctx context.Context, stack *v1beta1.Stack) error {

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

func (r *LedgerReconciler) createK8SService(ctx context.Context, stack *v1beta1.Stack, owner *v1beta1.Ledger, name string) error {
	_, _, err := CreateOrUpdate[*corev1.Service](ctx, r.Client, types.NamespacedName{
		Name:      name,
		Namespace: stack.Name,
	},
		ConfigureHTTPService(name),
		WithController[*corev1.Service](r.Scheme, owner),
	)
	return err
}

func (r *LedgerReconciler) createDeployment(ctx context.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger,
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

func (r *LedgerReconciler) setCommonContainerConfiguration(ctx context.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger,
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

func (r *LedgerReconciler) createBaseLedgerContainerV2() *corev1.Container {
	return &corev1.Container{
		Env: []corev1.EnvVar{
			Env("BIND", ":8080"),
		},
		Name: "ledger",
	}
}

func (r *LedgerReconciler) createLedgerContainerV2Full(ctx context.Context, stack *v1beta1.Stack) (*corev1.Container, error) {
	container := r.createBaseLedgerContainerV2()
	needPublisher, err := TopicExists(ctx, r.Client, stack, "ledger")
	if err != nil {
		return nil, err
	}
	if needPublisher {
		brokerEnvVars, err := GetBrokerEnvVarsIfEnabled(ctx, r.Client, stack.Name, "ledger")
		if err != nil {
			return nil, err
		}
		container.Env = append(container.Env, brokerEnvVars...)
		container.Env = append(container.Env, Env("PUBLISHER_TOPIC_MAPPING", "*:"+GetObjectName(stack.Name, "ledger")))
	}
	return container, nil
}

func (r *LedgerReconciler) createLedgerContainerV2WriteOnly(ctx context.Context, stack *v1beta1.Stack) (*corev1.Container, error) {
	return r.createLedgerContainerV2Full(ctx, stack)
}

func (r *LedgerReconciler) createLedgerContainerV2ReadOnly() *corev1.Container {
	container := r.createBaseLedgerContainerV2()
	container.Env = append(container.Env, Env("READ_ONLY", "true"))
	return container
}

func (r *LedgerReconciler) createGatewayDeployment(ctx context.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger) error {

	tpl := template.Must(template.New("ledger-caddyfile").Parse(ledgerCaddyfile))
	buf := bytes.NewBufferString("")

	openTelemetryEnabled, err := IsOpenTelemetryEnabled(ctx, r.Client, stack.Name)
	if err != nil {
		return err
	}

	if err := tpl.Execute(buf, map[string]any{
		"EnableOpenTelemetry": openTelemetryEnabled,
		"Debug":               stack.Spec.Debug || ledger.Spec.Debug,
	}); err != nil {
		return err
	}

	caddyfileConfigMap, _, err := CreateOrUpdate[*corev1.ConfigMap](ctx, r.Client, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "ledger-gateway",
	},
		func(t *corev1.ConfigMap) {
			t.Data = map[string]string{
				"Caddyfile": buf.String(),
			}
		},
		WithController[*corev1.ConfigMap](r.Scheme, ledger),
	)

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
func (r *LedgerReconciler) SetupWithManager(mgr ctrl.Manager) error {

	indexer := mgr.GetFieldIndexer()
	if err := indexer.IndexField(context.Background(), &v1beta1.Ledger{}, ".spec.stack", func(rawObj client.Object) []string {
		return []string{rawObj.(*v1beta1.Ledger).Spec.Stack}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Ledger{}).
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
		Owns(&v1beta1.Database{}).
		Complete(r)
}

func NewLedgerReconciler(client client.Client, scheme *runtime.Scheme) *LedgerReconciler {
	return &LedgerReconciler{
		Client: client,
		Scheme: scheme,
	}
}
