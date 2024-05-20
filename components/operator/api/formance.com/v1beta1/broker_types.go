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

// BrokerSpec defines the desired state of Broker
type BrokerSpec struct {
	StackDependency `json:",inline"`
}

// Mode defined how streams are created on the broker (mainly nats)
type Mode string

const (
	ModeOneStreamByService = "OneStreamByService"
	ModeOneStreamByStack   = "OneStreamByStack"
)

// BrokerStatus defines the observed state of Broker
type BrokerStatus struct {
	Status `json:",inline"`
	//+optional
	URI *URI `json:"uri,omitempty"`
	//+optional
	Mode Mode `json:"mode"`
	// Streams created when Mode == ModeOneStreamByService
	//+optional
	Streams []string `json:"streams,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Mode",type=string,JSONPath=".status.mode",description="Mode"
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Ready"
//+kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"

// Broker is the Schema for the brokers API
type Broker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BrokerSpec   `json:"spec,omitempty"`
	Status BrokerStatus `json:"status,omitempty"`
}

func (in *Broker) SetReady(ready bool) {
	in.Status.SetReady(ready)
}

func (in *Broker) IsReady() bool {
	return in.Status.Ready
}

func (in *Broker) SetError(s string) {
	in.Status.SetError(s)
}

func (in *Broker) GetStack() string {
	return in.Spec.GetStack()
}

func (in *Broker) GetConditions() *Conditions {
	return &in.Status.Conditions
}

//+kubebuilder:object:root=true

// BrokerList contains a list of Broker
type BrokerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Broker `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Broker{}, &BrokerList{})
}
