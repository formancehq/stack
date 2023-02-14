// +kubebuilder:object:generate=true
package v1beta1

import (
	v1 "k8s.io/api/core/v1"
)

type MonitoringSpec struct {
	// +optional
	Traces *TracesSpec `json:"traces,omitempty"`
}

type TracesOtlpSpec struct {
	// +optional
	Endpoint string `json:"endpoint,omitempty"`
	// +optional
	EndpointFrom *v1.EnvVarSource `json:"endpointFrom,omitempty"`
	// +optional
	Port int32 `json:"port,omitempty"`
	// +optional
	PortFrom *v1.EnvVarSource `json:"portFrom,omitempty"`
	// +optional
	Insecure bool `json:"insecure,omitempty"`
	// +kubebuilder:validation:Enum:={grpc,http}
	// +kubebuilder:validation:default:=grpc
	// +optional
	Mode string `json:"mode,omitempty"`
	// +optional
	ResourceAttributes string `json:"resourceAttributes,omitempty"`
}

type TracesSpec struct {
	// +optional
	Otlp *TracesOtlpSpec `json:"otlp,omitempty"`
}
