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

type KafkaSASLConfig struct {
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	Mechanism    string `json:"mechanism"`
	ScramSHASize string `json:"scramSHASize"`
}

type KafkaConfig struct {
	Brokers []string `json:"brokers"`
	// +optional
	TLS bool `json:"tls"`
	// +optional
	SASL *KafkaSASLConfig `json:"sasl,omitempty"`
}

type NatsConfig struct {
	URL string `json:"url"`
	// +kubebuilder:default:=1
	// +optional
	Replicas int `json:"replicas"`
}

// BrokerConfigSpec defines the desired state of BrokerConfig
type BrokerConfigSpec struct {
	// +optional
	Kafka *KafkaConfig `json:"kafka,omitempty"`
	// +optional
	Nats *NatsConfig `json:"nats,omitempty"`
}

// BrokerConfigStatus defines the observed state of BrokerConfig
type BrokerConfigStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// BrokerConfiguration is the Schema for the brokerconfigurations API
type BrokerConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BrokerConfigSpec   `json:"spec,omitempty"`
	Status BrokerConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BrokerConfigurationList contains a list of BrokerConfig
type BrokerConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BrokerConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BrokerConfiguration{}, &BrokerConfigurationList{})
}
