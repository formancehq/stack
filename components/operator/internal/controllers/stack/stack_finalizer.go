package stack

import (
	"context"
	"time"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllers/stack/backup"
	"github.com/formancehq/operator/internal/controllers/stack/delete"
	"github.com/formancehq/operator/internal/controllers/stack/storage/s3"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	finalizerName = "backup-delete"
)

var (
	subjects = []string{"ledger", "payments"}
)

func (r *Reconciler) HandleFinalizer(ctx context.Context, log logr.Logger, stack *stackv1beta3.Stack, conf *stackv1beta3.Configuration, req ctrl.Request) (*ctrl.Result, error) {

	// examine DeletionTimestamp to determine if object is under deletion
	// The object is being deleted
	if !stack.ObjectMeta.DeletionTimestamp.IsZero() && controllerutil.ContainsFinalizer(stack, finalizerName) {
		// Make sur that the object is disable
		stack.Spec.Disabled = true
		if err := r.client.Update(ctx, stack); err != nil {
			return &ctrl.Result{}, err
		}

		// We also need to make sure that all deployements are Terminated,
		// So that no one is already accessing our database.
		found := &appsv1.DeploymentList{}
		opts := &client.ListOptions{
			Namespace: req.Name,
		}

		if err := r.client.List(ctx, found, opts); err != nil {
			return &ctrl.Result{}, err
		}

		if len(found.Items) > 0 {
			return &ctrl.Result{
				Requeue:      true,
				RequeueAfter: time.Second,
			}, nil
		}

		// S3 is optional
		if conf.Spec.S3 != nil {
			if err := r.backupStack(ctx, conf, stack.Name, stack.ObjectMeta.DeletionTimestamp, log); err != nil {
				return &ctrl.Result{
					Requeue:      true,
					RequeueAfter: time.Second,
				}, err
			}
			// Stack has been fully backuped, we can now delete it
		}

		// our finalizer is present, so lets handle any external dependency
		res, err := r.deleteStack(ctx, stack, conf, log)
		if err != nil {
			// if fail to delete the external dependency here, return with error
			// so that it can be retried
			return res, err
		}

		// remove our finalizer from the list and update it.
		controllerutil.RemoveFinalizer(stack, finalizerName)
		if err := r.client.Update(ctx, stack); err != nil {
			return &ctrl.Result{}, err
		}

		// Stop reconciliation as the item is being deleted
		return &ctrl.Result{}, nil

	}

	// The object is not being deleted, so if it does not have our finalizer,
	// then lets add the finalizer and update the object. This is equivalent
	// registering our finalizer.
	if !controllerutil.ContainsFinalizer(stack, finalizerName) {
		controllerutil.AddFinalizer(stack, finalizerName)
		if err := r.client.Update(ctx, stack); err != nil {
			return &ctrl.Result{}, err
		}
	}

	return nil, nil
}

func (r *Reconciler) deleteStack(ctx context.Context, stack *stackv1beta3.Stack, conf *stackv1beta3.Configuration, log logr.Logger) (*ctrl.Result, error) {
	log.Info("start deleting databases " + stack.Name)
	if err := delete.DeleteByService(conf, stack.Name, log); err != nil {
		return &ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, err
	}

	log.Info("start deleting brokers subjects " + stack.Name)
	if err := delete.DeleteByBrokers(conf, stack.Name, subjects, log); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Reconciler) backupStack(ctx context.Context, conf *stackv1beta3.Configuration, stackName string, t *v1.Time, log logr.Logger) error {
	session, err := s3.NewSession(
		conf.Spec.S3.S3SecretConfig.AccessKey,
		conf.Spec.S3.S3SecretConfig.SecretKey,
		conf.Spec.S3.Endpoint,
		conf.Spec.S3.Region,
		conf.Spec.S3.ForceStylePath,
		conf.Spec.S3.Insecure,
	)

	if err != nil {
		log.Error(err, "Cannot create s3 client")
		return err
	}

	storage := s3.NewS3Storage(session, conf.Spec.S3.Bucket)

	log.Info("start backup for " + stackName)
	if err := backup.BackupServices(conf, stackName, storage, t, log); err != nil {
		return err
	}

	return nil
}
