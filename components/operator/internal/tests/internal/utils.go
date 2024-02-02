package internal

import (
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func RandObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name: uuid.NewString(),
	}
}

func LoadResource(ns, name string, object client.Object) error {
	return Get(types.NamespacedName{
		Namespace: ns,
		Name:      name,
	}, object)
}
