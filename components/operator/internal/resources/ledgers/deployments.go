package ledgers

import (
	"fmt"
	"strconv"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/brokerconfigurations"
	"github.com/formancehq/operator/internal/resources/brokertopics"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/opentelemetryconfigurations"
	"github.com/formancehq/operator/internal/resources/services"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v14 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func installLedger(ctx core.Context, stack *v1beta1.Stack,
	ledger *v1beta1.Ledger, database *v1beta1.Database, image string, isV2 bool) error {

	switch ledger.Spec.DeploymentStrategy {
	case v1beta1.DeploymentStrategyMonoWriterMultipleReader:
		if err := core.DeleteIfExists[*v1.Deployment](ctx, core.GetNamespacedResourceName(stack.Name, "ledger")); err != nil {
			return err
		}
		return installLedgerMonoWriterMultipleReader(ctx, stack, ledger, database, image, isV2)
	default:
		if err := uninstallLedgerMonoWriterMultipleReader(ctx, stack); err != nil {
			return err
		}
		return installLedgerSingleInstance(ctx, stack, ledger, database, image, isV2)
	}
}

func installLedgerSingleInstance(ctx core.Context, stack *v1beta1.Stack,
	ledger *v1beta1.Ledger, database *v1beta1.Database, version string, v2 bool) error {
	container, err := createLedgerContainerFull(ctx, stack, v2)
	if err != nil {
		return err
	}

	err = setCommonContainerConfiguration(ctx, stack, ledger, version, database, container, v2)
	if err != nil {
		return err
	}

	if !v2 && ledger.Spec.Locking.Strategy == "redis" {
		container.Env = append(container.Env,
			core.Env("NUMARY_LOCK_STRATEGY", "redis"),
			core.Env("NUMARY_LOCK_STRATEGY_REDIS_URL", ledger.Spec.Locking.Redis.Uri),
			core.Env("NUMARY_LOCK_STRATEGY_REDIS_TLS_ENABLED", strconv.FormatBool(ledger.Spec.Locking.Redis.TLS)),
			core.Env("NUMARY_LOCK_STRATEGY_REDIS_TLS_INSECURE", strconv.FormatBool(ledger.Spec.Locking.Redis.InsecureTLS)),
		)

		if ledger.Spec.Locking.Redis.Duration != 0 {
			container.Env = append(container.Env, core.Env("NUMARY_LOCK_STRATEGY_REDIS_DURATION", ledger.Spec.Locking.Redis.Duration.String()))
		}

		if ledger.Spec.Locking.Redis.Retry != 0 {
			container.Env = append(container.Env, core.Env("NUMARY_LOCK_STRATEGY_REDIS_RETRY", ledger.Spec.Locking.Redis.Retry.String()))
		}
	}

	if err := createDeployment(ctx, ledger, "ledger", *container,
		deployments.WithReplicas(1),
		setInitContainer(database, version, v2),
	); err != nil {
		return err
	}

	return nil
}

func getUpgradeContainer(database *v1beta1.Database, image string) corev1.Container {
	return databases.MigrateDatabaseContainer(
		image,
		database.Status.Configuration.DatabaseConfigurationSpec,
		database.Status.Configuration.Database,
		func(m *databases.MigrationConfiguration) {
			m.Command = []string{"buckets", "upgrade-all"}
			m.AdditionalEnv = []corev1.EnvVar{
				core.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"),
			}
		},
	)
}

func setInitContainer(database *v1beta1.Database, image string, v2 bool) func(t *v1.Deployment) {
	return func(t *v1.Deployment) {
		if !v2 {
			t.Spec.Template.Spec.InitContainers = []corev1.Container{}
			return
		}
		t.Spec.Template.Spec.InitContainers = []corev1.Container{getUpgradeContainer(database, image)}
	}
}

func installLedgerMonoWriterMultipleReader(ctx core.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger, database *v1beta1.Database, image string, v2 bool) error {

	createDeployment := func(name string, container corev1.Container, mutators ...core.ObjectMutator[*v1.Deployment]) error {
		err := setCommonContainerConfiguration(ctx, stack, ledger, image, database, &container, v2)
		if err != nil {
			return err
		}

		if err := createDeployment(ctx, ledger, name, container, mutators...); err != nil {
			return err
		}

		if _, err := services.Create(ctx, ledger, name); err != nil {
			return err
		}

		return nil
	}

	container, err := createLedgerContainerWriteOnly(ctx, stack, v2)
	if err != nil {
		return err
	}
	if err := createDeployment("ledger-write", *container,
		deployments.WithReplicas(1),
		setInitContainer(database, image, v2),
	); err != nil {
		return err
	}

	container = createLedgerContainerReadOnly(v2)
	if err := createDeployment("ledger-read", *container); err != nil {
		return err
	}

	if err := createGatewayDeployment(ctx, stack, ledger); err != nil {
		return err
	}

	return nil
}

func uninstallLedgerMonoWriterMultipleReader(ctx core.Context, stack *v1beta1.Stack) error {

	remove := func(name string) error {
		if err := core.DeleteIfExists[*v1.Deployment](ctx, core.GetNamespacedResourceName(stack.Name, name)); err != nil {
			return err
		}
		if err := core.DeleteIfExists[*corev1.Service](ctx, core.GetNamespacedResourceName(stack.Name, name)); err != nil {
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

	if err := core.DeleteIfExists[*v1.Deployment](ctx, core.GetNamespacedResourceName(stack.Name, "ledger-gateway")); err != nil {
		return err
	}

	return nil
}

func createDeployment(ctx core.Context, ledger *v1beta1.Ledger,
	name string, container corev1.Container, mutators ...core.ObjectMutator[*v1.Deployment]) error {
	mutators = append([]core.ObjectMutator[*v1.Deployment]{
		deployments.WithContainers(container),
		deployments.WithMatchingLabels(name),
	}, mutators...)

	_, err := deployments.CreateOrUpdate(ctx, ledger, name, mutators...)
	return err
}

func setCommonContainerConfiguration(ctx core.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger, image string, database *v1beta1.Database, container *corev1.Container, v2 bool) error {

	prefix := ""
	if !v2 {
		prefix = "NUMARY_"
	}
	env := make([]corev1.EnvVar, 0)
	otlpEnv, err := opentelemetryconfigurations.EnvVarsIfEnabledWithPrefix(ctx, stack.Name, core.GetModuleName(ctx, ledger), prefix)
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)

	gatewayEnv, err := gateways.EnvVarsIfEnabledWithPrefix(ctx, stack.Name, prefix)
	if err != nil {
		return err
	}
	env = append(env, gatewayEnv...)
	env = append(env, core.GetDevEnvVarsWithPrefix(stack, ledger, prefix)...)

	authEnvVars, err := auths.ProtectedAPIEnvVarsWithPrefix(ctx, stack, "ledger", ledger.Spec.Auth, prefix)
	if err != nil {
		return err
	}
	env = append(env, authEnvVars...)

	container.Resources = core.GetResourcesRequirementsWithDefault(ledger.Spec.ResourceRequirements, core.ResourceSizeSmall())
	container.Image = image
	container.Env = append(container.Env, env...)
	container.Env = append(container.Env, databases.PostgresEnvVarsWithPrefix(
		database.Status.Configuration.DatabaseConfigurationSpec, database.Status.Configuration.Database, prefix)...)
	container.Env = append(container.Env, core.Env(fmt.Sprintf("%sSTORAGE_POSTGRES_CONN_STRING", prefix), fmt.Sprintf("$(%sPOSTGRES_URI)", prefix)))
	container.Env = append(container.Env, core.Env(fmt.Sprintf("%sSTORAGE_DRIVER", prefix), "postgres"))
	container.Ports = []corev1.ContainerPort{deployments.StandardHTTPPort()}
	container.LivenessProbe = deployments.DefaultLiveness("http")

	return nil
}

func createBaseLedgerContainer(v2 bool) *corev1.Container {
	ret := &corev1.Container{
		Name: "ledger",
	}
	var bindFlag = "BIND"
	if !v2 {
		bindFlag = "NUMARY_SERVER_HTTP_BIND_ADDRESS"
	}
	ret.Env = append(ret.Env, core.Env(bindFlag, ":8080"))

	return ret
}

func createLedgerContainerFull(ctx core.Context, stack *v1beta1.Stack, v2 bool) (*corev1.Container, error) {
	container := createBaseLedgerContainer(v2)
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
		container.Env = append(container.Env, core.Env(fmt.Sprintf("%sPUBLISHER_TOPIC_MAPPING", prefix), "*:"+core.GetObjectName(stack.Name, "ledger")))
	}

	return container, nil
}

func createLedgerContainerWriteOnly(ctx core.Context, stack *v1beta1.Stack, v2 bool) (*corev1.Container, error) {
	return createLedgerContainerFull(ctx, stack, v2)
}

func createLedgerContainerReadOnly(v2 bool) *corev1.Container {
	container := createBaseLedgerContainer(v2)
	container.Env = append(container.Env, core.Env("READ_ONLY", "true"))
	return container
}

func createGatewayDeployment(ctx core.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger) error {

	caddyfileConfigMap, err := core.CreateCaddyfileConfigMap(ctx, stack, "ledger", Caddyfile, map[string]any{
		"Debug": stack.Spec.Debug || ledger.Spec.Debug,
	}, core.WithController[*corev1.ConfigMap](ctx.GetScheme(), ledger))
	if err != nil {
		return err
	}

	env := make([]corev1.EnvVar, 0)
	otlpEnv, err := opentelemetryconfigurations.EnvVarsIfEnabled(ctx, stack.Name, core.GetModuleName(ctx, ledger))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)
	env = append(env, core.GetDevEnvVars(stack, ledger)...)

	_, err = deployments.CreateOrUpdate(ctx, ledger, "ledger-gateway",
		core.ConfigureCaddy(caddyfileConfigMap, "caddy:2.7.6-alpine", env, nil),
		deployments.WithMatchingLabels("ledger"),
	)
	return err
}

func migrateToLedgerV2(ctx core.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger, database *v1beta1.Database, image string) error {

	job := &v13.Job{}
	err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "migrate-v2",
	}, job)
	if client.IgnoreNotFound(err) != nil {
		return err
	}
	if err == nil {
		if job.Status.Succeeded == 0 {
			return core.ErrPending
		}
		return nil
	}

	list := &v1.DeploymentList{}
	if err := ctx.GetClient().List(ctx, list, client.InNamespace(stack.Name)); err != nil {
		return err
	}

	for _, item := range list.Items {
		if controller := v14.GetControllerOf(&item); controller != nil && controller.UID == ledger.GetUID() {
			if err := ctx.GetClient().Delete(ctx, &item); err != nil {
				return err
			}
		}
	}

	return ctx.GetClient().Create(ctx, &v13.Job{
		ObjectMeta: v14.ObjectMeta{
			Namespace: stack.Name,
			Name:      "migrate-v2",
		},
		Spec: v13.JobSpec{
			BackoffLimit:            pointer.For(int32(10000)),
			TTLSecondsAfterFinished: pointer.For(int32(30)),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers:    []corev1.Container{getUpgradeContainer(database, image)},
				},
			},
		},
	})
}