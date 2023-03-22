# clients

### Available Operations

* [addScopeToClient](#addscopetoclient) - Add scope to client
* [createClient](#createclient) - Create client
* [createSecret](#createsecret) - Add a secret to a client
* [deleteClient](#deleteclient) - Delete client
* [deleteScopeFromClient](#deletescopefromclient) - Delete scope from client
* [deleteSecret](#deletesecret) - Delete a secret from a client
* [listClients](#listclients) - List clients
* [readClient](#readclient) - Read client
* [updateClient](#updateclient) - Update client

## addScopeToClient

Add scope to client

### Example Usage

```typescript
import { SDK } from "openapi";
import { AddScopeToClientResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.clients.addScopeToClient({
  clientId: "velit",
  scopeId: "error",
}).then((res: AddScopeToClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createClient

Create client

### Example Usage

```typescript
import { SDK } from "openapi";
import { CreateClientResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.clients.createClient({
  description: "quia",
  metadata: {
    "vitae": "laborum",
    "animi": "enim",
  },
  name: "Angelica Dietrich",
  postLogoutRedirectUris: [
    "possimus",
    "aut",
    "quasi",
  ],
  public: false,
  redirectUris: [
    "temporibus",
    "laborum",
    "quasi",
  ],
  trusted: false,
}).then((res: CreateClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createSecret

Add a secret to a client

### Example Usage

```typescript
import { SDK } from "openapi";
import { CreateSecretResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.clients.createSecret({
  createSecretRequest: {
    metadata: {
      "voluptatibus": "vero",
      "nihil": "praesentium",
      "voluptatibus": "ipsa",
      "omnis": "voluptate",
    },
    name: "Thomas Batz",
  },
  clientId: "maiores",
}).then((res: CreateSecretResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## deleteClient

Delete client

### Example Usage

```typescript
import { SDK } from "openapi";
import { DeleteClientResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.clients.deleteClient({
  clientId: "dicta",
}).then((res: DeleteClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## deleteScopeFromClient

Delete scope from client

### Example Usage

```typescript
import { SDK } from "openapi";
import { DeleteScopeFromClientResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.clients.deleteScopeFromClient({
  clientId: "corporis",
  scopeId: "dolore",
}).then((res: DeleteScopeFromClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## deleteSecret

Delete a secret from a client

### Example Usage

```typescript
import { SDK } from "openapi";
import { DeleteSecretResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.clients.deleteSecret({
  clientId: "iusto",
  secretId: "dicta",
}).then((res: DeleteSecretResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listClients

List clients

### Example Usage

```typescript
import { SDK } from "openapi";
import { ListClientsResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.clients.listClients().then((res: ListClientsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## readClient

Read client

### Example Usage

```typescript
import { SDK } from "openapi";
import { ReadClientResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.clients.readClient({
  clientId: "harum",
}).then((res: ReadClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## updateClient

Update client

### Example Usage

```typescript
import { SDK } from "openapi";
import { UpdateClientResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.clients.updateClient({
  updateClientRequest: {
    description: "enim",
    metadata: {
      "commodi": "repudiandae",
      "quae": "ipsum",
      "quidem": "molestias",
      "excepturi": "pariatur",
    },
    name: "Irma Ledner DVM",
    postLogoutRedirectUris: [
      "veritatis",
      "itaque",
      "incidunt",
    ],
    public: false,
    redirectUris: [
      "consequatur",
      "est",
    ],
    trusted: false,
  },
  clientId: "quibusdam",
}).then((res: UpdateClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
