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

type TemporalTLSConfig struct {
	// +optional
	CRT string `json:"crt"`
	// +optional
	Key string `json:"key"`
	// +optional
	SecretName string `json:"secretName"`
}

// TemporalConfigurationSpec defines the desired state of TemporalConfiguration
type TemporalConfigurationSpec struct {
	Address   string            `json:"address"`
	Namespace string            `json:"namespace"`
	TLS       TemporalTLSConfig `json:"tls"`
}

// TemporalConfigurationStatus defines the observed state of TemporalConfiguration
type TemporalConfigurationStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// TemporalConfiguration is the Schema for the temporalconfigurations API
type TemporalConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TemporalConfigurationSpec   `json:"spec,omitempty"`
	Status TemporalConfigurationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TemporalConfigurationList contains a list of TemporalConfiguration
type TemporalConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TemporalConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TemporalConfiguration{}, &TemporalConfigurationList{})
}
