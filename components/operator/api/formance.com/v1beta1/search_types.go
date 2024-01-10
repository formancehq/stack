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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SearchStreamProcessorSpec struct {
	//+optional
	ResourceRequirements *v1.ResourceRequirements `json:"resourceRequirements,omitempty"`
}

// SearchSpec defines the desired state of Search
type SearchSpec struct {
	StackDependency         `json:",inline"`
	CommonServiceProperties `json:",inline"`
	//+optional
	Batching *Batching `json:"batching,omitempty"`
	// +optional
	StreamProcessor *SearchStreamProcessorSpec `json:"streamProcessor,omitempty"`
	//+optional
	Service *ServiceConfiguration `json:"service,omitempty"`
	// +optional
	Auth *AuthConfig `json:"auth,omitempty"`
}

// SearchStatus defines the observed state of Search
type SearchStatus struct {
	CommonStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
//+kubebuilder:printcolumn:name="Error",type=string,JSONPath=".status.error",description="Error"

// Search is the Schema for the searches API
type Search struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SearchSpec   `json:"spec,omitempty"`
	Status SearchStatus `json:"status,omitempty"`
}

func (a Search) GetStack() string {
	return a.Spec.Stack
}

func (a *Search) SetCondition(condition Condition) {
	a.Status.SetCondition(condition)
}

//+kubebuilder:object:root=true

// SearchList contains a list of Search
type SearchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Search `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Search{}, &SearchList{})
}
