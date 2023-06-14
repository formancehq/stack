# webhooks

### Available Operations

* [activateConfig](#activateconfig) - Activate one config
* [changeConfigSecret](#changeconfigsecret) - Change the signing secret of a config
* [deactivateConfig](#deactivateconfig) - Deactivate one config
* [deleteConfig](#deleteconfig) - Delete one config
* [getManyConfigs](#getmanyconfigs) - Get many configs
* [insertConfig](#insertconfig) - Insert a new config
* [testConfig](#testconfig) - Test one config

## activateConfig

Activate a webhooks config by ID, to start receiving webhooks to its endpoint.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ActivateConfigResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.webhooks.activateConfig({
  id: "4997257d-dfb6-445b-929c-cbe2ab182818",
}).then((res: ActivateConfigResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## changeConfigSecret

Change the signing secret of the endpoint of a webhooks config.

If not passed or empty, a secret is automatically generated.
The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)


### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ChangeConfigSecretResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.webhooks.changeConfigSecret({
  configChangeSecret: {
    secret: "V0bivxRWveaoz08afqjU6Ko/jwO0Cb+3",
  },
  id: "4997257d-dfb6-445b-929c-cbe2ab182818",
}).then((res: ChangeConfigSecretResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## deactivateConfig

Deactivate a webhooks config by ID, to stop receiving webhooks to its endpoint.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DeactivateConfigResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.webhooks.deactivateConfig({
  id: "4997257d-dfb6-445b-929c-cbe2ab182818",
}).then((res: DeactivateConfigResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## deleteConfig

Delete a webhooks config by ID.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DeleteConfigResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.webhooks.deleteConfig({
  id: "4997257d-dfb6-445b-929c-cbe2ab182818",
}).then((res: DeleteConfigResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getManyConfigs

Sorted by updated date descending

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetManyConfigsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.webhooks.getManyConfigs({
  endpoint: "https://example.com",
  id: "4997257d-dfb6-445b-929c-cbe2ab182818",
}).then((res: GetManyConfigsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## insertConfig

Insert a new webhooks config.

The endpoint should be a valid https URL and be unique.

The secret is the endpoint's verification secret.
If not passed or empty, a secret is automatically generated.
The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)

All eventTypes are converted to lower-case when inserted.


### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { InsertConfigResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.webhooks.insertConfig({
  endpoint: "https://example.com",
  eventTypes: [
    "TYPE1",
    "TYPE1",
  ],
  secret: "V0bivxRWveaoz08afqjU6Ko/jwO0Cb+3",
}).then((res: InsertConfigResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## testConfig

Test a config by sending a webhook to its endpoint.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { TestConfigResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.webhooks.testConfig({
  id: "4997257d-dfb6-445b-929c-cbe2ab182818",
}).then((res: TestConfigResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
