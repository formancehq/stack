<!-- Start SDK Example Usage -->
```typescript
import { SDK } from "@formance/formance-sdk";
import { GetVersionsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.getVersions().then((res: GetVersionsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
<!-- End SDK Example Usage -->