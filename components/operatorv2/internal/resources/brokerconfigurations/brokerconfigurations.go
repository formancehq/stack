package brokerconfigurations

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func Require(ctx core.Context, stackName string) (*v1beta1.BrokerConfiguration, error) {
	brokerConfiguration, err := Get(ctx, stackName)
	if err != nil {
		return nil, err
	}
	if brokerConfiguration == nil {
		return nil, errors.New("no broker configuration found")
	}
	return brokerConfiguration, nil
}

func Get(ctx core.Context, stackName string) (*v1beta1.BrokerConfiguration, error) {

	stackSelectorRequirement, err := labels.NewRequirement("formance.com/stack", selection.In, []string{"any", stackName})
	if err != nil {
		return nil, err
	}

	brokerConfigurationList := &v1beta1.BrokerConfigurationList{}
	if err := ctx.GetClient().List(ctx, brokerConfigurationList, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*stackSelectorRequirement),
	}); err != nil {
		return nil, err
	}

	switch len(brokerConfigurationList.Items) {
	case 0:
		return nil, nil
	case 1:
		return &brokerConfigurationList.Items[0], nil
	default:
		return nil, errors.New("found multiple broker config")
	}
}

func GetEnvVars(ctx core.Context, stackName, serviceName string) ([]v1.EnvVar, error) {
	configuration, err := Get(ctx, stackName)
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
