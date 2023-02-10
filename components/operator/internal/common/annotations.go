package common

import (
	"github.com/formancehq/operator/internal/controllerutils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const ReloaderAnnotationKey = "reloader.stakater.com/auto"

func WithReloaderAnnotations[T client.Object]() controllerutils.ObjectMutator[T] {
	return controllerutils.WithAnnotations[T](map[string]string{
		ReloaderAnnotationKey: "true",
	})
}
