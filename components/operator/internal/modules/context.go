package modules

import (
	"crypto/sha256"
	"encoding/base64"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	corev1 "k8s.io/api/core/v1"
)

type ContainerResolutionConfiguration struct {
	ServiceInstallConfiguration
	Configs ConfigHandles
	Secrets SecretHandles
}

func (ctx ContainerResolutionConfiguration) GetConfig(name string) ConfigHandle {
	return ctx.Configs[name]
}

func (ctx ContainerResolutionConfiguration) GetSecret(name string) SecretHandle {
	return ctx.Secrets[name]
}

func (ctx ContainerResolutionConfiguration) volumes(serviceName string) ([]corev1.Volume, string) {
	ret := make([]corev1.Volume, 0)
	hash := sha256.New()
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
		hash.Write([]byte(config.Hash))
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
		hash.Write([]byte(secret.Hash))
	}

	return ret, base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

type RegisteredService struct {
	Port int32
	Service
}

type RegisteredModule struct {
	Module   Module
	Services map[string]RegisteredService
}

type RegisteredModules map[string]RegisteredModule

type ServiceInstallConfiguration struct {
	ReconciliationConfig
	RegisteredModules RegisteredModules
	PostgresConfig    *stackv1beta3.PostgresConfig
}
