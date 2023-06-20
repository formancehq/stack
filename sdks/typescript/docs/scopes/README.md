# scopes

### Available Operations

* [addTransientScope](#addtransientscope) - Add a transient scope to a scope
* [createScope](#createscope) - Create scope
* [deleteScope](#deletescope) - Delete scope
* [deleteTransientScope](#deletetransientscope) - Delete a transient scope from a scope
* [listScopes](#listscopes) - List scopes
* [readScope](#readscope) - Read scope
* [updateScope](#updatescope) - Update scope

## addTransientScope

Add a transient scope to a scope

### Example Usage

```typescript
import { SDK } from "openapi";
import { AddTransientScopeResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.scopes.addTransientScope({
  scopeId: "aspernatur",
  transientScopeId: "architecto",
}).then((res: AddTransientScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createScope

Create scope

### Example Usage

```typescript
import { SDK } from "openapi";
import { CreateScopeResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.scopes.createScope({
  label: "magnam",
  metadata: {
    "excepturi": "ullam",
  },
}).then((res: CreateScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## deleteScope

Delete scope

### Example Usage

```typescript
import { SDK } from "openapi";
import { DeleteScopeResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.scopes.deleteScope({
  scopeId: "provident",
}).then((res: DeleteScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## deleteTransientScope

Delete a transient scope from a scope

### Example Usage

```typescript
import { SDK } from "openapi";
import { DeleteTransientScopeResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.scopes.deleteTransientScope({
  scopeId: "quos",
  transientScopeId: "sint",
}).then((res: DeleteTransientScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listScopes

List Scopes

### Example Usage

```typescript
import { SDK } from "openapi";
import { ListScopesResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.scopes.listScopes().then((res: ListScopesResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## readScope

Read scope

### Example Usage

```typescript
import { SDK } from "openapi";
import { ReadScopeResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.scopes.readScope({
  scopeId: "accusantium",
}).then((res: ReadScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## updateScope

Update scope

### Example Usage

```typescript
import { SDK } from "openapi";
import { UpdateScopeResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.scopes.updateScope({
  updateScopeRequest: {
    label: "mollitia",
    metadata: {
      "mollitia": "ad",
      "eum": "dolor",
      "necessitatibus": "odit",
      "nemo": "quasi",
    },
  },
  scopeId: "iure",
}).then((res: UpdateScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
