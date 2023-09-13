# balances

### Available Operations

* [getBalances](#getbalances) - Get the balances from a ledger's account
* [getBalancesAggregated](#getbalancesaggregated) - Get the aggregated balances from selected accounts

## getBalances

Get the balances from a ledger's account

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetBalancesRequest;
import com.formance.formance_sdk.models.operations.GetBalancesResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("doloribus") {{
                    authorization = "";
                }})
                .build();

            GetBalancesRequest req = new GetBalancesRequest("ledger001") {{
                address = "users:001";
                after = "users:003";
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                pageSize = 958950L;
            }};            

            GetBalancesResponse res = sdk.balances.getBalances(req);

            if (res.balancesCursorResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                       | Type                                                                                                            | Required                                                                                                        | Description                                                                                                     |
| --------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                       | [com.formance.formance_sdk.models.operations.GetBalancesRequest](../../models/operations/GetBalancesRequest.md) | :heavy_check_mark:                                                                                              | The request object to use for the request.                                                                      |


### Response

**[com.formance.formance_sdk.models.operations.GetBalancesResponse](../../models/operations/GetBalancesResponse.md)**


## getBalancesAggregated

Get the aggregated balances from selected accounts

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetBalancesAggregatedRequest;
import com.formance.formance_sdk.models.operations.GetBalancesAggregatedResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("architecto") {{
                    authorization = "";
                }})
                .build();

            GetBalancesAggregatedRequest req = new GetBalancesAggregatedRequest("ledger001") {{
                address = "users:001";
            }};            

            GetBalancesAggregatedResponse res = sdk.balances.getBalancesAggregated(req);

            if (res.aggregateBalancesResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                                           | Type                                                                                                                                | Required                                                                                                                            | Description                                                                                                                         |
| ----------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                                           | [com.formance.formance_sdk.models.operations.GetBalancesAggregatedRequest](../../models/operations/GetBalancesAggregatedRequest.md) | :heavy_check_mark:                                                                                                                  | The request object to use for the request.                                                                                          |


### Response

**[com.formance.formance_sdk.models.operations.GetBalancesAggregatedResponse](../../models/operations/GetBalancesAggregatedResponse.md)**

