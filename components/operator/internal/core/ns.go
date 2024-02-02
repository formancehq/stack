package core

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func WithTemporaryNamespace(ctx Context, name string, callback func() error) error {
	ns, _, err := CreateOrUpdate[*v1.Namespace](ctx, types.NamespacedName{
		Name: name,
	})
	if err != nil {
		return err
	}

	if err := callback(); err != nil {
		return err
	}

	if err := ctx.GetClient().Delete(ctx, ns); err != nil {
		return err
	}

	return nil
}
