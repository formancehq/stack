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
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/jobs"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/resourcereferences"
	"github.com/formancehq/operator/internal/resources/settings"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	databaseFinalizer = "delete-database"
)

//+kubebuilder:rbac:groups=formance.com,resources=databases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=databases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=databases/finalizers,verbs=update

func Reconcile(ctx core.Context, stack *v1beta1.Stack, database *v1beta1.Database) error {

	databaseURL, err := settings.RequireURL(ctx, database.Spec.Stack, "postgres", database.Spec.Service, "uri")
	if err != nil {
		return errors.Wrap(err, "retrieving database configuration")
	}

	if secret := databaseURL.Query().Get("secret"); secret != "" {
		_, err = resourcereferences.Create(ctx, database, "postgres", secret, &v1.Secret{})
	} else {
		err = resourcereferences.Delete(ctx, database, "postgres")
	}
	if err != nil {
		return err
	}

	if awsRole := databaseURL.Query().Get("awsRole"); awsRole != "" {
		_, err = resourcereferences.Create(ctx, database, "database", awsRole, &v1.ServiceAccount{})
	} else {
		err = resourcereferences.Delete(ctx, database, "database")
	}
	if err != nil {
		return err
	}

	switch {
	case !database.Status.Ready:
		// Some job fields are immutable (env vars for example)
		// So, if the configuration has changed, wee need to delete the job,
		// Then recreate a new one

		if database.Status.URI != nil && database.Status.URI.String() != databaseURL.String() {
			object := &batchv1.Job{}
			object.SetName(fmt.Sprintf("%s-create-database", database.Spec.Service))
			object.SetNamespace(database.Spec.Stack)
			if err := ctx.GetClient().Delete(ctx, object, &client.DeleteOptions{
				GracePeriodSeconds: pointer.For(int64(0)),
			}); client.IgnoreNotFound(err) != nil {
				return err
			}
		}

		dbName := core.GetObjectName(database.Spec.Stack, database.Spec.Service)
		database.Status.URI = databaseURL
		database.Status.Database = dbName

		if err := handleDatabaseJob(ctx, stack, database, "create-database", "db", "create"); err != nil {
			return err
		}
	case database.Status.URI.String() != databaseURL.String():
		database.Status.OutOfSync = true
	}

	return nil
}

func Delete(ctx core.Context, database *v1beta1.Database) error {
	if database.Status.URI == nil {
		return nil
	}

	clearDatabase, err := settings.GetBoolOrFalse(ctx, database.Spec.Stack, "clear-database")
	if err != nil {
		return err
	}
	if !clearDatabase {
		return nil
	}

	logger := log.FromContext(ctx)
	logger = logger.WithValues("name", database.Name)
	logger.Info("Deleting database")

	stack := &v1beta1.Stack{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: database.Spec.Stack,
	}, stack); err != nil {
		return err
	}

	if err := handleDatabaseJob(ctx, stack, database, "drop-database", "db", "drop"); err != nil {
		return err
	}

	database.Status.URI = nil
	logger.Info("Database deleted.")

	return nil
}

func handleDatabaseJob(ctx core.Context, stack *v1beta1.Stack, database *v1beta1.Database, name string, args ...string) error {

	operatorUtilsImageVersion, err := core.GetImageVersionForStack(ctx, stack, "operator-utils")
	if err != nil {
		return err
	}

	operatorUtilsImage, err := registries.GetImage(ctx, stack, "operator-utils", operatorUtilsImageVersion)
	if err != nil {
		return err
	}

	annotations := make(map[string]string)
	secretReference := &v1beta1.ResourceReference{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: fmt.Sprintf("%s-postgres", database.Name),
	}, secretReference); client.IgnoreNotFound(err) != nil {
		return err
	} else if err == nil {
		annotations["secret-hash"] = secretReference.Status.Hash
	}

	env := GetPostgresEnvVars(database)
	if database.Spec.Debug {
		env = append(env, core.Env("DEBUG", "true"))
	}

	return jobs.Handle(ctx, database, name, v1.Container{
		Name:  name,
		Image: operatorUtilsImage,
		Args:  args,
		Env:   env,
	},
		jobs.Mutator(core.WithAnnotations[*batchv1.Job](annotations)),
		jobs.WithServiceAccount(database.Status.URI.Query().Get("awsRole")),
	)
}

func init() {
	core.Init(
		core.WithResourceReconciler(Reconcile,
			core.WithOwn[*v1beta1.Database](&batchv1.Job{}),
			core.WithOwn[*v1beta1.Database](&v1beta1.ResourceReference{}),
			core.WithWatchSettings[*v1beta1.Database](),
			core.WithFinalizer(databaseFinalizer, Delete),
		),
	)
}
