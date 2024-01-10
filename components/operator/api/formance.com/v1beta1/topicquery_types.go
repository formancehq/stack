/*
Copyright 2022.

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

// TopicQuerySpec defines the desired state of TopicQuery
type TopicQuerySpec struct {
	StackDependency `json:",inline"`
	Service         string `json:"service"`
	QueriedBy       string `json:"queriedBy"`
}

// TopicQueryStatus defines the observed state of TopicQuery
type TopicQueryStatus struct {
	CommonStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Ready"
//+kubebuilder:printcolumn:name="Error",type=string,JSONPath=".status.error",description="Error"

// TopicQuery is the Schema for the topicqueries API
type TopicQuery struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TopicQuerySpec   `json:"spec,omitempty"`
	Status TopicQueryStatus `json:"status,omitempty"`
}

func (a TopicQuery) GetStack() string {
	return a.Spec.Stack
}

//+kubebuilder:object:root=true

// TopicQueryList contains a list of TopicQuery
type TopicQueryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TopicQuery `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TopicQuery{}, &TopicQueryList{})
}
