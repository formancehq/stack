package common

import (
	"strings"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllerutils"
	corev1 "k8s.io/api/core/v1"
)

func BrokerEnvVars(broker v1beta3.Broker, serviceName, prefix string) []corev1.EnvVar {
	ret := make([]corev1.EnvVar, 0)

	if broker.Kafka != nil {
		ret = append(ret,
			controllerutils.EnvWithPrefix(prefix, "PUBLISHER_KAFKA_ENABLED", "true"),
			controllerutils.EnvWithPrefix(prefix, "PUBLISHER_KAFKA_BROKER",
				strings.Join(broker.Kafka.Brokers, ",")),
		)
		if broker.Kafka.SASL != nil {
			ret = append(ret,
				controllerutils.EnvWithPrefix(prefix, "PUBLISHER_KAFKA_SASL_ENABLED", "true"),
				controllerutils.EnvWithPrefix(prefix, "PUBLISHER_KAFKA_SASL_USERNAME", broker.Kafka.SASL.Username),
				controllerutils.EnvWithPrefix(prefix, "PUBLISHER_KAFKA_SASL_PASSWORD", broker.Kafka.SASL.Password),
				controllerutils.EnvWithPrefix(prefix, "PUBLISHER_KAFKA_SASL_MECHANISM", broker.Kafka.SASL.Mechanism),
				controllerutils.EnvWithPrefix(prefix, "PUBLISHER_KAFKA_SASL_SCRAM_SHA_SIZE", broker.Kafka.SASL.ScramSHASize),
			)
		}
		if broker.Kafka.TLS {
			ret = append(ret, controllerutils.EnvWithPrefix(prefix, "PUBLISHER_KAFKA_TLS_ENABLED", "true"))
		}
	} else {
		ret = append(ret,
			controllerutils.EnvWithPrefix(prefix, "PUBLISHER_NATS_ENABLED", "true"),
			controllerutils.EnvWithPrefix(prefix, "PUBLISHER_NATS_URL", broker.Nats.URL),
			controllerutils.EnvWithPrefix(prefix, "PUBLISHER_NATS_CLIENT_ID", serviceName),
		)
	}
	return ret
}
