package core

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
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

// TODO: The stack reconciler can create a config map container env var for dev and debug
// This way, we avoid the need to fetch the stack object at each reconciliation loop
func GetDevEnvVars(stack *v1beta1.Stack, service interface {
	IsDebug() bool
	IsDev() bool
}) []corev1.EnvVar {
	return GetDevEnvVarsWithPrefix(stack, service, "")
}

func GetDevEnvVarsWithPrefix(stack *v1beta1.Stack, service interface {
	IsDebug() bool
	IsDev() bool
}, prefix string) []corev1.EnvVar {
	return []corev1.EnvVar{
		EnvFromBool(fmt.Sprintf("%sDEBUG", prefix), stack.Spec.Debug || service.IsDebug()),
		EnvFromBool(fmt.Sprintf("%sDEV", prefix), stack.Spec.Dev || service.IsDev()),
		Env(fmt.Sprintf("%sSTACK", prefix), stack.Name),
	}
}
