package core

import (
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func MapObjectToReconcileRequests[T client.Object](items ...T) []reconcile.Request {
	return collectionutils.Map(items, func(gateway T) reconcile.Request {
		return reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      gateway.GetName(),
				Namespace: gateway.GetNamespace(),
			},
		}
	})
}
