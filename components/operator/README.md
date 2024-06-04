# operator

## Description

The operator allow to install formance components on a k8s cluster.

## Getting Started

### Prerequisites
- go version v1.20.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- k3d

### To Deploy on the cluster

**Create a dev cluster**

```sh
k3d cluster create internal
```

**Install CRDs**

```sh
make install
```

**Run operator**

```sh
make run
```

**Install a PostgresServer**

```sh
helm install postgres oci://registry-1.docker.io/bitnamicharts/postgresql \
  --set auth.username=formance \
  --set auth.password=formance \
  --set auth.database=formance
```

**Create a Settings object for database connection**

```sh
cat <<EOF | kubectl create -f -
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: databaseconfiguration0
spec:
  stacks:
  - "*"
  key: postgres.*.uri
  value: postgresql://formance:formance@postgres-postgresql.default:5432?disableSSLMode=true
EOF
```

**Create a stack**

```sh
cat <<EOF | kubectl create -f -
apiVersion: formance.com/v1beta1
kind: Stack
metadata:
  name: stack0
spec: {}
EOF
```

**Create a ledger**

```sh
cat <<EOF | kubectl create -f -
apiVersion: formance.com/v1beta1
kind: Ledger
metadata:
  name: stack0-ledger
spec:
  stack: stack0
EOF
```

**Create a Gateway**

```sh
cat <<EOF | kubectl create -f -
apiVersion: formance.com/v1beta1
kind: Gateway
metadata:
  name: gateway0
spec:
  stack: stack0
#  ingress:
#    host: testing.formance.dev
#    scheme: https
EOF
```

**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified. 
And it is required to have access to pull the image from the working environment. 
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/operator:tag
```

**Run locally without building/pushing image**

```sh
make run
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin 
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Settings

Settings allow to configure some parts of the deployments.
Settings are encoded as string, but under the hood, each settings can be unmarshalled on a dedicated type.

While we have some basic types (string, number, bool ...), we also have some complex structures : 
* Maps: maps are just one level dictionary with values as string. Repeat `<key>=<value>` pattern for each entry, while separating with comma.
* URIs: URIs are used each time we need to address an external resources (postgres, kafka ...). URIs are convenient to encode a lot of information in a simple, normalized format.

Available settings:

| Key                                                                                      | Type   | Example             | Description                                                          |
|------------------------------------------------------------------------------------------|--------|---------------------|----------------------------------------------------------------------|
| aws.service-account                                                                      | string |                     | AWS Role                                                             |
| postgres.`<module-name>`.uri                                                             | URI    |                     | Postgres database configuration                                      |
| elasticsearch.dsn                                                                        | URI    |                     | Elasticsearch connection URI                                         |
| temporal.dsn                                                                             | URI    |                     | Temporal URI                                                         |
| temporal.tls.crt                                                                         | string |                     | Temporal certificate                                                 |
| temporal.tls.key                                                                         | string |                     | Temporal certificate key                                             |
| broker.dsn                                                                               | URI    |                     | Broker URI                                                           |
| opentelemetry.traces.dsn                                                                 | URI    |                     | OpenTelemetry collector URI                                          |
| clear-database                                                                           | bool   | true                | Whether to remove databases on stack deletion                        |
| ledger.deployment-strategy                                                               | string | single              | Ledger deployment type                                               |
| payments.encryption-key                                                                  | string |                     | Payments data encryption key                                         |
| deployments.`<deployment-name>`.init-containers.`<container-name>`.resource-requirements | Map    | cpu=X, mem=X        |                                                                      |
| deployments.`<deployment-name>`.containers.`<container-name>`.resource-requirements      | Map    | cpu=X, mem=X        |                                                                      |
| deployments.`<deployment-name>`.init-containers.`<container-name>`.run-as                | Map    | user=X, group=X     |                                                                      |
| deployments.`<deployment-name>`.containers.`<container-name>`.run-as                     | Map    | user=X, group=X     |                                                                      |
| deployments.`<deployment-name>`.replicas                                                 | string | 2                   |                                                                      |
| caddy.image                                                                              | string |                     | Caddy image                                                          |
| registries.`<name>`.endpoint                                                             | string |                     | Specify a custom endpoint for a specific docker repository           |
| registries.`<name>`.images.`<path>`.rewrite                                              | string | formancehq/example  | Allow to rewrite the image path                                      |
| search.batching                                                                          | Map    | period=1s, count=10 | Override default batching parameters                                 |
| services.`<service-name>`.annotations                                                    | Map    |                     | Allow to specify custom annotations to apply on created k8s services |
| gateway.ingress.annotations                                                              | Map    |                     | Allow to specify custom annotations to apply on the gateway ingress  |
| logging.json                                                                             | bool   |                     | Configure services to log as json                                    |


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


## Contributing

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

