// +kubebuilder:object:generate=true
package v1beta3

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

type ConfigSource struct {
	// Selects a key of a ConfigMap.
	// +optional
	ConfigMapKeyRef *corev1.ConfigMapKeySelector `json:"configMapKeyRef,omitempty" protobuf:"bytes,3,opt,name=configMapKeyRef"`
	// Selects a key of a secret in the pod's namespace
	// +optional
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty" protobuf:"bytes,4,opt,name=secretKeyRef"`
}

func (c *ConfigSource) Env() *corev1.EnvVarSource {
	return &corev1.EnvVarSource{
		ConfigMapKeyRef: c.ConfigMapKeyRef,
		SecretKeyRef:    c.SecretKeyRef,
	}
}

type PostgresConfig struct {
	Port     int    `json:"port"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	// +optional
	DisableSSLMode bool `json:"disableSSLMode"`
}

func (in *PostgresConfig) DSN() string {
	queryParams := ""
	if in.DisableSSLMode {
		queryParams = "?sslmode=disable"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/postgres%s", in.Username, in.Password, in.Host, in.Port, queryParams)
}
