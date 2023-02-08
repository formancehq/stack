package v1beta3

import (
	pkgapisv1beta2 "github.com/formancehq/operator/pkg/apis/v1beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type PostgresConfigCreateDatabase struct {
	pkgapisv1beta2.PostgresConfigWithDatabase `json:",inline"`
	CreateDatabase                            bool `json:"createDatabase"`
}

func (in *PostgresConfigCreateDatabase) CreateDatabaseInitContainer() corev1.Container {
	return corev1.Container{
		Name:            "init-db",
		Image:           "postgres:13",
		ImagePullPolicy: corev1.PullIfNotPresent,
		Env:             in.Env(""),
		Command: []string{
			"sh",
			"-c",
			`psql -Atx ${POSTGRES_NO_DATABASE_URI}/postgres -c "SELECT 1 FROM pg_database WHERE datname = '${POSTGRES_DATABASE}'" | grep -q 1 && echo "Base already exists" || psql -Atx ${POSTGRES_NO_DATABASE_URI}/postgres -c "CREATE DATABASE \"${POSTGRES_DATABASE}\""`,
		},
	}
}

type Broker struct {
	// +optional
	Kafka *pkgapisv1beta2.KafkaConfig `json:"kafka,omitempty"`
	// +optional
	Nats *pkgapisv1beta2.NatsConfig `json:"nats,omitempty"`
}

func (b Broker) Validate() []*field.Error {
	if b.Kafka == nil && b.Nats == nil {
		return field.ErrorList{
			field.Invalid(
				field.NewPath("kafka"),
				nil,
				"either 'nats' or 'kafka' config must be specified",
			),
		}
	}
	return nil
}

func (b Broker) Env(serviceName string, prefix string) []corev1.EnvVar {
	switch {
	case b.Kafka != nil:
		return b.Kafka.Env(prefix)
	case b.Nats != nil:
		return b.Nats.Env(serviceName, prefix)
	default:
		panic("should not happen")
	}
}

type CollectorConfig struct {
	Broker `json:",inline"`
	Topic  string `json:"topic"`
}

func (c CollectorConfig) Env(serviceName, prefix string) []corev1.EnvVar {
	ret := c.Broker.Env(serviceName, prefix)
	return append(ret, pkgapisv1beta2.EnvWithPrefix(prefix, "PUBLISHER_TOPIC_MAPPING", "*:"+c.Topic))
}
