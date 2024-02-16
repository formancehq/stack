# API Reference

## Packages
- [formance.com/v1beta1](#formancecomv1beta1)


## formance.com/v1beta1

Package v1beta1 contains API Schema definitions for the formance v1beta1 API group

### Resource Types
- [Auth](#auth)
- [AuthClient](#authclient)
- [Benthos](#benthos)
- [BenthosStream](#benthosstream)
- [BrokerTopic](#brokertopic)
- [BrokerTopicConsumer](#brokertopicconsumer)
- [Database](#database)
- [Gateway](#gateway)
- [GatewayHTTPAPI](#gatewayhttpapi)
- [Ledger](#ledger)
- [Orchestration](#orchestration)
- [Payments](#payments)
- [Reconciliation](#reconciliation)
- [ResourceReference](#resourcereference)
- [Search](#search)
- [Settings](#settings)
- [Stack](#stack)
- [Stargate](#stargate)
- [Versions](#versions)
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
| `stack` _string_ |  |
| `delegatedOIDCServer` _[DelegatedOIDCServerConfiguration](#delegatedoidcserverconfiguration)_ |  |
| `signingKey` _string_ |  |
| `signingKeyFromSecret` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#secretkeyselector-v1-core)_ |  |
| `enableScopes` _boolean_ |  |




#### Batching





_Appears in:_
- [BenthosSpec](#benthosspec)
- [SearchSpec](#searchspec)

| Field | Description |
| --- | --- |
| `count` _integer_ |  |
| `period` _string_ |  |


#### Benthos



Benthos is the Schema for the benthos API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Benthos`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[BenthosSpec](#benthosspec)_ |  |


#### BenthosSpec



BenthosSpec defines the desired state of Benthos

_Appears in:_
- [Benthos](#benthos)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |
| `batching` _[Batching](#batching)_ |  |
| `initContainers` _[Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#container-v1-core) array_ |  |




#### BenthosStream



BenthosStream is the Schema for the benthosstreams API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `BenthosStream`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[BenthosStreamSpec](#benthosstreamspec)_ |  |


#### BenthosStreamSpec



BenthosStreamSpec defines the desired state of BenthosStream

_Appears in:_
- [BenthosStream](#benthosstream)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `data` _string_ |  |
| `name` _string_ |  |




#### BrokerTopic



BrokerTopic is the Schema for the brokertopics API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `BrokerTopic`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[BrokerTopicSpec](#brokertopicspec)_ |  |


#### BrokerTopicConsumer



BrokerTopicConsumer is the Schema for the brokertopicconsumers API



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




#### CommonStatus





_Appears in:_
- [AuthClientStatus](#authclientstatus)
- [AuthStatus](#authstatus)
- [BenthosStatus](#benthosstatus)
- [BenthosStreamStatus](#benthosstreamstatus)
- [BrokerTopicConsumerStatus](#brokertopicconsumerstatus)
- [BrokerTopicStatus](#brokertopicstatus)
- [DatabaseStatus](#databasestatus)
- [GatewayHTTPAPIStatus](#gatewayhttpapistatus)
- [GatewayStatus](#gatewaystatus)
- [LedgerStatus](#ledgerstatus)
- [ModuleStatus](#modulestatus)
- [OrchestrationStatus](#orchestrationstatus)
- [PaymentsStatus](#paymentsstatus)
- [ReconciliationStatus](#reconciliationstatus)
- [ResourceReferenceStatus](#resourcereferencestatus)
- [SearchStatus](#searchstatus)
- [StackStatus](#stackstatus)
- [StargateStatus](#stargatestatus)
- [StatusWithConditions](#statuswithconditions)
- [WalletsStatus](#walletsstatus)
- [WebhooksStatus](#webhooksstatus)

| Field | Description |
| --- | --- |
| `ready` _boolean_ |  |
| `info` _string_ |  |


#### Condition



Condition contains details for one aspect of the current state of this API Resource. --- This struct is intended for direct use as an array at the field path .status.conditions.  For example, 
 type FooStatus struct{ // Represents the observations of a foo's current state. // Known .status.conditions.type are: "Available", "Progressing", and "Degraded" // +patchMergeKey=type // +patchStrategy=merge // +listType=map // +listMapKey=type CommonStatus []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"` 
 // other fields }

_Appears in:_
- [AuthStatus](#authstatus)
- [BenthosStatus](#benthosstatus)
- [GatewayStatus](#gatewaystatus)
- [LedgerStatus](#ledgerstatus)
- [ModuleStatus](#modulestatus)
- [OrchestrationStatus](#orchestrationstatus)
- [PaymentsStatus](#paymentsstatus)
- [ReconciliationStatus](#reconciliationstatus)
- [SearchStatus](#searchstatus)
- [StargateStatus](#stargatestatus)
- [StatusWithConditions](#statuswithconditions)
- [WalletsStatus](#walletsstatus)
- [WebhooksStatus](#webhooksstatus)

| Field | Description |
| --- | --- |
| `type` _string_ | type of condition in CamelCase or in foo.example.com/CamelCase. --- Many .condition.type values are consistent across resources like Available, but because arbitrary conditions can be useful (see .node.status.conditions), the ability to deconflict is important. The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt) |
| `observedGeneration` _integer_ | observedGeneration represents the .metadata.generation that the condition was set based upon. For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date with respect to the current state of the instance. |
| `lastTransitionTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#time-v1-meta)_ | lastTransitionTime is the last time the condition transitioned from one status to another. This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable. |
| `message` _string_ | message is a human readable message indicating details about the transition. This may be an empty string. |
| `reason` _string_ | reason contains a programmatic identifier indicating the reason for the condition's last transition. Producers of specific condition types may define expected values and meanings for this field, and whether the values are considered a guaranteed API. The value should be a CamelCase string. This field may not be empty. |


#### Database



Database is the Schema for the databases API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Database`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[DatabaseSpec](#databasespec)_ |  |


#### DatabaseSpec



DatabaseSpec defines the desired state of Database

_Appears in:_
- [Database](#database)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `service` _string_ |  |
| `debug` _boolean_ |  |




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
- [BenthosSpec](#benthosspec)
- [GatewaySpec](#gatewayspec)
- [LedgerSpec](#ledgerspec)
- [ModuleProperties](#moduleproperties)
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




#### Gateway



Gateway is the Schema for the gateways API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Gateway`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[GatewaySpec](#gatewayspec)_ |  |


#### GatewayHTTPAPI



GatewayHTTPAPI is the Schema for the HTTPAPIs API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `GatewayHTTPAPI`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[GatewayHTTPAPISpec](#gatewayhttpapispec)_ |  |


#### GatewayHTTPAPIRule





_Appears in:_
- [GatewayHTTPAPISpec](#gatewayhttpapispec)

| Field | Description |
| --- | --- |
| `path` _string_ |  |
| `methods` _string array_ |  |
| `secured` _boolean_ |  |


#### GatewayHTTPAPISpec



GatewayHTTPAPISpec defines the desired state of GatewayHTTPAPI

_Appears in:_
- [GatewayHTTPAPI](#gatewayhttpapi)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `name` _string_ | Name indicates prefix api |
| `rules` _[GatewayHTTPAPIRule](#gatewayhttpapirule) array_ | Rules |
| `healthCheckEndpoint` _string_ | Health check endpoint |




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
| `ingress` _[GatewayIngress](#gatewayingress)_ |  |




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
| `stack` _string_ |  |
| `auth` _[AuthConfig](#authconfig)_ |  |
| `deploymentStrategy` _[DeploymentStrategy](#deploymentstrategy)_ |  |
| `locking` _[LockingStrategy](#lockingstrategy)_ | Locking is intended for ledger v1 only |




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




#### ModuleProperties





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


#### ModuleStatus





_Appears in:_
- [AuthStatus](#authstatus)
- [GatewayStatus](#gatewaystatus)
- [LedgerStatus](#ledgerstatus)
- [OrchestrationStatus](#orchestrationstatus)
- [PaymentsStatus](#paymentsstatus)
- [ReconciliationStatus](#reconciliationstatus)
- [SearchStatus](#searchstatus)
- [StargateStatus](#stargatestatus)
- [WalletsStatus](#walletsstatus)
- [WebhooksStatus](#webhooksstatus)

| Field | Description |
| --- | --- |
| `ready` _boolean_ |  |
| `info` _string_ |  |
| `conditions` _[Condition](#condition) array_ |  |




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
| `auth` _[AuthConfig](#authconfig)_ |  |




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
| `encryptionKey` _string_ |  |
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
| `auth` _[AuthConfig](#authconfig)_ |  |






#### ResourceReference



ResourceReference is the Schema for the resourcereferences API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `ResourceReference`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[ResourceReferenceSpec](#resourcereferencespec)_ |  |


#### ResourceReferenceSpec



ResourceReferenceSpec defines the desired state of ResourceReference

_Appears in:_
- [ResourceReference](#resourcereference)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `gvk` _[GroupVersionKind](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#groupversionkind-v1-meta)_ |  |
| `name` _string_ |  |




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
| `batching` _[Batching](#batching)_ |  |
| `auth` _[AuthConfig](#authconfig)_ |  |




#### Settings



Settings is the Schema for the settings API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Settings`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[SettingsSpec](#settingsspec)_ |  |


#### SettingsSpec



SettingsSpec defines the desired state of Settings

_Appears in:_
- [Settings](#settings)

| Field | Description |
| --- | --- |
| `stacks` _string array_ |  |
| `key` _string_ |  |
| `value` _string_ |  |




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
- [BenthosSpec](#benthosspec)
- [BenthosStreamSpec](#benthosstreamspec)
- [BrokerTopicConsumerSpec](#brokertopicconsumerspec)
- [BrokerTopicSpec](#brokertopicspec)
- [DatabaseSpec](#databasespec)
- [GatewayHTTPAPISpec](#gatewayhttpapispec)
- [GatewaySpec](#gatewayspec)
- [LedgerSpec](#ledgerspec)
- [OrchestrationSpec](#orchestrationspec)
- [PaymentsSpec](#paymentsspec)
- [ReconciliationSpec](#reconciliationspec)
- [ResourceReferenceSpec](#resourcereferencespec)
- [SearchSpec](#searchspec)
- [StargateSpec](#stargatespec)
- [WalletsSpec](#walletsspec)
- [WebhooksSpec](#webhooksspec)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |


#### StackSpec



StackSpec defines the desired state of Stack The version of the stack can be specified using either the field `version` or the `versionsFromFile` field. The `version` field will have priority over `versionFromFile` If `versions` and `versionsFromFile` are not specified, "latest" will be used.

_Appears in:_
- [Stack](#stack)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `version` _string_ | Version allow to specify the version of the components Must be a valid docker tag |
| `versionsFromFile` _string_ | VersionsFromFile allow to specify a formance.com/Versions object which contains individual versions for each component. Must reference a valid formance.com/Versions object |
| `enableAudit` _boolean_ |  |
| `disabled` _boolean_ |  |




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
| `stack` _string_ |  |
| `serverURL` _string_ |  |
| `organizationID` _string_ |  |
| `stackID` _string_ |  |
| `auth` _[StargateAuthSpec](#stargateauthspec)_ |  |




#### StatusWithConditions





_Appears in:_
- [AuthStatus](#authstatus)
- [BenthosStatus](#benthosstatus)
- [GatewayStatus](#gatewaystatus)
- [LedgerStatus](#ledgerstatus)
- [ModuleStatus](#modulestatus)
- [OrchestrationStatus](#orchestrationstatus)
- [PaymentsStatus](#paymentsstatus)
- [ReconciliationStatus](#reconciliationstatus)
- [SearchStatus](#searchstatus)
- [StargateStatus](#stargatestatus)
- [WalletsStatus](#walletsstatus)
- [WebhooksStatus](#webhooksstatus)

| Field | Description |
| --- | --- |
| `ready` _boolean_ |  |
| `info` _string_ |  |
| `conditions` _[Condition](#condition) array_ |  |


#### URI





_Appears in:_
- [BenthosStatus](#benthosstatus)
- [BrokerTopicStatus](#brokertopicstatus)
- [DatabaseStatus](#databasestatus)
- [OrchestrationStatus](#orchestrationstatus)
- [SearchStatus](#searchstatus)



#### Versions



Versions is the Schema for the versions API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Versions`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _object (keys:string, values:string)_ |  |


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
| `stack` _string_ |  |
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
| `auth` _[AuthConfig](#authconfig)_ |  |




