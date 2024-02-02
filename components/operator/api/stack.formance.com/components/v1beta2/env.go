package v1beta2

import (
	corev1 "k8s.io/api/core/v1"
)

func Env(key, value string) corev1.EnvVar {
	return corev1.EnvVar{
		Name:  key,
		Value: value,
	}
}
