---
title: Accounts
---

Accounts are containers for assets. They can send and receive assets using transactions, the sum of which will determine their balances.

$$
balance(account) = \sum postings(destination) - \sum postings(source)
$$

The number of accounts in a ledger is unlimited.

Accounts don't have to represent a real account in your bank; they can represent anything that has meaning to your application. Once you grasp the idiomatic Formance way of building financial applications, you'll start using accounts to represent abstract concepts in your business logic, like a sale, a contract or a payment.

Accounts also do not need to be created prior being used, as submitting a transaction involving it will automatically make it exist in the ledger.

## Naming accounts

Accounts are identified by their address, which must match `^[a-zA-Z_0-9]+(:[a-zA-Z_0-9]+){0,}$`. It recommended to use colons in addresses to organize them in structured segments, i.e. `payments:001:authorization`, although this does not trigger any particular behavior on the ledger.

```
payments:001:authorizations:001
sales:001:contract
```

## Using Metadata

Accounts can bear metadata, which in a key-value format. This metadata can be used to store any information that is relevant to the account, like a reference to an external system, or a description of the account.

:::tip
Accounts metadata can also be used in Numscript transactions, see [metadata in Numscript](/ledger/reference/numscript/metadata) for more information.
:::
