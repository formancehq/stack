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

// BenthosStreamSpec defines the desired state of BenthosStream
type BenthosStreamSpec struct {
	StackDependency `json:",inline"`
	Data            string `json:"data"`
	Name            string `json:"name"`
}

// BenthosStreamStatus defines the observed state of BenthosStream
type BenthosStreamStatus struct {
	Status `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"

// BenthosStream is the Schema for the benthosstreams API
type BenthosStream struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BenthosStreamSpec   `json:"spec,omitempty"`
	Status BenthosStreamStatus `json:"status,omitempty"`
}

func (in *BenthosStream) SetReady(b bool) {
	in.Status.Ready = b
}

func (in *BenthosStream) IsReady() bool {
	return in.Status.Ready
}

func (in *BenthosStream) SetError(s string) {
	in.Status.Info = s
}

func (a BenthosStream) GetStack() string {
	return a.Spec.Stack
}

func (in *BenthosStream) GetConditions() *Conditions {
	return &in.Status.Conditions
}

//+kubebuilder:object:root=true

// BenthosStreamList contains a list of BenthosStream
type BenthosStreamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BenthosStream `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BenthosStream{}, &BenthosStreamList{})
}
