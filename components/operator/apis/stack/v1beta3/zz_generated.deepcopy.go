//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta3

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthSpec) DeepCopyInto(out *AuthSpec) {
	*out = *in
	out.Postgres = in.Postgres
	if in.StaticClients != nil {
		in, out := &in.StaticClients, &out.StaticClients
		*out = make([]StaticClient, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthSpec.
func (in *AuthSpec) DeepCopy() *AuthSpec {
	if in == nil {
		return nil
	}
	out := new(AuthSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Batching) DeepCopyInto(out *Batching) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Batching.
func (in *Batching) DeepCopy() *Batching {
	if in == nil {
		return nil
	}
	out := new(Batching)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Broker) DeepCopyInto(out *Broker) {
	*out = *in
	if in.Kafka != nil {
		in, out := &in.Kafka, &out.Kafka
		*out = new(KafkaConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.Nats != nil {
		in, out := &in.Nats, &out.Nats
		*out = new(NatsConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Broker.
func (in *Broker) DeepCopy() *Broker {
	if in == nil {
		return nil
	}
	out := new(Broker)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClientConfiguration) DeepCopyInto(out *ClientConfiguration) {
	*out = *in
	if in.Description != nil {
		in, out := &in.Description, &out.Description
		*out = new(string)
		**out = **in
	}
	if in.RedirectUris != nil {
		in, out := &in.RedirectUris, &out.RedirectUris
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.PostLogoutRedirectUris != nil {
		in, out := &in.PostLogoutRedirectUris, &out.PostLogoutRedirectUris
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Scopes != nil {
		in, out := &in.Scopes, &out.Scopes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClientConfiguration.
func (in *ClientConfiguration) DeepCopy() *ClientConfiguration {
	if in == nil {
		return nil
	}
	out := new(ClientConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CollectorConfig) DeepCopyInto(out *CollectorConfig) {
	*out = *in
	in.Broker.DeepCopyInto(&out.Broker)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CollectorConfig.
func (in *CollectorConfig) DeepCopy() *CollectorConfig {
	if in == nil {
		return nil
	}
	out := new(CollectorConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonServiceProperties) DeepCopyInto(out *CommonServiceProperties) {
	*out = *in
	out.DevProperties = in.DevProperties
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonServiceProperties.
func (in *CommonServiceProperties) DeepCopy() *CommonServiceProperties {
	if in == nil {
		return nil
	}
	out := new(CommonServiceProperties)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Condition) DeepCopyInto(out *Condition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Condition.
func (in *Condition) DeepCopy() *Condition {
	if in == nil {
		return nil
	}
	out := new(Condition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Conditions) DeepCopyInto(out *Conditions) {
	{
		in := &in
		*out = make(Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Conditions.
func (in Conditions) DeepCopy() Conditions {
	if in == nil {
		return nil
	}
	out := new(Conditions)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigSource) DeepCopyInto(out *ConfigSource) {
	*out = *in
	if in.ConfigMapKeyRef != nil {
		in, out := &in.ConfigMapKeyRef, &out.ConfigMapKeyRef
		*out = new(v1.ConfigMapKeySelector)
		(*in).DeepCopyInto(*out)
	}
	if in.SecretKeyRef != nil {
		in, out := &in.SecretKeyRef, &out.SecretKeyRef
		*out = new(v1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigSource.
func (in *ConfigSource) DeepCopy() *ConfigSource {
	if in == nil {
		return nil
	}
	out := new(ConfigSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Configuration) DeepCopyInto(out *Configuration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Configuration.
func (in *Configuration) DeepCopy() *Configuration {
	if in == nil {
		return nil
	}
	out := new(Configuration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Configuration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigurationList) DeepCopyInto(out *ConfigurationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Configuration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigurationList.
func (in *ConfigurationList) DeepCopy() *ConfigurationList {
	if in == nil {
		return nil
	}
	out := new(ConfigurationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ConfigurationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigurationServicesSpec) DeepCopyInto(out *ConfigurationServicesSpec) {
	*out = *in
	out.Control = in.Control
	out.Ledger = in.Ledger
	out.Payments = in.Payments
	out.Webhooks = in.Webhooks
	out.Wallets = in.Wallets
	out.Orchestration = in.Orchestration
	in.Search.DeepCopyInto(&out.Search)
	in.Auth.DeepCopyInto(&out.Auth)
	out.Stargate = in.Stargate
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigurationServicesSpec.
func (in *ConfigurationServicesSpec) DeepCopy() *ConfigurationServicesSpec {
	if in == nil {
		return nil
	}
	out := new(ConfigurationServicesSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigurationSpec) DeepCopyInto(out *ConfigurationSpec) {
	*out = *in
	in.Services.DeepCopyInto(&out.Services)
	in.Broker.DeepCopyInto(&out.Broker)
	if in.Monitoring != nil {
		in, out := &in.Monitoring, &out.Monitoring
		*out = new(MonitoringSpec)
		(*in).DeepCopyInto(*out)
	}
	in.Ingress.DeepCopyInto(&out.Ingress)
	out.Temporal = in.Temporal
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigurationSpec.
func (in *ConfigurationSpec) DeepCopy() *ConfigurationSpec {
	if in == nil {
		return nil
	}
	out := new(ConfigurationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControlAuthentication) DeepCopyInto(out *ControlAuthentication) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControlAuthentication.
func (in *ControlAuthentication) DeepCopy() *ControlAuthentication {
	if in == nil {
		return nil
	}
	out := new(ControlAuthentication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControlSpec) DeepCopyInto(out *ControlSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControlSpec.
func (in *ControlSpec) DeepCopy() *ControlSpec {
	if in == nil {
		return nil
	}
	out := new(ControlSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseSpec) DeepCopyInto(out *DatabaseSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseSpec.
func (in *DatabaseSpec) DeepCopy() *DatabaseSpec {
	if in == nil {
		return nil
	}
	out := new(DatabaseSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DelegatedOIDCServerConfiguration) DeepCopyInto(out *DelegatedOIDCServerConfiguration) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DelegatedOIDCServerConfiguration.
func (in *DelegatedOIDCServerConfiguration) DeepCopy() *DelegatedOIDCServerConfiguration {
	if in == nil {
		return nil
	}
	out := new(DelegatedOIDCServerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevProperties) DeepCopyInto(out *DevProperties) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevProperties.
func (in *DevProperties) DeepCopy() *DevProperties {
	if in == nil {
		return nil
	}
	out := new(DevProperties)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSearchBasicAuthConfig) DeepCopyInto(out *ElasticSearchBasicAuthConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSearchBasicAuthConfig.
func (in *ElasticSearchBasicAuthConfig) DeepCopy() *ElasticSearchBasicAuthConfig {
	if in == nil {
		return nil
	}
	out := new(ElasticSearchBasicAuthConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSearchConfig) DeepCopyInto(out *ElasticSearchConfig) {
	*out = *in
	out.TLS = in.TLS
	if in.BasicAuth != nil {
		in, out := &in.BasicAuth, &out.BasicAuth
		*out = new(ElasticSearchBasicAuthConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSearchConfig.
func (in *ElasticSearchConfig) DeepCopy() *ElasticSearchConfig {
	if in == nil {
		return nil
	}
	out := new(ElasticSearchConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSearchTLSConfig) DeepCopyInto(out *ElasticSearchTLSConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSearchTLSConfig.
func (in *ElasticSearchTLSConfig) DeepCopy() *ElasticSearchTLSConfig {
	if in == nil {
		return nil
	}
	out := new(ElasticSearchTLSConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressConfig) DeepCopyInto(out *IngressConfig) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressConfig.
func (in *IngressConfig) DeepCopy() *IngressConfig {
	if in == nil {
		return nil
	}
	out := new(IngressConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressGlobalConfig) DeepCopyInto(out *IngressGlobalConfig) {
	*out = *in
	in.IngressConfig.DeepCopyInto(&out.IngressConfig)
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(IngressTLS)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressGlobalConfig.
func (in *IngressGlobalConfig) DeepCopy() *IngressGlobalConfig {
	if in == nil {
		return nil
	}
	out := new(IngressGlobalConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressTLS) DeepCopyInto(out *IngressTLS) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressTLS.
func (in *IngressTLS) DeepCopy() *IngressTLS {
	if in == nil {
		return nil
	}
	out := new(IngressTLS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaConfig) DeepCopyInto(out *KafkaConfig) {
	*out = *in
	if in.Brokers != nil {
		in, out := &in.Brokers, &out.Brokers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SASL != nil {
		in, out := &in.SASL, &out.SASL
		*out = new(KafkaSASLConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaConfig.
func (in *KafkaConfig) DeepCopy() *KafkaConfig {
	if in == nil {
		return nil
	}
	out := new(KafkaConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaSASLConfig) DeepCopyInto(out *KafkaSASLConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaSASLConfig.
func (in *KafkaSASLConfig) DeepCopy() *KafkaSASLConfig {
	if in == nil {
		return nil
	}
	out := new(KafkaSASLConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LedgerSpec) DeepCopyInto(out *LedgerSpec) {
	*out = *in
	out.Postgres = in.Postgres
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LedgerSpec.
func (in *LedgerSpec) DeepCopy() *LedgerSpec {
	if in == nil {
		return nil
	}
	out := new(LedgerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LockingStrategyRedisConfig) DeepCopyInto(out *LockingStrategyRedisConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LockingStrategyRedisConfig.
func (in *LockingStrategyRedisConfig) DeepCopy() *LockingStrategyRedisConfig {
	if in == nil {
		return nil
	}
	out := new(LockingStrategyRedisConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Migration) DeepCopyInto(out *Migration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Migration.
func (in *Migration) DeepCopy() *Migration {
	if in == nil {
		return nil
	}
	out := new(Migration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Migration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MigrationList) DeepCopyInto(out *MigrationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Migration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MigrationList.
func (in *MigrationList) DeepCopy() *MigrationList {
	if in == nil {
		return nil
	}
	out := new(MigrationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MigrationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MigrationSpec) DeepCopyInto(out *MigrationSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MigrationSpec.
func (in *MigrationSpec) DeepCopy() *MigrationSpec {
	if in == nil {
		return nil
	}
	out := new(MigrationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MigrationStatus) DeepCopyInto(out *MigrationStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MigrationStatus.
func (in *MigrationStatus) DeepCopy() *MigrationStatus {
	if in == nil {
		return nil
	}
	out := new(MigrationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MonitoringSpec) DeepCopyInto(out *MonitoringSpec) {
	*out = *in
	if in.Traces != nil {
		in, out := &in.Traces, &out.Traces
		*out = new(TracesSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MonitoringSpec.
func (in *MonitoringSpec) DeepCopy() *MonitoringSpec {
	if in == nil {
		return nil
	}
	out := new(MonitoringSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NatsConfig) DeepCopyInto(out *NatsConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NatsConfig.
func (in *NatsConfig) DeepCopy() *NatsConfig {
	if in == nil {
		return nil
	}
	out := new(NatsConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OrchestrationSpec) DeepCopyInto(out *OrchestrationSpec) {
	*out = *in
	out.DevProperties = in.DevProperties
	out.Postgres = in.Postgres
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OrchestrationSpec.
func (in *OrchestrationSpec) DeepCopy() *OrchestrationSpec {
	if in == nil {
		return nil
	}
	out := new(OrchestrationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PaymentsSpec) DeepCopyInto(out *PaymentsSpec) {
	*out = *in
	out.Postgres = in.Postgres
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PaymentsSpec.
func (in *PaymentsSpec) DeepCopy() *PaymentsSpec {
	if in == nil {
		return nil
	}
	out := new(PaymentsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgresConfig) DeepCopyInto(out *PostgresConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgresConfig.
func (in *PostgresConfig) DeepCopy() *PostgresConfig {
	if in == nil {
		return nil
	}
	out := new(PostgresConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchSpec) DeepCopyInto(out *SearchSpec) {
	*out = *in
	in.ElasticSearchConfig.DeepCopyInto(&out.ElasticSearchConfig)
	out.Batching = in.Batching
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchSpec.
func (in *SearchSpec) DeepCopy() *SearchSpec {
	if in == nil {
		return nil
	}
	out := new(SearchSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Stack) DeepCopyInto(out *Stack) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Stack.
func (in *Stack) DeepCopy() *Stack {
	if in == nil {
		return nil
	}
	out := new(Stack)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Stack) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StackAuthSpec) DeepCopyInto(out *StackAuthSpec) {
	*out = *in
	out.DelegatedOIDCServer = in.DelegatedOIDCServer
	if in.StaticClients != nil {
		in, out := &in.StaticClients, &out.StaticClients
		*out = make([]StaticClient, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StackAuthSpec.
func (in *StackAuthSpec) DeepCopy() *StackAuthSpec {
	if in == nil {
		return nil
	}
	out := new(StackAuthSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StackList) DeepCopyInto(out *StackList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Stack, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StackList.
func (in *StackList) DeepCopy() *StackList {
	if in == nil {
		return nil
	}
	out := new(StackList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StackList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StackSpec) DeepCopyInto(out *StackSpec) {
	*out = *in
	out.DevProperties = in.DevProperties
	in.Auth.DeepCopyInto(&out.Auth)
	if in.Stargate != nil {
		in, out := &in.Stargate, &out.Stargate
		*out = new(StackStargateConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StackSpec.
func (in *StackSpec) DeepCopy() *StackSpec {
	if in == nil {
		return nil
	}
	out := new(StackSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StackStargateConfig) DeepCopyInto(out *StackStargateConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StackStargateConfig.
func (in *StackStargateConfig) DeepCopy() *StackStargateConfig {
	if in == nil {
		return nil
	}
	out := new(StackStargateConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StackStatus) DeepCopyInto(out *StackStatus) {
	*out = *in
	in.Status.DeepCopyInto(&out.Status)
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make(map[string]map[string]int32, len(*in))
		for key, val := range *in {
			var outVal map[string]int32
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make(map[string]int32, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
			}
			(*out)[key] = outVal
		}
	}
	if in.StaticAuthClients != nil {
		in, out := &in.StaticAuthClients, &out.StaticAuthClients
		*out = make(map[string]StaticClient, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StackStatus.
func (in *StackStatus) DeepCopy() *StackStatus {
	if in == nil {
		return nil
	}
	out := new(StackStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StargateSpec) DeepCopyInto(out *StargateSpec) {
	*out = *in
	out.DevProperties = in.DevProperties
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StargateSpec.
func (in *StargateSpec) DeepCopy() *StargateSpec {
	if in == nil {
		return nil
	}
	out := new(StargateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StaticClient) DeepCopyInto(out *StaticClient) {
	*out = *in
	in.ClientConfiguration.DeepCopyInto(&out.ClientConfiguration)
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StaticClient.
func (in *StaticClient) DeepCopy() *StaticClient {
	if in == nil {
		return nil
	}
	out := new(StaticClient)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Status) DeepCopyInto(out *Status) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Status.
func (in *Status) DeepCopy() *Status {
	if in == nil {
		return nil
	}
	out := new(Status)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TemporalConfig) DeepCopyInto(out *TemporalConfig) {
	*out = *in
	out.TLS = in.TLS
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TemporalConfig.
func (in *TemporalConfig) DeepCopy() *TemporalConfig {
	if in == nil {
		return nil
	}
	out := new(TemporalConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TemporalTLSConfig) DeepCopyInto(out *TemporalTLSConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TemporalTLSConfig.
func (in *TemporalTLSConfig) DeepCopy() *TemporalTLSConfig {
	if in == nil {
		return nil
	}
	out := new(TemporalTLSConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TracesOtlpSpec) DeepCopyInto(out *TracesOtlpSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TracesOtlpSpec.
func (in *TracesOtlpSpec) DeepCopy() *TracesOtlpSpec {
	if in == nil {
		return nil
	}
	out := new(TracesOtlpSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TracesSpec) DeepCopyInto(out *TracesSpec) {
	*out = *in
	if in.Otlp != nil {
		in, out := &in.Otlp, &out.Otlp
		*out = new(TracesOtlpSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TracesSpec.
func (in *TracesSpec) DeepCopy() *TracesSpec {
	if in == nil {
		return nil
	}
	out := new(TracesSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Versions) DeepCopyInto(out *Versions) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Versions.
func (in *Versions) DeepCopy() *Versions {
	if in == nil {
		return nil
	}
	out := new(Versions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Versions) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VersionsList) DeepCopyInto(out *VersionsList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Versions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VersionsList.
func (in *VersionsList) DeepCopy() *VersionsList {
	if in == nil {
		return nil
	}
	out := new(VersionsList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VersionsList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VersionsSpec) DeepCopyInto(out *VersionsSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VersionsSpec.
func (in *VersionsSpec) DeepCopy() *VersionsSpec {
	if in == nil {
		return nil
	}
	out := new(VersionsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VersionsStatus) DeepCopyInto(out *VersionsStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VersionsStatus.
func (in *VersionsStatus) DeepCopy() *VersionsStatus {
	if in == nil {
		return nil
	}
	out := new(VersionsStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WalletsSpec) DeepCopyInto(out *WalletsSpec) {
	*out = *in
	out.DevProperties = in.DevProperties
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WalletsSpec.
func (in *WalletsSpec) DeepCopy() *WalletsSpec {
	if in == nil {
		return nil
	}
	out := new(WalletsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhooksSpec) DeepCopyInto(out *WebhooksSpec) {
	*out = *in
	out.DevProperties = in.DevProperties
	out.Postgres = in.Postgres
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhooksSpec.
func (in *WebhooksSpec) DeepCopy() *WebhooksSpec {
	if in == nil {
		return nil
	}
	out := new(WebhooksSpec)
	in.DeepCopyInto(out)
	return out
}
