package core

import (
	"fmt"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	corev1 "k8s.io/api/core/v1"
)

func Env(name string, value string) corev1.EnvVar {
	return corev1.EnvVar{
		Name:  name,
		Value: value,
	}
}

func EnvFromBool(name string, vb bool) corev1.EnvVar {
	value := "true"
	if !vb {
		value = "false"
	}
	return Env(name, value)
}

func EnvFromConfig(name, configName, key string) corev1.EnvVar {
	return corev1.EnvVar{
		Name: name,
		ValueFrom: &corev1.EnvVarSource{
			ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
				Key: key,
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configName,
				},
			},
		},
	}
}

func EnvFromSecret(name, secretName, key string) corev1.EnvVar {
	return corev1.EnvVar{
		Name: name,
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				Key: key,
				LocalObjectReference: corev1.LocalObjectReference{
					Name: secretName,
				},
			},
		},
	}
}

func EnvVarPlaceholder(key string) string {
	return fmt.Sprintf("$(%s)", key)
}

func ComputeEnvVar(format string, keys ...string) string {
	return fmt.Sprintf(format,
		collectionutils.Map(keys, func(key string) any {
			return EnvVarPlaceholder(key)
		})...,
	)
}
