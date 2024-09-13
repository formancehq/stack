package ledgers

import (
	"fmt"
	"strconv"

	"github.com/formancehq/operator/internal/resources/brokers"
	"github.com/formancehq/operator/internal/resources/brokertopics"
	"github.com/formancehq/operator/internal/resources/caddy"
	"k8s.io/apimachinery/pkg/types"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/applications"
	"github.com/formancehq/operator/internal/resources/auths"
	"github.com/formancehq/operator/internal/resources/databases"
	"github.com/formancehq/operator/internal/resources/gateways"
	"github.com/formancehq/operator/internal/resources/jobs"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/services"
	"github.com/formancehq/operator/internal/resources/settings"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func installLedger(ctx core.Context, stack *v1beta1.Stack,
	ledger *v1beta1.Ledger, database *v1beta1.Database, image string, isV2 bool) error {

	deploymentStrategySettings, err := settings.GetStringOrDefault(ctx, stack.Name, v1beta1.DeploymentStrategySingle, "ledger", "deployment-strategy")
	if err != nil {
		return err
	}

	isSingle := true
	if deploymentStrategySettings == v1beta1.DeploymentStrategyMonoWriterMultipleReader {
		isSingle = false
	}
	if ledger.Spec.DeploymentStrategy == v1beta1.DeploymentStrategyMonoWriterMultipleReader {
		isSingle = false
	}

	if isSingle {
		if err := uninstallLedgerMonoWriterMultipleReader(ctx, stack); err != nil {
			return err
		}
		return installLedgerSingleInstance(ctx, stack, ledger, database, image, isV2)
	}

	if err := core.DeleteIfExists[*appsv1.Deployment](ctx, core.GetNamespacedResourceName(stack.Name, "ledger")); err != nil {
		return err
	}
	return installLedgerMonoWriterMultipleReader(ctx, stack, ledger, database, image, isV2)
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

	if !v2 && ledger.Spec.Locking != nil && ledger.Spec.Locking.Strategy == "redis" {
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

	if err := createDeployment(ctx, stack, ledger, "ledger", *container, v2, 1); err != nil {
		return err
	}

	return nil
}

func getUpgradeContainer(ctx core.Context, stack *v1beta1.Stack, database *v1beta1.Database, image, version string) (corev1.Container, error) {
	return databases.MigrateDatabaseContainer(ctx, stack, image, database,
		func(m *databases.MigrationConfiguration) {
			if core.IsLower(version, "v2.0.0-rc.6") {
				m.Command = []string{"buckets", "upgrade-all"}
			}
			m.AdditionalEnv = []corev1.EnvVar{
				core.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"),
			}
		},
	)
}

func installLedgerMonoWriterMultipleReader(ctx core.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger, database *v1beta1.Database, image string, v2 bool) error {

	createDeployment := func(name string, container corev1.Container, replicas uint64) error {
		err := setCommonContainerConfiguration(ctx, stack, ledger, image, database, &container, v2)
		if err != nil {
			return err
		}

		if err := createDeployment(ctx, stack, ledger, name, container, v2, replicas); err != nil {
			return err
		}

		if _, err := services.Create(ctx, ledger, name, services.WithDefault(name)); err != nil {
			return err
		}

		return nil
	}

	container, err := createLedgerContainerWriteOnly(ctx, stack, v2)
	if err != nil {
		return err
	}
	if err := createDeployment("ledger-write", *container, 1); err != nil {
		return err
	}

	container = createLedgerContainerReadOnly(v2)
	if err := createDeployment("ledger-read", *container, 0); err != nil {
		return err
	}

	if err := createGatewayDeployment(ctx, stack, ledger); err != nil {
		return err
	}

	return nil
}

func uninstallLedgerMonoWriterMultipleReader(ctx core.Context, stack *v1beta1.Stack) error {

	remove := func(name string) error {
		if err := core.DeleteIfExists[*appsv1.Deployment](ctx, core.GetNamespacedResourceName(stack.Name, name)); err != nil {
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

	if err := core.DeleteIfExists[*appsv1.Deployment](ctx, core.GetNamespacedResourceName(stack.Name, "ledger-gateway")); err != nil {
		return err
	}

	return nil
}

func createDeployment(ctx core.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger, name string, container corev1.Container, v2 bool, replicas uint64) error {
	serviceAccountName, err := settings.GetAWSServiceAccount(ctx, stack.Name)
	if err != nil {
		return err
	}

	var volumes []corev1.Volume
	if !v2 {
		volumes = []corev1.Volume{{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		}}
	}

	tpl := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers:         []corev1.Container{container},
					Volumes:            volumes,
					ServiceAccountName: serviceAccountName,
				},
			},
		},
	}

	return applications.
		New(ledger, tpl).
		WithStateful(replicas > 0).
		Install(ctx)
}

func setCommonContainerConfiguration(ctx core.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger, image string, database *v1beta1.Database, container *corev1.Container, v2 bool) error {

	prefix := ""
	if !v2 {
		prefix = "NUMARY_"
	}
	env := make([]corev1.EnvVar, 0)
	otlpEnv, err := settings.GetOTELEnvVarsWithPrefix(ctx, stack.Name, core.LowerCamelCaseKind(ctx, ledger), prefix)
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

	postgresEnvVar, err := databases.PostgresEnvVarsWithPrefix(ctx, stack, database, prefix)
	if err != nil {
		return err
	}
	env = append(env, postgresEnvVar...)

	container.Image = image
	container.Env = append(container.Env, env...)
	container.Env = append(container.Env, core.Env(fmt.Sprintf("%sSTORAGE_POSTGRES_CONN_STRING", prefix), fmt.Sprintf("$(%sPOSTGRES_URI)", prefix)))
	container.Env = append(container.Env, core.Env(fmt.Sprintf("%sSTORAGE_DRIVER", prefix), "postgres"))
	container.Ports = []corev1.ContainerPort{applications.StandardHTTPPort()}
	container.LivenessProbe = applications.DefaultLiveness("http")

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
	if !v2 {
		ret.VolumeMounts = []corev1.VolumeMount{{
			Name:      "config",
			ReadOnly:  false,
			MountPath: "/root/.numary",
		}}
	}

	return ret
}

func createLedgerContainerFull(ctx core.Context, stack *v1beta1.Stack, v2 bool) (*corev1.Container, error) {
	container := createBaseLedgerContainer(v2)

	var broker *v1beta1.Broker
	if t, err := brokertopics.Find(ctx, stack, "ledger"); err != nil {
		return nil, err
	} else if t != nil && t.Status.Ready {
		broker = &v1beta1.Broker{}
		if err := ctx.GetClient().Get(ctx, types.NamespacedName{
			Name: stack.Name,
		}, broker); err != nil {
			return nil, err
		}
	}

	if broker != nil {
		if !broker.Status.Ready {
			return nil, core.NewPendingError().WithMessage("broker not ready")
		}
		prefix := ""
		if !v2 {
			prefix = "NUMARY_"
		}

		brokerEnvVar, err := brokers.GetEnvVarsWithPrefix(ctx, broker.Status.URI, stack.Name, "ledger", prefix)
		if err != nil {
			return nil, err
		}

		container.Env = append(container.Env, brokerEnvVar...)
		container.Env = append(container.Env, brokers.GetPublisherEnvVars(stack, broker, "ledger", prefix)...)
	}

	if v2 {
		hasDependency, err := core.HasDependency(ctx, stack.Name, &v1beta1.Analytics{})
		if err != nil {
			return nil, err
		}
		if hasDependency {
			container.Env = append(container.Env, core.Env("EMIT_LOGS", "true"))
		}
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

	caddyfileConfigMap, err := caddy.CreateCaddyfileConfigMap(ctx, stack, "ledger", Caddyfile, map[string]any{
		"Debug": stack.Spec.Debug || ledger.Spec.Debug,
	}, core.WithController[*corev1.ConfigMap](ctx.GetScheme(), ledger))
	if err != nil {
		return err
	}

	env := make([]corev1.EnvVar, 0)
	env = append(env, core.GetDevEnvVars(stack, ledger)...)

	caddyImage, err := registries.GetCaddyImage(ctx, stack, "2.7.6-alpine")
	if err != nil {
		return err
	}

	tpl, err := caddy.DeploymentTemplate(ctx, stack, ledger, caddyfileConfigMap, caddyImage, env)
	if err != nil {
		return err
	}

	tpl.Name = "ledger-gateway"
	return applications.
		New(ledger, tpl).
		Install(ctx)
}

func migrate(ctx core.Context, stack *v1beta1.Stack, ledger *v1beta1.Ledger, database *v1beta1.Database, image, version string) error {
	serviceAccountName, err := settings.GetAWSServiceAccount(ctx, stack.Name)
	if err != nil {
		return err
	}

	upgradeContainer, err := getUpgradeContainer(ctx, stack, database, image, version)
	if err != nil {
		return err
	}

	return jobs.Handle(ctx, ledger, "migrate-v2", upgradeContainer,
		jobs.PreCreate(func() error {
			list := &appsv1.DeploymentList{}
			if err := ctx.GetClient().List(ctx, list, client.InNamespace(stack.Name)); err != nil {
				return err
			}

			for _, item := range list.Items {
				if controller := metav1.GetControllerOf(&item); controller != nil && controller.UID == ledger.GetUID() {
					if err := ctx.GetClient().Delete(ctx, &item); err != nil {
						return err
					}
				}
			}
			return nil
		}),
		jobs.WithServiceAccount(serviceAccountName),
	)
}
