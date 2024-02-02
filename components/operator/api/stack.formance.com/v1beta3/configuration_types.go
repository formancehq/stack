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

package v1beta3

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigurationServicesSpec define all existing services for a stack.
// Fields order is important.
// For example, auth must be defined later as other services create static auth clients which must be used by auth.
type ConfigurationServicesSpec struct {
	Control        ControlSpec        `json:"control,omitempty"`
	Ledger         LedgerSpec         `json:"ledger,omitempty"`
	Payments       PaymentsSpec       `json:"payments,omitempty"`
	Reconciliation ReconciliationSpec `json:"reconciliation,omitempty"`
	Webhooks       WebhooksSpec       `json:"webhooks,omitempty"`
	Wallets        WalletsSpec        `json:"wallets,omitempty"`
	Orchestration  OrchestrationSpec  `json:"orchestration,omitempty"`
	Search         SearchSpec         `json:"search,omitempty"`
	Auth           AuthSpec           `json:"auth,omitempty"`

	// +optional
	Gateway  GatewaySpec  `json:"gateway,omitempty"`
	Stargate StargateSpec `json:"stargate,omitempty"`
}

func (in *ConfigurationServicesSpec) getDisabledProperty(service string) *bool {
	valueOf := reflect.ValueOf(in).Elem().FieldByName(strcase.ToCamel(service))
	if valueOf.IsValid() {
		_, ok := valueOf.Type().FieldByName("CommonServiceProperties")
		if !ok {
			return nil
		}
		commonServiceProperties := valueOf.FieldByName("CommonServiceProperties")
		properties := commonServiceProperties.Interface().(CommonServiceProperties)
		return properties.Disabled
	}
	return nil
}

func (in *ConfigurationServicesSpec) IsExplicitlyDisabled(service string) bool {
	disabled := in.getDisabledProperty(service)
	return disabled != nil && *disabled
}

func (in *ConfigurationServicesSpec) IsExplicitlyEnabled(service string) bool {
	disabled := in.getDisabledProperty(service)
	return disabled != nil && !*disabled
}

// TODO: Handle validation
type TemporalTLSConfig struct {
	// +optional
	CRT string `json:"crt"`
	// +optional
	Key string `json:"key"`
	// +optional
	SecretName string `json:"secretName"`
}

type TemporalConfig struct {
	Address   string            `json:"address"`
	Namespace string            `json:"namespace"`
	TLS       TemporalTLSConfig `json:"tls,omitempty"`
}

type CollectorConfig struct {
	Broker `json:",inline"`
	Topic  string `json:"topic"`
}

type Broker struct {
	// +optional
	Kafka *KafkaConfig `json:"kafka,omitempty"`
	// +optional
	Nats *NatsConfig `json:"nats,omitempty"`
}

type MonitoringSpec struct {
	// +optional
	Traces *TracesSpec `json:"traces,omitempty"`
	// +optional
	Metrics *MetricsSpec `json:"metrics,omitempty"`
}

type AnnotationsServicesSpec struct {
	// +optional
	Service map[string]string `json:"service,omitempty"`
}

type OtlpSpec struct {
	// +optional
	Endpoint string `json:"endpoint,omitempty"`
	// +optional
	Port int32 `json:"port,omitempty"`
	// +optional
	Insecure bool `json:"insecure,omitempty"`
	// +kubebuilder:validation:Enum:={grpc,http}
	// +kubebuilder:validation:default:=grpc
	// +optional
	Mode string `json:"mode,omitempty"`
	// +optional
	ResourceAttributes string `json:"resourceAttributes,omitempty"`
}

type TracesSpec struct {
	// +optional
	Otlp *OtlpSpec `json:"otlp,omitempty"`
}

type MetricsSpec struct {
	// +optional
	Otlp *OtlpSpec `json:"otlp,omitempty"`
}

type RegistryConfig struct {
	Endpoint string `json:"endpoint"`
}

type ConfigurationSpec struct {
	Services ConfigurationServicesSpec `json:"services"`
	Broker   Broker                    `json:"broker"`
	// +optional
	Monitoring *MonitoringSpec `json:"monitoring,omitempty"`

	// +optional
	Auth *AuthConfig `json:"auth,omitempty"`

	// +optional
	Ingress  IngressGlobalConfig `json:"ingress,omitempty"`
	Temporal TemporalConfig      `json:"temporal"`
	// LightMode is experimental and indicate we want monopods
	// +optional
	LightMode bool `json:"light,omitempty"`
	// +optional
	Registries map[string]RegistryConfig `json:"registries,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// Configuration is the Schema for the configurations API
type Configuration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigurationSpec `json:"spec,omitempty"`
	Status Status            `json:"status,omitempty"`
}

func (c *Configuration) Default() {}

func (*Configuration) Hub() {}

func (c Configuration) Validate() error {
	if c.Spec.Broker.Kafka == nil && c.Spec.Broker.Nats == nil {
		return errors.New("either 'kafka' or 'nats' is required")
	}
	return nil
}

func (c *Configuration) ResolveImage(image string) string {
	parts := strings.Split(image, ":")
	repository := parts[0]
	repositoryParts := strings.SplitN(repository, "/", 2)
	var (
		registry, path string
	)
	if len(repositoryParts) == 1 {
		registry = "docker.io"
		path = repository
	} else {
		registry = repositoryParts[0]
		path = repositoryParts[1]
	}
	if config, ok := c.Spec.Registries[registry]; ok && config.Endpoint != "" {
		return fmt.Sprintf("%s/%s:%s", config.Endpoint, path, parts[1])
	}
	return image
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
