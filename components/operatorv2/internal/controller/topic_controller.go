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

package controller

import (
	"context"

	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/internal"
	"github.com/nats-io/nats.go"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TopicReconciler reconciles a Topic object
type TopicReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=formance.com,resources=topics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=topics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=topics/finalizers,verbs=update

func (r *TopicReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx, "topic", req.NamespacedName)
	log.Info("Starting reconciliation")

	topic := &v1beta1.Topic{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name: req.Name,
	}, topic); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if len(topic.Spec.Queries) == 0 {
		if err := r.Client.Delete(ctx, topic); err != nil {
			return ctrl.Result{}, nil
		}
	}

	stackSelectorRequirement, err := labels.NewRequirement("formance.com/stack", selection.In, []string{"any", topic.Spec.Stack})
	if err != nil {
		return ctrl.Result{}, err
	}

	brokerConfigList := &v1beta1.BrokerConfigurationList{}
	if err := r.Client.List(ctx, brokerConfigList, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*stackSelectorRequirement),
	}); err != nil {
		return ctrl.Result{}, err
	}

	switch len(brokerConfigList.Items) {
	case 0:
		topic.Status.Error = "unable to find a broker configuration"
		topic.Status.Ready = false
	case 1:
		configuration := brokerConfigList.Items[0]
		switch {
		case configuration.Spec.Nats != nil:
			if err := r.createNATSStream(GetObjectName(topic.Spec.Stack, topic.Name), configuration); err != nil {
				topic.Status.Error = err.Error()
				topic.Status.Ready = false
			} else {
				topic.Status.Error = ""
				topic.Status.Ready = true
			}
		default:
			topic.Status.Error = ""
			topic.Status.Ready = true
		}
	default:
		topic.Status.Error = "multiple broker configuration object found"
		topic.Status.Ready = false
	}

	if err := r.Client.Status().Update(ctx, topic); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *TopicReconciler) createNATSStream(topicName string, configuration v1beta1.BrokerConfiguration) error {
	nc, err := nats.Connect(configuration.Spec.Nats.URL)
	if err != nil {
		return err
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		return err
	}
	streamConfig := nats.StreamConfig{
		Name:      topicName,
		Subjects:  []string{topicName},
		Retention: nats.InterestPolicy,
		Replicas:  configuration.Spec.Nats.Replicas,
	}

	_, err = js.StreamInfo(topicName)
	if err != nil {
		_, err := js.AddStream(&streamConfig)
		return err
	}

	_, err = js.UpdateStream(&streamConfig)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *TopicReconciler) SetupWithManager(mgr ctrl.Manager) error {
	indexer := mgr.GetFieldIndexer()
	if err := indexer.IndexField(context.Background(), &v1beta1.Topic{}, ".spec.service", func(rawObj client.Object) []string {
		return []string{rawObj.(*v1beta1.Topic).Spec.Service}
	}); err != nil {
		return err
	}

	if err := indexer.IndexField(context.Background(), &v1beta1.Topic{}, ".spec.stack", func(rawObj client.Object) []string {
		return []string{rawObj.(*v1beta1.Topic).Spec.Stack}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Topic{}).
		Complete(r)
}

func NewTopicReconciler(client client.Client, scheme *runtime.Scheme) *TopicReconciler {
	return &TopicReconciler{
		Client: client,
		Scheme: scheme,
	}
}