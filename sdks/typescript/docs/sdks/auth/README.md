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
    authorization: "",
  },
});

sdk.auth.addScopeToClient({
  clientId: "vel",
  scopeId: "error",
}).then((res: AddScopeToClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `request`                                                                                | [operations.AddScopeToClientRequest](../../models/operations/addscopetoclientrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |
| `config`                                                                                 | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                             | :heavy_minus_sign:                                                                       | Available config options for making requests.                                            |


### Response

**Promise<[operations.AddScopeToClientResponse](../../models/operations/addscopetoclientresponse.md)>**


## addTransientScope

Add a transient scope to a scope

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { AddTransientScopeResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.addTransientScope({
  scopeId: "deserunt",
  transientScopeId: "suscipit",
}).then((res: AddTransientScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                                  | Type                                                                                       | Required                                                                                   | Description                                                                                |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ |
| `request`                                                                                  | [operations.AddTransientScopeRequest](../../models/operations/addtransientscoperequest.md) | :heavy_check_mark:                                                                         | The request object to use for the request.                                                 |
| `config`                                                                                   | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                               | :heavy_minus_sign:                                                                         | Available config options for making requests.                                              |


### Response

**Promise<[operations.AddTransientScopeResponse](../../models/operations/addtransientscoperesponse.md)>**


## createClient

Create client

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateClientResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.createClient({
  description: "iure",
  metadata: {
    "debitis": "ipsa",
    "delectus": "tempora",
  },
  name: "Minnie Schiller",
  postLogoutRedirectUris: [
    "excepturi",
    "nisi",
  ],
  public: false,
  redirectUris: [
    "temporibus",
    "ab",
    "quis",
    "veritatis",
  ],
  trusted: false,
}).then((res: CreateClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `request`                                                                | [shared.CreateClientRequest](../../models/shared/createclientrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |
| `config`                                                                 | [AxiosRequestConfig](https://axios-http.com/docs/req_config)             | :heavy_minus_sign:                                                       | Available config options for making requests.                            |


### Response

**Promise<[operations.CreateClientResponse](../../models/operations/createclientresponse.md)>**


## createScope

Create scope

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateScopeResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.createScope({
  label: "deserunt",
  metadata: {
    "ipsam": "repellendus",
  },
}).then((res: CreateScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                              | Type                                                                   | Required                                                               | Description                                                            |
| ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| `request`                                                              | [shared.CreateScopeRequest](../../models/shared/createscoperequest.md) | :heavy_check_mark:                                                     | The request object to use for the request.                             |
| `config`                                                               | [AxiosRequestConfig](https://axios-http.com/docs/req_config)           | :heavy_minus_sign:                                                     | Available config options for making requests.                          |


### Response

**Promise<[operations.CreateScopeResponse](../../models/operations/createscoperesponse.md)>**


## createSecret

Add a secret to a client

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateSecretResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.createSecret({
  createSecretRequest: {
    metadata: {
      "quo": "odit",
      "at": "at",
      "maiores": "molestiae",
      "quod": "quod",
    },
    name: "Deanna Sauer MD",
  },
  clientId: "officia",
}).then((res: CreateSecretResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `request`                                                                        | [operations.CreateSecretRequest](../../models/operations/createsecretrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |
| `config`                                                                         | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                     | :heavy_minus_sign:                                                               | Available config options for making requests.                                    |


### Response

**Promise<[operations.CreateSecretResponse](../../models/operations/createsecretresponse.md)>**


## deleteClient

Delete client

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DeleteClientResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.deleteClient({
  clientId: "occaecati",
}).then((res: DeleteClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `request`                                                                        | [operations.DeleteClientRequest](../../models/operations/deleteclientrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |
| `config`                                                                         | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                     | :heavy_minus_sign:                                                               | Available config options for making requests.                                    |


### Response

**Promise<[operations.DeleteClientResponse](../../models/operations/deleteclientresponse.md)>**


## deleteScope

Delete scope

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DeleteScopeResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.deleteScope({
  scopeId: "fugit",
}).then((res: DeleteScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `request`                                                                      | [operations.DeleteScopeRequest](../../models/operations/deletescoperequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |
| `config`                                                                       | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                   | :heavy_minus_sign:                                                             | Available config options for making requests.                                  |


### Response

**Promise<[operations.DeleteScopeResponse](../../models/operations/deletescoperesponse.md)>**


## deleteScopeFromClient

Delete scope from client

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DeleteScopeFromClientResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.deleteScopeFromClient({
  clientId: "deleniti",
  scopeId: "hic",
}).then((res: DeleteScopeFromClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                                          | Type                                                                                               | Required                                                                                           | Description                                                                                        |
| -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `request`                                                                                          | [operations.DeleteScopeFromClientRequest](../../models/operations/deletescopefromclientrequest.md) | :heavy_check_mark:                                                                                 | The request object to use for the request.                                                         |
| `config`                                                                                           | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                                       | :heavy_minus_sign:                                                                                 | Available config options for making requests.                                                      |


### Response

**Promise<[operations.DeleteScopeFromClientResponse](../../models/operations/deletescopefromclientresponse.md)>**


## deleteSecret

Delete a secret from a client

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DeleteSecretResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.deleteSecret({
  clientId: "optio",
  secretId: "totam",
}).then((res: DeleteSecretResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `request`                                                                        | [operations.DeleteSecretRequest](../../models/operations/deletesecretrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |
| `config`                                                                         | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                     | :heavy_minus_sign:                                                               | Available config options for making requests.                                    |


### Response

**Promise<[operations.DeleteSecretResponse](../../models/operations/deletesecretresponse.md)>**


## deleteTransientScope

Delete a transient scope from a scope

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DeleteTransientScopeResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.deleteTransientScope({
  scopeId: "beatae",
  transientScopeId: "commodi",
}).then((res: DeleteTransientScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                                        | Type                                                                                             | Required                                                                                         | Description                                                                                      |
| ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| `request`                                                                                        | [operations.DeleteTransientScopeRequest](../../models/operations/deletetransientscoperequest.md) | :heavy_check_mark:                                                                               | The request object to use for the request.                                                       |
| `config`                                                                                         | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                                     | :heavy_minus_sign:                                                                               | Available config options for making requests.                                                    |


### Response

**Promise<[operations.DeleteTransientScopeResponse](../../models/operations/deletetransientscoperesponse.md)>**


## getServerInfo

Get server info

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetServerInfoResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.getServerInfo().then((res: GetServerInfoResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `config`                                                     | [AxiosRequestConfig](https://axios-http.com/docs/req_config) | :heavy_minus_sign:                                           | Available config options for making requests.                |


### Response

**Promise<[operations.GetServerInfoResponse](../../models/operations/getserverinforesponse.md)>**


## listClients

List clients

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListClientsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.listClients().then((res: ListClientsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `config`                                                     | [AxiosRequestConfig](https://axios-http.com/docs/req_config) | :heavy_minus_sign:                                           | Available config options for making requests.                |


### Response

**Promise<[operations.ListClientsResponse](../../models/operations/listclientsresponse.md)>**


## listScopes

List Scopes

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListScopesResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.listScopes().then((res: ListScopesResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `config`                                                     | [AxiosRequestConfig](https://axios-http.com/docs/req_config) | :heavy_minus_sign:                                           | Available config options for making requests.                |


### Response

**Promise<[operations.ListScopesResponse](../../models/operations/listscopesresponse.md)>**


## listUsers

List users

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListUsersResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.listUsers().then((res: ListUsersResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `config`                                                     | [AxiosRequestConfig](https://axios-http.com/docs/req_config) | :heavy_minus_sign:                                           | Available config options for making requests.                |


### Response

**Promise<[operations.ListUsersResponse](../../models/operations/listusersresponse.md)>**


## readClient

Read client

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ReadClientResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.readClient({
  clientId: "molestiae",
}).then((res: ReadClientResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `request`                                                                    | [operations.ReadClientRequest](../../models/operations/readclientrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |
| `config`                                                                     | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                 | :heavy_minus_sign:                                                           | Available config options for making requests.                                |


### Response

**Promise<[operations.ReadClientResponse](../../models/operations/readclientresponse.md)>**


## readScope

Read scope

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ReadScopeResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.readScope({
  scopeId: "modi",
}).then((res: ReadScopeResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                  | Type                                                                       | Required                                                                   | Description                                                                |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `request`                                                                  | [operations.ReadScopeRequest](../../models/operations/readscoperequest.md) | :heavy_check_mark:                                                         | The request object to use for the request.                                 |
| `config`                                                                   | [AxiosRequestConfig](https://axios-http.com/docs/req_config)               | :heavy_minus_sign:                                                         | Available config options for making requests.                              |


### Response

**Promise<[operations.ReadScopeResponse](../../models/operations/readscoperesponse.md)>**


## readUser

Read user

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ReadUserResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.readUser({
  userId: "qui",
}).then((res: ReadUserResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

### Parameters

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `request`                                                                | [operations.ReadUserRequest](../../models/operations/readuserrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |
| `config`                                                                 | [AxiosRequestConfig](https://axios-http.com/docs/req_config)             | :heavy_minus_sign:                                                       | Available config options for making requests.                            |


### Response

**Promise<[operations.ReadUserResponse](../../models/operations/readuserresponse.md)>**


## updateClient

Update client

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { UpdateClientResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
  },
});

sdk.auth.updateClient({
  updateClientRequest: {
    description: "impedit",
    metadata: {
      "esse": "ipsum",
      "excepturi": "aspernatur",
      "perferendis": "ad",
    },
    name: "Louis Moore",
    postLogoutRedirectUris: [
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

### Parameters

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `request`                                                                        | [operations.UpdateClientRequest](../../models/operations/updateclientrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |
| `config`                                                                         | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                     | :heavy_minus_sign:                                                               | Available config options for making requests.                                    |


### Response

**Promise<[operations.UpdateClientResponse](../../models/operations/updateclientresponse.md)>**


## updateScope

Update scope

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { UpdateScopeResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "",
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

### Parameters

| Parameter                                                                      | Type                                                                           | Required                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `request`                                                                      | [operations.UpdateScopeRequest](../../models/operations/updatescoperequest.md) | :heavy_check_mark:                                                             | The request object to use for the request.                                     |
| `config`                                                                       | [AxiosRequestConfig](https://axios-http.com/docs/req_config)                   | :heavy_minus_sign:                                                             | Available config options for making requests.                                  |


### Response

**Promise<[operations.UpdateScopeResponse](../../models/operations/updatescoperesponse.md)>**

