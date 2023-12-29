package core

import (
	"fmt"

	"k8s.io/apimachinery/pkg/types"
)

func GetObjectName(stack, name string) string {
	return fmt.Sprintf("%s-%s", stack, name)
}

func GetNamespacedResourceName(namespace, name string) types.NamespacedName {
	return types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
}

func GetResourceName(name string) types.NamespacedName {
	return types.NamespacedName{
		Name: name,
	}
}
