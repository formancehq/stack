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

// ResourceReferenceSpec defines the desired state of ResourceReference
type ResourceReferenceSpec struct {
	StackDependency  `json:",inline"`
	GroupVersionKind *metav1.GroupVersionKind `json:"gvk"`
	Name             string                   `json:"name"`
}

// ResourceReferenceStatus defines the observed state of ResourceReference
type ResourceReferenceStatus struct {
	CommonStatus `json:",inline"`
	//+optional
	SyncedResource string `json:"syncedResource,omitempty"`
	//+optional
	Hash string `json:"hash,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// ResourceReference is the Schema for the resourcereferences API
type ResourceReference struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceReferenceSpec   `json:"spec,omitempty"`
	Status ResourceReferenceStatus `json:"status,omitempty"`
}

func (in *ResourceReference) SetReady(b bool) {
	in.Status.SetReady(b)
}

func (in *ResourceReference) IsReady() bool {
	return in.Status.Ready
}

func (in *ResourceReference) SetError(s string) {
	in.Status.SetError(s)
}

func (in *ResourceReference) GetStack() string {
	return in.Spec.Stack
}

//+kubebuilder:object:root=true

// ResourceReferenceList contains a list of ResourceReference
type ResourceReferenceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceReference `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceReference{}, &ResourceReferenceList{})
}
