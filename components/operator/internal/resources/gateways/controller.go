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

package gateways

import (
	_ "embed"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/brokertopics"
	"github.com/formancehq/operator/internal/resources/services"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

//+kubebuilder:rbac:groups=formance.com,resources=gateways,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=gateways/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=gateways/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack, gateway *v1beta1.Gateway, version string) error {

	httpAPIs := make([]*v1beta1.HTTPAPI, 0)
	err := GetAllStackDependencies(ctx, gateway.Spec.Stack, &httpAPIs)
	if err != nil {
		return err
	}

	auth := &v1beta1.Auth{}
	ok, err := GetIfExists(ctx, stack.Name, auth)
	if err != nil {
		return err
	}
	if !ok {
		auth = nil
	}

	topic, err := createAuditTopic(ctx, stack, gateway)
	if err != nil {
		return err
	}

	configMap, err := createConfigMap(ctx, stack, gateway, httpAPIs, auth, topic)
	if err != nil {
		return err
	}

	if err := createDeployment(ctx, stack, gateway, configMap, topic, version); err != nil {
		return err
	}

	if _, err := services.Create(ctx, gateway, "gateway", func(t *corev1.Service) {
		if gateway.Spec.Service != nil {
			t.Annotations = gateway.Spec.Service.Annotations
		}
	}); err != nil {
		return err
	}

	if err := createIngress(ctx, stack, gateway); err != nil {
		return err
	}

	gateway.Status.SyncHTTPAPIs = Map(httpAPIs, func(from *v1beta1.HTTPAPI) string {
		return from.Spec.Name
	})
	gateway.Status.AuthEnabled = auth != nil

	return nil
}

func createAuditTopic(ctx Context, stack *v1beta1.Stack, gateway *v1beta1.Gateway) (*v1beta1.BrokerTopic, error) {
	if stack.Spec.EnableAudit && gateway.Spec.CompareVersion(stack, "v0.2.0") > 0 {
		topic, err := brokertopics.CreateOrUpdate(ctx, stack, gateway, "gateway", "audit")
		if err != nil {
			return nil, err
		}
		if !topic.Status.Ready {
			return nil, ErrPending
		}
		return topic, nil
	}
	return nil, nil
}

func init() {
	Init(
		WithModuleReconciler(Reconcile,
			WithOwn(&corev1.ConfigMap{}),
			WithOwn(&appsv1.Deployment{}),
			WithOwn(&corev1.Service{}),
			WithOwn(&networkingv1.Ingress{}),
			WithWatchConfigurationObject(&v1beta1.OpenTelemetryConfiguration{}),
			WithWatchDependency(&v1beta1.HTTPAPI{}),
			WithWatchDependency(&v1beta1.Auth{}),
			WithWatch[*v1beta1.BrokerTopic](brokertopics.Watch[*v1beta1.Gateway]("gateway")),
			WithWatchStack(),
		),
	)
}
