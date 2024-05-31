:::warning
This Module is subject to a user license.
:::

Wallets is a fully managed, white-label wallet service to materialize and spend users funds. It comes with built-in support for multi-currency balances and temporary holds capabilities (and upcoming support for reserved funds and expirable fungibles). It is built on top of the Formance Ledger service and is designed to provide an easy way to add wallets capabilities to your application without having to worry about the underlying transactions structure, providing an opinionated model implementation.

## Requirements

Formance Wallets requires:
- **Ledger**: See configuration guide [here](../04-Modules/03-Ledger.md).

## Wallets Object

:::info
You can find all the available parameters in [the comprehensive CRD documentation](../09-Configuration%20reference/02-Custom%20Resource%20Definitions.md#wallets).
:::

```yaml
apiVersion: formance.com/v1beta1
kind: Wallets
metadata:
  name: formance-dev
spec:
  stack: formance-dev
```
