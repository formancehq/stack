# mapping

### Available Operations

* [getMapping](#getmapping) - Get the mapping of a ledger
* [updateMapping](#updatemapping) - Update the mapping of a ledger

## getMapping

Get the mapping of a ledger

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetMappingRequest;
import com.formance.formance_sdk.models.operations.GetMappingResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("commodi") {{
                    authorization = "";
                }})
                .build();

            GetMappingRequest req = new GetMappingRequest("ledger001");            

            GetMappingResponse res = sdk.mapping.getMapping(req);

            if (res.mappingResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                     | Type                                                                                                          | Required                                                                                                      | Description                                                                                                   |
| ------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                     | [com.formance.formance_sdk.models.operations.GetMappingRequest](../../models/operations/GetMappingRequest.md) | :heavy_check_mark:                                                                                            | The request object to use for the request.                                                                    |


### Response

**[com.formance.formance_sdk.models.operations.GetMappingResponse](../../models/operations/GetMappingResponse.md)**


## updateMapping

Update the mapping of a ledger

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.UpdateMappingRequest;
import com.formance.formance_sdk.models.operations.UpdateMappingResponse;
import com.formance.formance_sdk.models.shared.Contract;
import com.formance.formance_sdk.models.shared.ContractExpr;
import com.formance.formance_sdk.models.shared.Mapping;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("quam") {{
                    authorization = "";
                }})
                .build();

            UpdateMappingRequest req = new UpdateMappingRequest(                new Mapping(                new com.formance.formance_sdk.models.shared.Contract[]{{
                                                add(new Contract(                new ContractExpr();) {{
                                                    account = "users:001";
                                                    expr = new ContractExpr() {{}};
                                                }}),
                                                add(new Contract(                new ContractExpr();) {{
                                                    account = "users:001";
                                                    expr = new ContractExpr() {{}};
                                                }}),
                                            }});, "ledger001");            

            UpdateMappingResponse res = sdk.mapping.updateMapping(req);

            if (res.mappingResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                           | Type                                                                                                                | Required                                                                                                            | Description                                                                                                         |
| ------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                           | [com.formance.formance_sdk.models.operations.UpdateMappingRequest](../../models/operations/UpdateMappingRequest.md) | :heavy_check_mark:                                                                                                  | The request object to use for the request.                                                                          |


### Response

**[com.formance.formance_sdk.models.operations.UpdateMappingResponse](../../models/operations/UpdateMappingResponse.md)**

