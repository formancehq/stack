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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Batching struct {
	Count  int    `json:"count"`
	Period string `json:"period"`
}

// StreamProcessorSpec defines the desired state of StreamProcessor
type StreamProcessorSpec struct {
	StackDependency `json:",inline"`
	DevProperties   `json:",inline"`
	//+optional
	ResourceProperties *ResourceProperties `json:"resourceProperties,omitempty"`
	//+optional
	Batching *Batching `json:"batching,omitempty"`
	//+optional
	InitContainers []corev1.Container `json:"initContainers"`
}

// StreamProcessorStatus defines the observed state of StreamProcessor
type StreamProcessorStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// StreamProcessor is the Schema for the streamprocessors API
type StreamProcessor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StreamProcessorSpec   `json:"spec,omitempty"`
	Status StreamProcessorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StreamProcessorList contains a list of StreamProcessor
type StreamProcessorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StreamProcessor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StreamProcessor{}, &StreamProcessorList{})
}
