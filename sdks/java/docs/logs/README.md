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
                .setSecurity(new Security("enim") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListLogsRequest req = new ListLogsRequest("ledger001") {{
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                endTime = OffsetDateTime.parse("2022-03-22T21:41:36.666Z");
                pageSize = 196582L;
                startTime = OffsetDateTime.parse("2021-11-23T05:54:08.890Z");
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
