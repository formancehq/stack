package core

import (
	corev1 "k8s.io/api/core/v1"
)

func NewVolumeFromConfigMap(name string, configMap *corev1.ConfigMap) corev1.Volume {
	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configMap.Name,
				},
			},
		},
	}
}

func NewVolumeMount(name, mountPath string, readOnly bool) corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      name,
		ReadOnly:  readOnly,
		MountPath: mountPath,
	}
}
