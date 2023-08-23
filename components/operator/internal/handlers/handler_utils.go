package handlers

import (
	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func getResourcesWithDefault(
	resourceProperties *stackv1beta3.ResourceProperties,
	defaultResources v1.ResourceRequirements,
) v1.ResourceRequirements {
	if resourceProperties == nil {
		return defaultResources
	}

	resources := defaultResources
	if resourceProperties.Request != nil {
		if resourceProperties.Request.Cpu != "" {
			resources.Requests[v1.ResourceCPU] = resource.MustParse(resourceProperties.Request.Cpu)
		}

		if resourceProperties.Request.Memory != "" {
			resources.Requests[v1.ResourceMemory] = resource.MustParse(resourceProperties.Request.Memory)
		}
	}

	if resourceProperties.Limits != nil {
		if resourceProperties.Limits.Cpu != "" {
			resources.Limits[v1.ResourceCPU] = resource.MustParse(resourceProperties.Limits.Cpu)
		}

		if resourceProperties.Limits.Memory != "" {
			resources.Limits[v1.ResourceMemory] = resource.MustParse(resourceProperties.Limits.Memory)
		}
	}

	return resources
}
