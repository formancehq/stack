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

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ActivateConfigRequest;
import com.formance.formance_sdk.models.operations.ActivateConfigResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("iure") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ActivateConfigRequest req = new ActivateConfigRequest("4997257d-dfb6-445b-929c-cbe2ab182818");            

            ActivateConfigResponse res = sdk.webhooks.activateConfig(req);

            if (res.configResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## changeConfigSecret

Change the signing secret of the endpoint of a webhooks config.

If not passed or empty, a secret is automatically generated.
The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)


### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ChangeConfigSecretRequest;
import com.formance.formance_sdk.models.operations.ChangeConfigSecretResponse;
import com.formance.formance_sdk.models.shared.ConfigChangeSecret;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("odio") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ChangeConfigSecretRequest req = new ChangeConfigSecretRequest("4997257d-dfb6-445b-929c-cbe2ab182818") {{
                configChangeSecret = new ConfigChangeSecret("V0bivxRWveaoz08afqjU6Ko/jwO0Cb+3");;
            }};            

            ChangeConfigSecretResponse res = sdk.webhooks.changeConfigSecret(req);

            if (res.configResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## deactivateConfig

Deactivate a webhooks config by ID, to stop receiving webhooks to its endpoint.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.DeactivateConfigRequest;
import com.formance.formance_sdk.models.operations.DeactivateConfigResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("quaerat") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            DeactivateConfigRequest req = new DeactivateConfigRequest("4997257d-dfb6-445b-929c-cbe2ab182818");            

            DeactivateConfigResponse res = sdk.webhooks.deactivateConfig(req);

            if (res.configResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## deleteConfig

Delete a webhooks config by ID.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.DeleteConfigRequest;
import com.formance.formance_sdk.models.operations.DeleteConfigResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("accusamus") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            DeleteConfigRequest req = new DeleteConfigRequest("4997257d-dfb6-445b-929c-cbe2ab182818");            

            DeleteConfigResponse res = sdk.webhooks.deleteConfig(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getManyConfigs

Sorted by updated date descending

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetManyConfigsRequest;
import com.formance.formance_sdk.models.operations.GetManyConfigsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("quidem") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetManyConfigsRequest req = new GetManyConfigsRequest() {{
                endpoint = "https://example.com";
                id = "4997257d-dfb6-445b-929c-cbe2ab182818";
            }};            

            GetManyConfigsResponse res = sdk.webhooks.getManyConfigs(req);

            if (res.configsResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## insertConfig

Insert a new webhooks config.

The endpoint should be a valid https URL and be unique.

The secret is the endpoint's verification secret.
If not passed or empty, a secret is automatically generated.
The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)

All eventTypes are converted to lower-case when inserted.


### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.InsertConfigResponse;
import com.formance.formance_sdk.models.shared.ConfigUser;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("voluptatibus") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            com.formance.formance_sdk.models.shared.ConfigUser req = new ConfigUser("https://example.com",                 new String[]{{
                                add("TYPE1"),
                                add("TYPE1"),
                            }}) {{
                secret = "V0bivxRWveaoz08afqjU6Ko/jwO0Cb+3";
            }};            

            InsertConfigResponse res = sdk.webhooks.insertConfig(req);

            if (res.configResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## testConfig

Test a config by sending a webhook to its endpoint.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.TestConfigRequest;
import com.formance.formance_sdk.models.operations.TestConfigResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("natus") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            TestConfigRequest req = new TestConfigRequest("4997257d-dfb6-445b-929c-cbe2ab182818");            

            TestConfigResponse res = sdk.webhooks.testConfig(req);

            if (res.attemptResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```
