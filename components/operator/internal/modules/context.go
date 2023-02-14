package modules

import (
	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	corev1 "k8s.io/api/core/v1"
)

type ContainerResolutionContext struct {
	InstallContext
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

type PrepareContext struct {
	Context
	Postgres *stackv1beta3.PostgresConfig
	Module   string
}

type RegisteredModule struct {
	Module   Module
	Services Services
}

type RegisteredModules map[string]RegisteredModule

type InstallContext struct {
	PrepareContext
	RegisteredModules RegisteredModules
}
