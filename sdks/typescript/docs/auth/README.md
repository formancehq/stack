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
  clientId: "corrupti",
  scopeId: "provident",
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
  scopeId: "distinctio",
  transientScopeId: "quibusdam",
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
  description: "unde",
  metadata: {
    "corrupti": "illum",
    "vel": "error",
    "deserunt": "suscipit",
    "iure": "magnam",
  },
  name: "Larry Windler",
  postLogoutRedirectUris: [
    "minus",
    "placeat",
  ],
  public: false,
  redirectUris: [
    "iusto",
    "excepturi",
    "nisi",
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
  label: "recusandae",
  metadata: {
    "ab": "quis",
    "veritatis": "deserunt",
    "perferendis": "ipsam",
    "repellendus": "sapiente",
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
      "odit": "at",
      "at": "maiores",
      "molestiae": "quod",
      "quod": "esse",
    },
    name: "Miss Lowell Parisian",
  },
  clientId: "occaecati",
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
  clientId: "fugit",
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
  scopeId: "deleniti",
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
  clientId: "hic",
  scopeId: "optio",
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
  clientId: "totam",
  secretId: "beatae",
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
  scopeId: "commodi",
  transientScopeId: "molestiae",
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
  clientId: "modi",
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
  scopeId: "qui",
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
  userId: "impedit",
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
    description: "cum",
    metadata: {
      "ipsum": "excepturi",
      "aspernatur": "perferendis",
    },
    name: "Faye Cormier",
    postLogoutRedirectUris: [
      "laboriosam",
      "hic",
      "saepe",
    ],
    public: false,
    redirectUris: [
      "in",
      "corporis",
      "iste",
    ],
    trusted: false,
  },
  clientId: "iure",
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
    label: "saepe",
    metadata: {
      "architecto": "ipsa",
      "reiciendis": "est",
      "mollitia": "laborum",
    },
  },
  scopeId: "dolores",
}).then((res: UpdateScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
