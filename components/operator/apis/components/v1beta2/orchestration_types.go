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

package v1beta2

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TemporalTLSConfig struct {
	CRT string `json:"crt"`
	Key string `json:"key"`
}

type TemporalConfig struct {
	Address   string            `json:"address"`
	Namespace string            `json:"namespace"`
	TaskQueue string            `json:"taskQueue"`
	TLS       TemporalTLSConfig `json:"tls"`
}

func (in *TemporalConfig) Env() []corev1.EnvVar {
	return []corev1.EnvVar{
		Env("TEMPORAL_ADDRESS", in.Address),
		Env("TEMPORAL_NAMESPACE", in.Namespace),
		Env("TEMPORAL_SSL_CLIENT_KEY", in.TLS.Key),
		Env("TEMPORAL_SSL_CLIENT_CERT", in.TLS.CRT),
		Env("TEMPORAL_TASK_QUEUE", in.TaskQueue),
	}
}

// OrchestrationSpec defines the desired state of Orchestration
type OrchestrationSpec struct {
	CommonServiceProperties `json:",inline"`
	Scalable                `json:",inline"`

	// +optional
	Postgres PostgresConfigCreateDatabase `json:"postgres"`
	// +optional
	Monitoring *MonitoringSpec `json:"monitoring"`
	// +optional
	Collector *CollectorConfig `json:"collector"`

	Auth     OAuth2ClientConfiguration `json:"auth"`
	StackURL string                    `json:"stackUrl"`
	Temporal TemporalConfig            `json:"temporal"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector
//+kubebuilder:storageversion

// Orchestration is the Schema for the orchestrations API
type Orchestration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec OrchestrationSpec `json:"spec"`
	// +optional
	Status Status `json:"status"`
}

//+kubebuilder:object:root=true

// OrchestrationList contains a list of Orchestration
type OrchestrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Orchestration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Orchestration{}, &OrchestrationList{})
}
