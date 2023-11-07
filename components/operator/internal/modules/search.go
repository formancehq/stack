package modules

import (
	"fmt"
	"strings"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
)

func SearchEnvVars(rc ReconciliationConfig) ContainerEnv {
	env := ElasticSearchEnvVars(rc.Stack, rc.Configuration, rc.Versions).
		Append(
			Env("OPEN_SEARCH_SERVICE", fmt.Sprintf("%s:%d%s",
				rc.Configuration.Spec.Services.Search.ElasticSearchConfig.Host,
				rc.Configuration.Spec.Services.Search.ElasticSearchConfig.Port,
				rc.Configuration.Spec.Services.Search.ElasticSearchConfig.PathPrefix)),
			Env("OPEN_SEARCH_SCHEME", rc.Configuration.Spec.Services.Search.ElasticSearchConfig.Scheme),
			Env("MAPPING_INIT_DISABLED", "true"),
		)
	if rc.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth != nil {
		env = env.Append(
			Env("OPEN_SEARCH_USERNAME", rc.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Username),
			Env("OPEN_SEARCH_PASSWORD", rc.Configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Password),
		)
	}
	if rc.Versions.IsLower("search", "v0.7.0") {
		env = env.Append(Env("ES_INDICES", rc.Stack.Name))
	} else {
		env = env.Append(Env("ES_INDICES", stackv1beta3.DefaultESIndex))
	}

	return env
}

func ElasticSearchEnvVars(stack *stackv1beta3.Stack, configuration *stackv1beta3.Configuration, versions *stackv1beta3.Versions) ContainerEnv {
	ret := ContainerEnv{
		Env("OPENSEARCH_URL", configuration.Spec.Services.Search.ElasticSearchConfig.Endpoint()),
		Env("OPENSEARCH_BATCHING_COUNT", fmt.Sprint(configuration.Spec.Services.Search.Batching.Count)),
		Env("OPENSEARCH_BATCHING_PERIOD", configuration.Spec.Services.Search.Batching.Period),
		Env("TOPIC_PREFIX", stack.Name+"-"),
	}
	if versions.IsLower("search", "v0.7.0") {
		ret = append(ret, Env("OPENSEARCH_INDEX", stack.Name))
	} else {
		ret = append(ret, Env("OPENSEARCH_INDEX", stackv1beta3.DefaultESIndex))
	}
	if configuration.Spec.Broker.Kafka != nil {
		ret = ret.Append(
			Env("KAFKA_ADDRESS", strings.Join(configuration.Spec.Broker.Kafka.Brokers, ",")),
		)
		if configuration.Spec.Broker.Kafka.TLS {
			ret = ret.Append(
				Env("KAFKA_TLS_ENABLED", "true"),
			)
		}
		if configuration.Spec.Broker.Kafka.SASL != nil {
			ret = ret.Append(
				Env("KAFKA_SASL_USERNAME", configuration.Spec.Broker.Kafka.SASL.Username),
				Env("KAFKA_SASL_PASSWORD", configuration.Spec.Broker.Kafka.SASL.Password),
				Env("KAFKA_SASL_MECHANISM", configuration.Spec.Broker.Kafka.SASL.Mechanism),
			)
		}
	}
	if configuration.Spec.Broker.Nats != nil {
		ret = ret.Append(
			Env("NATS_URL", configuration.Spec.Broker.Nats.URL),
		)
	}
	if configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth != nil {
		ret = ret.Append(
			Env("BASIC_AUTH_ENABLED", "true"),
		)
		if configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.SecretName == "" {
			ret = ret.Append(
				Env("BASIC_AUTH_USERNAME", configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Username),
				Env("BASIC_AUTH_PASSWORD", configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.Password),
			)
		} else {
			ret = ret.Append(
				EnvFromSecret("BASIC_AUTH_USERNAME", configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.SecretName, "username"),
				EnvFromSecret("BASIC_AUTH_PASSWORD", configuration.Spec.Services.Search.ElasticSearchConfig.BasicAuth.SecretName, "password"),
			)
		}
	}
	return ret
}
