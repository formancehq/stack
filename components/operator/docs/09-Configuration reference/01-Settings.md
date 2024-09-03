The Settings CRD is one of the most important CRDs in Operator v2. It enables all the necessary adjustments so that the Operator can adapt to your usage and environment.

Settings are encoded as string, but under the hood, each setting can be unmarshalled to a specific type.

While we have some basic types (string, number, bool ...), we also have some complex structures:
* Maps: maps are just one level dictionary with values as string. Repeat `<key>=<value>` pattern for each entry, while separating with comma.
* URIs: URIs are used each time we need to address an external resource (postgres, kafka ...). URIs are convenient to encode a lot of information in a simple, normalized format.

## Available settings

| Key                                                                                      | Type   | Example                                    | Description                                                                                                                    |
|------------------------------------------------------------------------------------------|--------|--------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------|
| aws.service-account                                                                      | string |                                            | AWS Role                                                                                                                       |
| postgres.`<module-name>`.uri                                                             | URI    |                                            | Postgres database configuration                                                                                                |
| elasticsearch.dsn                                                                        | URI    |                                            | Elasticsearch connection URI                                                                                                   |
| temporal.dsn                                                                             | URI    |                                            | Temporal URI                                                                                                                   |
| temporal.tls.crt                                                                         | string |                                            | Temporal certificate                                                                                                           |
| temporal.tls.key                                                                         | string |                                            | Temporal certificate key                                                                                                       |
| broker.dsn                                                                               | URI    |                                            | Broker URI                                                                                                                     |
| opentelemetry.traces.dsn                                                                 | URI    |                                            | OpenTelemetry collector URI                                                                                                    |
| opentelemetry.traces.resource-attributes                                                 | Map | key1=value1,key2=value2 | Opentelemetry additional resource attributes                                                                                   |
| clear-database                                                                           | bool   | true                                       | Whether to remove databases on stack deletion                                                                                  |
| ledger.deployment-strategy                                                               | string | single                                     | Ledger deployment type                                                                                                         |
| ledger.logs.max-batch-size                                                               | Int    | 1024                                       | Ledger logs batching max size                                                                                                  |
| payments.encryption-key                                                                  | string |                                            | Payments data encryption key                                                                                                   |
| deployments.`<deployment-name>`.init-containers.`<container-name>`.resource-requirements | Map    | cpu=X, mem=X                               |                                                                                                                                |
| deployments.`<deployment-name>`.containers.`<container-name>`.resource-requirements      | Map    | cpu=X, mem=X                               |                                                                                                                                |
| deployments.`<deployment-name>`.init-containers.`<container-name>`.run-as                | Map    | user=X, group=X                            |                                                                                                                                |
| deployments.`<deployment-name>`.containers.`<container-name>`.run-as                     | Map    | user=X, group=X                            |                                                                                                                                |
| deployments.`<deployment-name>`.replicas                                                 | string | 2                                          |                                                                                                                                |
| caddy.image                                                                              | string |                                            | Caddy image                                                                                                                    |
| registries.`<name>`.endpoint                                                             | string |                                            | Specify a custom endpoint for a specific docker repository                                                                     |
| registries.`<name>`.images.`<path>`.rewrite                                              | string | formancehq/example                         | Allow to rewrite the image path                                                                                                |
| search.batching                                                                          | Map    | period=1s, count=10                        | Override default batching parameters                                                                                           |
| services.`<service-name>`.annotations                                                    | Map    |                                            | Allow to specify custom annotations to apply on created k8s services                                                           |
| gateway.ingress.annotations                                                              | Map    |                                            | Allow to specify custom annotations to apply on the gateway ingress                                                            |
| logging.json                                                                             | bool   |                                            | Configure services to log as json                                                                                              |
| modules.`<module-name>`.database.connection-pool                                         | Map    | max-idle=10, max-idle-time=10, max-open=10 | Configure database connection pool for each module. See [Golang documentation](https://go.dev/doc/database/manage-connections) |
| orchestration.max-parallel-activities                                                    | Int    | 10                                         | Configure max parallel temporal activities on orchestration workers                                                            |


### Postgres URI format

Scheme: postgresql

Query params :

| Name           | Type   | Default | Description                                          |
|----------------|--------|---------|------------------------------------------------------|
| secret         | string |         | Specify a secret where credentials are defined       |
| disableSSLMode | bool   | false   | Disable SSL on Postgres connection                   |

### ElasticSearch URI format

Scheme: elasticsearch

Query params :

| Name   | Type   | Default | Description                                    |
|--------|--------|---------|------------------------------------------------|
| secret | string |         | Specify a secret where credentials are defined |

### Temporal URI format

Scheme : temporal

Path : Match the temporal namespace

Query params :

| Name   | Type   | Default | Description                                              |
|--------|--------|---------|----------------------------------------------------------|
| secret | string |         | Specify a secret where temporal certificates are defined |

### Broker URI format

Scheme : nats | kafka

#### Broker URI format (nats)

Scheme: nats

Query params :

| Name     | Type   | Default | Description                                                               |
|----------|--------|---------|---------------------------------------------------------------------------|
| replicas | number | 1       | Specify the number of replicas to configure on newly created nats streams |

#### Broker URI format (kafka)

Scheme: kafka

Query params :

| Name             | Type   | Default | Description                                    |
|------------------|--------|---------|------------------------------------------------|
| saslEnabled      | bool   | false   | Specify is sasl authentication must be enabled |
| saslUsername     | string |         | Username on sasl authentication                |
| saslPassword     | string |         | Password on sasl authentication                |
| saslMechanism    | string |         | Mechanism on sasl authentication               |
| saslSCRAMSHASize | string |         | SCRAM SHA size on sasl authentication          |
| tls              | bool   | false   | Whether enable ssl or not                      |



The process is always the same: you create a YAML file, submit it to Kubernetes, and the Operator takes care of the rest.
All the values present in the `Metadata` section are not used by the Operator. Conversely, the `Spec` section is used to define the Operator's parameters.
You will always find 3 parameters there:
- **stacks**: defines the stacks that should use this configuration (you can put a `*` to indicate that all stacks should use this configuration)
- **key**: defines the key of the configuration (you can put a `*` so that it applies to all services)
- **value**: defines the value of the configuration

## Examples
### Define PostgreSQL clusters
In this example, you will set up a configuration for a PostgreSQL cluster that will be used only by the `formance-dev` stack but will apply to all the modules of this stack.
Thus, the different modules of the Stack will use this PostgreSQL cluster while being isolated in their own database.

:::info
This database is created following the format: `{stackName}-{module}`
:::

```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-postgres-uri
spec:
  key: postgres.*.uri
  stacks:
    - 'formance-dev'
  value: postgresql://formance:formance@postgresql.formance-system.svc:5432?disableSSLMode=true
```

### Use AWS IAM Role
In this example, you'll use an AWS IAM role to connect to the database. The `formance-dev` stack will use this configuration.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: aws-rds-access-role
  namespace: formance-system
  labels:
    formance.com/stack: any
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::AWS_ACCOUNT_ID:role/AWS_ROLE_NAME
---
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-postgres-uri
spec:
  key: postgres.*.uri
  stacks:
    - 'formance-dev'
  value: postgresql://formance@postgresql.formance-system.svc:5432
 ```

### Define module resource requests
In this example, you'll set up a configuration for the resource requests of the `formance-dev` stack. This configuration will apply to all the modules of this stack.

```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-resource-requests
spec:
  key: deployments.*.containers.*.resource-requirements.requests
  stacks:
    - 'formance-dev'
  value: cpu=10m,memory=100Mi
```

### Define a Broker
In this example, you'll set up a configuration for the Broker of the `formance-dev` stack. This configuration will apply to all the modules of this stack.

```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-broker
spec:
  key: broker.dsn
  stacks:
    - 'formance-dev'
  value: nats://nats.formance-system.svc:4222?replicas=3
```

### Define a OpenTelemetry Collector

In this example, you'll set up a configuration to send traces and metrics to an OpenTelemetry collector. This configuration will apply to all modules in this stack.

```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: stacks-otel-collector
spec:
  key: opentelemetry.*.dsn
  stacks:
    - "formance-dev"
  value: grpc://opentelemetry-collector.formance-system.svc:4317?insecure=true
```

<!-- ### Define a Replicas -->
<!-- In this example, we'll set up a configuration to define the number of replicas for the `formance-dev` stack. This configuration will apply to all modules in this stack. -->

<!-- ```yaml -->
<!-- apiVersion: formance.com/v1beta1 -->
<!-- kind: Settings -->
<!-- metadata: -->
<!--   name: stacks-replicas -->
<!-- spec: -->
<!--   key: replicas -->
<!--   stacks: -->
<!--     - "formance-dev" -->
<!--   value: "2" -->
<!-- ``` -->