# Message Broker

The broker is the messaging system that the Formance stack uses to communicate between its modules. The Formance stack supports both [NATS](https://nats.io/) and [Kafka](https://kafka.apache.org/) as brokers.

The broker sends messages between the different modules of the Formance stack. The producers are Ledger, Gateway, and Payments. The consumers are Benthos, Orchestration and Webhooks. Benthos is used to transform messages so that they can be ingested by the ElasticSearch / OpenSearch cluster to power the Search module.

:::info
This stream is created following the format: `{stackName}-{module}`
:::

## Option 1: NATS

_Version 2.6 or higher with Jetstream is required._

The recommended way to spin up a NATS deployment is through the official NATS helm [chart](https://artifacthub.io/packages/helm/nats/nats).

:::info
Depending on your setup, you may need to activate Jetsream mode on your NATS deployment manually. Jetstream is required for the resources deployed by the Formance Operator to function properly.
:::

### Create the NATS settings

In this example, you'll set up a configuration for the Broker of the `formance-dev` stack. This configuration will apply to all the modules of this stack.

```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-nats
spec:
  key: broker.dsn
  stacks:
    - 'formance-dev'
  value: nats://nats.formance-system.svc:4222?replicas=3
```

## Option 2: Kafka

The Formance stack also supports Kafka as a broker. To use Kafka, you need to set up a Kafka cluster and configure the Formance Operator to use it.

### Create the Kafka settings

In this example, you'll set up a configuration for the Broker of the `formance-dev` stack. This configuration will apply to all the modules of this stack.

```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-nats
spec:
  key: broker.dsn
  stacks:
    - 'formance-dev'
  value: kafka://kafka.formance-system.svc:9092
```