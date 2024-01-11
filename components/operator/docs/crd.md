# API Reference

## Packages
- [formance.com/v1beta1](#formancecomv1beta1)


## formance.com/v1beta1

Package v1beta1 contains API Schema definitions for the formance v1beta1 API group

### Resource Types
- [Auth](#auth)
- [AuthClient](#authclient)
- [BrokerConfiguration](#brokerconfiguration)
- [BrokerTopic](#brokertopic)
- [BrokerTopicConsumer](#brokertopicconsumer)
- [Database](#database)
- [DatabaseConfiguration](#databaseconfiguration)
- [ElasticSearchConfiguration](#elasticsearchconfiguration)
- [Gateway](#gateway)
- [HTTPAPI](#httpapi)
- [Ledger](#ledger)
- [OpenTelemetryConfiguration](#opentelemetryconfiguration)
- [Orchestration](#orchestration)
- [Payments](#payments)
- [Reconciliation](#reconciliation)
- [RegistriesConfiguration](#registriesconfiguration)
- [Search](#search)
- [Stack](#stack)
- [Stargate](#stargate)
- [Stream](#stream)
- [StreamProcessor](#streamprocessor)
- [Wallets](#wallets)
- [Webhooks](#webhooks)



#### Auth



Auth is the Schema for the auths API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Auth`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[AuthSpec](#authspec)_ |  |


#### AuthClient



AuthClient is the Schema for the authclients API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `AuthClient`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[AuthClientSpec](#authclientspec)_ |  |


#### AuthClientSpec



AuthClientSpec defines the desired state of AuthClient

_Appears in:_
- [AuthClient](#authclient)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `id` _string_ |  |
| `public` _boolean_ |  |
| `description` _string_ |  |
| `redirectUris` _string array_ |  |
| `postLogoutRedirectUris` _string array_ |  |
| `scopes` _string array_ |  |
| `secret` _string_ |  |




#### AuthConfig





_Appears in:_
- [LedgerSpec](#ledgerspec)
- [OrchestrationSpec](#orchestrationspec)
- [PaymentsSpec](#paymentsspec)
- [ReconciliationSpec](#reconciliationspec)
- [SearchSpec](#searchspec)
- [WalletsSpec](#walletsspec)
- [WebhooksSpec](#webhooksspec)

| Field | Description |
| --- | --- |
| `readKeySetMaxRetries` _integer_ |  |
| `checkScopes` _boolean_ |  |


#### AuthSpec



AuthSpec defines the desired state of Auth

_Appears in:_
- [Auth](#auth)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `version` _string_ |  |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |
| `stack` _string_ |  |
| `delegatedOIDCServer` _[DelegatedOIDCServerConfiguration](#delegatedoidcserverconfiguration)_ |  |
| `signingKey` _string_ |  |
| `signingKeyFromSecret` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#secretkeyselector-v1-core)_ |  |
| `enableScopes` _boolean_ |  |
| `service` _[ServiceConfiguration](#serviceconfiguration)_ |  |




#### Batching





_Appears in:_
- [SearchSpec](#searchspec)
- [StreamProcessorSpec](#streamprocessorspec)

| Field | Description |
| --- | --- |
| `count` _integer_ |  |
| `period` _string_ |  |


#### BrokerConfiguration



BrokerConfiguration is the Schema for the brokerconfigurations API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `BrokerConfiguration`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[BrokerConfigurationSpec](#brokerconfigurationspec)_ |  |


#### BrokerConfigurationSpec



BrokerConfigurationSpec defines the desired state of BrokerConfig

_Appears in:_
- [BrokerConfiguration](#brokerconfiguration)
- [BrokerTopicStatus](#brokertopicstatus)

| Field | Description |
| --- | --- |
| `kafka` _[KafkaConfig](#kafkaconfig)_ |  |
| `nats` _[NatsConfig](#natsconfig)_ |  |




#### BrokerTopic



BrokerTopic is the Schema for the topics API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `BrokerTopic`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[BrokerTopicSpec](#brokertopicspec)_ |  |


#### BrokerTopicConsumer



BrokerTopicConsumer is the Schema for the topicqueries API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `BrokerTopicConsumer`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[BrokerTopicConsumerSpec](#brokertopicconsumerspec)_ |  |


#### BrokerTopicConsumerSpec



BrokerTopicConsumerSpec defines the desired state of BrokerTopicConsumer

_Appears in:_
- [BrokerTopicConsumer](#brokertopicconsumer)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `service` _string_ |  |
| `queriedBy` _string_ |  |




#### BrokerTopicSpec



BrokerTopicSpec defines the desired state of BrokerTopic

_Appears in:_
- [BrokerTopic](#brokertopic)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `service` _string_ |  |




#### CommonServiceProperties





_Appears in:_
- [AuthSpec](#authspec)
- [GatewaySpec](#gatewayspec)
- [LedgerSpec](#ledgerspec)
- [OrchestrationSpec](#orchestrationspec)
- [PaymentsSpec](#paymentsspec)
- [ReconciliationSpec](#reconciliationspec)
- [SearchSpec](#searchspec)
- [StargateSpec](#stargatespec)
- [WalletsSpec](#walletsspec)
- [WebhooksSpec](#webhooksspec)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `version` _string_ |  |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |


#### CommonStatus





_Appears in:_
- [AuthStatus](#authstatus)
- [BrokerTopicConsumerStatus](#brokertopicconsumerstatus)
- [BrokerTopicStatus](#brokertopicstatus)
- [GatewayStatus](#gatewaystatus)
- [HTTPAPIStatus](#httpapistatus)
- [LedgerStatus](#ledgerstatus)
- [OrchestrationStatus](#orchestrationstatus)
- [PaymentsStatus](#paymentsstatus)
- [ReconciliationStatus](#reconciliationstatus)
- [SearchStatus](#searchstatus)
- [StackStatus](#stackstatus)
- [StargateStatus](#stargatestatus)
- [StreamProcessorStatus](#streamprocessorstatus)
- [WalletsStatus](#walletsstatus)
- [WebhooksStatus](#webhooksstatus)

| Field | Description |
| --- | --- |
| `conditions` _[Condition](#condition) array_ |  |
| `ready` _boolean_ |  |
| `error` _string_ |  |


#### Condition



Condition contains details for one aspect of the current state of this API Resource. --- This struct is intended for direct use as an array at the field path .status.conditions.  For example, 
 type FooStatus struct{ // Represents the observations of a foo's current state. // Known .status.conditions.type are: "Available", "Progressing", and "Degraded" // +patchMergeKey=type // +patchStrategy=merge // +listType=map // +listMapKey=type CommonStatus []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"` 
 // other fields }

_Appears in:_
- [AuthStatus](#authstatus)
- [BrokerTopicConsumerStatus](#brokertopicconsumerstatus)
- [BrokerTopicStatus](#brokertopicstatus)
- [CommonStatus](#commonstatus)
- [GatewayStatus](#gatewaystatus)
- [HTTPAPIStatus](#httpapistatus)
- [LedgerStatus](#ledgerstatus)
- [OrchestrationStatus](#orchestrationstatus)
- [PaymentsStatus](#paymentsstatus)
- [ReconciliationStatus](#reconciliationstatus)
- [SearchStatus](#searchstatus)
- [StackStatus](#stackstatus)
- [StargateStatus](#stargatestatus)
- [StreamProcessorStatus](#streamprocessorstatus)
- [WalletsStatus](#walletsstatus)
- [WebhooksStatus](#webhooksstatus)

| Field | Description |
| --- | --- |
| `type` _string_ | type of condition in CamelCase or in foo.example.com/CamelCase. --- Many .condition.type values are consistent across resources like Available, but because arbitrary conditions can be useful (see .node.status.conditions), the ability to deconflict is important. The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt) |
| `observedGeneration` _integer_ | observedGeneration represents the .metadata.generation that the condition was set based upon. For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date with respect to the current state of the instance. |
| `lastTransitionTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#time-v1-meta)_ | lastTransitionTime is the last time the condition transitioned from one status to another. This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable. |
| `message` _string_ | message is a human readable message indicating details about the transition. This may be an empty string. |


#### CreatedDatabase





_Appears in:_
- [DatabaseStatus](#databasestatus)

| Field | Description |
| --- | --- |
| `port` _integer_ |  |
| `host` _string_ |  |
| `username` _string_ |  |
| `password` _string_ |  |
| `credentialsFromSecret` _string_ |  |
| `disableSSLMode` _boolean_ |  |
| `database` _string_ |  |


#### Database



Database is the Schema for the databases API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Database`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[DatabaseSpec](#databasespec)_ |  |


#### DatabaseConfiguration



DatabaseConfiguration is the Schema for the databaseconfigurations API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `DatabaseConfiguration`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[DatabaseConfigurationSpec](#databaseconfigurationspec)_ |  |


#### DatabaseConfigurationSpec



DatabaseConfigurationSpec defines the desired state of DatabaseConfiguration

_Appears in:_
- [CreatedDatabase](#createddatabase)
- [DatabaseConfiguration](#databaseconfiguration)

| Field | Description |
| --- | --- |
| `port` _integer_ |  |
| `host` _string_ |  |
| `username` _string_ |  |
| `password` _string_ |  |
| `credentialsFromSecret` _string_ |  |
| `disableSSLMode` _boolean_ |  |




#### DatabaseSpec



DatabaseSpec defines the desired state of Database

_Appears in:_
- [Database](#database)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `service` _string_ |  |




#### DelegatedOIDCServerConfiguration





_Appears in:_
- [AuthSpec](#authspec)

| Field | Description |
| --- | --- |
| `issuer` _string_ |  |
| `clientID` _string_ |  |
| `clientSecret` _string_ |  |


#### DeploymentStrategy

_Underlying type:_ _string_



_Appears in:_
- [LedgerSpec](#ledgerspec)



#### DevProperties





_Appears in:_
- [AuthSpec](#authspec)
- [CommonServiceProperties](#commonserviceproperties)
- [GatewaySpec](#gatewayspec)
- [LedgerSpec](#ledgerspec)
- [OrchestrationSpec](#orchestrationspec)
- [PaymentsSpec](#paymentsspec)
- [ReconciliationSpec](#reconciliationspec)
- [SearchSpec](#searchspec)
- [StackSpec](#stackspec)
- [StargateSpec](#stargatespec)
- [StreamProcessorSpec](#streamprocessorspec)
- [WalletsSpec](#walletsspec)
- [WebhooksSpec](#webhooksspec)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |


#### ElasticSearchBasicAuthConfig





_Appears in:_
- [ElasticSearchConfigurationSpec](#elasticsearchconfigurationspec)

| Field | Description |
| --- | --- |
| `username` _string_ |  |
| `password` _string_ |  |
| `secretName` _string_ |  |


#### ElasticSearchConfiguration



ElasticSearchConfiguration is the Schema for the elasticsearchconfigs API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `ElasticSearchConfiguration`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[ElasticSearchConfigurationSpec](#elasticsearchconfigurationspec)_ |  |


#### ElasticSearchConfigurationSpec



ElasticSearchConfigurationSpec defines the desired state of ElasticSearchConfiguration

_Appears in:_
- [ElasticSearchConfiguration](#elasticsearchconfiguration)

| Field | Description |
| --- | --- |
| `scheme` _string_ |  |
| `host` _string_ |  |
| `port` _integer_ |  |
| `tls` _[ElasticSearchTLSConfig](#elasticsearchtlsconfig)_ |  |
| `basicAuth` _[ElasticSearchBasicAuthConfig](#elasticsearchbasicauthconfig)_ |  |




#### ElasticSearchTLSConfig





_Appears in:_
- [ElasticSearchConfigurationSpec](#elasticsearchconfigurationspec)

| Field | Description |
| --- | --- |
| `enabled` _boolean_ |  |
| `skipCertVerify` _boolean_ |  |




#### Gateway



Gateway is the Schema for the gateways API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Gateway`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[GatewaySpec](#gatewayspec)_ |  |


#### GatewayIngress





_Appears in:_
- [GatewaySpec](#gatewayspec)

| Field | Description |
| --- | --- |
| `host` _string_ |  |
| `scheme` _string_ |  |
| `annotations` _object (keys:string, values:string)_ |  |
| `tls` _[GatewayIngressTLS](#gatewayingresstls)_ |  |


#### GatewayIngressTLS





_Appears in:_
- [GatewayIngress](#gatewayingress)

| Field | Description |
| --- | --- |
| `secretName` _string_ |  |


#### GatewaySpec



GatewaySpec defines the desired state of Gateway

_Appears in:_
- [Gateway](#gateway)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `version` _string_ |  |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |
| `ingress` _[GatewayIngress](#gatewayingress)_ |  |
| `service` _[ServiceConfiguration](#serviceconfiguration)_ |  |




#### HTTPAPI



HTTPAPI is the Schema for the HTTPAPIs API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `HTTPAPI`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[HTTPAPISpec](#httpapispec)_ |  |


#### HTTPAPIRule





_Appears in:_
- [HTTPAPISpec](#httpapispec)

| Field | Description |
| --- | --- |
| `path` _string_ |  |
| `methods` _string array_ |  |
| `secured` _boolean_ |  |


#### HTTPAPISpec



HTTPAPISpec defines the desired state of HTTPAPI

_Appears in:_
- [HTTPAPI](#httpapi)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `name` _string_ | Name indicates prefix api |
| `rules` _[HTTPAPIRule](#httpapirule) array_ | Rules |
| `service` _[ServiceConfiguration](#serviceconfiguration)_ | ServiceConfiguration |




#### KafkaConfig





_Appears in:_
- [BrokerConfigurationSpec](#brokerconfigurationspec)

| Field | Description |
| --- | --- |
| `brokers` _string array_ |  |
| `tls` _boolean_ |  |
| `sasl` _[KafkaSASLConfig](#kafkasaslconfig)_ |  |


#### KafkaSASLConfig





_Appears in:_
- [KafkaConfig](#kafkaconfig)

| Field | Description |
| --- | --- |
| `username` _string_ |  |
| `password` _string_ |  |
| `mechanism` _string_ |  |
| `scramSHASize` _string_ |  |


#### Ledger



Ledger is the Schema for the ledgers API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Ledger`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[LedgerSpec](#ledgerspec)_ |  |


#### LedgerSpec



LedgerSpec defines the desired state of Ledger

_Appears in:_
- [Ledger](#ledger)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `version` _string_ |  |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |
| `stack` _string_ |  |
| `auth` _[AuthConfig](#authconfig)_ |  |
| `deploymentStrategy` _[DeploymentStrategy](#deploymentstrategy)_ |  |
| `service` _[ServiceConfiguration](#serviceconfiguration)_ |  |




#### MetricsSpec





_Appears in:_
- [OpenTelemetryConfigurationSpec](#opentelemetryconfigurationspec)

| Field | Description |
| --- | --- |
| `otlp` _[OtlpSpec](#otlpspec)_ |  |


#### NatsConfig





_Appears in:_
- [BrokerConfigurationSpec](#brokerconfigurationspec)

| Field | Description |
| --- | --- |
| `url` _string_ |  |
| `replicas` _integer_ |  |


#### OpenTelemetryConfiguration



OpenTelemetryConfiguration is the Schema for the opentelemetrytraces API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `OpenTelemetryConfiguration`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[OpenTelemetryConfigurationSpec](#opentelemetryconfigurationspec)_ |  |


#### OpenTelemetryConfigurationSpec



OpenTelemetryConfigurationSpec defines the desired state of OpenTelemetryConfiguration

_Appears in:_
- [OpenTelemetryConfiguration](#opentelemetryconfiguration)

| Field | Description |
| --- | --- |
| `traces` _[TracesSpec](#tracesspec)_ |  |
| `metrics` _[MetricsSpec](#metricsspec)_ |  |




#### Orchestration



Orchestration is the Schema for the orchestrations API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Orchestration`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[OrchestrationSpec](#orchestrationspec)_ |  |


#### OrchestrationSpec



OrchestrationSpec defines the desired state of Orchestration

_Appears in:_
- [Orchestration](#orchestration)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `version` _string_ |  |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |
| `temporal` _[TemporalConfig](#temporalconfig)_ |  |
| `service` _[ServiceConfiguration](#serviceconfiguration)_ |  |
| `auth` _[AuthConfig](#authconfig)_ |  |




#### OtlpSpec





_Appears in:_
- [MetricsSpec](#metricsspec)
- [TracesSpec](#tracesspec)

| Field | Description |
| --- | --- |
| `endpoint` _string_ |  |
| `port` _integer_ |  |
| `insecure` _boolean_ |  |
| `mode` _string_ |  |
| `resourceAttributes` _object (keys:string, values:string)_ |  |


#### Payments



Payments is the Schema for the payments API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Payments`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[PaymentsSpec](#paymentsspec)_ |  |


#### PaymentsSpec



PaymentsSpec defines the desired state of Payments

_Appears in:_
- [Payments](#payments)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `version` _string_ |  |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |
| `encryptionKey` _string_ |  |
| `service` _[ServiceConfiguration](#serviceconfiguration)_ |  |
| `auth` _[AuthConfig](#authconfig)_ |  |




#### Reconciliation



Reconciliation is the Schema for the reconciliations API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Reconciliation`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[ReconciliationSpec](#reconciliationspec)_ |  |


#### ReconciliationSpec



ReconciliationSpec defines the desired state of Reconciliation

_Appears in:_
- [Reconciliation](#reconciliation)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `version` _string_ |  |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |
| `service` _[ServiceConfiguration](#serviceconfiguration)_ |  |
| `auth` _[AuthConfig](#authconfig)_ |  |




#### RegistriesConfiguration



RegistriesConfiguration is the Schema for the registries API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `RegistriesConfiguration`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[RegistriesConfigurationSpec](#registriesconfigurationspec)_ |  |


#### RegistriesConfigurationSpec



RegistriesConfigurationSpec defines the desired state of RegistriesConfiguration

_Appears in:_
- [RegistriesConfiguration](#registriesconfiguration)

| Field | Description |
| --- | --- |
| `registries` _object (keys:string, values:[RegistryConfigurationSpec](#registryconfigurationspec))_ |  |




#### RegistryConfigurationSpec





_Appears in:_
- [RegistriesConfigurationSpec](#registriesconfigurationspec)

| Field | Description |
| --- | --- |
| `endpoint` _string_ |  |


#### Search



Search is the Schema for the searches API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Search`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[SearchSpec](#searchspec)_ |  |


#### SearchSpec



SearchSpec defines the desired state of Search

_Appears in:_
- [Search](#search)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `version` _string_ |  |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |
| `batching` _[Batching](#batching)_ |  |
| `streamProcessor` _[SearchStreamProcessorSpec](#searchstreamprocessorspec)_ |  |
| `service` _[ServiceConfiguration](#serviceconfiguration)_ |  |
| `auth` _[AuthConfig](#authconfig)_ |  |




#### SearchStreamProcessorSpec





_Appears in:_
- [SearchSpec](#searchspec)

| Field | Description |
| --- | --- |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |


#### ServiceConfiguration





_Appears in:_
- [AuthSpec](#authspec)
- [GatewaySpec](#gatewayspec)
- [HTTPAPISpec](#httpapispec)
- [LedgerSpec](#ledgerspec)
- [OrchestrationSpec](#orchestrationspec)
- [PaymentsSpec](#paymentsspec)
- [ReconciliationSpec](#reconciliationspec)
- [SearchSpec](#searchspec)
- [StargateSpec](#stargatespec)
- [WalletsSpec](#walletsspec)
- [WebhooksSpec](#webhooksspec)

| Field | Description |
| --- | --- |
| `annotations` _object (keys:string, values:string)_ |  |


#### Stack



Stack is the Schema for the stacks API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Stack`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[StackSpec](#stackspec)_ |  |


#### StackDependency





_Appears in:_
- [AuthClientSpec](#authclientspec)
- [AuthSpec](#authspec)
- [BrokerTopicConsumerSpec](#brokertopicconsumerspec)
- [BrokerTopicSpec](#brokertopicspec)
- [DatabaseSpec](#databasespec)
- [GatewaySpec](#gatewayspec)
- [HTTPAPISpec](#httpapispec)
- [LedgerSpec](#ledgerspec)
- [OrchestrationSpec](#orchestrationspec)
- [PaymentsSpec](#paymentsspec)
- [ReconciliationSpec](#reconciliationspec)
- [SearchSpec](#searchspec)
- [StargateSpec](#stargatespec)
- [StreamProcessorSpec](#streamprocessorspec)
- [StreamSpec](#streamspec)
- [WalletsSpec](#walletsspec)
- [WebhooksSpec](#webhooksspec)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |


#### StackSpec



StackSpec defines the desired state of Stack

_Appears in:_
- [Stack](#stack)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `version` _string_ |  |
| `enableAudit` _boolean_ |  |




#### Stargate



Stargate is the Schema for the stargates API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Stargate`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[StargateSpec](#stargatespec)_ |  |


#### StargateAuthSpec





_Appears in:_
- [StargateSpec](#stargatespec)

| Field | Description |
| --- | --- |
| `clientID` _string_ |  |
| `clientSecret` _string_ |  |
| `issuer` _string_ |  |


#### StargateSpec



StargateSpec defines the desired state of Stargate

_Appears in:_
- [Stargate](#stargate)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `version` _string_ |  |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |
| `stack` _string_ |  |
| `serverURL` _string_ |  |
| `organizationID` _string_ |  |
| `stackID` _string_ |  |
| `auth` _[StargateAuthSpec](#stargateauthspec)_ |  |
| `service` _[ServiceConfiguration](#serviceconfiguration)_ |  |




#### Stream



Stream is the Schema for the streams API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Stream`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[StreamSpec](#streamspec)_ |  |


#### StreamProcessor



StreamProcessor is the Schema for the streamprocessors API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `StreamProcessor`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[StreamProcessorSpec](#streamprocessorspec)_ |  |


#### StreamProcessorSpec



StreamProcessorSpec defines the desired state of StreamProcessor

_Appears in:_
- [StreamProcessor](#streamprocessor)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |
| `batching` _[Batching](#batching)_ |  |
| `initContainers` _[Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#container-v1-core) array_ |  |




#### StreamSpec



StreamSpec defines the desired state of Stream

_Appears in:_
- [Stream](#stream)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `data` _string_ |  |




#### TemporalConfig





_Appears in:_
- [OrchestrationSpec](#orchestrationspec)

| Field | Description |
| --- | --- |
| `address` _string_ |  |
| `namespace` _string_ |  |
| `tls` _[TemporalTLSConfig](#temporaltlsconfig)_ |  |


#### TemporalTLSConfig





_Appears in:_
- [TemporalConfig](#temporalconfig)

| Field | Description |
| --- | --- |
| `crt` _string_ |  |
| `key` _string_ |  |
| `secretName` _string_ |  |


#### TracesSpec





_Appears in:_
- [OpenTelemetryConfigurationSpec](#opentelemetryconfigurationspec)

| Field | Description |
| --- | --- |
| `otlp` _[OtlpSpec](#otlpspec)_ |  |


#### Wallets



Wallets is the Schema for the wallets API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Wallets`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[WalletsSpec](#walletsspec)_ |  |


#### WalletsSpec



WalletsSpec defines the desired state of Wallets

_Appears in:_
- [Wallets](#wallets)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `version` _string_ |  |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |
| `stack` _string_ |  |
| `service` _[ServiceConfiguration](#serviceconfiguration)_ |  |
| `auth` _[AuthConfig](#authconfig)_ |  |




#### Webhooks



Webhooks is the Schema for the webhooks API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Webhooks`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[WebhooksSpec](#webhooksspec)_ |  |


#### WebhooksSpec



WebhooksSpec defines the desired state of Webhooks

_Appears in:_
- [Webhooks](#webhooks)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `version` _string_ |  |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |
| `service` _[ServiceConfiguration](#serviceconfiguration)_ |  |
| `auth` _[AuthConfig](#authconfig)_ |  |




