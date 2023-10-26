package modules

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/common"
	"github.com/formancehq/operator/internal/controllerutils"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/nats-io/nats.go"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

type Config struct {
	Data  map[string]string
	Mount bool
}

func (c Config) create(ctx context.Context, deployer Deployer, serviceName, configName string) (*ConfigHandle, error) {

	hash := sha256.New()
	dataKeys := collectionutils.Keys(c.Data)
	sort.Strings(dataKeys)
	for _, k := range dataKeys {
		hash.Write([]byte(k))
		hash.Write([]byte(c.Data[k]))
	}

	configMap, _, err := deployer.
		ConfigMaps().
		CreateOrUpdate(ctx, serviceName+"-"+configName, func(t *corev1.ConfigMap) {
			t.Data = c.Data
		})
	if err != nil {
		return nil, err
	}
	var mountPath string
	if c.Mount {
		mountPath = fmt.Sprintf("/config/%s", configName)
	}
	h := NewConfigHandle(configMap.Name, mountPath, base64.URLEncoding.EncodeToString(hash.Sum(nil)))
	return &h, nil
}

type Configs map[string]Config

func (c Configs) create(ctx context.Context, deployer Deployer, serviceName string) (ConfigHandles, error) {
	configHandles := ConfigHandles{}
	for configName, configDefinition := range c {
		configHandle, err := configDefinition.create(ctx, deployer, serviceName, configName)
		if err != nil {
			return nil, err
		}
		configHandles[configName] = *configHandle
	}
	return configHandles, nil
}

type Secret struct {
	Data  map[string][]byte
	Mount bool
}

func (s Secret) create(ctx context.Context, deployer *scopedResourceDeployer, serviceName, secretName string) (*SecretHandle, error) {

	hash := sha256.New()
	dataKeys := collectionutils.Keys(s.Data)
	sort.Strings(dataKeys)
	for _, k := range dataKeys {
		hash.Write([]byte(k))
		hash.Write(s.Data[k])
	}
	secret, _, err := deployer.
		Secrets().
		CreateOrUpdate(ctx, serviceName+"-"+secretName, func(t *corev1.Secret) {
			// Only create secret if it does not exist
			if t.Data != nil {
				for k := range s.Data {
					if _, ok := t.Data[k]; !ok {
						goto apply
					}
				}
				return
			}
		apply:
			t.Data = s.Data
		})
	if err != nil {
		return nil, err
	}
	var mountPath string
	if s.Mount {
		mountPath = fmt.Sprintf("/secret/%s", secretName)
	}
	h := NewSecretHandle(secret.Name, mountPath, base64.URLEncoding.EncodeToString(hash.Sum(nil)))
	return &h, nil
}

type Secrets map[string]Secret

func (s Secrets) create(ctx context.Context, deployer *scopedResourceDeployer, serviceName string) (SecretHandles, error) {
	secretHandles := SecretHandles{}
	for secretName, secretDefinition := range s {
		secretHandle, err := secretDefinition.create(ctx, deployer, serviceName, secretName)
		if err != nil {
			return nil, err
		}
		secretHandles[secretName] = *secretHandle
	}
	return secretHandles, nil
}

type EnvVar struct {
	Name       string
	FromString *string
	FromConfig *string
	FromSecret *string
	Key        string
}

func Env(name string, value string) EnvVar {
	return EnvVar{
		Name:       name,
		FromString: &value,
	}
}

func EnvFromBool(name string, vb bool) EnvVar {
	value := "true"
	if !vb {
		value = "false"
	}
	return EnvVar{
		Name:       name,
		FromString: &value,
	}
}

func EnvFromConfig(name, configName, key string) EnvVar {
	return EnvVar{
		Name:       name,
		FromConfig: &configName,
		Key:        key,
	}
}

func EnvFromSecret(name, secretName, key string) EnvVar {
	return EnvVar{
		Name:       name,
		FromSecret: &secretName,
		Key:        key,
	}
}

type Liveness int

const (
	LivenessDefault = iota
	LivenessLegacy
	LivenessDisable
)

type ContainerEnv []EnvVar

func (env ContainerEnv) Append(v ...EnvVar) ContainerEnv {
	return append(env, v...)
}

func NewEnv() ContainerEnv {
	return ContainerEnv{}
}

func (e ContainerEnv) ToCoreEnv() []corev1.EnvVar {
	ret := make([]corev1.EnvVar, 0)
	for _, envVar := range e {
		switch {
		case envVar.FromSecret != nil:
			ret = append(ret, corev1.EnvVar{
				Name: envVar.Name,
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						Key: envVar.Key,
						LocalObjectReference: corev1.LocalObjectReference{
							Name: *envVar.FromSecret,
						},
					},
				},
			})
		case envVar.FromConfig != nil:
			ret = append(ret, corev1.EnvVar{
				Name: envVar.Name,
				ValueFrom: &corev1.EnvVarSource{
					ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
						Key: envVar.Key,
						LocalObjectReference: corev1.LocalObjectReference{
							Name: *envVar.FromSecret,
						},
					},
				},
			})
		case envVar.FromString != nil:
			ret = append(ret, corev1.EnvVar{
				Name:  envVar.Name,
				Value: *envVar.FromString,
			})
		}
	}
	return ret
}

type SecretHandle struct {
	MountPath string
	Name      string
	Hash      string
}

func NewSecretHandle(name, mountPath, hash string) SecretHandle {
	return SecretHandle{
		MountPath: mountPath,
		Name:      name,
		Hash:      hash,
	}
}

func (h SecretHandle) GetName() string {
	return h.Name
}

func (h SecretHandle) GetMountPath() string {
	return h.MountPath
}

type ConfigHandle struct {
	MountPath string
	Name      string
	Hash      string
}

func NewConfigHandle(name, mountPath, hash string) ConfigHandle {
	return ConfigHandle{
		MountPath: mountPath,
		Name:      name,
		Hash:      hash,
	}
}

func (h ConfigHandle) GetName() string {
	return h.Name
}

func (h ConfigHandle) GetMountPath() string {
	if h.MountPath == "" {
		panic("config not defined as mountable")
	}
	return h.MountPath
}

type ConfigHandles map[string]ConfigHandle

func (h ConfigHandles) sort() []string {
	ret := make([]string, 0)
	for key := range h {
		ret = append(ret, key)
	}
	sort.Strings(ret)
	return ret
}

type SecretHandles map[string]SecretHandle

func (h SecretHandles) sort() []string {
	ret := make([]string, 0)
	for key := range h {
		ret = append(ret, key)
	}
	sort.Strings(ret)
	return ret
}

type ExposeHTTP struct {
	// Path indicates the path used to expose the service using an ingress
	Path    string
	Methods []string
	Name    string
}

type Path struct {
	Path    string
	Methods []string
	Name    string
}

var DefaultExposeHTTP = &ExposeHTTP{}

type Service struct {
	Name string
	// Secured indicate if the service is able to handle security
	Secured bool
	// ExposeHTTP indicate the service expose a http endpoint
	ExposeHTTP *ExposeHTTP
	// Paths indicates the paths exposed by the service
	Paths []Path
	// ListenEnvVar indicate the flag used to configure the http service address
	// TODO(gfyrag): Remove this in a future version when all services implements --listen
	ListenEnvVar string
	// Port indicate the listening port of the service.
	// Deprecated
	// All services should have the --listen flag to allow the operator to specify the port
	Port int32
	// Annotations indicates the annotations apply to the Service
	Annotations             map[string]string
	InjectPostgresVariables bool
	HasVersionEndpoint      bool
	Liveness                Liveness
	AuthConfiguration       func(config ReconciliationConfig) stackv1beta3.ClientConfiguration
	Configs                 func(resolveContext ServiceInstallConfiguration) Configs
	Secrets                 func(resolveContext ServiceInstallConfiguration) Secrets
	Container               func(resolveContext ContainerResolutionConfiguration) Container
	InitContainer           func(resolveContext ContainerResolutionConfiguration) []Container
	NeedTopic               bool
	Replicas                *int32

	EnvPrefix string
}

type Services []*Service

func (services Services) Len() int {
	return len(services)
}

func (services Services) Less(i, j int) bool {
	return strings.Compare(services[i].Name, services[j].Name) < 0
}

func (services Services) Swap(i, j int) {
	services[i], services[j] = services[j], services[i]
}

type serviceReconciler struct {
	*moduleReconciler
	name     string
	service  Service
	usedPort int32
}

func (r *serviceReconciler) prepare() {
	if r.service.AuthConfiguration != nil {
		_ = r.Stack.GetOrCreateClient(r.name, r.service.AuthConfiguration(r.ReconciliationConfig))
	}
	if r.service.ExposeHTTP != nil || r.service.Liveness != LivenessDisable {
		r.allocatePort()
	}
}

func (r *serviceReconciler) reconcile(ctx context.Context, config ServiceInstallConfiguration) error {

	// TODO: Use a job
	if r.Configuration.Spec.Broker.Nats != nil && r.service.NeedTopic {
		r.configureNats()
	}

	configHandles, err := r.installConfigs(ctx, config)
	if err != nil {
		return err
	}

	secretHandles, err := r.installSecrets(ctx, config)
	if err != nil {
		return err
	}

	if r.service.ExposeHTTP != nil {
		if err := r.install(ctx); err != nil {
			return err
		}
		if r.service.ExposeHTTP.Path != "" {
			if err := r.createIngress(ctx); err != nil {
				return err
			}
		}
	}

	err = r.createDeployment(ctx, ContainerResolutionConfiguration{
		ServiceInstallConfiguration: config,
		Configs:                     configHandles,
		Secrets:                     secretHandles,
	})
	if err != nil {
		return err
	}

	return nil
}

// TODO: Properly handle errors
func (r *serviceReconciler) configureNats() {
	topicName := r.Stack.GetServiceNamespacedName(r.name).Name
	streamConfig := nats.StreamConfig{
		Name:      topicName,
		Subjects:  []string{topicName},
		Retention: nats.InterestPolicy,
	}
	nc, err := nats.Connect(r.Configuration.Spec.Broker.Nats.URL)
	if err != nil {
		logging.Error(err)
	}
	js, err := nc.JetStream()
	if err != nil {
		logging.Error(err)
	}
	_, err = js.StreamInfo(topicName)
	if err != nil {
		_, err := js.AddStream(&streamConfig)
		if err != nil {
			logging.Error(err)
		}
	} else {
		_, err = js.UpdateStream(&streamConfig)
		if err != nil {
			logging.Error(fmt.Sprintf("%s: %s", topicName, err))
		}
	}
}

func (r *serviceReconciler) allocatePort() {
	r.usedPort = r.service.Port
	if r.usedPort == 0 {
		r.usedPort = r.portAllocator.NextPort()
	}

	if r.Stack.Status.Ports == nil {
		r.Stack.Status.Ports = make(map[string]map[string]int32)
	}

	if r.Stack.Status.Ports[r.module.Name()] == nil {
		r.Stack.Status.Ports[r.module.Name()] = make(map[string]int32)
	}

	r.Stack.Status.Ports[r.module.Name()][r.name] = r.usedPort
}

func (r *serviceReconciler) installConfigs(ctx context.Context, config ServiceInstallConfiguration) (ConfigHandles, error) {
	if r.service.Configs == nil {
		return ConfigHandles{}, nil
	}
	return r.service.Configs(config).create(ctx, r.namespacedResourceDeployer, r.name)
}

func (r *serviceReconciler) installSecrets(ctx context.Context, config ServiceInstallConfiguration) (SecretHandles, error) {
	if r.service.Secrets == nil {
		return SecretHandles{}, nil
	}
	return r.service.Secrets(config).create(ctx, r.namespacedResourceDeployer, r.name)
}

func (r *serviceReconciler) createDeployment(ctx context.Context, config ContainerResolutionConfiguration) error {
	container := r.service.Container(config)
	volumes, volumesHash := config.volumes(r.name)
	return r.podDeployer.deploy(ctx, pod{
		name:                 r.name,
		moduleName:           r.module.Name(),
		volumes:              volumes,
		initContainers:       r.initContainers(config, r.name),
		containers:           r.containers(config, container, r.name),
		disableRollingUpdate: container.DisableRollingUpdate,
		replicas:             r.service.Replicas,
		annotations: map[string]string{
			"stack.formance.cloud/volumes-hash": volumesHash,
		},
	})
}

func (r *serviceReconciler) initContainers(ctx ContainerResolutionConfiguration, serviceName string) []corev1.Container {
	ret := make([]corev1.Container, 0)
	if r.service.InitContainer != nil {
		for _, c := range r.service.InitContainer(ctx) {
			ret = append(ret, r.createContainer(ctx, c, "init-"+serviceName, true))
		}
	}
	return ret
}

func (r *serviceReconciler) containers(ctx ContainerResolutionConfiguration, container Container, serviceName string) []corev1.Container {
	return []corev1.Container{
		r.createContainer(ctx, container, serviceName, false),
	}
}

func (r *serviceReconciler) createContainer(ctx ContainerResolutionConfiguration, container Container, serviceName string, init bool) corev1.Container {
	c := corev1.Container{
		Name: func() string {
			if container.Name != "" {
				return container.Name
			}
			return serviceName
		}(),
		Image:           container.Image,
		ImagePullPolicy: GetPullPolicy(container.Image),
		Command:         container.Command,
		Args:            container.Args,
		Resources:       container.Resources,
	}
	env := NewEnv()
	if r.service.InjectPostgresVariables {
		env = env.Append(
			DefaultPostgresEnvVarsWithPrefix(*ctx.PostgresConfig, r.Stack.GetServiceName(r.module.Name()), r.service.EnvPrefix)...,
		)
	}
	if r.service.ListenEnvVar != "" {
		env = env.Append(
			Env(fmt.Sprintf("%s%s", r.service.EnvPrefix, r.service.ListenEnvVar), fmt.Sprintf(":%d", r.usedPort)),
		)
	}

	if r.Configuration.Spec.Monitoring != nil {
		if r.Configuration.Spec.Monitoring.Traces != nil {
			env = env.Append(
				MonitoringTracesEnvVars(r.Configuration.Spec.Monitoring.Traces, r.service.EnvPrefix)...,
			)
		}
		if r.Configuration.Spec.Monitoring.Metrics != nil {
			env = env.Append(
				MonitoringMetricsEnvVars(r.Configuration.Spec.Monitoring.Metrics, r.service.EnvPrefix)...,
			)
		}
	}

	if !init {
		env = env.Append(
			Env(fmt.Sprintf("%sDEBUG", r.service.EnvPrefix), fmt.Sprintf("%v", r.Stack.Spec.Debug)),
			Env(fmt.Sprintf("%sDEV", r.service.EnvPrefix), fmt.Sprintf("%v", r.Stack.Spec.Dev)),
			// TODO: the stack url is a full url, we can target the gateway. Need to find how to generalize this
			// as the gateway is a component like another
			Env(fmt.Sprintf("%sSTACK_URL", r.service.EnvPrefix), r.Stack.URL()),
			Env(fmt.Sprintf("%sOTEL_SERVICE_NAME", r.service.EnvPrefix), serviceName),
			Env("STACK", r.Stack.Name),
		)
	}

	for _, envVar := range container.Env {
		envVar.Name = fmt.Sprintf("%s%s", r.service.EnvPrefix, envVar.Name)
		env = append(env, envVar)
	}

	c.Env = env.ToCoreEnv()

	if !init {
		ret := make([]corev1.VolumeMount, 0)
		for _, configName := range ctx.Configs.sort() {
			config := ctx.Configs[configName]
			if config.MountPath == "" {
				continue
			}
			ret = append(ret, corev1.VolumeMount{
				Name:      configName,
				ReadOnly:  true,
				MountPath: fmt.Sprintf("/config/%s", configName),
			})
		}
		for _, secretName := range ctx.Secrets.sort() {
			secret := ctx.Secrets[secretName]
			if secret.MountPath == "" {
				continue
			}
			ret = append(ret, corev1.VolumeMount{
				Name:      secretName,
				ReadOnly:  true,
				MountPath: fmt.Sprintf("/secret/%s", secretName),
			})
		}
		c.VolumeMounts = ret

		switch r.service.Liveness {
		case LivenessDefault:
			c.LivenessProbe = common.DefaultLiveness(r.GetUsedPort())
		case LivenessLegacy:
			c.LivenessProbe = common.LegacyLiveness(r.GetUsedPort())
		}
		if r.usedPort != 0 {
			c.Ports = common.SinglePort("http", r.usedPort)
		}
	}
	return c
}

func (r *serviceReconciler) install(ctx context.Context) error {
	return controllerutils.JustError(r.namespacedResourceDeployer.Services().CreateOrUpdate(ctx, r.name, func(t *corev1.Service) {
		annotations := r.service.Annotations
		if annotations == nil {
			annotations = map[string]string{}
		} else {
			annotations = collectionutils.CopyMap(annotations)
		}
		t.ObjectMeta.Annotations = annotations

		selector := r.name
		if r.Configuration.Spec.LightMode {
			selector = r.Stack.Name
		}
		t.Labels = map[string]string{
			"app.kubernetes.io/service-name": r.name,
		}
		t.Spec = corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:        "http",
				Port:        r.usedPort,
				Protocol:    "TCP",
				AppProtocol: pointer.String("http"),
				TargetPort:  intstr.FromInt(int(r.usedPort)),
			}},
			Selector: map[string]string{
				"app.kubernetes.io/name": selector,
			},
		}
	}))
}

func (r *serviceReconciler) createIngress(ctx context.Context) error {
	return controllerutils.JustError(r.namespacedResourceDeployer.Ingresses().CreateOrUpdate(ctx, r.name, func(ingress *networkingv1.Ingress) {
		annotations := r.Configuration.Spec.Ingress.Annotations
		if annotations == nil {
			annotations = map[string]string{}
		} else {
			annotations = collectionutils.CopyMap(annotations)
		}

		pathType := networkingv1.PathTypePrefix
		ingress.ObjectMeta.Annotations = annotations
		ingress.Spec = networkingv1.IngressSpec{
			TLS: func() []networkingv1.IngressTLS {
				if r.Configuration.Spec.Ingress.TLS == nil {
					return nil
				}
				return []networkingv1.IngressTLS{{
					SecretName: r.Configuration.Spec.Ingress.TLS.SecretName,
				}}
			}(),
			Rules: []networkingv1.IngressRule{
				{
					Host: r.Stack.Spec.Host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     r.service.ExposeHTTP.Path,
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: r.name,
											Port: networkingv1.ServiceBackendPort{
												Name: "http",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
	}))
}

func (r *serviceReconciler) GetUsedPort() int32 {
	return r.usedPort
}

func newServiceReconciler(moduleReconciler *moduleReconciler, service Service, name string) *serviceReconciler {
	return &serviceReconciler{
		moduleReconciler: moduleReconciler,
		name:             name,
		service:          service,
	}
}
