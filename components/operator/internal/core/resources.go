package core

import (
	"dario.cat/mergo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func GetResourcesRequirementsWithDefault(
	resourceProperties *corev1.ResourceRequirements,
	defaultResources corev1.ResourceRequirements,
) corev1.ResourceRequirements {
	if resourceProperties == nil {
		return defaultResources
	}

	if err := mergo.Merge(&defaultResources, resourceProperties); err != nil {
		panic(err)
	}

	return defaultResources
}

func ResourceSizeSmall() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("100m"),
			corev1.ResourceMemory: resource.MustParse("50Mi"),
		},
	}
}
