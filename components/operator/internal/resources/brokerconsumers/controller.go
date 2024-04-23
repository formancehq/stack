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

package brokerconsumers

import (
	"fmt"
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/jobs"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"
)

const (
	ConditionTypeReady                      = "Ready"
	ConditionTypeBrokerTopicCreated         = "BrokerTopicCreated"
	ConditionTypeNatsStackConsumerCreated   = "NatsStackConsumerCreated"
	ConditionTypeNatsServiceConsumerCreated = "NatsServiceConsumerCreated"
)

//+kubebuilder:rbac:groups=formance.com,resources=brokerconsumers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=brokerconsumers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=brokerconsumers/finalizers,verbs=update

func Reconcile(ctx core.Context, stack *v1beta1.Stack, consumer *v1beta1.BrokerConsumer) error {

	for _, service := range consumer.Spec.Services {
		topic := &v1beta1.BrokerTopic{}
		if err := ctx.GetClient().Get(ctx, types.NamespacedName{
			Name: core.GetObjectName(consumer.Spec.Stack, service),
		}, topic); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
			topic = &v1beta1.BrokerTopic{
				ObjectMeta: ctrl.ObjectMeta{
					Name: core.GetObjectName(consumer.Spec.Stack, service),
				},
				Spec: v1beta1.BrokerTopicSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: consumer.Spec.Stack,
					},
					Service: service,
				},
			}
			if err := controllerutil.SetOwnerReference(consumer, topic, ctx.GetScheme()); err != nil {
				return err
			}

			if err := controllerutil.SetOwnerReference(stack, topic, ctx.GetScheme()); err != nil {
				return err
			}
			if err := ctx.GetClient().Create(ctx, topic); err != nil {
				return err
			}
			return nil
		} else {
			patch := client.MergeFromWithOptions(topic.DeepCopy(), client.MergeFromWithOptimisticLock{})
			if err := controllerutil.SetOwnerReference(consumer, topic, ctx.GetScheme()); err != nil {
				return err
			}
			if err := ctx.GetClient().Patch(ctx, topic, patch); err != nil {
				return err
			}
		}

		if !topic.Status.Ready {
			consumer.Status.SetCondition(v1beta1.Condition{
				Type:               ConditionTypeReady,
				Status:             "False",
				ObservedGeneration: consumer.Generation,
				LastTransitionTime: metav1.Now(),
				Message:            fmt.Sprintf("BrokerTopic %s not yet ready", topic.Name),
			})
			consumer.Status.SetCondition(v1beta1.Condition{
				Type:               ConditionTypeBrokerTopicCreated,
				Status:             "False",
				ObservedGeneration: consumer.Generation,
				LastTransitionTime: metav1.Now(),
				Message:            fmt.Sprintf("BrokerTopic %s not yet ready", topic.Name),
			})
			return core.NewPendingError()
		}
	}

	consumer.Status.SetCondition(v1beta1.Condition{
		Type:               ConditionTypeBrokerTopicCreated,
		Status:             "True",
		ObservedGeneration: consumer.Generation,
		LastTransitionTime: metav1.Now(),
		Message:            "All topics created",
	})

	broker := &v1beta1.Broker{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: stack.Name,
	}, broker); err != nil {
		return err
	}

	if !broker.Status.Ready {
		return core.NewPendingError().WithMessage("broker not ready")
	}

	if broker.Status.URI.Scheme == "nats" {
		switch broker.Status.Mode {
		case v1beta1.ModeOneStreamByStack:
			if !consumer.Status.CheckCondition(v1beta1.CheckCondition{
				Type:       ConditionTypeNatsStackConsumerCreated,
				Generation: consumer.Generation,
			}) {
				if err := createStackNatsConsumer(ctx, stack, consumer, broker); err != nil {
					return err
				}
			}
		case v1beta1.ModeOneStreamByService:
			for _, service := range consumer.Spec.Services {
				if !consumer.Status.CheckCondition(v1beta1.CheckCondition{
					Type:       ConditionTypeNatsServiceConsumerCreated,
					Generation: consumer.Generation,
					Reason:     service,
				}) {
					if err := createServiceNatsConsumer(ctx, stack, consumer, broker, service); err != nil {
						return err
					}
				}
			}
		}
	}

	consumer.Status.SetCondition(v1beta1.Condition{
		Type:               ConditionTypeReady,
		Status:             "True",
		ObservedGeneration: consumer.Generation,
		LastTransitionTime: metav1.Now(),
		Message:            "Consumer completely configured",
	})

	return nil
}

func createServiceNatsConsumer(ctx core.Context, stack *v1beta1.Stack, consumer *v1beta1.BrokerConsumer, broker *v1beta1.Broker, service string) error {
	const script = `
	exists=$(nats consumer ls $STACK-$SERVICE -n | grep $NAME)
	[[ -z "$exists" ]] || {
		nats --server $NATS_URI consumer add $STACK-$SERVICE $NAME \
			--deliver-group $NAME \
			--deliver all \
			--max-pending 1024 \
			--ack explicit \
			--target $STACK-$NAME \
			--replay instant \
			--filter $STACK-$SERVICE \
			--defaults
	}`

	natsBoxImage, err := registries.GetNatsBoxImage(ctx, stack, "0.14.1")
	if err != nil {
		return err
	}

	err = jobs.Handle(ctx, consumer, "cc-"+service, corev1.Container{
		Image: natsBoxImage,
		Name:  "create-consumer",
		Args:  core.ShellScript(script),
		Env: []corev1.EnvVar{
			core.Env("NATS_URI", fmt.Sprintf("nats://%s", broker.Status.URI.Host)),
			core.Env("STACK", stack.Name),
			core.Env("NAME", consumer.Spec.QueriedBy),
			core.Env("SERVICE", service),
		},
	})

	if err != nil {
		consumer.Status.SetCondition(v1beta1.Condition{
			Type:               ConditionTypeReady,
			Status:             "False",
			ObservedGeneration: consumer.Generation,
			LastTransitionTime: metav1.Now(),
			Message:            fmt.Sprintf("Error creating consumer on nats: %s", err),
		})
		consumer.Status.SetCondition(v1beta1.Condition{
			Type:               ConditionTypeNatsServiceConsumerCreated,
			Status:             "False",
			ObservedGeneration: consumer.Generation,
			LastTransitionTime: metav1.Now(),
			Message:            fmt.Sprintf("Error creating consumer on nats: %s", err),
			Reason:             service,
		})
	} else {
		consumer.Status.SetCondition(v1beta1.Condition{
			Type:               ConditionTypeNatsServiceConsumerCreated,
			Status:             "True",
			ObservedGeneration: consumer.Generation,
			LastTransitionTime: metav1.Now(),
			Message:            "Nats consumer created",
			Reason:             service,
		})
	}
	return err
}

func createStackNatsConsumer(ctx core.Context, stack *v1beta1.Stack, consumer *v1beta1.BrokerConsumer, broker *v1beta1.Broker) error {
	const script = `
	filters=""
	for f in $SUBJECTS; do
		filters="$filters --filter $f"
	done
	nats --server $NATS_URI consumer add $STREAM $NAME \
		--deliver-group $DELIVER \
		--deliver all \
		--max-pending 1024 \
		--ack explicit \
		--target $STREAM-$NAME \
		--replay instant \
		--defaults $filters
	`

	natsBoxImage, err := registries.GetNatsBoxImage(ctx, stack, "0.14.1")
	if err != nil {
		return err
	}

	consumerName := consumer.Spec.QueriedBy
	if consumer.Spec.Name != "" {
		consumerName += "_" + consumer.Spec.Name
	}

	err = jobs.Handle(ctx, consumer, "create-consumer", corev1.Container{
		Image: natsBoxImage,
		Name:  "create-consumer",
		Args:  core.ShellScript(script),
		Env: []corev1.EnvVar{
			core.Env("NATS_URI", fmt.Sprintf("nats://%s", broker.Status.URI.Host)),
			core.Env("STREAM", stack.Name),
			core.Env("NAME", consumerName),
			core.Env("DELIVER", consumer.Spec.QueriedBy),
			core.Env("SUBJECTS", strings.Join(
				collectionutils.Map(consumer.Spec.Services, func(from string) string {
					return fmt.Sprintf("%s.%s", stack.Name, from)
				}), " ",
			)),
		},
	})
	if err != nil {
		consumer.Status.SetCondition(v1beta1.Condition{
			Type:               ConditionTypeReady,
			Status:             "False",
			ObservedGeneration: consumer.Generation,
			LastTransitionTime: metav1.Now(),
			Message:            fmt.Sprintf("Error creating consumer on nats: %s", err),
		})
		consumer.Status.SetCondition(v1beta1.Condition{
			Type:               ConditionTypeNatsStackConsumerCreated,
			Status:             "False",
			ObservedGeneration: consumer.Generation,
			LastTransitionTime: metav1.Now(),
			Message:            fmt.Sprintf("Error creating consumer on nats: %s", err),
		})
	} else {
		consumer.Status.SetCondition(v1beta1.Condition{
			Type:               ConditionTypeNatsStackConsumerCreated,
			Status:             "True",
			ObservedGeneration: consumer.Generation,
			LastTransitionTime: metav1.Now(),
			Message:            "Nats consumer created",
		})
	}
	return err
}
