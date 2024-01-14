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

// WebhooksSpec defines the desired state of Webhooks
type WebhooksSpec struct {
	StackDependency  `json:",inline"`
	ModuleProperties `json:",inline"`
	//+optional
	Service *ServiceConfiguration `json:"service,omitempty"`
	// +optional
	Auth *AuthConfig `json:"auth,omitempty"`
}

// WebhooksStatus defines the observed state of Webhooks
type WebhooksStatus struct {
	CommonStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
//+kubebuilder:printcolumn:name="Error",type=string,JSONPath=".status.error",description="Error"

// Webhooks is the Schema for the webhooks API
type Webhooks struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebhooksSpec   `json:"spec,omitempty"`
	Status WebhooksStatus `json:"status,omitempty"`
}

func (a Webhooks) GetStack() string {
	return a.Spec.Stack
}

func (a Webhooks) IsDebug() bool {
	return a.Spec.Debug
}

func (a Webhooks) IsDev() bool {
	return a.Spec.Dev
}

func (a *Webhooks) SetCondition(condition Condition) {
	a.Status.SetCondition(condition)
}

//+kubebuilder:object:root=true

// WebhooksList contains a list of Webhooks
type WebhooksList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Webhooks `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Webhooks{}, &WebhooksList{})
}
