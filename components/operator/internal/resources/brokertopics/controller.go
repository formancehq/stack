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
	"github.com/formancehq/operator/internal/resources/settings"
	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
)

//+kubebuilder:rbac:groups=formance.com,resources=brokertopics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=brokertopics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=brokertopics/finalizers,verbs=update

func Reconcile(ctx core.Context, stack *v1beta1.Stack, topic *v1beta1.BrokerTopic) error {

	if len(topic.GetOwnerReferences()) == 1 { // Remains only the stack
		return ctx.GetClient().Delete(ctx, topic)
	}

	if topic.Status.Ready {
		return nil
	}

	brokerURI, err := settings.RequireURL(ctx, stack.Name, "broker.dsn")
	if err != nil {
		return err
	}
	if brokerURI == nil {
		return errors.New("broker configuration not found")
	}

	topic.Status.URI = brokerURI

	switch {
	case brokerURI.Scheme == "nats":
		if err := createJob(ctx, topic, brokerURI); err != nil {
			return err
		}
	}

	return nil
}

func clear(ctx core.Context, topic *v1beta1.BrokerTopic) error {
	if topic.Status.Ready && topic.Status.URI != nil {
		switch {
		case topic.Status.URI.Scheme == "nats":
			if err := deleteJob(ctx, topic); err != nil {
				return err
			}
		}
	}

	return nil
}

func init() {
	core.Init(
		core.WithResourceReconciler(Reconcile,
			core.WithOwn[*v1beta1.BrokerTopic](&batchv1.Job{}),
			core.WithWatchSettings[*v1beta1.BrokerTopic](),
			core.WithFinalizer[*v1beta1.BrokerTopic]("clear", clear),
		),
		core.WithSimpleIndex[*v1beta1.BrokerTopic](".spec.service", func(t *v1beta1.BrokerTopic) string {
			return t.Spec.Service
		}),
	)
}
