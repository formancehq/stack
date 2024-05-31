:::warning
This Module is subject to a user license.
:::

Formance Flows is a handy service that lets you quickly set up end-to-end money flows, without the headache of piecing together APIs and untangling complex systems.

With a clever compatibility model, you can easily move value between different ledgers, wallets, and payment processors. Plus, Formance Flows takes care of translating and interpreting the transactions for you.

On top of that, Formance Flows comes with flexible workflow capabilities, so you can create complex flows that account for delays or external events, as well as retry and fallback options. All of this helps make your financial management smoother and stress-free.


## Requirements

Formance Flows requires:
- **PostgreSQL**: See configuration guide [here](../05-Infrastructure%20services/02-Message%20broker.md).
- **Temporal**: See configuration guide [here](../05-Infrastructure%20services/04-Temporal.md).
- **Broker**: See configuration guide [here](/next/operator/infra/broker).

## Orchestration Object

:::info
You can find all the available parameters in [the comprehensive CRD documentation](../09-Configuration%20reference/02-Custom%20Resource%20Definitions.md#orchestration).
:::

```yaml
apiVersion: formance.com/v1beta1
kind: Orchestration
metadata:
  name: formance-dev
spec:
  stack: formance-dev
```
