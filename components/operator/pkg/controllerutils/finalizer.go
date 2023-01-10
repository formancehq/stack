package controllerutils

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func IsDeleted(meta client.Object) bool {
	return meta.GetDeletionTimestamp() != nil && !meta.GetDeletionTimestamp().IsZero()
}

type Finalizer struct {
	name string
}

func (f *Finalizer) logger(ctx context.Context) logr.Logger {
	return log.FromContext(ctx, "finalizer", f.name)
}

func (f *Finalizer) Add(ctx context.Context, client client.Client, object client.Object) error {
	controllerutil.AddFinalizer(object, f.name)
	return client.Update(ctx, object)
}

func (f *Finalizer) Remove(ctx context.Context, client client.Client, object client.Object) error {
	controllerutil.RemoveFinalizer(object, f.name)
	return client.Update(ctx, object)
}

func (f *Finalizer) IsPresent(ob client.Object) bool {
	return controllerutil.ContainsFinalizer(ob, f.name)
}

func (f *Finalizer) Handle(ctx context.Context, client client.Client, ob client.Object, fn func() error) (bool, error) {
	if IsDeleted(ob) {
		if !f.IsPresent(ob) {
			f.logger(ctx).Info("Resource deleted and finalizer not present")
			return true, nil
		}
		if err := fn(); err != nil {
			f.logger(ctx).Error(err, "Resource deleted and finalizer callback return error")
			return true, err
		}
		if err := f.Remove(ctx, client, ob); err != nil {
			f.logger(ctx).Error(err, "Resource deleted, callback was called, but an error occurred removing the finalizer")
			return true, err
		}
		f.logger(ctx).Info("Resource deleted, callback was called and finalizer removed")
		return true, nil
	} else {
		// Assert finalizer is properly installed on the object
		if err := f.AssertIsInstalled(ctx, client, ob); err != nil {
			return false, err
		}
	}
	return false, nil
}

func (f *Finalizer) AssertIsInstalled(ctx context.Context, client client.Client, ob client.Object) error {
	if !f.IsPresent(ob) {
		f.logger(ctx).Info("Install finalizer")
		return f.Add(ctx, client, ob)
	}
	return nil
}

func New(name string) *Finalizer {
	return &Finalizer{
		name: name,
	}
}
