:::warning
This Module is subject to a user license.
:::

## Overview

Formance Search is a service that provides a search engine for your Formance stack. It is based on Elasticsearch and provides a powerful search engine for your data.

Search comes with with two components:
- **Search**: The public search service itself used to search your data.
- **[Benthos](https://www.benthos.dev/)**: The data pipeline that indexes your data into the search engine.

Benthos is a stream processor that is used to ingest data from your Formance stack into the Eleasticsearch search engine. Benthos listens to the [messages broker](../05-Infrastructure%20services/02-Message%20broker.md), transforms the data and indexes it into the search engine.

When you install the Search module with the Formance Operator, it will automatically install and configure Benthos on your behalf.

## Requirements

Formance Search requires:
- **Elasticsearch / OpenSearch**: See configuration guide [here](../05-Infrastructure%20services/03-Elasticsearch.md).

## Search Object

:::info
You can find all the available parameters in [the comprehensive CRD documentation](../09-Configuration%20reference/02-Custom%20Resource%20Definitions.md).
:::

```yaml
apiVersion: formance.com/v1beta1
kind: Search
metadata:
  name: formance-dev
spec:
  stack: formance-dev
```
