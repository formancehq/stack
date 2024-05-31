:::warning
This Module is subject to a user license.
:::

## Requirements

Formance Reconciliation requires:
- **PostgreSQL**: See configuration guide [here](../05-Infrastructure%20services/01-PostgreSQL.md).
- **Ledger**: See configuration guide [here](../04-Modules/03-Ledger.md).
- **Payments**: See configuration guide [here](../04-Modules/04-Payments.md).

## Reconciliation Object

:::info
You can find all the available parameters in [the comprehensive CRD documentation](../09-Configuration%20reference/02-Custom%20Resource%20Definitions.md#reconciliation).
:::

```yaml
apiVersion: formance.com/v1beta1
kind: Reconciliation
metadata:
  name: formance-dev
spec:
  stack: formance-dev
```
