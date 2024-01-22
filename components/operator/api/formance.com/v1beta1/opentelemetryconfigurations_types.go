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
	ResourceAttributes map[string]string `json:"resourceAttributes,omitempty"`
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
type OpenTelemetryConfiguration struct {
	ConfigurationProperties `json:",inline"`
	// +optional
	Traces *TracesSpec `json:"traces,omitempty"`
	// +optional
	Metrics *MetricsSpec `json:"metrics,omitempty"`
}
