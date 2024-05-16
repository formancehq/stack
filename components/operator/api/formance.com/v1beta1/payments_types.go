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

// PaymentsSpec defines the desired state of Payments
type PaymentsSpec struct {
	StackDependency  `json:",inline"`
	ModuleProperties `json:",inline"`
	// +optional
	EncryptionKey string `json:"encryptionKey"`
	// +optional
	Auth *AuthConfig `json:"auth,omitempty"`
}

// PaymentsStatus defines the observed state of Payments
type PaymentsStatus struct {
	ModuleStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
// +kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"
// +kubebuilder:metadata:labels=formance.com/kind=module
// Payments is the Schema for the payments API
type Payments struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PaymentsSpec   `json:"spec,omitempty"`
	Status PaymentsStatus `json:"status,omitempty"`
}

func (in *Payments) IsEE() bool {
	return false
}

func (in *Payments) SetReady(b bool) {
	in.Status.Ready = b
}

func (in *Payments) IsReady() bool {
	return in.Status.Ready
}

func (in *Payments) SetError(s string) {
	in.Status.Info = s
}

func (in *Payments) GetConditions() []Condition {
	return in.Status.Conditions
}

func (in *Payments) GetVersion() string {
	return in.Spec.Version
}

func (a Payments) isEventPublisher() {}

func (a Payments) GetStack() string {
	return a.Spec.Stack
}

func (a Payments) IsDebug() bool {
	return a.Spec.Debug
}

func (a Payments) IsDev() bool {
	return a.Spec.Dev
}

func (a *Payments) SetCondition(condition Condition) {
	a.Status.SetCondition(condition)
}

//+kubebuilder:object:root=true

// PaymentsList contains a list of Payments
type PaymentsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Payments `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Payments{}, &PaymentsList{})
}
