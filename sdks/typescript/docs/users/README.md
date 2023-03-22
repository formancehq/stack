# users

### Available Operations

* [listUsers](#listusers) - List users
* [readUser](#readuser) - Read user

## listUsers

List users

### Example Usage

```typescript
import { SDK } from "openapi";
import { ListUsersResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.users.listUsers().then((res: ListUsersResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## readUser

Read user

### Example Usage

```typescript
import { SDK } from "openapi";
import { ReadUserResponse } from "openapi/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.users.readUser({
  userId: "dolores",
}).then((res: ReadUserResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
