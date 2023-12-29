package internal

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"io/fs"
	"k8s.io/apimachinery/pkg/types"
	"path/filepath"
	"strings"

	"github.com/formancehq/operator/v2/api/v1beta1"
	pkgError "github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetDevEnvVars(stack *v1beta1.Stack, service interface {
	IsDebug() bool
	IsDev() bool
}) []corev1.EnvVar {
	return []corev1.EnvVar{
		EnvFromBool("DEBUG", stack.Spec.Debug || service.IsDebug()),
		EnvFromBool("DEV", stack.Spec.Dev || service.IsDev()),
		Env("STACK", stack.Name),
	}
}

func GetCommonServicesEnvVars(ctx context.Context, _client client.Client, stack *v1beta1.Stack, serviceName string, service interface {
	IsDebug() bool
	IsDev() bool
}) ([]corev1.EnvVar, error) {
	ret := make([]corev1.EnvVar, 0)
	configuration, err := GetOpenTelemetryConfiguration(ctx, _client, stack.Name)
	if err != nil {
		return nil, err
	}
	if configuration != nil {
		ret = append(ret, MonitoringEnvVars(configuration, serviceName)...)
	}

	env, err := GetURLSAsEnvVarsIfGatewayEnabled(ctx, _client, stack.Name)
	if err != nil {
		return nil, err
	}
	ret = append(ret, env...)
	ret = append(ret, GetDevEnvVars(stack, service)...)

	return ret, nil
}

func GetOpenTelemetryConfiguration(ctx context.Context, _client client.Client, stackName string) (*v1beta1.OpenTelemetryConfiguration, error) {
	stackSelectorRequirement, err := labels.NewRequirement("formance.com/stack", selection.In, []string{"any", stackName})
	if err != nil {
		return nil, err
	}

	openTelemetryTracesList := &v1beta1.OpenTelemetryConfigurationList{}
	if err := _client.List(ctx, openTelemetryTracesList, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*stackSelectorRequirement),
	}); err != nil {
		return nil, err
	}

	switch len(openTelemetryTracesList.Items) {
	case 0:
		return nil, nil
	case 1:
		return &openTelemetryTracesList.Items[0], nil
	default:
		return nil, pkgError.New("found multiple opentelemetry config")
	}
}

func IsOpenTelemetryEnabled(ctx context.Context, _client client.Client, stackName string) (bool, error) {
	configuration, err := GetOpenTelemetryConfiguration(ctx, _client, stackName)
	if err == nil {
		return false, err
	}
	if configuration == nil {
		return false, nil
	}
	return true, nil
}

type monitoringType string

const (
	monitoringTypeTraces  monitoringType = "TRACES"
	monitoringTypeMetrics monitoringType = "METRICS"
)

func MonitoringEnvVars(config *v1beta1.OpenTelemetryConfiguration, serviceName string) []corev1.EnvVar {
	ret := make([]corev1.EnvVar, 0)
	if config.Spec.Traces != nil {
		if config.Spec.Traces.Otlp != nil {
			ret = append(ret, MonitoringOTLPEnvVars(config.Spec.Traces.Otlp, monitoringTypeTraces, serviceName)...)
		}
	}
	if config.Spec.Metrics != nil {
		if config.Spec.Metrics.Otlp != nil {
			ret = append(ret, MonitoringOTLPEnvVars(config.Spec.Metrics.Otlp, monitoringTypeMetrics, serviceName)...)
		}
	}
	return nil
}

func MonitoringOTLPEnvVars(otlp *v1beta1.OtlpSpec, monitoringType monitoringType, serviceName string) []corev1.EnvVar {
	return []corev1.EnvVar{
		Env(fmt.Sprintf("OTEL_%s", string(monitoringType)), "true"),
		Env(fmt.Sprintf("OTEL_%s_EXPORTER", string(monitoringType)), "otlp"),
		EnvFromBool(fmt.Sprintf("OTEL_%s_EXPORTER_OTLP_INSECURE", string(monitoringType)), otlp.Insecure),
		Env(fmt.Sprintf("OTEL_%s_EXPORTER_OTLP_MODE", string(monitoringType)), otlp.Mode),
		Env(fmt.Sprintf("OTEL_%s_PORT", string(monitoringType)), fmt.Sprint(otlp.Port)),
		Env(fmt.Sprintf("OTEL_%s_ENDPOINT", string(monitoringType)), otlp.Endpoint),
		Env(fmt.Sprintf("OTEL_%s_EXPORTER_OTLP_ENDPOINT", string(monitoringType)), ComputeEnvVar("%s:%s", fmt.Sprintf("OTEL_%s_ENDPOINT", string(monitoringType)), fmt.Sprintf("OTEL_%s_PORT", string(monitoringType)))),
		Env("OTEL_RESOURCE_ATTRIBUTES", otlp.ResourceAttributes),
		Env("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", ComputeEnvVar("http://%s", fmt.Sprintf("OTEL_TRACES_EXPORTER_OTLP_ENDPOINT"))),
		Env("OTEL_SERVICE_NAME", serviceName),
	}
}

func RequireBrokerConfiguration(ctx context.Context, _client client.Client, stackName string) (*v1beta1.BrokerConfiguration, error) {
	brokerConfiguration, err := GetBrokerConfiguration(ctx, _client, stackName)
	if err != nil {
		return nil, err
	}
	if brokerConfiguration == nil {
		return nil, errors.New("no broker configuration found")
	}
	return brokerConfiguration, nil
}

func GetBrokerConfiguration(ctx context.Context, _client client.Client, stackName string) (*v1beta1.BrokerConfiguration, error) {

	stackSelectorRequirement, err := labels.NewRequirement("formance.com/stack", selection.In, []string{"any", stackName})
	if err != nil {
		return nil, err
	}

	brokerConfigurationList := &v1beta1.BrokerConfigurationList{}
	if err := _client.List(ctx, brokerConfigurationList, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*stackSelectorRequirement),
	}); err != nil {
		return nil, err
	}

	switch len(brokerConfigurationList.Items) {
	case 0:
		return nil, nil
	case 1:
		return &brokerConfigurationList.Items[0], nil
	default:
		return nil, pkgError.New("found multiple broker config")
	}
}

func RequireElasticSearchConfiguration(ctx context.Context, _client client.Client, stackName string) (*v1beta1.ElasticSearchConfiguration, error) {

	elasticSearchConfiguration, err := GetElasticSearchConfiguration(ctx, _client, stackName)
	if err != nil {
		return nil, err
	}
	if elasticSearchConfiguration == nil {
		return nil, ErrNoConfigurationFound
	}

	return elasticSearchConfiguration, nil
}

func GetElasticSearchConfiguration(ctx context.Context, _client client.Client, stackName string) (*v1beta1.ElasticSearchConfiguration, error) {

	stackSelectorRequirement, err := labels.NewRequirement("formance.com/stack", selection.In, []string{"any", stackName})
	if err != nil {
		return nil, err
	}

	elasticSearchConfigurationList := &v1beta1.ElasticSearchConfigurationList{}
	if err := _client.List(ctx, elasticSearchConfigurationList, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*stackSelectorRequirement),
	}); err != nil {
		return nil, err
	}

	switch len(elasticSearchConfigurationList.Items) {
	case 0:
		return nil, nil
	case 1:
		return &elasticSearchConfigurationList.Items[0], nil
	default:
		return nil, pkgError.New("found multiple elasticsearch config")
	}
}

var ErrNoConfigurationFound = pkgError.New("no configuration found")

func GetBrokerEnvVars(ctx context.Context, _client client.Client, stackName, serviceName string) ([]corev1.EnvVar, error) {
	configuration, err := GetBrokerConfiguration(ctx, _client, stackName)
	if err != nil {
		return nil, err
	}
	if configuration == nil {
		return nil, ErrNoConfigurationFound
	}

	return BrokerEnvVars(*configuration, serviceName), nil
}

func BrokerEnvVars(broker v1beta1.BrokerConfiguration, serviceName string) []corev1.EnvVar {
	ret := make([]corev1.EnvVar, 0)

	if broker.Spec.Kafka != nil {
		ret = append(ret,
			Env("BROKER", "kafka"),
			Env("PUBLISHER_KAFKA_ENABLED", "true"),
			Env("PUBLISHER_KAFKA_BROKER", strings.Join(broker.Spec.Kafka.Brokers, " ")),
		)
		if broker.Spec.Kafka.SASL != nil {
			ret = append(ret,
				Env("PUBLISHER_KAFKA_SASL_ENABLED", "true"),
				Env("PUBLISHER_KAFKA_SASL_USERNAME", broker.Spec.Kafka.SASL.Username),
				Env("PUBLISHER_KAFKA_SASL_PASSWORD", broker.Spec.Kafka.SASL.Password),
				Env("PUBLISHER_KAFKA_SASL_MECHANISM", broker.Spec.Kafka.SASL.Mechanism),
				Env("PUBLISHER_KAFKA_SASL_SCRAM_SHA_SIZE", broker.Spec.Kafka.SASL.ScramSHASize),
			)
		}
		if broker.Spec.Kafka.TLS {
			ret = append(ret,
				Env("PUBLISHER_KAFKA_TLS_ENABLED", "true"),
			)
		}
	} else {
		ret = append(ret,
			Env("BROKER", "nats"),
			Env("PUBLISHER_NATS_ENABLED", "true"),
			Env("PUBLISHER_NATS_URL", broker.Spec.Nats.URL),
			Env("PUBLISHER_NATS_CLIENT_ID", serviceName),
		)
	}
	return ret
}

func HashFromConfigMap(configMap *corev1.ConfigMap) string {
	digest := sha256.New()
	if err := json.NewEncoder(digest).Encode(configMap.Data); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(digest.Sum(nil))
}

func ConfigureCaddy(caddyfile *corev1.ConfigMap, image string, env []corev1.EnvVar,
	resourceProperties *v1beta1.ResourceProperties) []ObjectMutator[*appsv1.Deployment] {
	return []ObjectMutator[*appsv1.Deployment]{
		func(t *appsv1.Deployment) {
			t.Spec.Template.Annotations = MergeMaps(t.Spec.Template.Annotations, map[string]string{
				"caddyfile-hash": HashFromConfigMap(caddyfile),
			})
			t.Spec.Template.Spec.Volumes = []corev1.Volume{
				volumeFromConfigMap("caddyfile", caddyfile),
			}
			t.Spec.Template.Spec.Containers = []corev1.Container{
				{
					Name:    "gateway",
					Command: []string{"/usr/bin/caddy"},
					Args: []string{
						"run",
						"--config", "/gateway/Caddyfile",
						"--adapter", "caddyfile",
					},
					Image:     image,
					Env:       env,
					Resources: GetResourcesWithDefault(resourceProperties, ResourceSizeSmall()),
					VolumeMounts: []corev1.VolumeMount{
						volumeMount("caddyfile", "/gateway"),
					},
					Ports: []corev1.ContainerPort{{
						Name:          "http",
						ContainerPort: 8080,
					}},
				},
			}
		},
	}
}

func volumeFromConfigMap(name string, cm *corev1.ConfigMap) corev1.Volume {
	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cm.Name,
				},
			},
		},
	}
}

func volumeMount(name, mountPath string) corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      name,
		ReadOnly:  true,
		MountPath: mountPath,
	}
}

func ConfigureK8SService(name string, options ...func(service *corev1.Service)) func(service *corev1.Service) {
	return func(t *corev1.Service) {
		t.Labels = map[string]string{
			"app.kubernetes.io/service-name": name,
		}
		t.Spec = corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:       "http",
				Port:       8080,
				Protocol:   "TCP",
				TargetPort: intstr.FromInt32(8080),
			}},
			Selector: map[string]string{
				"app.kubernetes.io/name": name,
			},
		}
		for _, option := range options {
			option(t)
		}
	}
}

func IsAuthEnabled(ctx context.Context, _client client.Client, stackName string) (bool, error) {
	list := &v1beta1.AuthList{}
	if err := _client.List(ctx, list, client.MatchingFields{
		".spec.stack": stackName,
	}); err != nil {
		return false, err
	}

	return len(list.Items) > 0, nil
}

type MigrationConfiguration struct {
	Command       []string
	AdditionalEnv []corev1.EnvVar
}

func MigrateDatabaseContainer(serviceName string, database v1beta1.DatabaseConfigurationSpec, databaseName, version string, options ...func(m *MigrationConfiguration)) corev1.Container {
	m := &MigrationConfiguration{}
	for _, option := range options {
		option(m)
	}
	args := m.Command
	if len(args) == 0 {
		args = []string{"migrate"}
	}
	env := PostgresEnvVars(database, databaseName)
	if m.AdditionalEnv != nil {
		env = append(env, m.AdditionalEnv...)
	}
	return corev1.Container{
		Name:  "migrate",
		Image: GetImage(serviceName, version),
		Args:  args,
		Env:   env,
	}
}

func StandardHTTPPort() corev1.ContainerPort {
	return corev1.ContainerPort{
		Name:          "http",
		ContainerPort: 8080,
	}
}

func GetVersion(stack *v1beta1.Stack, defaultVersion string) string {
	if defaultVersion == "" {
		return stack.GetVersion()
	}
	return defaultVersion
}

func GetGatewayIfEnabled(ctx context.Context, _client client.Client, stackName string) (*v1beta1.Gateway, error) {
	gatewayList := &v1beta1.GatewayList{}
	if err := _client.List(ctx, gatewayList, client.MatchingFields{
		".spec.stack": stackName,
	}); err != nil {
		return nil, err
	}

	switch len(gatewayList.Items) {
	case 0:
		return nil, nil
	case 1:
		return &gatewayList.Items[0], nil
	default:
		return nil, pkgError.New("found multiple gateway")
	}
}

func GetURLSAsEnvVarsIfGatewayEnabled(ctx context.Context, _client client.Client, stackName string) ([]corev1.EnvVar, error) {
	gateway, err := GetGatewayIfEnabled(ctx, _client, stackName)
	if err != nil {
		return nil, err
	}
	if gateway == nil {
		return nil, nil
	}
	ret := []corev1.EnvVar{{
		Name:  "STACK_URL",
		Value: "http://gateway:8080",
	}}
	if gateway.Spec.Ingress != nil {
		ret = append(ret, corev1.EnvVar{
			Name:  "STACK_PUBLIC_URL",
			Value: fmt.Sprintf("%s://%s", gateway.Spec.Ingress.Scheme, gateway.Spec.Ingress.Host),
		})
	}

	return ret, nil
}

func GetAuthClientEnvVars(authClient *v1beta1.AuthClient) []corev1.EnvVar {
	return []corev1.EnvVar{
		EnvFromSecret("STACK_CLIENT_ID", fmt.Sprintf("auth-client-%s", authClient.Name), "id"),
		EnvFromSecret("STACK_CLIENT_SECRET", fmt.Sprintf("auth-client-%s", authClient.Name), "secret"),
	}
}

func GetStack(ctx context.Context, client client.Client, spec interface {
	GetStack() string
}) (*v1beta1.Stack, error) {
	stack := &v1beta1.Stack{}
	if err := client.Get(ctx, types.NamespacedName{
		Name: spec.GetStack(),
	}, stack); err != nil {
		return nil, err
	}

	return stack, nil
}

var ErrMultipleInstancesFound = errors.New("multiple resources found")

//func GetSingleStackDependencyObject[T client.Object](ctx context.Context, scheme *runtime.Scheme, _client client.Client, stackName string) (T, error) {
//
//	fmt.Println("find object")
//	fmt.Println("find object")
//	fmt.Println("find object")
//	fmt.Println("find object")
//	fmt.Println("find object")
//
//	var t T
//	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
//	kinds, _, err := scheme.ObjectKinds(t)
//	if err != nil {
//		return t, err
//	}
//	spew.Dump(kinds)
//
//	list := &unstructured.UnstructuredList{}
//	list.SetGroupVersionKind(kinds[0])
//	err = _client.List(ctx, list, client.MatchingFields{
//		".spec.stack": stackName,
//	})
//	if err != nil {
//		return t, err
//	}
//
//	switch len(list.Items) {
//	case 0:
//		return t, nil
//	case 1:
//		if err := runtime.DefaultUnstructuredConverter.
//			FromUnstructured(list.Items[0].UnstructuredContent(), t); err != nil {
//			return t, err
//		}
//		return t, nil
//	default:
//		return t, ErrMultipleInstancesFound
//	}
//}

func GetAuthIfEnabled(ctx context.Context, _client client.Client, stackName string) (*v1beta1.Auth, error) {
	authList := &v1beta1.AuthList{}
	if err := _client.List(ctx, authList, client.MatchingFields{
		".spec.stack": stackName,
	}); err != nil {
		return nil, err
	}

	switch len(authList.Items) {
	case 0:
		return nil, nil
	case 1:
		return &authList.Items[0], nil
	default:
		return nil, pkgError.New("found multiple auth")
	}
}

func GetLedgerIfEnabled(ctx context.Context, _client client.Client, stackName string) (*v1beta1.Ledger, error) {
	LedgerList := &v1beta1.LedgerList{}
	if err := _client.List(ctx, LedgerList, client.MatchingFields{
		".spec.stack": stackName,
	}); err != nil {
		return nil, err
	}

	switch len(LedgerList.Items) {
	case 0:
		return nil, nil
	case 1:
		return &LedgerList.Items[0], nil
	default:
		return nil, pkgError.New("found multiple Ledger")
	}
}

func GetPaymentsIfEnabled(ctx context.Context, _client client.Client, stackName string) (*v1beta1.Payments, error) {
	PaymentsList := &v1beta1.PaymentsList{}
	if err := _client.List(ctx, PaymentsList, client.MatchingFields{
		".spec.stack": stackName,
	}); err != nil {
		return nil, err
	}

	switch len(PaymentsList.Items) {
	case 0:
		return nil, nil
	case 1:
		return &PaymentsList.Items[0], nil
	default:
		return nil, pkgError.New("found multiple Payments")
	}
}

func GetWalletsIfEnabled(ctx context.Context, _client client.Client, stackName string) (*v1beta1.Wallets, error) {

	WalletsList := &v1beta1.WalletsList{}
	if err := _client.List(ctx, WalletsList, client.MatchingFields{
		".spec.stack": stackName,
	}); err != nil {
		return nil, err
	}

	switch len(WalletsList.Items) {
	case 0:
		return nil, nil
	case 1:
		return &WalletsList.Items[0], nil
	default:
		return nil, pkgError.New("found multiple Wallets")
	}
}

func CreateTopicQuery(ctx context.Context, _client client.Client, stack *v1beta1.Stack, service, queriedBy string) error {
	_, _, err := CreateOrUpdate[*v1beta1.TopicQuery](ctx, _client, types.NamespacedName{
		Name: GetObjectName(stack.Name, fmt.Sprintf("%s-%s", queriedBy, service)),
	}, func(t *v1beta1.TopicQuery) {
		t.Spec.QueriedBy = queriedBy
		t.Spec.Stack = stack.Name
		t.Spec.Service = service
	})
	return err
}

func CopyDir(f fs.FS, root, path string, ret *map[string]string) {
	dirEntries, err := fs.ReadDir(f, path)
	if err != nil {
		panic(err)
	}
	for _, dirEntry := range dirEntries {
		dirEntryPath := filepath.Join(path, dirEntry.Name())
		if dirEntry.IsDir() {
			CopyDir(f, root, dirEntryPath, ret)
		} else {
			fileContent, err := fs.ReadFile(f, dirEntryPath)
			if err != nil {
				panic(err)
			}
			sanitizedPath := strings.TrimPrefix(dirEntryPath, root)
			sanitizedPath = strings.TrimPrefix(sanitizedPath, "/")
			(*ret)[sanitizedPath] = string(fileContent)
		}
	}
}

func LoadStreams(ctx context.Context, client client.Client, fs embed.FS,
	stackName string, streamDirectory string) error {
	streamFiles, err := fs.ReadDir(streamDirectory)
	if err != nil {
		return err
	}

	// TODO: Only if search enabled
	for _, file := range streamFiles {
		streamContent, err := fs.ReadFile(streamDirectory + "/" + file.Name())
		if err != nil {
			return err
		}

		sanitizedName := strings.ReplaceAll(file.Name(), "_", "-")

		_, _, err = CreateOrUpdate[*v1beta1.Stream](ctx, client, types.NamespacedName{
			Name: fmt.Sprintf("%s-%s", stackName, sanitizedName),
		}, func(t *v1beta1.Stream) {
			t.Spec.Data = string(streamContent)
			t.Spec.Stack = stackName
		})
		if err != nil {
			return pkgError.Wrap(err, "creating stream")
		}
	}

	return nil
}
