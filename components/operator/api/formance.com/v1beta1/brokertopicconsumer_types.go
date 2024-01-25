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

// BrokerTopicConsumerSpec defines the desired state of BrokerTopicConsumer
type BrokerTopicConsumerSpec struct {
	StackDependency `json:",inline"`
	Service         string `json:"service"`
	QueriedBy       string `json:"queriedBy"`
}

// BrokerTopicConsumerStatus defines the observed state of BrokerTopicConsumer
type BrokerTopicConsumerStatus struct {
	CommonStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Ready"
//+kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"

// BrokerTopicConsumer is the Schema for the brokertopicconsumers API
type BrokerTopicConsumer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BrokerTopicConsumerSpec   `json:"spec,omitempty"`
	Status BrokerTopicConsumerStatus `json:"status,omitempty"`
}

func (in *BrokerTopicConsumer) SetReady(b bool) {
	in.Status.Ready = b
}

func (in *BrokerTopicConsumer) IsReady() bool {
	return in.Status.Ready
}

func (in *BrokerTopicConsumer) SetError(s string) {
	in.Status.Info = s
}

func (a BrokerTopicConsumer) GetStack() string {
	return a.Spec.Stack
}

//+kubebuilder:object:root=true

// BrokerTopicConsumerList contains a list of BrokerTopicConsumer
type BrokerTopicConsumerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BrokerTopicConsumer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BrokerTopicConsumer{}, &BrokerTopicConsumerList{})
}
