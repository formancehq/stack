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

// BrokerConsumerSpec defines the desired state of BrokerConsumer
type BrokerConsumerSpec struct {
	StackDependency `json:",inline"`
	Services        []string `json:"services"`
	QueriedBy       string   `json:"queriedBy"`
}

// BrokerConsumerStatus defines the observed state of BrokerConsumer
type BrokerConsumerStatus struct {
	CommonStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Services",type=string,JSONPath=".spec.services",description="Listened services"
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Ready"
//+kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"

// BrokerConsumer is the Schema for the brokerconsumers API
type BrokerConsumer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BrokerConsumerSpec   `json:"spec,omitempty"`
	Status BrokerConsumerStatus `json:"status,omitempty"`
}

func (in *BrokerConsumer) SetReady(b bool) {
	in.Status.SetReady(b)
}

func (in *BrokerConsumer) IsReady() bool {
	return in.Status.Ready
}

func (in *BrokerConsumer) SetError(s string) {
	in.Status.SetError(s)
}

func (in *BrokerConsumer) GetStack() string {
	return in.Spec.Stack
}

//+kubebuilder:object:root=true

// BrokerConsumerList contains a list of BrokerConsumer
type BrokerConsumerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BrokerConsumer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BrokerConsumer{}, &BrokerConsumerList{})
}
