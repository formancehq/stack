---
title: Concurrency model
---

Transactions commited to the ledger are fully atomic and serialized, supported by two separate and ordered concurrency-control mechanisms preventing effectively race-conditions to happen:

0. Pre-commit in-memory or Redis-based locking
1. Optimistic locking

## Transaction commit flow

![flow](/img/advanced/concurrency-model.png)

## Recommendation for multi-instances deployments

Due to the sequential nature of the ledger, multiple instances deployments should be carefully configured.

We recommended to use the Redis-based shared pre-commit lock using the [redis flags](/ledger/operations/env-vars). While the optimistic lock will ultimately be there to prevent race-conditions on commit, using the shared lock will reduce such commit attempts in the first place, yielding better performance for write heavy workloads.

Also, you can find a complete docker-compose example using 3 instances of the ledger, and a simple reverse proxy to handle traffic [there](https://github.com/formancehq/ledger/blob/main/examples/multi-node/docker-compose.yml).

Should the optimistic locking prevent a conflict on commit, it will surface it to the API consumer with a ([Response 409](/api#tag/transactions/operation/createTransaction)) - it is the responsibility of the client to retry the transaction in this case.
