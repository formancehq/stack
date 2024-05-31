# Elasticsearch

_Version 7.10 or higher is required._

The recommended way to spin up an Elasticsearch / OpenSearch deployment is to use your cloud provider's managed service, or Elastic Cloud itself.

:::info
If you are using Elastic Cloud, make sure to use a deployment in a network with low latency to your cluster.
:::

## Create the ElasticSearch / OpenSearch settings

In this example, you'll set up a configuration for the Broker of the `formance-dev` stack. This configuration will apply to all the modules of this stack.

```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-es
spec:
  key: elasticsearch.dsn
  stacks:
    - 'formance-dev'
  value: https://es.formance-system.svc:443?
```
