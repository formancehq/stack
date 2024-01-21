package brokerconfigurations

import (
	"fmt"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	v1 "k8s.io/api/core/v1"
)

func GetEnvVars(ctx core.Context, stackName, serviceName string) ([]v1.EnvVar, error) {
	configuration, err := core.GetConfigurationObject[*v1beta1.BrokerConfiguration](ctx, stackName)
	if err != nil {
		return nil, err
	}
	if configuration == nil {
		return nil, core.ErrNotFound
	}

	return BrokerEnvVars(configuration.Spec, stackName, serviceName), nil
}

func BrokerEnvVars(broker v1beta1.BrokerConfigurationSpec, stackName, serviceName string) []v1.EnvVar {
	return BrokerEnvVarsWithPrefix(broker, stackName, serviceName, "")
}

func BrokerEnvVarsWithPrefix(broker v1beta1.BrokerConfigurationSpec, stackName, serviceName, prefix string) []v1.EnvVar {
	ret := make([]v1.EnvVar, 0)

	switch {
	case broker.Kafka != nil:
		ret = append(ret,
			core.Env(fmt.Sprintf("%sBROKER", prefix), "kafka"),
			core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_ENABLED", prefix), "true"),
			core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_BROKER", prefix), strings.Join(broker.Kafka.Brokers, " ")),
		)
		if broker.Kafka.SASL != nil {
			ret = append(ret,
				core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_SASL_ENABLED", prefix), "true"),
				core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_SASL_USERNAME", prefix), broker.Kafka.SASL.Username),
				core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_SASL_PASSWORD", prefix), broker.Kafka.SASL.Password),
				core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_SASL_MECHANISM", prefix), broker.Kafka.SASL.Mechanism),
				core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_SASL_SCRAM_SHA_SIZE", prefix), broker.Kafka.SASL.ScramSHASize),
			)
		}
		if broker.Kafka.TLS {
			ret = append(ret,
				core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_TLS_ENABLED", prefix), "true"),
			)
		}
	case broker.Nats != nil:
		ret = append(ret,
			core.Env(fmt.Sprintf("%sBROKER", prefix), "nats"),
			core.Env(fmt.Sprintf("%sPUBLISHER_NATS_ENABLED", prefix), "true"),
			core.Env(fmt.Sprintf("%sPUBLISHER_NATS_URL", prefix), broker.Nats.URL),
			core.Env(fmt.Sprintf("%sPUBLISHER_NATS_CLIENT_ID", prefix), fmt.Sprintf("%s-%s", stackName, serviceName)),
		)
	}

	return ret
}
