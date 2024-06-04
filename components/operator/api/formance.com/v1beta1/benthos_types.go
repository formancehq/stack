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

// Batching allow to define custom batching configuration
type Batching struct {
	// Count indicates the number of messages that can be kept in memory before being flushed to ElasticSearch
	Count int `json:"count"`
	// Period indicates the maximum duration messages can be kept in memory before being flushed to ElasticSearch
	Period string `json:"period"`
}

type BenthosSpec struct {
	StackDependency `json:",inline"`
	DevProperties   `json:",inline"`
	//+optional
	ResourceProperties *corev1.ResourceRequirements `json:"resourceRequirements,omitempty"`
	//+optional
	Batching *Batching `json:"batching,omitempty"`
	//+optional
	InitContainers []corev1.Container `json:"initContainers"`
}

type BenthosStatus struct {
	Status `json:",inline"`
	//+optional
	ElasticSearchURI *URI `json:"elasticSearchURI"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
//+kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"

// Benthos is the Schema for the benthos API
type Benthos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BenthosSpec   `json:"spec,omitempty"`
	Status BenthosStatus `json:"status,omitempty"`
}

func (in *Benthos) SetReady(b bool) {
	in.Status.Ready = b
}

func (in *Benthos) IsReady() bool {
	return in.Status.Ready
}

func (in *Benthos) SetError(s string) {
	in.Status.Info = s
}

func (a Benthos) GetStack() string {
	return a.Spec.Stack
}

func (in *Benthos) GetConditions() *Conditions {
	return &in.Status.Conditions
}

//+kubebuilder:object:root=true

// BenthosList contains a list of Benthos
type BenthosList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Benthos `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Benthos{}, &BenthosList{})
}
