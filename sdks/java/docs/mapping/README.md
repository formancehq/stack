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
                .setSecurity(new Security("laborum") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
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

## updateMapping

Update the mapping of a ledger

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.UpdateMappingRequest;
import com.formance.formance_sdk.models.operations.UpdateMappingResponse;
import com.formance.formance_sdk.models.shared.Contract;
import com.formance.formance_sdk.models.shared.Mapping;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("quasi") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            UpdateMappingRequest req = new UpdateMappingRequest(                new Mapping(                new com.formance.formance_sdk.models.shared.Contract[]{{
                                                add(new Contract(                new java.util.HashMap<String, Object>() {{
                                                                    put("doloremque", "reprehenderit");
                                                                }}) {{
                                                    account = "users:001";
                                                    expr = new java.util.HashMap<String, Object>() {{
                                                        put("vero", "nihil");
                                                        put("praesentium", "voluptatibus");
                                                        put("ipsa", "omnis");
                                                        put("voluptate", "cum");
                                                    }};
                                                }}),
                                                add(new Contract(                new java.util.HashMap<String, Object>() {{
                                                                    put("dicta", "harum");
                                                                    put("enim", "accusamus");
                                                                }}) {{
                                                    account = "users:001";
                                                    expr = new java.util.HashMap<String, Object>() {{
                                                        put("maiores", "dicta");
                                                        put("corporis", "dolore");
                                                    }};
                                                }}),
                                                add(new Contract(                new java.util.HashMap<String, Object>() {{
                                                                    put("excepturi", "pariatur");
                                                                    put("modi", "praesentium");
                                                                    put("rem", "voluptates");
                                                                }}) {{
                                                    account = "users:001";
                                                    expr = new java.util.HashMap<String, Object>() {{
                                                        put("repudiandae", "quae");
                                                        put("ipsum", "quidem");
                                                    }};
                                                }}),
                                                add(new Contract(                new java.util.HashMap<String, Object>() {{
                                                                    put("itaque", "incidunt");
                                                                }}) {{
                                                    account = "users:001";
                                                    expr = new java.util.HashMap<String, Object>() {{
                                                        put("repudiandae", "sint");
                                                    }};
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
