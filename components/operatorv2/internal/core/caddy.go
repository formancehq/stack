package core

import (
	"github.com/formancehq/operator/v2/api/formance.com/v1beta1"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
)

func ConfigureCaddy(caddyfile *v1.ConfigMap, image string, env []v1.EnvVar,
	resourceProperties *v1beta1.ResourceProperties) []ObjectMutator[*appsv1.Deployment] {
	return []ObjectMutator[*appsv1.Deployment]{
		func(t *appsv1.Deployment) {
			t.Spec.Template.Annotations = collectionutils.MergeMaps(t.Spec.Template.Annotations, map[string]string{
				"caddyfile-hash": HashFromConfigMap(caddyfile),
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
					Resources: GetResourcesWithDefault(resourceProperties, ResourceSizeSmall()),
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
