package common

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	utils2 "github.com/formancehq/operator/v2/internal/utils"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	v12 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
)

func ConfigureCaddy(caddyfile *v1.ConfigMap, image string, env []v1.EnvVar,
	resourceProperties *v1beta1.ResourceProperties) []utils2.ObjectMutator[*v12.Deployment] {
	return []utils2.ObjectMutator[*v12.Deployment]{
		func(t *v12.Deployment) {
			t.Spec.Template.Annotations = collectionutils.MergeMaps(t.Spec.Template.Annotations, map[string]string{
				"caddyfile-hash": utils2.HashFromConfigMap(caddyfile),
			})
			t.Spec.Template.Spec.Volumes = []v1.Volume{
				volumeFromConfigMap("caddyfile", caddyfile),
			}
			t.Spec.Template.Spec.Containers = []v1.Container{
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
					Resources: utils2.GetResourcesWithDefault(resourceProperties, utils2.ResourceSizeSmall()),
					VolumeMounts: []v1.VolumeMount{
						volumeMount("caddyfile", "/gateway"),
					},
					Ports: []v1.ContainerPort{{
						Name:          "http",
						ContainerPort: 8080,
					}},
				},
			}
		},
	}
}

func volumeFromConfigMap(name string, cm *v1.ConfigMap) v1.Volume {
	return v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: cm.Name,
				},
			},
		},
	}
}

func volumeMount(name, mountPath string) v1.VolumeMount {
	return v1.VolumeMount{
		Name:      name,
		ReadOnly:  true,
		MountPath: mountPath,
	}
}
