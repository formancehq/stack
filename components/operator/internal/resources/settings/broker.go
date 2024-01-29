package settings

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	v1 "k8s.io/api/core/v1"
)

func ResolveBrokerEnvVars(ctx core.Context, stack *v1beta1.Stack, serviceName string) ([]v1.EnvVar, error) {
	uri, err := RequireURL(ctx, stack.Name, "broker.dsn")
	if err != nil {
		return nil, err
	}

	return GetBrokerEnvVars(uri, stack.Name, serviceName), nil
}

func GetBrokerEnvVars(brokerURI *v1beta1.URI, stackName, serviceName string) []v1.EnvVar {
	return GetBrokerEnvVarsWithPrefix(brokerURI, stackName, serviceName, "")
}

func GetBrokerEnvVarsWithPrefix(brokerURI *v1beta1.URI, stackName, serviceName, prefix string) []v1.EnvVar {
	ret := make([]v1.EnvVar, 0)

	ret = append(ret, core.Env(fmt.Sprintf("%sBROKER", prefix), brokerURI.Scheme))

	switch {
	case brokerURI.Scheme == "kafka":
		ret = append(ret,
			core.Env(fmt.Sprintf("%sBROKER", prefix), "kafka"),
			core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_ENABLED", prefix), "true"),
			core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_BROKER", prefix), brokerURI.Host),
		)
		if IsTrue(brokerURI.Query().Get("saslEnabled")) {
			ret = append(ret,
				core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_SASL_ENABLED", prefix), "true"),
				core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_SASL_USERNAME", prefix), brokerURI.Query().Get("saslUsername")),
				core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_SASL_PASSWORD", prefix), brokerURI.Query().Get("saslPassword")),
				core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_SASL_MECHANISM", prefix), brokerURI.Query().Get("saslMechanism")),
				core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_SASL_SCRAM_SHA_SIZE", prefix), brokerURI.Query().Get("saslSCRAMSHASize")),
			)
		}
		if IsTrue(brokerURI.Query().Get("tls")) {
			ret = append(ret,
				core.Env(fmt.Sprintf("%sPUBLISHER_KAFKA_TLS_ENABLED", prefix), "true"),
			)
		}
	case brokerURI.Scheme == "nats":
		ret = append(ret,
			core.Env(fmt.Sprintf("%sPUBLISHER_NATS_ENABLED", prefix), "true"),
			core.Env(fmt.Sprintf("%sPUBLISHER_NATS_URL", prefix), brokerURI.Host),
			core.Env(fmt.Sprintf("%sPUBLISHER_NATS_CLIENT_ID", prefix), fmt.Sprintf("%s-%s", stackName, serviceName)),
		)
	}

	return ret
}
