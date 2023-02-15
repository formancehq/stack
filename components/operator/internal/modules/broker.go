package modules

import (
	"strings"

	"github.com/formancehq/operator/apis/stack/v1beta3"
)

func BrokerEnvVars(broker v1beta3.Broker, serviceName string) ContainerEnv {
	return BrokerEnvVarsWithPrefix(broker, serviceName, "")
}

func BrokerEnvVarsWithPrefix(broker v1beta3.Broker, serviceName, prefix string) ContainerEnv {
	ret := ContainerEnv{}

	if broker.Kafka != nil {
		ret = ret.Append(
			Env(prefix+"BROKER", "kafka"),
			Env(prefix+"PUBLISHER_KAFKA_ENABLED", "true"),
			Env(prefix+"PUBLISHER_KAFKA_BROKER", strings.Join(broker.Kafka.Brokers, ",")),
		)
		if broker.Kafka.SASL != nil {
			ret = ret.Append(
				Env(prefix+"PUBLISHER_KAFKA_SASL_ENABLED", "true"),
				Env(prefix+"PUBLISHER_KAFKA_SASL_USERNAME", broker.Kafka.SASL.Username),
				Env(prefix+"PUBLISHER_KAFKA_SASL_PASSWORD", broker.Kafka.SASL.Password),
				Env(prefix+"PUBLISHER_KAFKA_SASL_MECHANISM", broker.Kafka.SASL.Mechanism),
				Env(prefix+"PUBLISHER_KAFKA_SASL_SCRAM_SHA_SIZE", broker.Kafka.SASL.ScramSHASize),
			)
		}
		if broker.Kafka.TLS {
			ret = ret.Append(
				Env(prefix+"PUBLISHER_KAFKA_TLS_ENABLED", "true"),
			)
		}
	} else {
		ret = ret.Append(
			Env(prefix+"BROKER", "nats"),
			Env(prefix+"PUBLISHER_NATS_ENABLED", "true"),
			Env(prefix+"PUBLISHER_NATS_URL", broker.Nats.URL),
			Env(prefix+"PUBLISHER_NATS_CLIENT_ID", serviceName),
		)
	}
	return ret
}
