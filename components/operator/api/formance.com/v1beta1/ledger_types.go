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
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type LockingStrategyRedisConfig struct {
	Uri string `json:"uri,omitempty"`
	// +optional
	TLS bool `json:"tls"`
	// +optional
	InsecureTLS bool `json:"insecure,omitempty"`
	// +optional
	Duration time.Duration `json:"duration,omitempty"`
	// +optional
	Retry time.Duration `json:"retry,omitempty"`
}

type LockingStrategy struct {
	// +kubebuilder:Enum:={memory,redis}
	// +kubebuilder:default:=memory
	// +optional
	Strategy string `json:"strategy,omitempty"`
	// +optional
	Redis *LockingStrategyRedisConfig `json:"redis"`
}

type DeploymentStrategy string

const (
	DeploymentStrategySingle                   = "single"
	DeploymentStrategyMonoWriterMultipleReader = "single-writer"
)

// LedgerSpec defines the desired state of Ledger
type LedgerSpec struct {
	ModuleProperties `json:",inline"`
	StackDependency  `json:",inline"`
	// +optional
	Auth *AuthConfig `json:"auth,omitempty"`
	//+optional
	DeploymentStrategy DeploymentStrategy `json:"deploymentStrategy,omitempty"`
	// Locking is intended for ledger v1 only
	//+optional
	Locking *LockingStrategy `json:"locking,omitempty"`
}

// LedgerStatus defines the observed state of Ledger
type LedgerStatus struct {
	ModuleStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Is ready"
//+kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"

// Ledger is the Schema for the ledgers API
type Ledger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LedgerSpec   `json:"spec,omitempty"`
	Status LedgerStatus `json:"status,omitempty"`
}

func (in *Ledger) IsEE() bool {
	return false
}

func (in *Ledger) IsReady() bool {
	return in.Status.Ready
}

func (in *Ledger) SetReady(b bool) {
	in.Status.Ready = b
}

func (in *Ledger) SetError(s string) {
	in.Status.Info = s
}

func (in *Ledger) GetConditions() []Condition {
	return in.Status.Conditions
}

func (in *Ledger) GetVersion() string {
	return in.Spec.Version
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
