# API Reference

## Packages
- [stack.formance.com/v1beta3](#stackformancecomv1beta3)


## stack.formance.com/v1beta3

Package v1beta3 contains API Schema definitions for the stack v1beta3 API group


### Resource Types
- [Configuration](#configuration)
- [Migration](#migration)
- [Stack](#stack)
- [Versions](#versions)



#### AnnotationsServicesSpec





_Appears in:_
- [AuthSpec](#authspec)
- [ControlSpec](#controlspec)
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
| `service` _object (keys:string, values:string)_ |  |


#### AuthConfig





_Appears in:_
- [ConfigurationSpec](#configurationspec)

| Field | Description |
| --- | --- |
| `readKeySetMaxRetries` _integer_ |  |
| `checkScopes` _boolean_ |  |


#### AuthSpec





_Appears in:_
- [ConfigurationServicesSpec](#configurationservicesspec)

| Field | Description |
| --- | --- |
| `postgres` _[PostgresConfig](#postgresconfig)_ |  |
| `staticClients` _[StaticClient](#staticclient) array_ |  |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `resourceProperties` _[ResourceProperties](#resourceproperties)_ |  |
| `annotations` _[AnnotationsServicesSpec](#annotationsservicesspec)_ |  |


#### Batching





_Appears in:_
- [SearchSpec](#searchspec)

| Field | Description |
| --- | --- |
| `count` _integer_ |  |
| `period` _string_ |  |


#### Broker





_Appears in:_
- [CollectorConfig](#collectorconfig)
- [ConfigurationSpec](#configurationspec)

| Field | Description |
| --- | --- |
| `kafka` _[KafkaConfig](#kafkaconfig)_ |  |
| `nats` _[NatsConfig](#natsconfig)_ |  |


#### ClientConfiguration





_Appears in:_
- [StaticClient](#staticclient)

| Field | Description |
| --- | --- |
| `public` _boolean_ |  |
| `description` _string_ |  |
| `redirectUris` _string array_ |  |
| `postLogoutRedirectUris` _string array_ |  |
| `scopes` _string array_ |  |




#### CommonServiceProperties





_Appears in:_
- [ControlSpec](#controlspec)
- [LedgerSpec](#ledgerspec)
- [OrchestrationSpec](#orchestrationspec)
- [PaymentsSpec](#paymentsspec)
- [ReconciliationSpec](#reconciliationspec)
- [SearchSpec](#searchspec)
- [WalletsSpec](#walletsspec)
- [WebhooksSpec](#webhooksspec)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `disabled` _boolean_ |  |


#### Condition





_Appears in:_
- [Conditions](#conditions)

| Field | Description |
| --- | --- |
| `type` _string_ | type of condition in CamelCase or in foo.example.com/CamelCase. --- Many .condition.type values are consistent across resources like Available, but because arbitrary conditions can be useful (see .node.status.conditions), the ability to deconflict is important. The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt) |
| `observedGeneration` _integer_ | observedGeneration represents the .metadata.generation that the condition was set based upon. For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date with respect to the current state of the instance. |
| `lastTransitionTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#time-v1-meta)_ | lastTransitionTime is the last time the condition transitioned from one status to another. This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable. |
| `message` _string_ | message is a human readable message indicating details about the transition. This may be an empty string. |


#### Conditions

_Underlying type:_ _[Condition](#condition)_



_Appears in:_
- [StackStatus](#stackstatus)
- [Status](#status)

| Field | Description |
| --- | --- |
| `type` _string_ | type of condition in CamelCase or in foo.example.com/CamelCase. --- Many .condition.type values are consistent across resources like Available, but because arbitrary conditions can be useful (see .node.status.conditions), the ability to deconflict is important. The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt) |
| `observedGeneration` _integer_ | observedGeneration represents the .metadata.generation that the condition was set based upon. For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date with respect to the current state of the instance. |
| `lastTransitionTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#time-v1-meta)_ | lastTransitionTime is the last time the condition transitioned from one status to another. This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable. |
| `message` _string_ | message is a human readable message indicating details about the transition. This may be an empty string. |


#### Configuration



Configuration is the Schema for the configurations API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `stack.formance.com/v1beta3`
| `kind` _string_ | `Configuration`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[ConfigurationSpec](#configurationspec)_ |  |


#### ConfigurationServicesSpec



ConfigurationServicesSpec define all existing services for a stack. Fields order is important. For example, auth must be defined later as other services create static auth clients which must be used by auth.

_Appears in:_
- [ConfigurationSpec](#configurationspec)

| Field | Description |
| --- | --- |
| `control` _[ControlSpec](#controlspec)_ |  |
| `ledger` _[LedgerSpec](#ledgerspec)_ |  |
| `payments` _[PaymentsSpec](#paymentsspec)_ |  |
| `reconciliation` _[ReconciliationSpec](#reconciliationspec)_ |  |
| `webhooks` _[WebhooksSpec](#webhooksspec)_ |  |
| `wallets` _[WalletsSpec](#walletsspec)_ |  |
| `orchestration` _[OrchestrationSpec](#orchestrationspec)_ |  |
| `search` _[SearchSpec](#searchspec)_ |  |
| `auth` _[AuthSpec](#authspec)_ |  |
| `gateway` _[GatewaySpec](#gatewayspec)_ |  |
| `stargate` _[StargateSpec](#stargatespec)_ |  |


#### ConfigurationSpec





_Appears in:_
- [Configuration](#configuration)

| Field | Description |
| --- | --- |
| `services` _[ConfigurationServicesSpec](#configurationservicesspec)_ |  |
| `broker` _[Broker](#broker)_ |  |
| `monitoring` _[MonitoringSpec](#monitoringspec)_ |  |
| `auth` _[AuthConfig](#authconfig)_ |  |
| `ingress` _[IngressGlobalConfig](#ingressglobalconfig)_ |  |
| `temporal` _[TemporalConfig](#temporalconfig)_ |  |
| `light` _boolean_ | LightMode is experimental and indicate we want monopods |
| `registries` _object (keys:string, values:[RegistryConfig](#registryconfig))_ |  |




#### ControlSpec





_Appears in:_
- [ConfigurationServicesSpec](#configurationservicesspec)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `disabled` _boolean_ |  |
| `resourceProperties` _[ResourceProperties](#resourceproperties)_ |  |
| `annotations` _[AnnotationsServicesSpec](#annotationsservicesspec)_ |  |




#### DelegatedOIDCServerConfiguration





_Appears in:_
- [StackAuthSpec](#stackauthspec)

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
- [ControlSpec](#controlspec)
- [LedgerSpec](#ledgerspec)
- [OrchestrationSpec](#orchestrationspec)
- [PaymentsSpec](#paymentsspec)
- [ReconciliationSpec](#reconciliationspec)
- [SearchSpec](#searchspec)
- [StackSpec](#stackspec)
- [StargateSpec](#stargatespec)
- [WalletsSpec](#walletsspec)
- [WebhooksSpec](#webhooksspec)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |


#### ElasticSearchBasicAuthConfig





_Appears in:_
- [ElasticSearchConfig](#elasticsearchconfig)

| Field | Description |
| --- | --- |
| `username` _string_ |  |
| `password` _string_ |  |
| `secretName` _string_ |  |


#### ElasticSearchConfig





_Appears in:_
- [SearchSpec](#searchspec)

| Field | Description |
| --- | --- |
| `scheme` _string_ |  |
| `host` _string_ |  |
| `port` _integer_ |  |
| `tls` _[ElasticSearchTLSConfig](#elasticsearchtlsconfig)_ |  |
| `basicAuth` _[ElasticSearchBasicAuthConfig](#elasticsearchbasicauthconfig)_ |  |
| `pathPrefix` _string_ |  |
| `useZinc` _boolean_ |  |


#### ElasticSearchTLSConfig





_Appears in:_
- [ElasticSearchConfig](#elasticsearchconfig)

| Field | Description |
| --- | --- |
| `enabled` _boolean_ |  |
| `skipCertVerify` _boolean_ |  |


#### GatewaySpec





_Appears in:_
- [ConfigurationServicesSpec](#configurationservicesspec)

| Field | Description |
| --- | --- |
| `resourceProperties` _[ResourceProperties](#resourceproperties)_ |  |
| `annotations` _[AnnotationsServicesSpec](#annotationsservicesspec)_ |  |
| `fallback` _string_ |  |
| `enableAuditPlugin` _boolean_ |  |
| `livenessEndpoint` _string_ |  |
| `enableScopes` _boolean_ |  |


#### IngressConfig





_Appears in:_
- [IngressGlobalConfig](#ingressglobalconfig)

| Field | Description |
| --- | --- |
| `annotations` _object (keys:string, values:string)_ |  |


#### IngressGlobalConfig





_Appears in:_
- [ConfigurationSpec](#configurationspec)

| Field | Description |
| --- | --- |
| `annotations` _object (keys:string, values:string)_ |  |
| `tls` _[IngressTLS](#ingresstls)_ |  |


#### IngressTLS





_Appears in:_
- [IngressGlobalConfig](#ingressglobalconfig)

| Field | Description |
| --- | --- |
| `secretName` _string_ | SecretName is the name of the secret used to terminate TLS traffic on port 443. Field is left optional to allow TLS routing based on SNI hostname alone. If the SNI host in a listener conflicts with the "Host" header field used by an IngressRule, the SNI host is used for termination and value of the Host header is used for routing. |


#### KafkaConfig





_Appears in:_
- [Broker](#broker)
- [CollectorConfig](#collectorconfig)

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


#### LedgerSpec





_Appears in:_
- [ConfigurationServicesSpec](#configurationservicesspec)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `disabled` _boolean_ |  |
| `postgres` _[PostgresConfig](#postgresconfig)_ |  |
| `locking` _[LockingStrategy](#lockingstrategy)_ |  |
| `allowPastTimestamps` _boolean_ |  |
| `resourceProperties` _[ResourceProperties](#resourceproperties)_ |  |
| `annotations` _[AnnotationsServicesSpec](#annotationsservicesspec)_ |  |
| `deploymentStrategy` _[DeploymentStrategy](#deploymentstrategy)_ |  |
| `disabled` _boolean_ |  |


#### LockingStrategy





_Appears in:_
- [LedgerSpec](#ledgerspec)

| Field | Description |
| --- | --- |
| `strategy` _string_ |  |
| `redis` _[LockingStrategyRedisConfig](#lockingstrategyredisconfig)_ |  |


#### LockingStrategyRedisConfig





_Appears in:_
- [LockingStrategy](#lockingstrategy)

| Field | Description |
| --- | --- |
| `uri` _string_ |  |
| `tls` _boolean_ |  |
| `insecure` _boolean_ |  |
| `duration` _[Duration](#duration)_ |  |
| `retry` _[Duration](#duration)_ |  |


#### MetricsSpec





_Appears in:_
- [MonitoringSpec](#monitoringspec)

| Field | Description |
| --- | --- |
| `otlp` _[OtlpSpec](#otlpspec)_ |  |


#### Migration



Migration is the Schema for the migrations API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `stack.formance.com/v1beta3`
| `kind` _string_ | `Migration`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[MigrationSpec](#migrationspec)_ |  |


#### MigrationSpec



MigrationSpec defines the desired state of Migration

_Appears in:_
- [Migration](#migration)

| Field | Description |
| --- | --- |
| `configuration` _string_ |  |
| `currentVersion` _string_ |  |
| `version` _string_ |  |
| `module` _string_ |  |
| `targetedVersion` _string_ |  |
| `postUpgrade` _boolean_ |  |




#### MonitoringSpec





_Appears in:_
- [ConfigurationSpec](#configurationspec)

| Field | Description |
| --- | --- |
| `traces` _[TracesSpec](#tracesspec)_ |  |
| `metrics` _[MetricsSpec](#metricsspec)_ |  |


#### NatsConfig





_Appears in:_
- [Broker](#broker)
- [CollectorConfig](#collectorconfig)

| Field | Description |
| --- | --- |
| `url` _string_ |  |
| `replicas` _integer_ |  |


#### OrchestrationSpec





_Appears in:_
- [ConfigurationServicesSpec](#configurationservicesspec)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `disabled` _boolean_ |  |
| `postgres` _[PostgresConfig](#postgresconfig)_ |  |
| `resourceProperties` _[ResourceProperties](#resourceproperties)_ |  |
| `annotations` _[AnnotationsServicesSpec](#annotationsservicesspec)_ |  |


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
| `resourceAttributes` _string_ |  |


#### PaymentsSpec





_Appears in:_
- [ConfigurationServicesSpec](#configurationservicesspec)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `disabled` _boolean_ |  |
| `encryptionKey` _string_ |  |
| `postgres` _[PostgresConfig](#postgresconfig)_ |  |
| `resourceProperties` _[ResourceProperties](#resourceproperties)_ |  |
| `annotations` _[AnnotationsServicesSpec](#annotationsservicesspec)_ |  |


#### PostgresConfig





_Appears in:_
- [AuthSpec](#authspec)
- [LedgerSpec](#ledgerspec)
- [OrchestrationSpec](#orchestrationspec)
- [PaymentsSpec](#paymentsspec)
- [ReconciliationSpec](#reconciliationspec)
- [WebhooksSpec](#webhooksspec)

| Field | Description |
| --- | --- |
| `port` _integer_ |  |
| `host` _string_ |  |
| `username` _string_ |  |
| `password` _string_ |  |
| `debug` _boolean_ |  |
| `credentialsFromSecret` _string_ |  |
| `disableSSLMode` _boolean_ |  |


#### ReconciliationSpec





_Appears in:_
- [ConfigurationServicesSpec](#configurationservicesspec)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `disabled` _boolean_ |  |
| `postgres` _[PostgresConfig](#postgresconfig)_ |  |
| `resourceProperties` _[ResourceProperties](#resourceproperties)_ |  |
| `annotations` _[AnnotationsServicesSpec](#annotationsservicesspec)_ |  |


#### RegistryConfig





_Appears in:_
- [ConfigurationSpec](#configurationspec)

| Field | Description |
| --- | --- |
| `endpoint` _string_ |  |


#### ResourceProperties





_Appears in:_
- [AuthSpec](#authspec)
- [ControlSpec](#controlspec)
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
| `request` _[resource](#resource)_ |  |
| `limits` _[resource](#resource)_ |  |


#### SearchSpec





_Appears in:_
- [ConfigurationServicesSpec](#configurationservicesspec)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `disabled` _boolean_ |  |
| `elasticSearch` _[ElasticSearchConfig](#elasticsearchconfig)_ |  |
| `batching` _[Batching](#batching)_ |  |
| `searchResourceProperties` _[ResourceProperties](#resourceproperties)_ |  |
| `benthosResourceProperties` _[ResourceProperties](#resourceproperties)_ |  |
| `annotations` _[AnnotationsServicesSpec](#annotationsservicesspec)_ |  |




#### Stack



Stack is the Schema for the stacks API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `stack.formance.com/v1beta3`
| `kind` _string_ | `Stack`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[StackSpec](#stackspec)_ |  |


#### StackAuthSpec





_Appears in:_
- [StackSpec](#stackspec)

| Field | Description |
| --- | --- |
| `delegatedOIDCServer` _[DelegatedOIDCServerConfiguration](#delegatedoidcserverconfiguration)_ |  |
| `staticClients` _[StaticClient](#staticclient) array_ |  |


#### StackServicePropertiesSpec





_Appears in:_
- [StackServicesSpec](#stackservicesspec)

| Field | Description |
| --- | --- |
| `disabled` _boolean_ |  |


#### StackServicesSpec





_Appears in:_
- [StackSpec](#stackspec)

| Field | Description |
| --- | --- |
| `ledger` _[StackServicePropertiesSpec](#stackservicepropertiesspec)_ |  |
| `orchestration` _[StackServicePropertiesSpec](#stackservicepropertiesspec)_ |  |
| `reconciliation` _[StackServicePropertiesSpec](#stackservicepropertiesspec)_ |  |
| `payments` _[StackServicePropertiesSpec](#stackservicepropertiesspec)_ |  |
| `wallets` _[StackServicePropertiesSpec](#stackservicepropertiesspec)_ |  |
| `webhooks` _[StackServicePropertiesSpec](#stackservicepropertiesspec)_ |  |
| `control` _[StackServicePropertiesSpec](#stackservicepropertiesspec)_ |  |


#### StackSpec



StackSpec defines the desired state of Stack

_Appears in:_
- [Stack](#stack)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `seed` _string_ |  |
| `host` _string_ |  |
| `auth` _[StackAuthSpec](#stackauthspec)_ |  |
| `stargate` _[StackStargateConfig](#stackstargateconfig)_ |  |
| `versions` _string_ |  |
| `scheme` _string_ |  |
| `disabled` _boolean_ |  |
| `services` _[StackServicesSpec](#stackservicesspec)_ |  |


#### StackStargateConfig





_Appears in:_
- [StackSpec](#stackspec)

| Field | Description |
| --- | --- |
| `stargateServerURL` _string_ |  |




#### StargateSpec





_Appears in:_
- [ConfigurationServicesSpec](#configurationservicesspec)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `resourceProperties` _[ResourceProperties](#resourceproperties)_ |  |
| `annotations` _[AnnotationsServicesSpec](#annotationsservicesspec)_ |  |


#### StaticClient





_Appears in:_
- [AuthSpec](#authspec)
- [StackAuthSpec](#stackauthspec)
- [StackStatus](#stackstatus)

| Field | Description |
| --- | --- |
| `public` _boolean_ |  |
| `description` _string_ |  |
| `redirectUris` _string array_ |  |
| `postLogoutRedirectUris` _string array_ |  |
| `scopes` _string array_ |  |
| `id` _string_ |  |
| `secrets` _string array_ |  |


#### Status





_Appears in:_
- [StackStatus](#stackstatus)

| Field | Description |
| --- | --- |
| `conditions` _[Conditions](#conditions)_ |  |


#### TemporalConfig





_Appears in:_
- [ConfigurationSpec](#configurationspec)

| Field | Description |
| --- | --- |
| `address` _string_ |  |
| `namespace` _string_ |  |
| `tls` _[TemporalTLSConfig](#temporaltlsconfig)_ |  |


#### TemporalTLSConfig



TODO: Handle validation

_Appears in:_
- [TemporalConfig](#temporalconfig)

| Field | Description |
| --- | --- |
| `crt` _string_ |  |
| `key` _string_ |  |
| `secretName` _string_ |  |


#### TracesSpec





_Appears in:_
- [MonitoringSpec](#monitoringspec)

| Field | Description |
| --- | --- |
| `otlp` _[OtlpSpec](#otlpspec)_ |  |




#### Versions



Versions is the Schema for the versions API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `stack.formance.com/v1beta3`
| `kind` _string_ | `Versions`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[VersionsSpec](#versionsspec)_ |  |


#### VersionsSpec



VersionsSpec defines the desired state of Versions

_Appears in:_
- [Versions](#versions)

| Field | Description |
| --- | --- |
| `control` _string_ |  |
| `ledger` _string_ |  |
| `payments` _string_ |  |
| `search` _string_ |  |
| `auth` _string_ |  |
| `webhooks` _string_ |  |
| `wallets` _string_ |  |
| `orchestration` _string_ |  |
| `gateway` _string_ |  |
| `stargate` _string_ |  |
| `reconciliation` _string_ |  |




#### WalletsSpec





_Appears in:_
- [ConfigurationServicesSpec](#configurationservicesspec)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `disabled` _boolean_ |  |
| `resourceProperties` _[ResourceProperties](#resourceproperties)_ |  |
| `annotations` _[AnnotationsServicesSpec](#annotationsservicesspec)_ |  |


#### WebhooksSpec





_Appears in:_
- [ConfigurationServicesSpec](#configurationservicesspec)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `disabled` _boolean_ |  |
| `postgres` _[PostgresConfig](#postgresconfig)_ |  |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `resourceProperties` _[ResourceProperties](#resourceproperties)_ |  |
| `annotations` _[AnnotationsServicesSpec](#annotationsservicesspec)_ |  |


