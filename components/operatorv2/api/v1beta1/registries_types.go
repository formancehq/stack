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

type RegistryConfiguration struct {
	Endpoint string `json:"endpoint"`
}

// RegistriesSpec defines the desired state of Registries
type RegistriesSpec struct {
	Registries map[string]RegistryConfiguration `json:"registries"`
}

// RegistriesStatus defines the observed state of Registries
type RegistriesStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// Registries is the Schema for the registries API
type Registries struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RegistriesSpec   `json:"spec,omitempty"`
	Status RegistriesStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RegistriesList contains a list of Registries
type RegistriesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Registries `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Registries{}, &RegistriesList{})
}
