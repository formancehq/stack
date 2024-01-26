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
	"github.com/formancehq/operator/internal/resources/secrets"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	databaseFinalizer = "finalize.databases.formance.com"
)

//+kubebuilder:rbac:groups=formance.com,resources=databases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=databases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=databases/finalizers,verbs=update

func Reconcile(ctx core.Context, _ *v1beta1.Stack, database *v1beta1.Database) error {

	databaseURL, err := settings.RequireURL(ctx, database.Spec.Stack, "postgres", database.Spec.Service, "uri")
	if err != nil {
		return errors.Wrap(err, "retrieving database configuration")
	}

	switch {
	// TODO: We have multiple occurrences of this type of code, we need to factorize
	case !database.Status.Ready || database.Status.URL().Host == databaseURL.Host:
		// Some job fields are immutable (env vars for example)
		// So, if the configuration has changed, wee need to delete the job,
		// Then recreate a new one
		if database.Status.DSN != databaseURL.String() {
			object := &batchv1.Job{}
			object.SetName(fmt.Sprintf("%s-create-database", database.Spec.Service))
			object.SetNamespace(database.Spec.Stack)
			if err := ctx.GetClient().Delete(ctx, object, &client.DeleteOptions{
				GracePeriodSeconds: pointer.For(int64(0)),
			}); client.IgnoreNotFound(err) != nil {
				return err
			}
		}

		if err := secrets.SyncFromURLs(ctx, database, database.Status.URL(), databaseURL); err != nil {
			return err
		}

		dbName := core.GetObjectName(database.Spec.Stack, database.Spec.Service)
		database.Status.DSN = databaseURL.String()
		database.Status.Database = dbName

		job, err := createDatabase(ctx, database)
		if err != nil {
			return err
		}

		if job.Status.Succeeded == 0 {
			return core.NewPendingError()
		}
	case !reflect.DeepEqual(database.Status.DSN, databaseURL.String()):
		database.Status.OutOfSync = true
	}

	return nil
}

func Delete(ctx core.Context, database *v1beta1.Database) error {
	if database.Status.DSN == "" {
		return nil
	}
	job, _, err := core.CreateOrUpdate[*batchv1.Job](ctx, types.NamespacedName{
		Namespace: database.Spec.Stack,
		Name:      fmt.Sprintf("%s-drop-database", database.Spec.Service),
	},
		func(t *batchv1.Job) error {
			dropDBCommand := `psql -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USERNAME} -c "DROP DATABASE \"${POSTGRES_DATABASE}\""`
			if isDisabledSSLMode(database.Status.URL()) {
				dropDBCommand += ` "sslmode=disable"`
			}

			t.Spec.BackoffLimit = pointer.For(int32(10000))
			t.Spec.TTLSecondsAfterFinished = pointer.For(int32(30))
			t.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyOnFailure
			t.Spec.Template.Spec.Containers = []v1.Container{{
				Name:  "drop-database",
				Image: "postgres:15-alpine",
				Args:  []string{"sh", "-c", dropDBCommand},
				Env:   psqlEnvVars(database),
			}}

			return nil
		},
		core.WithController[*batchv1.Job](ctx.GetScheme(), database),
	)
	if err != nil {
		return err
	}

	if job.Status.Succeeded == 0 {
		return core.NewPendingError()
	}

	return nil
}

func createDatabase(ctx core.Context, database *v1beta1.Database) (*batchv1.Job, error) {

	job, _, err := core.CreateOrUpdate[*batchv1.Job](ctx, types.NamespacedName{
		Namespace: database.Spec.Stack,
		Name:      fmt.Sprintf("%s-create-database", database.Spec.Service),
	},
		func(t *batchv1.Job) error {
			// PG does not support 'CREATE IF NOT EXISTS ' construct, emulate it with the above query
			createDBCommand := `echo SELECT \'CREATE DATABASE \"${POSTGRES_DATABASE}\"\' WHERE NOT EXISTS \(SELECT FROM pg_database WHERE datname = \'${POSTGRES_DATABASE}\'\)\\gexec | psql -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USERNAME}`
			if isDisabledSSLMode(database.Status.URL()) {
				createDBCommand += ` "sslmode=disable"`
			}

			t.Spec.BackoffLimit = pointer.For(int32(10000))
			t.Spec.TTLSecondsAfterFinished = pointer.For(int32(30))
			t.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyOnFailure
			t.Spec.Template.Spec.Containers = []v1.Container{{
				Name:  "create-database",
				Image: "postgres:15-alpine",
				Args:  []string{"sh", "-c", createDBCommand},
				Env:   psqlEnvVars(database),
			}}

			return nil
		},
		core.WithController[*batchv1.Job](ctx.GetScheme(), database),
	)
	return job, err
}

func psqlEnvVars(database *v1beta1.Database) []v1.EnvVar {
	return append(GetPostgresEnvVars(database),
		// psql use PGPASSWORD env var
		core.Env("PGPASSWORD", "$(POSTGRES_PASSWORD)"),
	)
}

func init() {
	core.Init(
		core.WithStackDependencyReconciler(Reconcile,
			core.WithOwn[*v1beta1.Database](&batchv1.Job{}),
			core.WithOwn[*v1beta1.Database](&v1.Secret{}),
			core.WithWatchSettings[*v1beta1.Database](),
			core.WithWatchSecrets[*v1beta1.Database](),
			core.WithFinalizer(databaseFinalizer, Delete),
		),
	)
}
