# script

### Available Operations

* [~~runScript~~](#runscript) - Execute a Numscript :warning: **Deprecated**

## ~~runScript~~

This route is deprecated, and has been merged into `POST /{ledger}/transactions`.


> :warning: **DEPRECATED**: this method will be removed in a future release, please migrate away from it as soon as possible.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.RunScriptRequest;
import com.formance.formance_sdk.models.operations.RunScriptResponse;
import com.formance.formance_sdk.models.shared.Script;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("dolores") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            RunScriptRequest req = new RunScriptRequest(                new Script("vars {
                            account $user
                            }
                            send [COIN 10] (
                            	source = @world
                            	destination = $user
                            )
                            ") {{
                                metadata = new java.util.HashMap<String, Object>() {{
                                    put("totam", "dignissimos");
                                    put("eaque", "quis");
                                }};
                                reference = "order_1234";
                                vars = new java.util.HashMap<String, Object>() {{
                                    put("eos", "perferendis");
                                }};
                            }};, "ledger001") {{
                preview = true;
            }};            

            RunScriptResponse res = sdk.script.runScript(req);

            if (res.scriptResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```
