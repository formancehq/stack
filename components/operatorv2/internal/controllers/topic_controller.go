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

package controllers

import (
	"context"
	"fmt"
	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/brokerconfigurations"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	types "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// TopicController reconciles a Topic object
type TopicController struct{}

//+kubebuilder:rbac:groups=formance.com,resources=topics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=topics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=topics/finalizers,verbs=update

func (r *TopicController) Reconcile(ctx Context, topic *v1beta1.Topic) error {

	if len(topic.Spec.Queries) == 0 {
		if err := ctx.GetClient().Delete(ctx, topic); err != nil {
			return nil
		}
	}

	stackSelectorRequirement, err := labels.NewRequirement("formance.com/stack", selection.In, []string{"any", topic.Spec.Stack})
	if err != nil {
		return err
	}

	brokerConfigList := &v1beta1.BrokerConfigurationList{}
	if err := ctx.GetClient().List(ctx, brokerConfigList, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*stackSelectorRequirement),
	}); err != nil {
		return err
	}

	switch len(brokerConfigList.Items) {
	case 0:
		topic.Status.Error = "unable to find a broker configuration"
		topic.Status.Ready = false
	case 1:
		configuration := brokerConfigList.Items[0]
		switch {
		case configuration.Spec.Nats != nil:
			job, err := r.createJob(ctx, topic, configuration)
			if err != nil {
				topic.Status.Error = err.Error()
				topic.Status.Ready = false
			} else {
				topic.Status.Error = ""
				topic.Status.Ready = job.Status.Succeeded > 0
			}
		default:
			topic.Status.Error = ""
			topic.Status.Ready = true
		}
	default:
		topic.Status.Error = "multiple broker configuration object found"
		topic.Status.Ready = false
	}

	return nil
}

func (r *TopicController) createJob(ctx Context,
	topic *v1beta1.Topic, configuration v1beta1.BrokerConfiguration) (*batchv1.Job, error) {

	job, _, err := CreateOrUpdate[*batchv1.Job](ctx, types.NamespacedName{
		Namespace: topic.Spec.Stack,
		Name:      fmt.Sprintf("%s-create-topic", topic.Spec.Service),
	},
		func(t *batchv1.Job) {
			args := []string{"nats", "stream", "add",
				"--server", fmt.Sprintf("nats://%s", configuration.Spec.Nats.URL),
				"--retention", "interest",
				"--subjects", topic.Name,
				"--defaults",
			}
			if configuration.Spec.Nats.Replicas > 0 {
				args = append(args, "--replicas", fmt.Sprint(configuration.Spec.Nats.Replicas))
			}
			args = append(args, topic.Name)

			t.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyOnFailure
			t.Spec.Template.Spec.Containers = []corev1.Container{{
				Image: "natsio/nats-box:0.14.1",
				Name:  "create-topic",
				Args:  args,
			}}
		},
		WithController[*batchv1.Job](ctx.GetScheme(), topic),
	)
	return job, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *TopicController) SetupWithManager(mgr Manager) (*builder.Builder, error) {

	indexer := mgr.GetFieldIndexer()
	if err := indexer.IndexField(context.Background(), &v1beta1.Topic{}, ".spec.service", func(rawObj client.Object) []string {
		return []string{rawObj.(*v1beta1.Topic).Spec.Service}
	}); err != nil {
		return nil, err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Topic{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(
			&v1beta1.BrokerConfiguration{},
			handler.EnqueueRequestsFromMapFunc(
				brokerconfigurations.Watch(mgr, &v1beta1.TopicList{})),
		).
		Owns(&batchv1.Job{}), nil
}

func ForTopic() *TopicController {
	return &TopicController{}
}
