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

package formance_com

import (
	_ "embed"
	"fmt"
	"golang.org/x/mod/semver"
	"strconv"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/brokerconfigurations"
	"github.com/formancehq/operator/internal/resources/brokertopics"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/httpapis"
	"github.com/formancehq/operator/internal/resources/ledgers"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/services"
	"github.com/formancehq/operator/internal/resources/streams"
	"github.com/formancehq/search/benthos"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// LedgerController reconciles a Ledger object
type LedgerController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=ledgers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=ledgers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=ledgers/finalizers,verbs=update

func (r *LedgerController) Reconcile(ctx Context, ledger *v1beta1.Ledger) error {

	stack, err := GetStack(ctx, ledger)
	if err != nil {
		return err
	}

	database, err := databases.Create(ctx, ledger)
	if err != nil {
		return err
	}

	image, err := registries.GetImage(ctx, stack, "ledger", ledger.Spec.Version)
	if err != nil {
		return err
	}

	if err := httpapis.Create(ctx, ledger,
		httpapis.WithServiceConfiguration(ledger.Spec.Service)); err != nil {
		return err
	}

	isV2 := false
	moduleVersion := GetModuleVersion(stack, ledger.Spec.Version)
	if !semver.IsValid(moduleVersion) || semver.Compare(moduleVersion, "v2.0.0-alpha") > 0 {
		isV2 = true
	}

	hasSearch, err := HasDependency[*v1beta1.Search](ctx, stack.Name)
	if err != nil {
		return err
	}
	if hasSearch {
		streamsVersion := "v1.0.0"
		if isV2 {
			streamsVersion = "v2.0.0"
		}
		if err := streams.LoadFromFileSystem(ctx, benthos.Streams, ledger.Spec.Stack, "streams/ledger/"+streamsVersion,
			WithController[*v1beta1.Stream](ctx.GetScheme(), ledger),
			WithLabels[*v1beta1.Stream](map[string]string{
				"service": "ledger",
			}),
		); err != nil {
			return err
		}
	} else {
		if err := ctx.GetClient().DeleteAllOf(ctx, &v1beta1.Stream{}, client.MatchingLabels{
			"service": "ledger",
		}); err != nil {
			return err
		}
	}

	if database.Status.Ready {
		err = r.installLedger(ctx, stack, ledger, database, image, isV2)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *LedgerController) installLedger(ctx Context, stack *v1beta1.Stack,
	ledger *v1beta1.Ledger, database *v1beta1.Database, image string, isV2 bool) error {

	switch ledger.Spec.DeploymentStrategy {
	case v1beta1.DeploymentStrategyMonoWriterMultipleReader:
		if err := DeleteIfExists[*appsv1.Deployment](ctx, GetNamespacedResourceName(stack.Name, "ledger")); err != nil {
			return err
		}
		return r.installLedgerMonoWriterMultipleReader(ctx, stack, ledger, database, image, isV2)
	default:
		if err := r.uninstallLedgerMonoWriterMultipleReader(ctx, stack); err != nil {
			return err
		}
		return r.installLedgerSingleInstance(ctx, stack, ledger, database, image, isV2)
	}
}

func (r *LedgerController) installLedgerSingleInstance(ctx Context, stack *v1beta1.Stack,
	ledger *v1beta1.Ledger, database *v1beta1.Database, version string, v2 bool) error {
	container, err := r.createLedgerContainerFull(ctx, stack, v2)
	if err != nil {
		return err
	}

	err = r.setCommonContainerConfiguration(ctx, stack, ledger, version, database, container, v2)
	if err != nil {
		return err
	}

	if !v2 && ledger.Spec.LockingStrategy.Strategy == "redis" {
		container.Env = append(container.Env,
			Env("NUMARY_LOCK_STRATEGY", "redis"),
			Env("NUMARY_LOCK_STRATEGY_REDIS_URL", ledger.Spec.LockingStrategy.Redis.Uri),
			Env("NUMARY_LOCK_STRATEGY_REDIS_TLS_ENABLED", strconv.FormatBool(ledger.Spec.LockingStrategy.Redis.TLS)),
			Env("NUMARY_LOCK_STRATEGY_REDIS_TLS_INSECURE", strconv.FormatBool(ledger.Spec.LockingStrategy.Redis.InsecureTLS)),
		)

		if ledger.Spec.LockingStrategy.Redis.Duration != 0 {
			container.Env = append(container.Env, Env("NUMARY_LOCK_STRATEGY_REDIS_DURATION", ledger.Spec.LockingStrategy.Redis.Duration.String()))
		}

		if ledger.Spec.LockingStrategy.Redis.Retry != 0 {
			container.Env = append(container.Env, Env("NUMARY_LOCK_STRATEGY_REDIS_RETRY", ledger.Spec.LockingStrategy.Redis.Retry.String()))
		}
	}

	if err := r.createDeployment(ctx, ledger, "ledger", *container,
		deployments.WithReplicas(1),
		r.setInitContainer(database, version, v2),
	); err != nil {
		return err
	}

	return nil
}

func (r *LedgerController) setInitContainer(database *v1beta1.Database, image string, v2 bool) func(t *appsv1.Deployment) {
	return func(t *appsv1.Deployment) {
		if !v2 {
			t.Spec.Template.Spec.InitContainers = []corev1.Container{}
			return
		}
		t.Spec.Template.Spec.InitContainers = []corev1.Container{
			databases.MigrateDatabaseContainer(
				image,
				database.Status.Configuration.DatabaseConfigurationSpec,
				database.Status.Configuration.Database,
				func(m *databases.MigrationConfiguration) {
					m.Command = []string{"buckets", "upgrade-all"}
					m.AdditionalEnv = []corev1.EnvVar{
						Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"),
					}
				},
			),
		}
	}
}

func (r *LedgerController) installLedgerMonoWriterMultipleReader(ctx Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger, database *v1beta1.Database, image string, v2 bool) error {

	createDeployment := func(name string, container corev1.Container, mutators ...ObjectMutator[*appsv1.Deployment]) error {
		err := r.setCommonContainerConfiguration(ctx, stack, ledger, image, database, &container, v2)
		if err != nil {
			return err
		}

		if err := r.createDeployment(ctx, ledger, name, container, mutators...); err != nil {
			return err
		}

		if err := r.createK8SService(ctx, stack, ledger, name); err != nil {
			return err
		}

		return nil
	}

	container, err := r.createLedgerContainerWriteOnly(ctx, stack, v2)
	if err != nil {
		return err
	}
	if err := createDeployment("ledger-write", *container,
		deployments.WithReplicas(1),
		r.setInitContainer(database, image, v2),
	); err != nil {
		return err
	}

	container = r.createLedgerContainerReadOnly(v2)
	if err := createDeployment("ledger-read", *container); err != nil {
		return err
	}

	if err := r.createGatewayDeployment(ctx, stack, ledger); err != nil {
		return err
	}

	return nil
}

func (r *LedgerController) uninstallLedgerMonoWriterMultipleReader(ctx Context, stack *v1beta1.Stack) error {

	remove := func(name string) error {
		if err := DeleteIfExists[*appsv1.Deployment](ctx, GetNamespacedResourceName(stack.Name, name)); err != nil {
			return err
		}
		if err := DeleteIfExists[*corev1.Service](ctx, GetNamespacedResourceName(stack.Name, name)); err != nil {
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

	if err := DeleteIfExists[*appsv1.Deployment](ctx, GetNamespacedResourceName(stack.Name, "ledger-gateway")); err != nil {
		return err
	}

	return nil
}

func (r *LedgerController) createK8SService(ctx Context, stack *v1beta1.Stack, owner *v1beta1.Ledger, name string) error {
	_, _, err := CreateOrUpdate[*corev1.Service](ctx, types.NamespacedName{
		Name:      name,
		Namespace: stack.Name,
	},
		services.ConfigureK8SService(name),
		WithController[*corev1.Service](ctx.GetScheme(), owner),
	)
	return err
}

func (r *LedgerController) createDeployment(ctx Context, ledger *v1beta1.Ledger,
	name string, container corev1.Container, mutators ...ObjectMutator[*appsv1.Deployment]) error {
	mutators = append([]ObjectMutator[*appsv1.Deployment]{
		deployments.WithContainers(container),
		deployments.WithMatchingLabels(name),
	}, mutators...)

	_, err := deployments.CreateOrUpdate(ctx, ledger, name, mutators...)
	return err
}

func (r *LedgerController) setCommonContainerConfiguration(ctx Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger, image string, database *v1beta1.Database, container *corev1.Container, v2 bool) error {

	prefix := ""
	if !v2 {
		prefix = "NUMARY_"
	}
	env, err := GetCommonModuleEnvVarsWithPrefix(ctx, stack, ledger, prefix)
	if err != nil {
		return err
	}

	authEnvVars, err := auths.EnvVarsWithPrefix(ctx, stack, "ledger", ledger.Spec.Auth, prefix)
	if err != nil {
		return err
	}
	env = append(env, authEnvVars...)

	container.Resources = GetResourcesRequirementsWithDefault(ledger.Spec.ResourceRequirements, ResourceSizeSmall())
	container.Image = image
	container.Env = append(container.Env, env...)
	container.Env = append(container.Env, databases.PostgresEnvVarsWithPrefix(
		database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database, prefix)...)
	container.Env = append(container.Env, Env(fmt.Sprintf("%sSTORAGE_POSTGRES_CONN_STRING", prefix), fmt.Sprintf("$(%sPOSTGRES_URI)", prefix)))
	container.Env = append(container.Env, Env(fmt.Sprintf("%sSTORAGE_DRIVER", prefix), "postgres"))
	container.Ports = []corev1.ContainerPort{deployments.StandardHTTPPort()}
	container.LivenessProbe = deployments.DefaultLiveness("http")

	return nil
}

func (r *LedgerController) createBaseLedgerContainer(v2 bool) *corev1.Container {
	ret := &corev1.Container{
		Name: "ledger",
	}
	var bindFlag = "BIND"
	if !v2 {
		bindFlag = "NUMARY_SERVER_HTTP_BIND_ADDRESS"
	}
	ret.Env = append(ret.Env, Env(bindFlag, ":8080"))

	return ret
}

func (r *LedgerController) createLedgerContainerFull(ctx Context, stack *v1beta1.Stack, v2 bool) (*corev1.Container, error) {
	container := r.createBaseLedgerContainer(v2)
	topic, err := brokertopics.Find(ctx, stack, "ledger")
	if err != nil {
		return nil, err
	}

	if topic != nil {
		if !topic.Status.Ready {
			return nil, fmt.Errorf("topic %s is not yet ready", topic.Name)
		}

		prefix := ""
		if !v2 {
			prefix = "NUMARY_"
		}

		container.Env = append(container.Env, brokerconfigurations.BrokerEnvVarsWithPrefix(*topic.Status.Configuration, stack.Name, "ledger", prefix)...)
		container.Env = append(container.Env, Env(fmt.Sprintf("%sPUBLISHER_TOPIC_MAPPING", prefix), "*:"+GetObjectName(stack.Name, "ledger")))
	}

	return container, nil
}

func (r *LedgerController) createLedgerContainerWriteOnly(ctx Context, stack *v1beta1.Stack, v2 bool) (*corev1.Container, error) {
	return r.createLedgerContainerFull(ctx, stack, v2)
}

func (r *LedgerController) createLedgerContainerReadOnly(v2 bool) *corev1.Container {
	container := r.createBaseLedgerContainer(v2)
	container.Env = append(container.Env, Env("READ_ONLY", "true"))
	return container
}

func (r *LedgerController) createGatewayDeployment(ctx Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger) error {

	caddyfileConfigMap, err := CreateCaddyfileConfigMap(ctx, stack, "ledger", ledgers.Caddyfile, map[string]any{
		"Debug": stack.Spec.Debug || ledger.Spec.Debug,
	}, WithController[*corev1.ConfigMap](ctx.GetScheme(), ledger))
	if err != nil {
		return err
	}

	env, err := GetCommonModuleEnvVars(ctx, stack, ledger)
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, ledger, "ledger-gateway",
		ConfigureCaddy(caddyfileConfigMap, "caddy:2.7.6-alpine", env, nil),
		deployments.WithMatchingLabels("ledger"),
	)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *LedgerController) SetupWithManager(mgr Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Ledger{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(&v1beta1.Stack{}, handler.EnqueueRequestsFromMapFunc(Watch[*v1beta1.Ledger](mgr))).
		Watches(
			&v1beta1.BrokerTopic{},
			handler.EnqueueRequestsFromMapFunc(
				brokertopics.Watch[*v1beta1.Ledger](mgr, "ledger")),
		).
		Watches(
			&v1beta1.Database{},
			handler.EnqueueRequestsFromMapFunc(
				databases.Watch[*v1beta1.Ledger](mgr, "ledger")),
		).
		Watches(
			&v1beta1.OpenTelemetryConfiguration{},
			handler.EnqueueRequestsFromMapFunc(WatchUsingLabels[*v1beta1.Ledger](mgr)),
		).
		Watches(
			&v1beta1.RegistriesConfiguration{},
			handler.EnqueueRequestsFromMapFunc(WatchUsingLabels[*v1beta1.Ledger](mgr)),
		).
		Watches(
			&v1beta1.Search{},
			handler.EnqueueRequestsFromMapFunc(WatchDependents[*v1beta1.Ledger](mgr)),
		).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&v1beta1.HTTPAPI{}).
		Owns(&v1beta1.Database{}), nil
}

func ForLedger() *LedgerController {
	return &LedgerController{}
}
