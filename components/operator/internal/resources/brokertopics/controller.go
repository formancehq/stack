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
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
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

	brokerConfiguration, err := core.RequireLabelledConfig[*v1beta1.BrokerConfiguration](ctx, topic.Spec.Stack)
	if err != nil {
		return err
	}

	topic.Status.Configuration = &brokerConfiguration.Spec

	switch {
	case brokerConfiguration.Spec.Nats != nil:
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
			core.WithWatchConfigurationObject(&v1beta1.BrokerConfiguration{}),
			core.WithWatchStack(),
		),
		core.WithIndex[*v1beta1.BrokerTopic](".spec.service", func(t *v1beta1.BrokerTopic) string {
			return t.Spec.Service
		}),
	)
}
