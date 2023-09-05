# search

### Available Operations

* [search](#search) - Search
* [searchgetServerInfo](#searchgetserverinfo) - Get server info

## search

ElasticSearch query engine

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.SearchResponse;
import com.formance.formance_sdk.models.shared.Query;
import com.formance.formance_sdk.models.shared.QueryRaw;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("magnam") {{
                    authorization = "";
                }})
                .build();

            com.formance.formance_sdk.models.shared.Query req = new Query() {{
                after = new String[]{{
                    add("users:002"),
                    add("users:002"),
                    add("users:002"),
                    add("users:002"),
                }};
                cursor = "YXVsdCBhbmQgYSBtYXhpbXVtIG1heF9yZXN1bHRzLol=";
                ledgers = new String[]{{
                    add("quickstart"),
                    add("quickstart"),
                    add("quickstart"),
                    add("quickstart"),
                }};
                pageSize = 411820L;
                policy = "OR";
                raw = new QueryRaw();;
                sort = "txid:asc";
                target = "aliquid";
                terms = new String[]{{
                    add("destination=central_bank1"),
                    add("destination=central_bank1"),
                    add("destination=central_bank1"),
                }};
            }};            

            SearchResponse res = sdk.search.search(req);

            if (res.response != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                     | Type                                                                          | Required                                                                      | Description                                                                   |
| ----------------------------------------------------------------------------- | ----------------------------------------------------------------------------- | ----------------------------------------------------------------------------- | ----------------------------------------------------------------------------- |
| `request`                                                                     | [com.formance.formance_sdk.models.shared.Query](../../models/shared/Query.md) | :heavy_check_mark:                                                            | The request object to use for the request.                                    |


### Response

**[com.formance.formance_sdk.models.operations.SearchResponse](../../models/operations/SearchResponse.md)**


## searchgetServerInfo

Get server info

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.SearchgetServerInfoResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("accusamus") {{
                    authorization = "";
                }})
                .build();

            SearchgetServerInfoResponse res = sdk.search.searchgetServerInfo();

            if (res.serverInfo != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```


### Response

**[com.formance.formance_sdk.models.operations.SearchgetServerInfoResponse](../../models/operations/SearchgetServerInfoResponse.md)**

