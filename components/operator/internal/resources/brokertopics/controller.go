/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package brokertopics

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	batchv1 "k8s.io/api/batch/v1"
)

//+kubebuilder:rbac:groups=formance.com,resources=brokertopics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=brokertopics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=brokertopics/finalizers,verbs=update

func Reconcile(ctx core.Context, stack *v1beta1.Stack, topic *v1beta1.BrokerTopic) error {

	if len(topic.GetOwnerReferences()) == 0 {
		if err := clear(ctx, topic); err != nil {
			return err
		}
		return core.ErrDeleted
	}

	if topic.Status.Ready {
		return nil
	}

	brokerConfiguration, err := FindBrokerConfiguration(ctx, stack)
	if err != nil {
		return err
	}

	topic.Status.Configuration = brokerConfiguration

	switch {
	case brokerConfiguration.Nats != nil:
		job, err := createJob(ctx, topic, *brokerConfiguration)
		if err != nil {
			return err
		}

		if job.Status.Succeeded == 0 {
			return core.ErrPending
		}
	}

	return nil
}

func FindBrokerConfiguration(ctx core.Context, stack *v1beta1.Stack) (*v1beta1.BrokerConfiguration, error) {
	brokerKind, err := settings.RequireString(ctx, stack.Name, "broker", "kind")
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
	url, err := settings.RequireString(ctx, stack.Name, "broker", "nats", "endpoint")
	if err != nil {
		return nil, err
	}

	replicas, err := settings.GetIntOrDefault(ctx, stack.Name, 1, "broker.nats.replicas")
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

	endpoints, err := settings.GetStringSlice(ctx, stack.Name, "broker", "kafka", "endpoints")
	if err != nil {
		return nil, err
	}

	tls, err := settings.GetBoolOrFalse(ctx, stack.Name, "broker", "kafka", "ssl", "enabled")
	if err != nil {
		return nil, err
	}

	saslEnabled, err := settings.GetBoolOrDefault(ctx, stack.Name, false, "kafka", "sasl", "enabled")
	if err != nil {
		return nil, err
	}

	var saslConfig *v1beta1.BrokerKafkaSASLConfig
	if saslEnabled {
		saslUsername, err := settings.GetStringOrEmpty(ctx, stack.Name, "kafka", "sasl", "username")
		if err != nil {
			return nil, err
		}
		saslPassword, err := settings.GetStringOrEmpty(ctx, stack.Name, "kafka", "sasl", "password")
		if err != nil {
			return nil, err
		}
		saslMechanism, err := settings.GetStringOrEmpty(ctx, stack.Name, "kafka", "sasl", "mechanism")
		if err != nil {
			return nil, err
		}
		saslScramSHASize, err := settings.GetStringOrEmpty(ctx, stack.Name, "kafka", "sasl", "scram-sha-size")
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

func clear(ctx core.Context, topic *v1beta1.BrokerTopic) error {
	if topic.Status.Ready && topic.Status.Configuration != nil {
		switch {
		case topic.Status.Configuration.Nats != nil:
			job, err := deleteJob(ctx, topic)
			if err != nil {
				return err
			}

			if job.Status.Succeeded == 0 {
				return core.ErrPending
			}
		}
	}

	return ctx.GetClient().Delete(ctx, topic)
}

func init() {
	core.Init(
		core.WithStackDependencyReconciler(Reconcile,
			core.WithOwn(&batchv1.Job{}),
			core.WithWatchConfigurationObject(&v1beta1.Settings{}),
			core.WithWatchStack(),
		),
		core.WithSimpleIndex[*v1beta1.BrokerTopic](".spec.service", func(t *v1beta1.BrokerTopic) string {
			return t.Spec.Service
		}),
	)
}
