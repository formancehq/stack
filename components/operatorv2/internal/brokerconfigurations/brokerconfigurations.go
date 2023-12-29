package brokerconfigurations

import (
	"errors"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	utils2 "github.com/formancehq/operator/v2/internal/utils"
	errors2 "github.com/pkg/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func Require(ctx reconcilers.Context, stackName string) (*v1beta1.BrokerConfiguration, error) {
	brokerConfiguration, err := Get(ctx, stackName)
	if err != nil {
		return nil, err
	}
	if brokerConfiguration == nil {
		return nil, errors.New("no broker configuration found")
	}
	return brokerConfiguration, nil
}

func Get(ctx reconcilers.Context, stackName string) (*v1beta1.BrokerConfiguration, error) {

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
		return nil, errors2.New("found multiple broker config")
	}
}

func GetEnvVars(ctx reconcilers.Context, stackName, serviceName string) ([]v1.EnvVar, error) {
	configuration, err := Get(ctx, stackName)
	if err != nil {
		return nil, err
	}
	if configuration == nil {
		return nil, utils2.ErrNotFound
	}

	return BrokerEnvVars(*configuration, serviceName), nil
}

func BrokerEnvVars(broker v1beta1.BrokerConfiguration, serviceName string) []v1.EnvVar {
	ret := make([]v1.EnvVar, 0)

	if broker.Spec.Kafka != nil {
		ret = append(ret,
			utils2.Env("BROKER", "kafka"),
			utils2.Env("PUBLISHER_KAFKA_ENABLED", "true"),
			utils2.Env("PUBLISHER_KAFKA_BROKER", strings.Join(broker.Spec.Kafka.Brokers, " ")),
		)
		if broker.Spec.Kafka.SASL != nil {
			ret = append(ret,
				utils2.Env("PUBLISHER_KAFKA_SASL_ENABLED", "true"),
				utils2.Env("PUBLISHER_KAFKA_SASL_USERNAME", broker.Spec.Kafka.SASL.Username),
				utils2.Env("PUBLISHER_KAFKA_SASL_PASSWORD", broker.Spec.Kafka.SASL.Password),
				utils2.Env("PUBLISHER_KAFKA_SASL_MECHANISM", broker.Spec.Kafka.SASL.Mechanism),
				utils2.Env("PUBLISHER_KAFKA_SASL_SCRAM_SHA_SIZE", broker.Spec.Kafka.SASL.ScramSHASize),
			)
		}
		if broker.Spec.Kafka.TLS {
			ret = append(ret,
				utils2.Env("PUBLISHER_KAFKA_TLS_ENABLED", "true"),
			)
		}
	} else {
		ret = append(ret,
			utils2.Env("BROKER", "nats"),
			utils2.Env("PUBLISHER_NATS_ENABLED", "true"),
			utils2.Env("PUBLISHER_NATS_URL", broker.Spec.Nats.URL),
			utils2.Env("PUBLISHER_NATS_CLIENT_ID", serviceName),
		)
	}
	return ret
}
