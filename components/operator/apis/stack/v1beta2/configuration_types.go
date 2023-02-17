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
	"encoding/json"

	"github.com/formancehq/operator/apis/components/v1beta2"
	"github.com/formancehq/operator/apis/stack/v1beta3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

type ConfigurationServicesSpec struct {
	Auth           AuthSpec           `json:"auth,omitempty"`
	Control        ControlSpec        `json:"control,omitempty"`
	Ledger         LedgerSpec         `json:"ledger,omitempty"`
	Payments       PaymentsSpec       `json:"payments,omitempty"`
	Search         SearchSpec         `json:"search,omitempty"`
	Webhooks       WebhooksSpec       `json:"webhooks,omitempty"`
	Wallets        WalletsSpec        `json:"wallets,omitempty"`
	Orchestration  OrchestrationSpec  `json:"orchestration,omitempty"`
	Counterparties CounterpartiesSpec `json:"counterparties,omitempty"`
}

type TemporalConfig struct {
	Address   string                    `json:"address"`
	Namespace string                    `json:"namespace"`
	TLS       v1beta2.TemporalTLSConfig `json:"tls"`
}

type ConfigurationSpec struct {
	Services ConfigurationServicesSpec `json:"services"`
	Kafka    v1beta2.KafkaConfig       `json:"kafka"`
	// +optional
	Monitoring *v1beta2.MonitoringSpec `json:"monitoring,omitempty"`
	// +optional
	Ingress  IngressGlobalConfig `json:"ingress,omitempty"`
	Temporal TemporalConfig      `json:"temporal"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// Configuration is the Schema for the configurations API
type Configuration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigurationSpec `json:"spec,omitempty"`
	Status v1beta2.Status    `json:"status,omitempty"`
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
		cfg := v1beta2.KafkaConfig{}
		if err := json.Unmarshal(kafkaConfigAsJSON, &cfg); err != nil {
			return err
		}

		config.Spec.Kafka = cfg
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
