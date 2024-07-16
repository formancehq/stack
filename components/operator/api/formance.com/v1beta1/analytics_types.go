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

// AnalyticsSpec defines the desired state of Analytics
type AnalyticsSpec struct {
	ModuleProperties `json:",inline"`
	StackDependency  `json:",inline"`
}

// AnalyticsStatus defines the observed state of Analytics
type AnalyticsStatus struct {
	Status `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
//+kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"
//+kubebuilder:metadata:labels=formance.com/kind=module
//+kubebuilder:metadata:labels=formance.com/is-ee=true

// Analytics is the Schema for the analytics API
type Analytics struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AnalyticsSpec   `json:"spec,omitempty"`
	Status AnalyticsStatus `json:"status,omitempty"`
}

func (in *Analytics) SetReady(b bool) {
	in.Status.SetReady(b)
}

func (in *Analytics) IsReady() bool {
	//TODO implement me

	return in.Status.Ready
}

func (in *Analytics) SetError(s string) {
	in.Status.SetError(s)
}

func (in *Analytics) GetConditions() *Conditions {
	return &in.Status.Conditions
}

func (in *Analytics) GetStack() string {
	//TODO implement me

	return in.Spec.Stack
}

func (in *Analytics) GetVersion() string {
	//TODO implement me
	return in.Spec.Version
}

func (in *Analytics) IsDebug() bool {
	//TODO implement me
	return in.Spec.Debug
}

func (in *Analytics) IsDev() bool {
	//TODO implement me
	return in.Spec.Dev
}

func (in *Analytics) IsEE() bool {
	//TODO implement me
	return true
}

//+kubebuilder:object:root=true

// AnalyticsList contains a list of Analytics
type AnalyticsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Analytics `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Analytics{}, &AnalyticsList{})
}
