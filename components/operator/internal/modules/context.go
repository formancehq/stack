package modules

import (
	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	corev1 "k8s.io/api/core/v1"
)

type ContainerResolutionContext struct {
	ServiceInstallContext
	Configs ConfigHandles
	Secrets SecretHandles
}

func (ctx ContainerResolutionContext) GetConfig(name string) ConfigHandle {
	return ctx.Configs[name]
}

func (ctx ContainerResolutionContext) GetSecret(name string) SecretHandle {
	return ctx.Secrets[name]
}

func (ctx ContainerResolutionContext) volumes(serviceName string) []corev1.Volume {
	ret := make([]corev1.Volume, 0)
	for _, configName := range ctx.Configs.sort() {
		config := ctx.Configs[configName]
		if config.MountPath == "" {
			continue
		}
		ret = append(ret, corev1.Volume{
			Name: configName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: serviceName + "-" + configName,
					},
				},
			},
		})
	}
	for _, secretName := range ctx.Secrets.sort() {
		secret := ctx.Secrets[secretName]
		if secret.MountPath == "" {
			continue
		}
		ret = append(ret, corev1.Volume{
			Name: secretName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: serviceName + "-" + secretName,
				},
			},
		})
	}

	return ret
}

type ModuleContext struct {
	Context
	PortAllocator PortAllocator
	Postgres      *stackv1beta3.PostgresConfig
	Module        string
}

func (ctx ModuleContext) HasVersionHigherOrEqual(version string) bool {
	return ctx.Versions.IsHigherOrEqual(ctx.Module, version)
}

func (ctx ModuleContext) HasVersionHigher(version string) bool {
	return ctx.Versions.IsHigher(ctx.Module, version)
}

func (ctx ModuleContext) HasVersionLower(version string) bool {
	return ctx.Versions.IsLower(ctx.Module, version)
}

type RegisteredModule struct {
	Module   Module
	Services Services
}

type RegisteredModules map[string]RegisteredModule

type ServiceInstallContext struct {
	ModuleContext
	RegisteredModules RegisteredModules
	PodDeployer       PodDeployer
}
