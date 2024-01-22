package settings

import (
	"fmt"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	v1 "k8s.io/api/core/v1"
)

func ResolveBrokerEnvVars(ctx core.Context, stack *v1beta1.Stack, serviceName string) ([]v1.EnvVar, error) {
	configuration, err := FindBrokerConfiguration(ctx, stack)
	if err != nil {
		return nil, err
	}
	if configuration == nil {
		return nil, core.ErrNotFound
	}

	return GetBrokerEnvVars(*configuration, stack.Name, serviceName), nil
}

func GetBrokerEnvVars(broker v1beta1.BrokerConfiguration, stackName, serviceName string) []v1.EnvVar {
	return GetBrokerEnvVarsWithPrefix(broker, stackName, serviceName, "")
}

func GetBrokerEnvVarsWithPrefix(broker v1beta1.BrokerConfiguration, stackName, serviceName, prefix string) []v1.EnvVar {
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

func FindBrokerConfiguration(ctx core.Context, stack *v1beta1.Stack) (*v1beta1.BrokerConfiguration, error) {
	brokerKind, err := RequireString(ctx, stack.Name, "broker", "kind")
	if err != nil {
		return nil, err
	}
	switch brokerKind {
	case "kafka":
		return resolveKafkaConfiguration(ctx, stack)
	case "nats":
		return resolveNatsConfiguration(ctx, stack)
	default:
		return nil, fmt.Errorf("broker kind '%s' unknown", brokerKind)
	}
}

func resolveNatsConfiguration(ctx core.Context, stack *v1beta1.Stack) (*v1beta1.BrokerConfiguration, error) {
	url, err := RequireString(ctx, stack.Name, "broker", "nats", "endpoint")
	if err != nil {
		return nil, err
	}

	replicas, err := GetIntOrDefault(ctx, stack.Name, 1, "broker.nats.replicas")
	if err != nil {
		return nil, err
	}

	return &v1beta1.BrokerConfiguration{
		Nats: &v1beta1.BrokerNatsConfig{
			URL:      url,
			Replicas: replicas,
		},
	}, nil
}

func resolveKafkaConfiguration(ctx core.Context, stack *v1beta1.Stack) (*v1beta1.BrokerConfiguration, error) {

	endpoints, err := GetStringSlice(ctx, stack.Name, "broker", "kafka", "endpoints")
	if err != nil {
		return nil, err
	}

	tls, err := GetBoolOrFalse(ctx, stack.Name, "broker", "kafka", "ssl", "enabled")
	if err != nil {
		return nil, err
	}

	saslEnabled, err := GetBoolOrDefault(ctx, stack.Name, false, "kafka", "sasl", "enabled")
	if err != nil {
		return nil, err
	}

	var saslConfig *v1beta1.BrokerKafkaSASLConfig
	if saslEnabled {
		saslUsername, err := GetStringOrEmpty(ctx, stack.Name, "kafka", "sasl", "username")
		if err != nil {
			return nil, err
		}
		saslPassword, err := GetStringOrEmpty(ctx, stack.Name, "kafka", "sasl", "password")
		if err != nil {
			return nil, err
		}
		saslMechanism, err := GetStringOrEmpty(ctx, stack.Name, "kafka", "sasl", "mechanism")
		if err != nil {
			return nil, err
		}
		saslScramSHASize, err := GetStringOrEmpty(ctx, stack.Name, "kafka", "sasl", "scram-sha-size")
		if err != nil {
			return nil, err
		}
		saslConfig = &v1beta1.BrokerKafkaSASLConfig{
			Username:     saslUsername,
			Password:     saslPassword,
			Mechanism:    saslMechanism,
			ScramSHASize: saslScramSHASize,
		}
	}

	return &v1beta1.BrokerConfiguration{
		Kafka: &v1beta1.BrokerKafkaConfig{
			Brokers: endpoints,
			TLS:     tls,
			SASL:    saslConfig,
		},
	}, nil
}
