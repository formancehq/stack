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

package databases

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pkg/errors"
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	databaseFinalizer = "finalize.databases.formance.com"
)

//+kubebuilder:rbac:groups=formance.com,resources=databases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=databases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=databases/finalizers,verbs=update

func Reconcile(ctx core.Context, stack *v1beta1.Stack, database *v1beta1.Database) error {

	if !database.DeletionTimestamp.IsZero() {
		if database.Status.Configuration != nil {
			job, err := deleteJob(ctx, database)
			if err != nil {
				return err
			}

			if job.Status.Succeeded == 0 {
				return core.ErrPending
			}
		}

		patch := client.MergeFrom(database.DeepCopy())
		if updated := controllerutil.RemoveFinalizer(database, databaseFinalizer); updated {
			if err := ctx.GetClient().Patch(ctx, database, patch); err != nil {
				return errors.Wrap(err, "removing finalizer")
			}
		}

		return core.ErrDeleted
	}

	databaseConfiguration, err := settings.FindDatabaseConfiguration(ctx, database)
	if err != nil {
		return err
	}

	spew.Dump(databaseConfiguration)

	if databaseConfiguration == nil {
		return fmt.Errorf("unable to find a database configuration")
	}

	switch {
	case !database.Status.Ready:
		// Some job fields are immutable (env vars for example)
		// So, if the configuration has changed, wee need to delete the job,
		// Then recreate a new one
		if database.Status.Configuration != nil {
			if !reflect.DeepEqual(
				database.Status.Configuration.DatabaseConfiguration,
				*databaseConfiguration,
			) {
				object := &batchv1.Job{}
				object.SetName(fmt.Sprintf("%s-create-database", database.Spec.Service))
				object.SetNamespace(database.Spec.Stack)
				if err := ctx.GetClient().Delete(ctx, object, &client.DeleteOptions{
					GracePeriodSeconds: pointer.For(int64(0)),
				}); client.IgnoreNotFound(err) != nil {
					return err
				}
			}
		}

		dbName := core.GetObjectName(database.Spec.Stack, database.Spec.Service)
		job, err := createJob(ctx, *databaseConfiguration, database, dbName)
		if err != nil {
			return err
		}

		patch := client.MergeFrom(database.DeepCopy())
		if updated := controllerutil.AddFinalizer(database, databaseFinalizer); updated {
			if err := ctx.GetClient().Patch(ctx, database, patch); err != nil {
				return err
			}
		}

		database.Status.Configuration = &v1beta1.CreatedDatabase{
			DatabaseConfiguration: *databaseConfiguration,
			Database:              dbName,
		}

		if job.Status.Succeeded == 0 {
			return core.ErrPending
		}
	case !reflect.DeepEqual(database.Status.Configuration.DatabaseConfiguration,
		*databaseConfiguration):
		database.Status.OutOfSync = true
	}

	return nil
}

func init() {
	core.Init(
		core.WithStackDependencyReconciler(Reconcile,
			core.WithOwn(&batchv1.Job{}),
			core.WithWatchConfigurationObject(&v1beta1.Settings{}),
		),
	)
}
