:::warning
This Module is subject to a user license.
:::

## Requirements

Formance Webhooks requires:
- **PostgreSQL**: See configuration guide [here](../05-Infrastructure%20services/02-Message%20broker.md).
- **Broker**: See configuration guide [here](/next/operator/infra/broker).

## Webhooks Object

:::info
You can find all the available parameters in [the comprehensive CRD documentation](../09-Configuration%20reference/02-Custom%20Resource%20Definitions.md#webhooks).
:::

```yaml
apiVersion: formance.com/v1beta1
kind: Webhooks
metadata:
  name: formance-dev
spec:
  stack: formance-dev
```
