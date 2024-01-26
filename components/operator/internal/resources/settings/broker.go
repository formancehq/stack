package settings

import (
	"fmt"
	"net/url"
	"strconv"
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
	brokerDSN, err := GetString(ctx, stack.Name, "broker.dsn")
	if err != nil {
		return nil, err
	}
	if brokerDSN == nil {
		return nil, nil
	}

	switch {
	case strings.HasPrefix(*brokerDSN, "kafka://"):
		return resolveKafkaConfiguration(*brokerDSN)
	case strings.HasPrefix(*brokerDSN, "nats://"):
		return resolveNatsConfiguration(*brokerDSN)
	default:
		return nil, fmt.Errorf("broker kind '%s' unknown", *brokerDSN)
	}
}

func resolveNatsConfiguration(natsDSN string) (*v1beta1.BrokerConfiguration, error) {
	natsURI, err := url.Parse(natsDSN)
	if err != nil {
		return nil, err
	}

	if natsURI.Scheme != "nats" {
		return nil, fmt.Errorf("invalid nats uri: %s", natsDSN)
	}

	replicas := uint64(1)
	if replicasValue := natsURI.Query().Get("replicas"); replicasValue != "" {
		replicas, err = strconv.ParseUint(replicasValue, 10, 16)
		if err != nil {
			return nil, err
		}
	}

	return &v1beta1.BrokerConfiguration{
		Nats: &v1beta1.BrokerNatsConfig{
			URL:      natsURI.Host,
			Replicas: int(replicas),
		},
	}, nil
}

func resolveKafkaConfiguration(kafkaDSN string) (*v1beta1.BrokerConfiguration, error) {

	kafkaURI, err := url.Parse(kafkaDSN)
	if err != nil {
		return nil, err
	}

	if kafkaURI.Scheme != "kafka" {
		return nil, fmt.Errorf("invalid kafka uri: %s", kafkaDSN)
	}

	var saslConfig *v1beta1.BrokerKafkaSASLConfig
	if IsTrue(kafkaURI.Query().Get("saslEnabled")) {
		saslConfig = &v1beta1.BrokerKafkaSASLConfig{
			Username:     kafkaURI.Query().Get("saslUsername"),
			Password:     kafkaURI.Query().Get("saslPassword"),
			Mechanism:    kafkaURI.Query().Get("saslMechanism"),
			ScramSHASize: kafkaURI.Query().Get("saslSCRAMSHASize"),
		}
	}

	return &v1beta1.BrokerConfiguration{
		Kafka: &v1beta1.BrokerKafkaConfig{
			Brokers: []string{kafkaURI.Host},
			TLS:     IsTrue(kafkaURI.Query().Get("tls")),
			SASL:    saslConfig,
		},
	}, nil
}
