# auth

### Available Operations

* [addScopeToClient](#addscopetoclient) - Add scope to client
* [addTransientScope](#addtransientscope) - Add a transient scope to a scope
* [createClient](#createclient) - Create client
* [createScope](#createscope) - Create scope
* [createSecret](#createsecret) - Add a secret to a client
* [deleteClient](#deleteclient) - Delete client
* [deleteScope](#deletescope) - Delete scope
* [deleteScopeFromClient](#deletescopefromclient) - Delete scope from client
* [deleteSecret](#deletesecret) - Delete a secret from a client
* [deleteTransientScope](#deletetransientscope) - Delete a transient scope from a scope
* [getServerInfo](#getserverinfo) - Get server info
* [listClients](#listclients) - List clients
* [listScopes](#listscopes) - List scopes
* [listUsers](#listusers) - List users
* [readClient](#readclient) - Read client
* [readScope](#readscope) - Read scope
* [readUser](#readuser) - Read user
* [updateClient](#updateclient) - Update client
* [updateScope](#updatescope) - Update scope

## addScopeToClient

Add scope to client

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { AddScopeToClientResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.addScopeToClient({
  clientId: "recusandae",
  scopeId: "temporibus",
}).then((res: AddScopeToClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## addTransientScope

Add a transient scope to a scope

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { AddTransientScopeResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.addTransientScope({
  scopeId: "ab",
  transientScopeId: "quis",
}).then((res: AddTransientScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createClient

Create client

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateClientResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.createClient({
  description: "veritatis",
  metadata: {
    "perferendis": "ipsam",
    "repellendus": "sapiente",
    "quo": "odit",
  },
  name: "Wilfred Wolff",
  postLogoutRedirectUris: [
    "esse",
    "totam",
    "porro",
    "dolorum",
  ],
  public: false,
  redirectUris: [
    "nam",
  ],
  trusted: false,
}).then((res: CreateClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createScope

Create scope

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateScopeResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.createScope({
  label: "officia",
  metadata: {
    "fugit": "deleniti",
    "hic": "optio",
    "totam": "beatae",
  },
}).then((res: CreateScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createSecret

Add a secret to a client

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateSecretResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.createSecret({
  createSecretRequest: {
    metadata: {
      "molestiae": "modi",
      "qui": "impedit",
    },
    name: "Cory Emmerich",
  },
  clientId: "perferendis",
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
import { SDK } from "@formance/formance-sdk";
import { DeleteClientResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.deleteClient({
  clientId: "ad",
}).then((res: DeleteClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## deleteScope

Delete scope

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DeleteScopeResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.deleteScope({
  scopeId: "natus",
}).then((res: DeleteScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## deleteScopeFromClient

Delete scope from client

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DeleteScopeFromClientResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.deleteScopeFromClient({
  clientId: "sed",
  scopeId: "iste",
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
import { SDK } from "@formance/formance-sdk";
import { DeleteSecretResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.deleteSecret({
  clientId: "dolor",
  secretId: "natus",
}).then((res: DeleteSecretResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## deleteTransientScope

Delete a transient scope from a scope

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DeleteTransientScopeResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.deleteTransientScope({
  scopeId: "laboriosam",
  transientScopeId: "hic",
}).then((res: DeleteTransientScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getServerInfo

Get server info

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetServerInfoResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.getServerInfo().then((res: GetServerInfoResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listClients

List clients

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListClientsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.listClients().then((res: ListClientsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listScopes

List Scopes

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListScopesResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.listScopes().then((res: ListScopesResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listUsers

List users

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListUsersResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.listUsers().then((res: ListUsersResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## readClient

Read client

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ReadClientResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.readClient({
  clientId: "saepe",
}).then((res: ReadClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## readScope

Read scope

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ReadScopeResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.readScope({
  scopeId: "fuga",
}).then((res: ReadScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## readUser

Read user

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ReadUserResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.readUser({
  userId: "in",
}).then((res: ReadUserResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## updateClient

Update client

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { UpdateClientResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.updateClient({
  updateClientRequest: {
    description: "corporis",
    metadata: {
      "iure": "saepe",
      "quidem": "architecto",
      "ipsa": "reiciendis",
    },
    name: "Shaun Osinski",
    postLogoutRedirectUris: [
      "explicabo",
      "nobis",
    ],
    public: false,
    redirectUris: [
      "omnis",
      "nemo",
    ],
    trusted: false,
  },
  clientId: "minima",
}).then((res: UpdateClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## updateScope

Update scope

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { UpdateScopeResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.auth.updateScope({
  updateScopeRequest: {
    label: "excepturi",
    metadata: {
      "iure": "culpa",
    },
  },
  scopeId: "doloribus",
}).then((res: UpdateScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
