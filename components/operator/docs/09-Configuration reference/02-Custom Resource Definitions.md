# API Reference

## Packages
- [formance.com/v1beta1](#formancecomv1beta1)


## formance.com/v1beta1

Package v1beta1 contains API Schema definitions for the formance v1beta1 API group.

It allow to configure a Formance stack.

A stack is composed of a [Stack](#stack) resource and some [modules](#modules).

Each module can create multiple resources following its needs. See [Other resources](#other-resources).

Various parts of the stack can be configured either using the CRD properties or using some [Settings](#settings).


Modules :
- [Auth](#auth)
- [Gateway](#gateway)
- [Ledger](#ledger)
- [Orchestration](#orchestration)
- [Payments](#payments)
- [Reconciliation](#reconciliation)
- [Search](#search)
- [Stargate](#stargate)
- [Wallets](#wallets)
- [Webhooks](#webhooks)

Other resources :
- [AuthClient](#authclient)
- [Benthos](#benthos)
- [BenthosStream](#benthosstream)
- [Broker](#broker)
- [BrokerConsumer](#brokerconsumer)
- [BrokerTopic](#brokertopic)
- [Database](#database)
- [GatewayHTTPAPI](#gatewayhttpapi)
- [ResourceReference](#resourcereference)
- [Versions](#versions)

### Main resources

#### Stack



Stack represents a formance stack.
A Stack is basically a container. It holds some global properties and
creates a namespace if not already existing.


To do more, you need to create some [modules](#modules).


The Stack resource allow to specify the version of the stack.


It can be specified using either the field `.spec.version` or the `.spec.versionsFromFile` field (Refer to the documentation of [Versions](#versions) resource.


The `version` field will have priority over `versionFromFile`.


If `versions` and `versionsFromFile` are not specified, "latest" will be used.















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Stack` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[StackSpec](#stackspec)_ |  |  |  |
| `status` _[StackStatus](#stackstatus)_ |  |  |  |



##### StackSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `debug` _boolean_ | Allow to enable debug mode on the module | false |  |
| `dev` _boolean_ | Allow to enable dev mode on the module<br />Dev mode is used to allow some application to do custom setup in development mode (allow insecure certificates for example) | false |  |
| `version` _string_ | Version allow to specify the version of the components<br />Must be a valid docker tag |  |  |
| `versionsFromFile` _string_ | VersionsFromFile allow to specify a formance.com/Versions object which contains individual versions<br />for each component.<br />Must reference a valid formance.com/Versions object |  |  |
| `enableAudit` _boolean_ | EnableAudit enable audit at the stack level.<br />Actually, it enables audit on [Gateway](#gateway) | false |  |
| `disabled` _boolean_ | Disabled indicate the stack is disabled.<br />A disabled stack disable everything<br />It just keeps the namespace and the [Database](#database) resources. | false |  |





##### StackStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |
| `modules` _string array_ | Modules register detected modules |  |  |


#### Settings



Settings represents a configurable piece of the stacks.


The purpose of this resource is to be able to configure some common settings between a set of stacks.


Example :
```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: postgres-uri
spec:
  key: postgres.ledger.uri
  stacks:
  - stack0
  value: postgresql://postgresql.formance.svc.cluster.local:5432
```


This example create a setting named `postgres-uri` targeting the stack named `stack0` and the service `ledger` (see the key `postgres.ledger.uri`).


Therefore, a [Database](#database) created for the stack `stack0` and the service named 'ledger' will use the uri `postgresql://postgresql.formance.svc.cluster.local:5432`.


Settings allow to use wildcards in keys and in stacks list.


For example, if you want to use the same database server for all the modules of a specific stack, you can write :
```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: postgres-uri
spec:
  key: postgres.*.uri # There, we use a wildcard to indicate we want to use that setting of all services of the stack `stack0`
  stacks:
  - stack0
  value: postgresql://postgresql.formance.svc.cluster.local:5432
```


Also, we could use that setting for all of our stacks using :
```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: postgres-uri
spec:
  key: postgres.*.uri # There, we use a wildcard to indicate we want to use that setting for all services of all stacks
  stacks:
  - * # There we select all the stacks
  value: postgresql://postgresql.formance.svc.cluster.local:5432
```


Some settings are really global, while some are used by specific module.


Refer to the documentation of each module and resource to discover available Settings.


##### Global settings
###### AWS account


A stack can use an AWS account for authentication.


It can be used to connect to any AWS service we could use.


It includes RDS, OpenSearch and MSK. To do so, you can create the following setting:
```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: aws-service-account
spec:
  key: aws.service-account
  stacks:
  - '*'
  value: aws-access
```
This setting instruct the operator than there is somewhere on the cluster a service account named `aws-access`.


So, each time a service has the capability to use AWS, the operator will use this service account.


The service account could look like that :
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::************:role/staging-eu-west-1-hosting-stack-access
  labels:
    formance.com/stack: any
  name: aws-access
```
You can note two things :
1. We have an annotation indicating the role arn used to connect to AWS. Refer to the AWS documentation to create this role
2. We have a label `formance.com/stack=any` indicating we are targeting all stacks.
   Refer to the documentation of [ResourceReference](#resourcereference) for further information.


###### JSON logging


You can use the setting `logging.json` with the value `true` to configure elligible service to log as json.
Example:
```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: json-logging
spec:
  key: logging.json
  stacks:
  - '*'
  value: "true"
```















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Settings` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[SettingsSpec](#settingsspec)_ |  |  |  |



##### SettingsSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stacks` _string array_ | Stacks on which the setting is applied. Can contain `*` to indicate a wildcard. |  |  |
| `key` _string_ | The setting Key. See the documentation of each module or [global settings](#global-settings) to discover them. |  |  |
| `value` _string_ | The value. It must have a specific format following the Key. |  |  |




### Modules

#### Auth



Auth represent the authentication module of a stack.


It is an OIDC compliant server.


Creating it for a stack automatically add authentication on all supported modules.


The auth service is basically a proxy to another OIDC compliant server.















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Auth` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[AuthSpec](#authspec)_ |  |  |  |
| `status` _[AuthStatus](#authstatus)_ |  |  |  |



##### AuthSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `debug` _boolean_ | Allow to enable debug mode on the module | false |  |
| `dev` _boolean_ | Allow to enable dev mode on the module<br />Dev mode is used to allow some application to do custom setup in development mode (allow insecure certificates for example) | false |  |
| `version` _string_ | Version allow to override global version defined at stack level for a specific module |  |  |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `delegatedOIDCServer` _[DelegatedOIDCServerConfiguration](#delegatedoidcserverconfiguration)_ | Contains information about a delegated authentication server to use to delegate authentication |  |  |
| `signingKey` _string_ | Allow to override the default signing key used to sign JWT tokens. |  |  |
| `signingKeyFromSecret` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#secretkeyselector-v1-core)_ | Allow to override the default signing key used to sign JWT tokens using a k8s secret |  |  |
| `enableScopes` _boolean_ | Allow to enable scopes usage on authentication.<br /><br />If not enabled, each service will check the authentication but will not restrict access following scopes.<br />in this case, if authenticated, it is ok. | false |  |

###### DelegatedOIDCServerConfiguration



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `issuer` _string_ | Issuer is the url of the delegated oidc server |  |  |
| `clientID` _string_ | ClientID is the client id to use for authentication |  |  |
| `clientSecret` _string_ | ClientSecret is the client secret to use for authentication |  |  |





##### AuthStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |
| `clients` _string array_ | Clients contains the list of clients created using [AuthClient](#authclient) |  |  |


#### Gateway



Gateway is the Schema for the gateways API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Gateway` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[GatewaySpec](#gatewayspec)_ |  |  |  |
| `status` _[GatewayStatus](#gatewaystatus)_ |  |  |  |



##### GatewaySpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `debug` _boolean_ | Allow to enable debug mode on the module | false |  |
| `dev` _boolean_ | Allow to enable dev mode on the module<br />Dev mode is used to allow some application to do custom setup in development mode (allow insecure certificates for example) | false |  |
| `version` _string_ | Version allow to override global version defined at stack level for a specific module |  |  |
| `ingress` _[GatewayIngress](#gatewayingress)_ | Allow to customize the generated ingress |  |  |

###### GatewayIngress



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `host` _string_ | Indicates the hostname on which the stack will be served.<br />Example : `formance.example.com` |  |  |
| `scheme` _string_ | Indicate the scheme.<br /><br />Actually, It should be `https` unless you know what you are doing. | https |  |
| `annotations` _object (keys:string, values:string)_ | Custom annotations to add on the ingress |  |  |
| `tls` _[GatewayIngressTLS](#gatewayingresstls)_ | Allow to customize the tls part of the ingress |  |  |

###### GatewayIngressTLS



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `secretName` _string_ | Specify the secret name used for the tls configuration on the ingress |  |  |





##### GatewayStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |
| `syncHTTPAPIs` _string array_ | Detected http apis. See [GatewayHTTPAPI](#gatewayhttpapi) |  |  |
| `authEnabled` _boolean_ | Indicates if a [Auth](#auth) module has been detected. | false |  |


#### Ledger



Ledger is the module allowing to install a ledger instance.


The ledger is actually a stateful application on the writer part.
So we cannot scale the ledger as we want without prior configuration.


So, the ledger can run in two modes :
* single instance: Only one instance will be deployed. We cannot scale in that mode.
* single writer / multiple reader: In this mode, we will have a single writer and multiple readers if needed.


Use setting `ledger.deployment-strategy` with either the value :
* single : For the single instance mode.
* single-writer: For the single writer / multiple reader mode.
  Under the hood, the operator create two deployments and force the scaling of the writer to stay at 1.
  Then you can scale the deployment of the reader to the value you want.















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Ledger` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[LedgerSpec](#ledgerspec)_ |  |  |  |
| `status` _[LedgerStatus](#ledgerstatus)_ |  |  |  |



##### LedgerSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `debug` _boolean_ | Allow to enable debug mode on the module | false |  |
| `dev` _boolean_ | Allow to enable dev mode on the module<br />Dev mode is used to allow some application to do custom setup in development mode (allow insecure certificates for example) | false |  |
| `version` _string_ | Version allow to override global version defined at stack level for a specific module |  |  |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `deploymentStrategy` _[DeploymentStrategy](#deploymentstrategy)_ | Deprecated. | single |  |
| `locking` _[LockingStrategy](#lockingstrategy)_ | Locking is intended for ledger v1 only |  |  |

###### DeploymentStrategy

_Underlying type:_ _string_


















###### LockingStrategy



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `strategy` _string_ |  | memory |  |
| `redis` _[LockingStrategyRedisConfig](#lockingstrategyredisconfig)_ |  |  |  |

###### LockingStrategyRedisConfig



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `uri` _string_ |  |  |  |
| `tls` _boolean_ |  | false |  |
| `insecure` _boolean_ |  | false |  |
| `duration` _string_ |  |  |  |
| `retry` _string_ |  |  |  |





##### LedgerStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |


#### Orchestration



Orchestration is the Schema for the orchestrations API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Orchestration` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[OrchestrationSpec](#orchestrationspec)_ |  |  |  |
| `status` _[OrchestrationStatus](#orchestrationstatus)_ |  |  |  |



##### OrchestrationSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `debug` _boolean_ | Allow to enable debug mode on the module | false |  |
| `dev` _boolean_ | Allow to enable dev mode on the module<br />Dev mode is used to allow some application to do custom setup in development mode (allow insecure certificates for example) | false |  |
| `version` _string_ | Version allow to override global version defined at stack level for a specific module |  |  |





##### OrchestrationStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |
| `temporalURI` _string_ |  |  | Type: string <br /> |


#### Payments



Payments is the Schema for the payments API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Payments` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[PaymentsSpec](#paymentsspec)_ |  |  |  |
| `status` _[PaymentsStatus](#paymentsstatus)_ |  |  |  |



##### PaymentsSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `debug` _boolean_ | Allow to enable debug mode on the module | false |  |
| `dev` _boolean_ | Allow to enable dev mode on the module<br />Dev mode is used to allow some application to do custom setup in development mode (allow insecure certificates for example) | false |  |
| `version` _string_ | Version allow to override global version defined at stack level for a specific module |  |  |
| `encryptionKey` _string_ |  |  |  |





##### PaymentsStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |


#### Reconciliation



Reconciliation is the Schema for the reconciliations API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Reconciliation` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ReconciliationSpec](#reconciliationspec)_ |  |  |  |
| `status` _[ReconciliationStatus](#reconciliationstatus)_ |  |  |  |



##### ReconciliationSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `debug` _boolean_ | Allow to enable debug mode on the module | false |  |
| `dev` _boolean_ | Allow to enable dev mode on the module<br />Dev mode is used to allow some application to do custom setup in development mode (allow insecure certificates for example) | false |  |
| `version` _string_ | Version allow to override global version defined at stack level for a specific module |  |  |





##### ReconciliationStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |


#### Search



Search is the Schema for the searches API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Search` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[SearchSpec](#searchspec)_ |  |  |  |
| `status` _[SearchStatus](#searchstatus)_ |  |  |  |



##### SearchSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `debug` _boolean_ | Allow to enable debug mode on the module | false |  |
| `dev` _boolean_ | Allow to enable dev mode on the module<br />Dev mode is used to allow some application to do custom setup in development mode (allow insecure certificates for example) | false |  |
| `version` _string_ | Version allow to override global version defined at stack level for a specific module |  |  |
| `batching` _[Batching](#batching)_ |  |  |  |

###### Batching



Batching allow to define custom batching configuration















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `count` _integer_ | Count indicates the number of messages that can be kept in memory before being flushed to ElasticSearch |  |  |
| `period` _string_ | Period indicates the maximum duration messages can be kept in memory before being flushed to ElasticSearch |  |  |





##### SearchStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |
| `elasticSearchURI` _string_ |  |  | Type: string <br /> |
| `topicCleaned` _boolean_ | TopicCleaned is used to flag stacks where the topics have been cleaned (still search-ledgerv2 and co consumers) | false |  |


#### Stargate



Stargate is the Schema for the stargates API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Stargate` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[StargateSpec](#stargatespec)_ |  |  |  |
| `status` _[StargateStatus](#stargatestatus)_ |  |  |  |



##### StargateSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `debug` _boolean_ | Allow to enable debug mode on the module | false |  |
| `dev` _boolean_ | Allow to enable dev mode on the module<br />Dev mode is used to allow some application to do custom setup in development mode (allow insecure certificates for example) | false |  |
| `version` _string_ | Version allow to override global version defined at stack level for a specific module |  |  |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `serverURL` _string_ |  |  |  |
| `organizationID` _string_ |  |  |  |
| `stackID` _string_ |  |  |  |
| `auth` _[StargateAuthSpec](#stargateauthspec)_ |  |  |  |

###### StargateAuthSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `clientID` _string_ |  |  |  |
| `clientSecret` _string_ |  |  |  |
| `issuer` _string_ |  |  |  |





##### StargateStatus



StargateStatus defines the observed state of Stargate















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |


#### Wallets



Wallets is the Schema for the wallets API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Wallets` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[WalletsSpec](#walletsspec)_ |  |  |  |
| `status` _[WalletsStatus](#walletsstatus)_ |  |  |  |



##### WalletsSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `debug` _boolean_ | Allow to enable debug mode on the module | false |  |
| `dev` _boolean_ | Allow to enable dev mode on the module<br />Dev mode is used to allow some application to do custom setup in development mode (allow insecure certificates for example) | false |  |
| `version` _string_ | Version allow to override global version defined at stack level for a specific module |  |  |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |





##### WalletsStatus



WalletsStatus defines the observed state of Wallets















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |


#### Webhooks



Webhooks is the Schema for the webhooks API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Webhooks` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[WebhooksSpec](#webhooksspec)_ |  |  |  |
| `status` _[WebhooksStatus](#webhooksstatus)_ |  |  |  |



##### WebhooksSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `debug` _boolean_ | Allow to enable debug mode on the module | false |  |
| `dev` _boolean_ | Allow to enable dev mode on the module<br />Dev mode is used to allow some application to do custom setup in development mode (allow insecure certificates for example) | false |  |
| `version` _string_ | Version allow to override global version defined at stack level for a specific module |  |  |





##### WebhooksStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |


### Other resources

#### AuthClient



AuthClient allow to create OAuth2/OIDC clients on the auth server (see [Auth](#auth))















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `AuthClient` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[AuthClientSpec](#authclientspec)_ |  |  |  |
| `status` _[AuthClientStatus](#authclientstatus)_ |  |  |  |



##### AuthClientSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `id` _string_ | ID indicates the client id<br />It must be used with oauth2 `client_id` parameter |  |  |
| `public` _boolean_ | Public indicate whether a client is confidential or not.<br />Confidential clients are clients which the secret can be kept secret...<br />As opposed to public clients which cannot have a secret (application single page for example) | false |  |
| `description` _string_ | Description represents an optional description of the client |  |  |
| `redirectUris` _string array_ | RedirectUris allow to list allowed redirect uris for the client |  |  |
| `postLogoutRedirectUris` _string array_ | RedirectUris allow to list allowed post logout redirect uris for the client |  |  |
| `scopes` _string array_ | Scopes allow to five some scope to the client |  |  |
| `secret` _string_ | Secret allow to configure a secret for the client.<br />It is not required as some client could use some oauth2 flows which does not requires a client secret |  |  |





##### AuthClientStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |


#### Benthos



Benthos is the Schema for the benthos API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Benthos` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[BenthosSpec](#benthosspec)_ |  |  |  |
| `status` _[BenthosStatus](#benthosstatus)_ |  |  |  |



##### BenthosSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `debug` _boolean_ | Allow to enable debug mode on the module | false |  |
| `dev` _boolean_ | Allow to enable dev mode on the module<br />Dev mode is used to allow some application to do custom setup in development mode (allow insecure certificates for example) | false |  |
| `resourceRequirements` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#resourcerequirements-v1-core)_ |  |  |  |
| `batching` _[Batching](#batching)_ |  |  |  |
| `initContainers` _[Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#container-v1-core) array_ |  |  |  |

###### Batching



Batching allow to define custom batching configuration















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `count` _integer_ | Count indicates the number of messages that can be kept in memory before being flushed to ElasticSearch |  |  |
| `period` _string_ | Period indicates the maximum duration messages can be kept in memory before being flushed to ElasticSearch |  |  |





##### BenthosStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |
| `elasticSearchURI` _string_ |  |  | Type: string <br /> |


#### BenthosStream



BenthosStream is the Schema for the benthosstreams API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `BenthosStream` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[BenthosStreamSpec](#benthosstreamspec)_ |  |  |  |
| `status` _[BenthosStreamStatus](#benthosstreamstatus)_ |  |  |  |



##### BenthosStreamSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `data` _string_ |  |  |  |
| `name` _string_ |  |  |  |





##### BenthosStreamStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |


#### Broker



Broker is the Schema for the brokers API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Broker` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[BrokerSpec](#brokerspec)_ |  |  |  |
| `status` _[BrokerStatus](#brokerstatus)_ |  |  |  |



##### BrokerSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |





##### BrokerStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |
| `uri` _string_ |  |  | Type: string <br /> |
| `mode` _[Mode](#mode)_ | Mode indicating the configuration of the nats streams<br />Two modes are defined :<br />* OneStreamByService: In this case, each service will have a dedicated stream created<br />* OneStreamByStack: In this case, a stream will be created for the stack and each service will use a specific subject inside this stream |  | Enum: [OneStreamByService OneStreamByStack] <br /> |
| `streams` _string array_ | Streams list streams created when Mode == ModeOneStreamByService |  |  |

###### Mode

_Underlying type:_ _string_

Mode defined how streams are created on the broker (mainly nats)

















#### BrokerConsumer



BrokerConsumer is the Schema for the brokerconsumers API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `BrokerConsumer` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[BrokerConsumerSpec](#brokerconsumerspec)_ |  |  |  |
| `status` _[BrokerConsumerStatus](#brokerconsumerstatus)_ |  |  |  |



##### BrokerConsumerSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `services` _string array_ |  |  |  |
| `queriedBy` _string_ |  |  |  |
| `name` _string_ | As the name is optional, if not provided, the name will be the QueriedBy property<br />This is only applied when using one stream by stack see Mode |  |  |





##### BrokerConsumerStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |


#### BrokerTopic



BrokerTopic is the Schema for the brokertopics API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `BrokerTopic` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[BrokerTopicSpec](#brokertopicspec)_ |  |  |  |
| `status` _[BrokerTopicStatus](#brokertopicstatus)_ |  |  |  |



##### BrokerTopicSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `service` _string_ |  |  |  |





##### BrokerTopicStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |


#### Database



Database represent a concrete database on a PostgreSQL server, it is created by modules requiring a database ([Ledger](#ledger) for example).


It uses the settings `postgres.<module-name>.uri` which must have the following uri format: `postgresql://[<username>@<password>]@<host>/<db-name>`
Additionally, the uri can define a query param `secret` indicating a k8s secret, than must be used to retrieve database credentials.


On creation, the reconciler behind the Database object will create the database on the postgresql server using a k8s job.
On Deletion, by default, the reconciler will let the database untouched.
You can allow the reconciler to drop the database on the server by using the [Settings](#settings) `clear-database` with the value `true`.
If you use that setting, the reconciler will use another job to drop the database.
Be careful, no backup are performed!


Database resource honors `aws.service-account` setting, so, you can create databases on an AWS server if you need.
See [AWS accounts](#aws-account)


Once a database is fully configured, it retains the postgres uri used.
If the setting indicating the server uri changed, the Database object will set the field `.status.outOfSync` to true
and will not change anything.


Therefore, to switch to a new server, you must change the setting value, then drop the Database object.
It will be recreated with correct uri.















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Database` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[DatabaseSpec](#databasespec)_ |  |  |  |
| `status` _[DatabaseStatus](#databasestatus)_ |  |  |  |



##### DatabaseSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `service` _string_ | Service is a discriminator for the created database.<br />Actually, it will be the module name (ledger, payments...).<br />Therefore, the created database will be named `<stack-name><service>` |  |  |
| `debug` _boolean_ |  | false |  |





##### DatabaseStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |
| `uri` _string_ |  |  | Type: string <br /> |
| `database` _string_ | The generated database name |  |  |
| `outOfSync` _boolean_ | OutOfSync indicates than a settings changed the uri of the postgres server<br />The Database object need to be removed to be recreated |  |  |


#### GatewayHTTPAPI



GatewayHTTPAPI is the Schema for the HTTPAPIs API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `GatewayHTTPAPI` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[GatewayHTTPAPISpec](#gatewayhttpapispec)_ |  |  |  |
| `status` _[GatewayHTTPAPIStatus](#gatewayhttpapistatus)_ |  |  |  |



##### GatewayHTTPAPISpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `name` _string_ | Name indicates prefix api |  |  |
| `rules` _[GatewayHTTPAPIRule](#gatewayhttpapirule) array_ | Rules |  |  |
| `healthCheckEndpoint` _string_ | Health check endpoint |  |  |

###### GatewayHTTPAPIRule



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `path` _string_ |  |  |  |
| `methods` _string array_ |  |  |  |
| `secured` _boolean_ |  | false |  |





##### GatewayHTTPAPIStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |
| `ready` _boolean_ |  |  |  |


#### ResourceReference



ResourceReference is a special resources used to refer to externally created resources.


It includes k8s service accounts and secrets.


Why? Because the operator create a namespace by stack, so, a stack does not have access to secrets and service
accounts created externally.


A ResourceReference is created by other resource who need to use a specific secret or service account.
For example, if you want to use a secret for your database connection (see [Database](#database), you will
create a setting indicating a secret name. You will need to create this secret yourself, and you will put this
secret inside the namespace you want (`default` maybe).


The Database reconciler will create a ResourceReference looking like that :
```
apiVersion: formance.com/v1beta1
kind: ResourceReference
metadata:
  name: jqkuffjxcezj-qlii-auth-postgres
  ownerReferences:
  - apiVersion: formance.com/v1beta1
    blockOwnerDeletion: true
    controller: true
    kind: Database
    name: jqkuffjxcezj-qlii-auth
    uid: 2cc4b788-3ffb-4e3d-8a30-07ed3941c8d2
spec:
  gvk:
    group: ""
    kind: Secret
    version: v1
  name: postgres
  stack: jqkuffjxcezj-qlii
status:
  ...
```
This reconciler behind this ResourceReference will search, in all namespaces, for a secret named "postgres".
The secret must have a label `formance.com/stack` with the value matching either a specific stack or `any` to target any stack.


Once the reconciler has found the secret, it will copy it inside the stack namespace, allowing the ResourceReconciler owner to use it.















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `ResourceReference` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ResourceReferenceSpec](#resourcereferencespec)_ |  |  |  |
| `status` _[ResourceReferenceStatus](#resourcereferencestatus)_ |  |  |  |



##### ResourceReferenceSpec



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `stack` _string_ | Stack indicates the stack on which the module is installed |  |  |
| `gvk` _[GroupVersionKind](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#groupversionkind-v1-meta)_ |  |  |  |
| `name` _string_ |  |  |  |





##### ResourceReferenceStatus



















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready indicates if the resource is seen as completely reconciled |  |  |
| `info` _string_ | Info can contain any additional like reconciliation errors |  |  |
| `syncedResource` _string_ |  |  |  |
| `hash` _string_ |  |  |  |


#### Versions



Versions is the Schema for the versions API















| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1` | | |
| `kind` _string_ | `Versions` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _object (keys:string, values:string)_ |  |  |  |





