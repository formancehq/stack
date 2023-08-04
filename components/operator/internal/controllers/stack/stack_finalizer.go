package stack

import (
	"context"
	"time"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllers/stack/delete"
	"github.com/formancehq/operator/internal/controllers/stack/storage/s3"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *Reconciler) HandleFinalizer(ctx context.Context, log logr.Logger, stack *stackv1beta3.Stack, conf *stackv1beta3.Configuration, req ctrl.Request) (*ctrl.Result, error) {
	finalizerName := stack.Name

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
			if err := r.backupStack(ctx, conf, stack, log); err != nil {
				return &ctrl.Result{}, err
			}
		}

		// our finalizer is present, so lets handle any external dependency
		if err := r.deleteStack(ctx, req.NamespacedName, stack, conf, log); err != nil {
			// if fail to delete the external dependency here, return with error
			// so that it can be retried
			return &ctrl.Result{}, err
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

// Neet to be able to be called multiple times !!!
// Need to be idempotent
func (r *Reconciler) deleteStack(ctx context.Context, key types.NamespacedName, stack *stackv1beta3.Stack, conf *stackv1beta3.Configuration, log logr.Logger) error {
	log.Info("start deleting databases " + stack.Name)
	if err := delete.DeleteServiceData(conf, stack.Name, log); err != nil {
		log.Error(err, "Error during deleting databases")
	}

	log.Info("start deleting brokers subjects " + stack.Name)
	if err := delete.DeleteBrokersData(conf, stack.Name, []string{"ledger", "payments"}, log); err != nil {
		log.Error(err, "Error during deleting brokers subjects")
	}

	return nil
}

func (r *Reconciler) backupStack(ctx context.Context, conf *stackv1beta3.Configuration, stack *stackv1beta3.Stack, log logr.Logger) error {
	s3Client, err := s3.NewClient(
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

	bucket := "backups"
	storage := s3.NewS3Storage(s3Client, bucket)

	log.Info("start backup for " + stack.Name)
	if err := delete.BackupServicesData(conf, stack, storage, log); err != nil {
		log.Error(err, "Error during backups")
	}

	return nil
}
