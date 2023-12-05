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

package stack

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
	pkgError "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// MigrationReconciler reconciles a Migration object
type MigrationReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	platform modules.Platform
}

// +kubebuilder:rbac:groups=stack.formance.com,resources=migrations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stack.formance.com,resources=migrations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stack.formance.com,resources=migrations/finalizers,verbs=update

func (r *MigrationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx, "migration", req.NamespacedName)
	log.Info("Starting reconciliation")

	migration := &stackv1beta3.Migration{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      req.Name,
	}, migration); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Object not found, skip")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, pkgError.Wrap(err, "Reading target")
	}

	if migration.Status.Terminated {
		return ctrl.Result{}, nil
	}

	configuration := &stackv1beta3.Configuration{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      migration.Spec.Configuration,
	}, configuration); err != nil {
		return ctrl.Result{}, nil
	}

	versions := &stackv1beta3.Versions{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      migration.Spec.Version,
	}, versions); err != nil {
		return ctrl.Result{}, nil
	}

	stack := &stackv1beta3.Stack{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      req.Namespace,
	}, stack); err != nil {
		return ctrl.Result{}, nil
	}

	module := modules.Get(migration.Spec.Module)
	version, ok := module.Versions()[migration.Spec.TargetedVersion]
	if !ok {
		return ctrl.Result{}, pkgError.New("migration not found")
	}

	rc := modules.ReconciliationConfig{
		Stack:         stack,
		Configuration: configuration,
		Versions:      versions,
		Platform:      r.platform,
	}

	var fn func(ctx context.Context, jobRunner modules.JobRunner, config modules.MigrationConfig) (bool, error)

	if migration.Spec.PostUpgrade {
		fn = version.PostUpgrade
	} else {
		fn = version.PreUpgrade
	}

	jobRunner := modules.NewJobRunner(r.Client, r.Scheme, stack, migration, fmt.Sprintf("%s-%s-%s-upgrade-",
		module.Name(), migration.Spec.TargetedVersion, migration.Discriminator()))

	migrationConfig := modules.MigrationConfig{
		ReconciliationConfig: rc,
		Version:              migration.Spec.CurrentVersion,
	}

	var returnedError error
	if terminated, err := fn(ctx, jobRunner, migrationConfig); err != nil {
		returnedError = err
		migration.Status.Err = returnedError.Error()
	} else {
		migration.Status.Terminated = terminated
		migration.Status.Err = ""
	}

	if err := r.Client.Status().Update(ctx, migration); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, returnedError
}

// SetupWithManager sets up the controller with the Manager.
func (r *MigrationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stackv1beta3.Migration{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&batchv1.Job{}).
		Complete(r)
}

func NewMigrationReconciler(client client.Client, scheme *runtime.Scheme, platform modules.Platform) *MigrationReconciler {
	return &MigrationReconciler{
		platform: platform,
		Client:   client,
		Scheme:   scheme,
	}
}
