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
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"reflect"
	"strconv"
	"strings"

	"github.com/formancehq/operator/internal/core"
	"github.com/pkg/errors"

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

	databaseConfiguration, err := findConfiguration(ctx, database)
	if err != nil {
		return err
	}

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

func GetSetting(ctx core.Context, stack string, keys ...string) (*string, error) {
	key := strings.Join(keys, ".")
	list := &v1beta1.SettingsList{}
	if err := ctx.GetClient().List(ctx, list, client.MatchingFields{
		"stack": stack,
		"key":   key,
	}); err != nil {
		return nil, err
	}

	if len(list.Items) == 0 {
		return nil, nil
	}
	if len(list.Items) > 1 {
		return nil, fmt.Errorf("found multiple matching setting with key '%s' and stack '%s'", key, stack)
	}

	return &list.Items[0].Spec.Value, nil
}

func GetStringSetting(ctx core.Context, stack string, keys ...string) (*string, error) {
	return GetSetting(ctx, stack, keys...)
}

func ValueOrDefault[T any](v *T, defaultValue T) T {
	if v == nil {
		return defaultValue
	}
	return *v
}

func GetInt64Setting(ctx core.Context, stack string, keys ...string) (*int64, error) {
	value, err := GetSetting(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	intValue, err := strconv.ParseInt(*value, 10, 64)
	if err != nil {
		return nil, err
	}

	return &intValue, nil
}

func GetBoolSetting(ctx core.Context, stack string, keys ...string) (*bool, error) {
	value, err := GetSetting(ctx, stack, keys...)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	return pointer.For(*value == "true"), nil
}

func findConfiguration(ctx core.Context, database *v1beta1.Database) (*v1beta1.DatabaseConfiguration, error) {
	host, err := GetStringSetting(ctx, database.Spec.Stack, "databases", database.Spec.Service, "host")
	if err != nil {
		return nil, err
	}
	if host == nil {
		host, err = GetStringSetting(ctx, database.Spec.Stack, "databases", "host")
		if err != nil {
			return nil, err
		}
	}

	if host == nil {
		return nil, errors.New("missing database host")
	}

	port, err := GetInt64Setting(ctx, database.Spec.Stack, "databases", database.Spec.Service, "port")
	if err != nil {
		return nil, err
	}
	if port == nil {
		port, err = GetInt64Setting(ctx, database.Spec.Stack, "databases", "port")
		if err != nil {
			return nil, err
		}
	}
	if port == nil {
		port = pointer.For(int64(5432)) // default postgres port
	}

	username, err := GetStringSetting(ctx, database.Spec.Stack, "databases", database.Spec.Service, "username")
	if err != nil {
		return nil, err
	}
	if username == nil {
		username, err = GetStringSetting(ctx, database.Spec.Stack, "databases", "username")
		if err != nil {
			return nil, err
		}
	}

	password, err := GetStringSetting(ctx, database.Spec.Stack, "databases", database.Spec.Service, "password")
	if err != nil {
		return nil, err
	}
	if password == nil {
		password, err = GetStringSetting(ctx, database.Spec.Stack, "databases", "password")
		if err != nil {
			return nil, err
		}
	}

	credentialsFromSecret, err := GetStringSetting(ctx, database.Spec.Stack, "databases", database.Spec.Service, "credentials-from-secret")
	if err != nil {
		return nil, err
	}
	if credentialsFromSecret == nil {
		credentialsFromSecret, err = GetStringSetting(ctx, database.Spec.Stack, "databases", "secret")
		if err != nil {
			return nil, err
		}
	}

	disableSSLMode, err := GetBoolSetting(ctx, database.Spec.Stack, "databases", database.Spec.Service, "ssl", "disable")
	if err != nil {
		return nil, err
	}
	if disableSSLMode == nil {
		disableSSLMode, err = GetBoolSetting(ctx, database.Spec.Stack, "databases", "disable-ssl-mode")
		if err != nil {
			return nil, err
		}
	}

	return &v1beta1.DatabaseConfiguration{
		Port:                  int(*port),
		Host:                  *host,
		Username:              ValueOrDefault(username, ""),
		Password:              ValueOrDefault(password, ""),
		CredentialsFromSecret: ValueOrDefault(credentialsFromSecret, ""),
		DisableSSLMode:        ValueOrDefault(disableSSLMode, false),
	}, nil
}

func init() {
	core.Init(
		core.WithStackDependencyReconciler(Reconcile,
			core.WithOwn(&batchv1.Job{}),
			core.WithWatchConfigurationObject(&v1beta1.RegistriesConfiguration{}),
			core.WithWatchConfigurationObject(&v1beta1.Settings{}),
		),
	)
}
