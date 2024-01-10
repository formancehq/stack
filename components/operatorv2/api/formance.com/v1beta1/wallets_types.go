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

// WalletsSpec defines the desired state of Wallets
type WalletsSpec struct {
	CommonServiceProperties `json:",inline"`
	StackDependency         `json:",inline"`
	//+optional
	Service *ServiceConfiguration `json:"service,omitempty"`
}

// WalletsStatus defines the observed state of Wallets
type WalletsStatus struct {
	CommonStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// Wallets is the Schema for the wallets API
type Wallets struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WalletsSpec   `json:"spec,omitempty"`
	Status WalletsStatus `json:"status,omitempty"`
}

func (c *Wallets) SetStatus(status bool, error string) {
	c.Status.SetStatus(status, error)
}

func (a Wallets) GetStack() string {
	return a.Spec.Stack
}

//+kubebuilder:object:root=true

// WalletsList contains a list of Wallets
type WalletsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Wallets `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Wallets{}, &WalletsList{})
}
