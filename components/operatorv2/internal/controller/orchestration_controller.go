/*
Copyright 2023.

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
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	formancev1beta1 "github.com/formancehq/operator/v2/api/v1beta1"
)

// OrchestrationController reconciles a Orchestration object
type OrchestrationController struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=formance.com,resources=orchestrations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=orchestrations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=orchestrations/finalizers,verbs=update

func (r *OrchestrationController) Reconcile(ctx context.Context, orchestration *v1beta1.Orchestration) error {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OrchestrationController) SetupWithManager(mgr ctrl.Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&formancev1beta1.Orchestration{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForOrchestration(client client.Client, scheme *runtime.Scheme) *OrchestrationController {
	return &OrchestrationController{
		Client: client,
		Scheme: scheme,
	}
}
