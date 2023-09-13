# logs

### Available Operations

* [listLogs](#listlogs) - List the logs from a ledger

## listLogs

List the logs from a ledger, sorted by ID in descending order.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListLogsRequest;
import com.formance.formance_sdk.models.operations.ListLogsResponse;
import com.formance.formance_sdk.models.shared.Security;
import java.time.OffsetDateTime;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("dolorem") {{
                    authorization = "";
                }})
                .build();

            ListLogsRequest req = new ListLogsRequest("ledger001") {{
                after = "culpa";
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                endTime = OffsetDateTime.parse("2022-01-02T17:10:32.894Z");
                pageSize = 653108L;
                startTime = OffsetDateTime.parse("2022-06-30T02:19:51.375Z");
            }};            

            ListLogsResponse res = sdk.logs.listLogs(req);

            if (res.logsCursorResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                 | Type                                                                                                      | Required                                                                                                  | Description                                                                                               |
| --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                 | [com.formance.formance_sdk.models.operations.ListLogsRequest](../../models/operations/ListLogsRequest.md) | :heavy_check_mark:                                                                                        | The request object to use for the request.                                                                |


### Response

**[com.formance.formance_sdk.models.operations.ListLogsResponse](../../models/operations/ListLogsResponse.md)**

