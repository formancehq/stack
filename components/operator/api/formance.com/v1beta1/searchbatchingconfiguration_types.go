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

// SearchBatchingConfigurationSpec defines the desired state of SearchBatchingConfiguration
type SearchBatchingConfigurationSpec struct {
	ConfigurationProperties `json:",inline"`
	Batching                Batching `json:"batching"`
}

// SearchBatchingConfigurationStatus defines the observed state of SearchBatchingConfiguration
type SearchBatchingConfigurationStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// SearchBatchingConfiguration is the Schema for the searchbatchingconfigurations API
type SearchBatchingConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SearchBatchingConfigurationSpec   `json:"spec,omitempty"`
	Status SearchBatchingConfigurationStatus `json:"status,omitempty"`
}

func (in *SearchBatchingConfiguration) GetStacks() []string {
	return in.Spec.Stacks
}

func (in *SearchBatchingConfiguration) IsWildcard() bool {
	return in.Spec.ApplyOnAllStacks
}

var _ ConfigurationObject = &SearchBatchingConfiguration{}

//+kubebuilder:object:root=true

// SearchBatchingConfigurationList contains a list of SearchBatchingConfiguration
type SearchBatchingConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SearchBatchingConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SearchBatchingConfiguration{}, &SearchBatchingConfigurationList{})
}
