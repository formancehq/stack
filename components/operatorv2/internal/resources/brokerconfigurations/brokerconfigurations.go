package brokerconfigurations

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	"k8s.io/api/core/v1"
	"strings"
)

func GetEnvVars(ctx core.Context, stackName, serviceName string) ([]v1.EnvVar, error) {
	configuration, err := stacks.GetByLabel[*v1beta1.BrokerConfiguration](ctx, stackName)
	if err != nil {
		return nil, err
	}
	if configuration == nil {
		return nil, stacks.ErrNotFound
	}

	return BrokerEnvVars(configuration.Spec, serviceName), nil
}

func BrokerEnvVars(broker v1beta1.BrokerConfigurationSpec, serviceName string) []v1.EnvVar {
	ret := make([]v1.EnvVar, 0)

	if broker.Kafka != nil {
		ret = append(ret,
			core.Env("BROKER", "kafka"),
			core.Env("PUBLISHER_KAFKA_ENABLED", "true"),
			core.Env("PUBLISHER_KAFKA_BROKER", strings.Join(broker.Kafka.Brokers, " ")),
		)
		if broker.Kafka.SASL != nil {
			ret = append(ret,
				core.Env("PUBLISHER_KAFKA_SASL_ENABLED", "true"),
				core.Env("PUBLISHER_KAFKA_SASL_USERNAME", broker.Kafka.SASL.Username),
				core.Env("PUBLISHER_KAFKA_SASL_PASSWORD", broker.Kafka.SASL.Password),
				core.Env("PUBLISHER_KAFKA_SASL_MECHANISM", broker.Kafka.SASL.Mechanism),
				core.Env("PUBLISHER_KAFKA_SASL_SCRAM_SHA_SIZE", broker.Kafka.SASL.ScramSHASize),
			)
		}
		if broker.Kafka.TLS {
			ret = append(ret,
				core.Env("PUBLISHER_KAFKA_TLS_ENABLED", "true"),
			)
		}
	} else {
		ret = append(ret,
			core.Env("BROKER", "nats"),
			core.Env("PUBLISHER_NATS_ENABLED", "true"),
			core.Env("PUBLISHER_NATS_URL", broker.Nats.URL),
			core.Env("PUBLISHER_NATS_CLIENT_ID", serviceName),
		)
	}
	return ret
}
