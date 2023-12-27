/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type OtlpSpec struct {
	// +optional
	Endpoint string `json:"endpoint,omitempty"`
	// +optional
	Port int32 `json:"port,omitempty"`
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
	Otlp *OtlpSpec `json:"otlp,omitempty"`
}

type MetricsSpec struct {
	// +optional
	Otlp *OtlpSpec `json:"otlp,omitempty"`
}

// OpenTelemetryConfigurationSpec defines the desired state of OpenTelemetryConfiguration
type OpenTelemetryConfigurationSpec struct {
	Stack string `json:"stack"`
	// +optional
	Traces *TracesSpec `json:"traces,omitempty"`
	// +optional
	Metrics *MetricsSpec `json:"metrics,omitempty"`
}

// OpenTelemetryConfigurationStatus defines the observed state of OpenTelemetryConfiguration
type OpenTelemetryConfigurationStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// OpenTelemetryConfiguration is the Schema for the opentelemetrytraces API
type OpenTelemetryConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenTelemetryConfigurationSpec   `json:"spec,omitempty"`
	Status OpenTelemetryConfigurationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OpenTelemetryConfigurationList contains a list of OpenTelemetryTraces
type OpenTelemetryConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenTelemetryConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenTelemetryConfiguration{}, &OpenTelemetryConfigurationList{})
}
