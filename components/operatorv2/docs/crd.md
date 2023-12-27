# API Reference

## Packages
- [formance.com/v1beta1](#formancecomv1beta1)


## formance.com/v1beta1

Package v1beta1 contains API Schema definitions for the formance v1beta1 API group

### Resource Types
- [Auth](#auth)
- [BrokerConfiguration](#brokerconfiguration)
- [Database](#database)
- [DatabaseConfiguration](#databaseconfiguration)
- [Gateway](#gateway)
- [HTTPAPI](#httpapi)
- [Ledger](#ledger)
- [OpenTelemetryConfiguration](#opentelemetryconfiguration)
- [Stack](#stack)
- [Topic](#topic)
- [TopicQuery](#topicquery)



#### Auth



Auth is the Schema for the auths API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Auth`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[AuthSpec](#authspec)_ |  |


#### AuthSpec



AuthSpec defines the desired state of Auth

_Appears in:_
- [Auth](#auth)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `stack` _string_ |  |
| `delegatedOIDCServer` _[DelegatedOIDCServerConfiguration](#delegatedoidcserverconfiguration)_ |  |
| `resourceProperties` _[ResourceProperties](#resourceproperties)_ |  |




#### BrokerConfigSpec



BrokerConfigSpec defines the desired state of BrokerConfig

_Appears in:_
- [BrokerConfiguration](#brokerconfiguration)

| Field | Description |
| --- | --- |
| `kafka` _[KafkaConfig](#kafkaconfig)_ |  |
| `nats` _[NatsConfig](#natsconfig)_ |  |




#### BrokerConfiguration



BrokerConfiguration is the Schema for the brokerconfigurations API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `BrokerConfiguration`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[BrokerConfigSpec](#brokerconfigspec)_ |  |




#### CreatedDatabase





_Appears in:_
- [DatabaseStatus](#databasestatus)

| Field | Description |
| --- | --- |
| `port` _integer_ |  |
| `host` _string_ |  |
| `username` _string_ |  |
| `password` _string_ |  |
| `debug` _boolean_ |  |
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
| `debug` _boolean_ |  |
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
- [LedgerSpec](#ledgerspec)
- [StackSpec](#stackspec)

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
| `ingress` _[GatewayIngress](#gatewayingress)_ |  |
| `enableAudit` _boolean_ |  |
| `resourceProperties` _[ResourceProperties](#resourceproperties)_ |  |




#### HTTPAPI



todo: rename to API HTTPAPI is the Schema for the HTTPAPIs API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `HTTPAPI`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[HTTPAPISpec](#httpapispec)_ |  |


#### HTTPAPISpec



HTTPAPISpec defines the desired state of HTTPAPI

_Appears in:_
- [HTTPAPI](#httpapi)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `name` _string_ | Name indicates prefix api |
| `secured` _boolean_ | Secured indicate if the service is able to handle security |
| `hasVersionEndpoint` _boolean_ | HasVersionEndpoint indicates if the service has a /_info endpoint |
| `liveness` _[Liveness](#liveness)_ | Liveness indicates if the service has a /_health(check) endpoint |
| `annotations` _object (keys:string, values:string)_ |  |
| `portName` _string_ | Port name of the container |




#### KafkaConfig





_Appears in:_
- [BrokerConfigSpec](#brokerconfigspec)

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
| `stack` _string_ |  |
| `deploymentStrategy` _[DeploymentStrategy](#deploymentstrategy)_ |  |
| `resourceProperties` _[ResourceProperties](#resourceproperties)_ |  |




#### Liveness

_Underlying type:_ _integer_



_Appears in:_
- [HTTPAPISpec](#httpapispec)



#### MetricsSpec





_Appears in:_
- [OpenTelemetryConfigurationSpec](#opentelemetryconfigurationspec)

| Field | Description |
| --- | --- |
| `otlp` _[OtlpSpec](#otlpspec)_ |  |


#### NatsConfig





_Appears in:_
- [BrokerConfigSpec](#brokerconfigspec)

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
| `stack` _string_ |  |
| `traces` _[TracesSpec](#tracesspec)_ |  |
| `metrics` _[MetricsSpec](#metricsspec)_ |  |




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




#### ResourceProperties





_Appears in:_
- [AuthSpec](#authspec)
- [GatewaySpec](#gatewayspec)
- [LedgerSpec](#ledgerspec)

| Field | Description |
| --- | --- |
| `request` _[resource](#resource)_ |  |
| `limits` _[resource](#resource)_ |  |


#### Stack



Stack is the Schema for the stacks API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Stack`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[StackSpec](#stackspec)_ |  |


#### StackSpec



StackSpec defines the desired state of Stack

_Appears in:_
- [Stack](#stack)

| Field | Description |
| --- | --- |
| `debug` _boolean_ |  |
| `dev` _boolean_ |  |
| `version` _string_ |  |




#### Topic



Topic is the Schema for the topics API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `Topic`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[TopicSpec](#topicspec)_ |  |


#### TopicQuery



TopicQuery is the Schema for the topicqueries API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `formance.com/v1beta1`
| `kind` _string_ | `TopicQuery`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[TopicQuerySpec](#topicqueryspec)_ |  |


#### TopicQuerySpec



TopicQuerySpec defines the desired state of TopicQuery

_Appears in:_
- [TopicQuery](#topicquery)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `service` _string_ |  |
| `queriedBy` _string_ |  |




#### TopicSpec



TopicSpec defines the desired state of Topic

_Appears in:_
- [Topic](#topic)

| Field | Description |
| --- | --- |
| `stack` _string_ |  |
| `service` _string_ |  |
| `queries` _string array_ |  |




#### TracesSpec





_Appears in:_
- [OpenTelemetryConfigurationSpec](#opentelemetryconfigurationspec)

| Field | Description |
| --- | --- |
| `otlp` _[OtlpSpec](#otlpspec)_ |  |


