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
	"encoding/json"
	"github.com/formancehq/operator/apis/stack/v1beta3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
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

type ConfigurationSpec struct {
	// +optional
	Monitoring *MonitoringSpec `json:"monitoring,omitempty"`
	// +optional
	Services ServicesSpec `json:"services,omitempty"`
	// +optional
	Auth *AuthSpec `json:"auth,omitempty"`
	// +optional
	Ingress IngressGlobalConfig `json:"ingress"`
	// +optional
	Kafka *KafkaConfig `json:"kafka"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// Configuration is the Schema for the configurations API
type Configuration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigurationSpec `json:"spec,omitempty"`
	Status Status            `json:"status,omitempty"`
}

func (config *Configuration) ConvertFrom(hubRaw conversion.Hub) error {
	hub := hubRaw.(*v1beta3.Configuration)
	specAsRaw, err := json.Marshal(hub.Spec)
	if err != nil {
		return err
	}
	*config = Configuration{
		ObjectMeta: hub.ObjectMeta,
		TypeMeta:   hub.TypeMeta,
	}
	config.APIVersion = GroupVersion.String()
	if err := json.Unmarshal(specAsRaw, &config.Spec); err != nil {
		return err
	}

	if hub.Spec.Broker.Kafka != nil {
		kafkaConfigAsJSON, err := json.Marshal(hub.Spec.Broker.Kafka)
		if err != nil {
			return err
		}
		cfg := KafkaConfig{}
		if err := json.Unmarshal(kafkaConfigAsJSON, &cfg); err != nil {
			return err
		}

		config.Spec.Kafka = &cfg
	}

	return nil
}

func (config *Configuration) ConvertTo(hubRaw conversion.Hub) error {

	hub := hubRaw.(*v1beta3.Configuration)
	specAsRaw, err := json.Marshal(config.Spec)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(specAsRaw, &hub.Spec); err != nil {
		return err
	}
	hub.ObjectMeta = config.ObjectMeta
	hub.TypeMeta = config.TypeMeta
	hub.APIVersion = v1beta3.GroupVersion.String()
	kafkaSpecAsJSON, err := json.Marshal(config.Spec.Kafka)
	if err != nil {
		return err
	}
	cfg := v1beta3.KafkaConfig{}
	if err := json.Unmarshal(kafkaSpecAsJSON, &cfg); err != nil {
		return err
	}
	hub.Spec.Broker.Kafka = &cfg
	hub.Spec.Services.Payments.EncryptionKey = "default-encryption-key"

	return nil
}

//+kubebuilder:object:root=true

// ConfigurationList contains a list of Configuration
type ConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Configuration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Configuration{}, &ConfigurationList{})
}
