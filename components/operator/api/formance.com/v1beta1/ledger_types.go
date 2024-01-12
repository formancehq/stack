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

type DeploymentStrategy string

const (
	DeploymentStrategySingle                   = "single"
	DeploymentStrategyMonoWriterMultipleReader = "single-writer"
)

// LedgerSpec defines the desired state of Ledger
type LedgerSpec struct {
	CommonServiceProperties `json:",inline"`
	StackDependency         `json:",inline"`
	// +optional
	Auth *AuthConfig `json:"auth,omitempty"`
	//+optional
	DeploymentStrategy DeploymentStrategy `json:"deploymentStrategy,omitempty"`
	//+optional
	Service *ServiceConfiguration `json:"service,omitempty"`
}

// LedgerStatus defines the observed state of Ledger
type LedgerStatus struct {
	CommonStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
//+kubebuilder:printcolumn:name="Error",type=string,JSONPath=".status.error",description="Error"

// Ledger is the Schema for the ledgers API
type Ledger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LedgerSpec   `json:"spec,omitempty"`
	Status LedgerStatus `json:"status,omitempty"`
}

func (a Ledger) isEventPublisher() {}

func (a Ledger) GetStack() string {
	return a.Spec.Stack
}

func (a Ledger) IsDebug() bool {
	return a.Spec.Debug
}

func (a Ledger) IsDev() bool {
	return a.Spec.Dev
}

func (a *Ledger) SetCondition(condition Condition) {
	a.Status.SetCondition(condition)
}

//+kubebuilder:object:root=true

// LedgerList contains a list of Ledger
type LedgerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ledger `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Ledger{}, &LedgerList{})
}
