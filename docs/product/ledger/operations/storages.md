---
title: Storages considerations
---
# Storages considerations

Formance abstracts storage behind a single interface, and storage plugins are responsible for implementing them in the most efficient way possible.

## Postgres

A Postgres storage backend is shipped with Formance out of the box and is the recommended storage for production use. It is not the default storage and can be enabled by updating the configuration variables `storage.driver` and `storage-postgres-conn-string`. Environments variables or file can be used.

```shell
NUMARY_STORAGE_DRIVER=postgres \
NUMARY_STORAGE_POSTGRES_CONN_STRING=postgresql://localhost/dbname \
numary server start
```

## SQLite

SQLite is a simple storage built into with Formance Ledger. It is ideal for testing purposes local usage of the ledger. It will use the data directory defined in config, which defaults to `$HOME/.numary/data`, or use the one defined in config file:

```yaml
storage:
  driver: postgres
  postgres:
    conn_string: postgresql://localhost/postgres
```

:::caution
SQLite is not recommended for production use and is on its way to be deprecated in favor of Postgres in the Ledger v2.x line.
:::
