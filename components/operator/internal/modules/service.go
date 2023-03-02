package modules

import (
	"context"
	"fmt"
	"sort"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/collectionutils"
	"github.com/formancehq/operator/internal/common"
	"github.com/formancehq/operator/internal/controllerutils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

type Config struct {
	Data  map[string]string
	Mount bool
}

func (c Config) create(ctx InstallContext, deployer Deployer, serviceName, configName string) (*ConfigHandle, error) {
	configMap, err := deployer.
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
	h := NewConfigHandle(configMap.Name, mountPath)
	return &h, nil
}

type Configs map[string]Config

func (c Configs) create(ctx InstallContext, deployer Deployer, serviceName string) (ConfigHandles, error) {
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

func (s Secret) create(ctx InstallContext, deployer *ComponentDeployer, serviceName, secretName string) (*SecretHandle, error) {
	secret, err := deployer.
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
	h := NewSecretHandle(secret.Name, mountPath)
	return &h, nil
}

type Secrets map[string]Secret

func (s Secrets) create(ctx InstallContext, deployer *ComponentDeployer, serviceName string) (SecretHandles, error) {
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

type Container struct {
	Command              []string
	Args                 []string
	Env                  ContainerEnv
	Image                string
	Name                 string
	Liveness             Liveness
	DisableRollingUpdate bool
}

type SecretHandle struct {
	MountPath string
	Name      string
}

func NewSecretHandle(name, mountPath string) SecretHandle {
	return SecretHandle{
		MountPath: mountPath,
		Name:      name,
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
}

func NewConfigHandle(name, mountPath string) ConfigHandle {
	return ConfigHandle{
		MountPath: mountPath,
		Name:      name,
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

type Context struct {
	context.Context
	Region        string
	Stack         *stackv1beta3.Stack
	Configuration *stackv1beta3.Configuration
	Versions      *stackv1beta3.Versions
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

type Service struct {
	Name                    string
	Secured                 bool
	Port                    int32
	Path                    string
	InjectPostgresVariables bool
	HasVersionEndpoint      bool
	AuthConfiguration       func(resolveContext PrepareContext) stackv1beta3.ClientConfiguration
	Configs                 func(resolveContext InstallContext) Configs
	Secrets                 func(resolveContext InstallContext) Secrets
	Container               func(resolveContext ContainerResolutionContext) Container
	InitContainer           func(resolveContext ContainerResolutionContext) []Container
}

func (service Service) Prepare(ctx PrepareContext, serviceName string) {
	if service.AuthConfiguration != nil {
		_ = ctx.Stack.GetOrCreateClient(serviceName, service.AuthConfiguration(ctx))
	}
}

func (service Service) installService(ctx InstallContext, deployer Deployer, serviceName string) error {
	return controllerutils.JustError(deployer.Services().CreateOrUpdate(ctx, serviceName, func(t *corev1.Service) {
		t.Spec = corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:        "http",
				Port:        service.Port,
				Protocol:    "TCP",
				AppProtocol: pointer.String("http"),
				TargetPort:  intstr.FromInt(int(service.Port)),
			}},
			Selector: collectionutils.CreateMap("app.kubernetes.io/name", serviceName),
		}
	}))
}

func (service Service) createIngress(ctx InstallContext, deployer *ComponentDeployer, serviceName string) error {
	return controllerutils.JustError(deployer.Ingresses().CreateOrUpdate(ctx, serviceName, func(ingress *networkingv1.Ingress) {
		annotations := ctx.Configuration.Spec.Ingress.Annotations
		if annotations == nil {
			annotations = map[string]string{}
		} else {
			annotations = collectionutils.CopyMap(annotations)
		}

		pathType := networkingv1.PathTypePrefix
		ingress.ObjectMeta.Annotations = annotations
		ingress.Spec = networkingv1.IngressSpec{
			TLS: func() []networkingv1.IngressTLS {
				if ctx.Configuration.Spec.Ingress.TLS == nil {
					return nil
				}
				return []networkingv1.IngressTLS{{
					SecretName: ctx.Configuration.Spec.Ingress.TLS.SecretName,
				}}
			}(),
			Rules: []networkingv1.IngressRule{
				{
					Host: ctx.Stack.Spec.Host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     service.Path,
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: serviceName,
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

func (service Service) installConfigs(ctx InstallContext, deployer Deployer, serviceName string) (ConfigHandles, error) {
	if service.Configs == nil {
		return ConfigHandles{}, nil
	}
	return service.Configs(ctx).create(ctx, deployer, serviceName)
}

func (service Service) installSecrets(ctx InstallContext, deployer *ComponentDeployer, serviceName string) (SecretHandles, error) {
	if service.Secrets == nil {
		return SecretHandles{}, nil
	}
	return service.Secrets(ctx).create(ctx, deployer, serviceName)
}

func (service Service) createDeployment(ctx ContainerResolutionContext, deployer *ComponentDeployer, serviceName string) error {
	container := service.Container(ctx)
	return controllerutils.JustError(deployer.
		Deployments().
		CreateOrUpdate(ctx, serviceName, func(t *appsv1.Deployment) {
			matchLabels := collectionutils.CreateMap("app.kubernetes.io/name", serviceName)
			strategy := appsv1.DeploymentStrategy{}
			if container.DisableRollingUpdate {
				strategy.Type = appsv1.RecreateDeploymentStrategyType
			}
			t.Spec = appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: matchLabels,
				},
				Strategy: strategy,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: matchLabels,
					},
					Spec: corev1.PodSpec{
						Volumes:        ctx.volumes(serviceName),
						InitContainers: service.initContainers(ctx, serviceName),
						Containers:     service.containers(ctx, container, serviceName),
					},
				},
			}
		}),
	)
}

func (service Service) initContainers(ctx ContainerResolutionContext, serviceName string) []corev1.Container {
	ret := make([]corev1.Container, 0)
	if service.InitContainer != nil {
		for _, c := range service.InitContainer(ctx) {
			ret = append(ret, service.createContainer(ctx, c, "init-"+serviceName, true))
		}
	}
	return ret
}

func (service Service) containers(ctx ContainerResolutionContext, container Container, serviceName string) []corev1.Container {
	return []corev1.Container{
		service.createContainer(ctx, container, serviceName, false),
	}
}

func (service Service) Install(ctx InstallContext, deployer *ComponentDeployer, serviceName string) error {
	configHandles, err := service.installConfigs(ctx, deployer, serviceName)
	if err != nil {
		return err
	}

	secretHandles, err := service.installSecrets(ctx, deployer, serviceName)
	if err != nil {
		return err
	}

	if service.Port != 0 {
		if err := service.installService(ctx, deployer, serviceName); err != nil {
			return err
		}
		if service.Path != "" {
			if err := service.createIngress(ctx, deployer, serviceName); err != nil {
				return err
			}
		}
	}

	err = service.createDeployment(ContainerResolutionContext{
		InstallContext: ctx,
		Configs:        configHandles,
		Secrets:        secretHandles,
	}, deployer, serviceName)
	if err != nil {
		return err
	}

	return nil
}

func (service Service) createContainer(ctx ContainerResolutionContext, container Container, serviceName string, init bool) corev1.Container {
	c := corev1.Container{
		Name: func() string {
			if container.Name != "" {
				return container.Name
			}
			return serviceName
		}(),
		Image:   container.Image,
		Command: container.Command,
		Args:    container.Args,
	}
	env := NewEnv()
	if service.InjectPostgresVariables {
		env = env.Append(
			DefaultPostgresEnvVarsWithPrefix(*ctx.Postgres, ctx.Stack.GetServiceName(ctx.Module), "")...,
		)
	}

	if ctx.Configuration.Spec.Monitoring != nil {
		env = env.Append(
			MonitoringEnvVarsWithPrefix(*ctx.Configuration.Spec.Monitoring)...,
		)
	}

	if !init {
		env = env.Append(
			Env("DEBUG", fmt.Sprintf("%v", ctx.Stack.Spec.Debug)),
			Env("DEV", fmt.Sprintf("%v", ctx.Stack.Spec.Dev)),
			Env("STACK_URL", ctx.Stack.URL()),
			Env("OTEL_SERVICE_NAME", serviceName),
		)
	}

	c.Env = env.Append(container.Env...).ToCoreEnv()

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

		switch container.Liveness {
		case LivenessDefault:
			c.LivenessProbe = common.DefaultLiveness()
		case LivenessLegacy:
			c.LivenessProbe = common.LegacyLiveness()
		}
		if service.Port != 0 {
			c.Ports = common.SinglePort("http", service.Port)
		}
	}
	return c
}
