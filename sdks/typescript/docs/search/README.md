# search

### Available Operations

* [search](#search) - Search
* [searchgetServerInfo](#searchgetserverinfo) - Get server info

## search

ElasticSearch query engine

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { SearchResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.search.search({
  after: [
    "users:002",
    "users:002",
    "users:002",
  ],
  cursor: "YXVsdCBhbmQgYSBtYXhpbXVtIG1heF9yZXN1bHRzLol=",
  ledgers: [
    "quickstart",
    "quickstart",
    "quickstart",
  ],
  pageSize: 725255,
  policy: "OR",
  raw: {
    "blanditiis": "deleniti",
    "sapiente": "amet",
    "deserunt": "nisi",
  },
  sort: "txid:asc",
  target: "vel",
  terms: [
    "destination=central_bank1",
    "destination=central_bank1",
    "destination=central_bank1",
  ],
}).then((res: SearchResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## searchgetServerInfo

Get server info

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { SearchgetServerInfoResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.search.searchgetServerInfo().then((res: SearchgetServerInfoResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
