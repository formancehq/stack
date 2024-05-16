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
	ModuleProperties `json:",inline"`
	StackDependency  `json:",inline"`
	// +optional
	Auth *AuthConfig `json:"auth,omitempty"`
}

// WalletsStatus defines the observed state of Wallets
type WalletsStatus struct {
	ModuleStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
// +kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"
// +kubebuilder:metadata:labels=formance.com/kind=module
// Wallets is the Schema for the wallets API
type Wallets struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WalletsSpec   `json:"spec,omitempty"`
	Status WalletsStatus `json:"status,omitempty"`
}

func (in *Wallets) IsEE() bool {
	return false
}

func (in *Wallets) SetReady(b bool) {
	in.Status.Ready = b
}

func (in *Wallets) IsReady() bool {
	return in.Status.Ready
}

func (in *Wallets) SetError(s string) {
	in.Status.Info = s
}

func (in *Wallets) GetConditions() []Condition {
	return in.Status.Conditions
}

func (in *Wallets) GetVersion() string {
	return in.Spec.Version
}

func (a Wallets) GetStack() string {
	return a.Spec.Stack
}

func (a *Wallets) SetCondition(condition Condition) {
	a.Status.SetCondition(condition)
}

func (a Wallets) IsDebug() bool {
	return a.Spec.Debug
}

func (a Wallets) IsDev() bool {
	return a.Spec.Dev
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
