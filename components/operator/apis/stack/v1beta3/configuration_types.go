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
	"reflect"
	"strings"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigurationServicesSpec define all existing services for a stack.
// Fields order is important.
// For example, auth must be defined later as other services create static auth clients which must be used by auth.
type ConfigurationServicesSpec struct {
	Control       ControlSpec       `json:"control,omitempty"`
	Ledger        LedgerSpec        `json:"ledger,omitempty"`
	Payments      PaymentsSpec      `json:"payments,omitempty"`
	Webhooks      WebhooksSpec      `json:"webhooks,omitempty"`
	Wallets       WalletsSpec       `json:"wallets,omitempty"`
	Orchestration OrchestrationSpec `json:"orchestration,omitempty"`
	Search        SearchSpec        `json:"search,omitempty"`
	Auth          AuthSpec          `json:"auth,omitempty"`

	// +optional
	Gateway  GatewaySpec  `json:"gateway,omitempty"`
	Stargate StargateSpec `json:"stargate,omitempty"`
}

func (in *ConfigurationServicesSpec) List() []string {
	valueOf := reflect.ValueOf(*in)
	ret := make([]string, 0)
	for i := 0; i < valueOf.Type().NumField(); i++ {
		ret = append(ret, strings.ToLower(valueOf.Type().Field(i).Name))
	}
	return ret
}

type TemporalTLSConfig struct {
	CRT string `json:"crt"`
	Key string `json:"key"`
}

type TemporalConfig struct {
	Address   string            `json:"address"`
	Namespace string            `json:"namespace"`
	TLS       TemporalTLSConfig `json:"tls"`
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

type S3SecretConfig struct {
	AccessKey string `json:"accessKey,omitempty"`
	SecretKey string `json:"secretKey,omitempty"`
}
type S3ConfigSpec struct {
	Endpoint string `json:"endpoint"`
	Bucket   string `json:"bucket"`

	// +optional
	Region string `json:"region,omitempty"`

	// +optional
	S3SecretConfig *S3SecretConfig `json:"secret,omitempty"`

	// +optional
	ForceStylePath bool `json:"forceStylePath,omitempty"`

	// +optional
	Insecure bool `json:"insecure,omitempty"`
}
type ConfigurationSpec struct {
	Services ConfigurationServicesSpec `json:"services"`
	Broker   Broker                    `json:"broker"`
	// +optional
	Monitoring *MonitoringSpec `json:"monitoring,omitempty"`

	// +optional
	S3 *S3ConfigSpec `json:"s3,omitempty"`

	// +optional
	Ingress  IngressGlobalConfig `json:"ingress,omitempty"`
	Temporal TemporalConfig      `json:"temporal"`
	// LightMode is experimental and indicate we want monopods
	// +optional
	LightMode bool `json:"light,omitempty"`
}

func (in *ConfigurationSpec) GetServices() []string {
	return in.Services.List()
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
