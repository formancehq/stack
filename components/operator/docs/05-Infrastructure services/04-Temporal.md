### Temporal

:::info
Using Temporal is only required for stacks using the flows service. It can be ommitted now if you don't plan on using it yet, and added at a later time.
:::

The recommended way to spin up a Temporal deployment is through Temporal Cloud, or by using the official Temporal helm [chart](https://github.com/temporalio/helm-charts).


## Create the Temporal settings

In this example, you'll set up a configuration for the Orchestration of the `formance-dev` stack. This configuration will apply to all the modules of this stack.

```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-temporal-dsn
spec:
  key: temporal.dsn
  stacks:
    - 'formance-dev'
  value: temporal://dev-eu-west-1.fsdfsdf.tmprl.cloud:7233/dev-eu-west-1.fsdfsdf?
---
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-temporal-tls-crt
spec:
  key: temporal.tls.crt
  stacks:
    - 'formance-dev'
  value: |
    -----BEGIN CERTIFICATE-----
    CERTIFICATE_CONTENT
    -----END CERTIFICATE-----

---
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-temporal-tls-key
spec:
  key: temporal.tls.key
  stacks:
    - 'formance-dev'
  value: |
    -----BEGIN PRIVATE KEY-----
    PRIVATE_KEY_CONTENT
    -----END PRIVATE KEY-----
```
