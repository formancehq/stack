package controllerutils

import "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

func JustError[T any](t T, result controllerutil.OperationResult, err error) error {
	return err
}
