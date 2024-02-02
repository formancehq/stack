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

// BrokerTopicSpec defines the desired state of BrokerTopic
type BrokerTopicSpec struct {
	StackDependency `json:",inline"`
	//+required
	Service string `json:"service"`
}

// BrokerTopicStatus defines the observed state of BrokerTopic
type BrokerTopicStatus struct {
	CommonStatus `json:",inline"`
	//+optional
	URI *URI `json:"uri,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Ready"
//+kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"

// BrokerTopic is the Schema for the brokertopics API
type BrokerTopic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BrokerTopicSpec   `json:"spec,omitempty"`
	Status BrokerTopicStatus `json:"status,omitempty"`
}

func (in *BrokerTopic) SetReady(b bool) {
	in.Status.Ready = b
}

func (in *BrokerTopic) IsReady() bool {
	return in.Status.Ready
}

func (in *BrokerTopic) SetError(s string) {
	in.Status.Info = s
}

func (a *BrokerTopic) GetStack() string {
	return a.Spec.Stack
}

func (a *BrokerTopic) isResource() {}

//+kubebuilder:object:root=true

// BrokerTopicList contains a list of BrokerTopic
type BrokerTopicList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BrokerTopic `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BrokerTopic{}, &BrokerTopicList{})
}
