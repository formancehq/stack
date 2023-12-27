package internal

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func GetResourcesWithDefault(
	resourceProperties *v1beta1.ResourceProperties,
	defaultResources corev1.ResourceRequirements,
) corev1.ResourceRequirements {
	if resourceProperties == nil {
		return defaultResources
	}

	resources := defaultResources
	if resourceProperties.Request != nil {
		if resources.Requests == nil {
			resources.Requests = make(corev1.ResourceList)
		}

		if resourceProperties.Request.Cpu != "" {
			resources.Requests[corev1.ResourceCPU] = resource.MustParse(resourceProperties.Request.Cpu)
		}

		if resourceProperties.Request.Memory != "" {
			resources.Requests[corev1.ResourceMemory] = resource.MustParse(resourceProperties.Request.Memory)
		}
	}

	if resourceProperties.Limits != nil {
		if resources.Limits == nil {
			resources.Limits = make(corev1.ResourceList)
		}

		if resourceProperties.Limits.Cpu != "" {
			resources.Limits[corev1.ResourceCPU] = resource.MustParse(resourceProperties.Limits.Cpu)
		}

		if resourceProperties.Limits.Memory != "" {
			resources.Limits[corev1.ResourceMemory] = resource.MustParse(resourceProperties.Limits.Memory)
		}
	}

	return resources
}

func ResourceSizeSmall() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("100m"),
			corev1.ResourceMemory: resource.MustParse("50Mi"),
		},
	}
}

func ResourceSizeMedium() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("100m"),
			corev1.ResourceMemory: resource.MustParse("512Mi"),
		},
	}
}
