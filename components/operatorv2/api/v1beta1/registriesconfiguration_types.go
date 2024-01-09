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

type RegistryConfigurationSpec struct {
	Endpoint string `json:"endpoint"`
}

// RegistriesConfigurationSpec defines the desired state of RegistriesConfiguration
type RegistriesConfigurationSpec struct {
	Registries map[string]RegistryConfigurationSpec `json:"registries"`
}

// RegistriesConfigurationStatus defines the observed state of RegistriesConfiguration
type RegistriesConfigurationStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// RegistriesConfiguration is the Schema for the registries API
type RegistriesConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RegistriesConfigurationSpec   `json:"spec,omitempty"`
	Status RegistriesConfigurationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RegistriesConfigurationList contains a list of RegistriesConfiguration
type RegistriesConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RegistriesConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RegistriesConfiguration{}, &RegistriesConfigurationList{})
}
