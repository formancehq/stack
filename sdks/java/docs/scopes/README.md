# scopes

### Available Operations

* [addTransientScope](#addtransientscope) - Add a transient scope to a scope
* [createScope](#createscope) - Create scope
* [deleteScope](#deletescope) - Delete scope
* [deleteTransientScope](#deletetransientscope) - Delete a transient scope from a scope
* [listScopes](#listscopes) - List scopes
* [readScope](#readscope) - Read scope
* [updateScope](#updatescope) - Update scope

## addTransientScope

Add a transient scope to a scope

### Example Usage

```java
package hello.world;

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.AddTransientScopeRequest;
import org.openapis.openapi.models.operations.AddTransientScopeResponse;
import org.openapis.openapi.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("ullam") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            AddTransientScopeRequest req = new AddTransientScopeRequest("provident", "quos");            

            AddTransientScopeResponse res = sdk.scopes.addTransientScope(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## createScope

Create scope

### Example Usage

```java
package hello.world;

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.CreateScopeResponse;
import org.openapis.openapi.models.shared.CreateScopeRequest;
import org.openapis.openapi.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("sint") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            org.openapis.openapi.models.shared.CreateScopeRequest req = new CreateScopeRequest("accusantium") {{
                metadata = new java.util.HashMap<String, Object>() {{
                    put("reiciendis", "mollitia");
                    put("ad", "eum");
                    put("dolor", "necessitatibus");
                }};
            }};            

            CreateScopeResponse res = sdk.scopes.createScope(req);

            if (res.createScopeResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## deleteScope

Delete scope

### Example Usage

```java
package hello.world;

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.DeleteScopeRequest;
import org.openapis.openapi.models.operations.DeleteScopeResponse;
import org.openapis.openapi.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("odit") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            DeleteScopeRequest req = new DeleteScopeRequest("nemo");            

            DeleteScopeResponse res = sdk.scopes.deleteScope(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## deleteTransientScope

Delete a transient scope from a scope

### Example Usage

```java
package hello.world;

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.DeleteTransientScopeRequest;
import org.openapis.openapi.models.operations.DeleteTransientScopeResponse;
import org.openapis.openapi.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("quasi") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            DeleteTransientScopeRequest req = new DeleteTransientScopeRequest("iure", "doloribus");            

            DeleteTransientScopeResponse res = sdk.scopes.deleteTransientScope(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listScopes

List Scopes

### Example Usage

```java
package hello.world;

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.ListScopesResponse;
import org.openapis.openapi.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("debitis") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListScopesResponse res = sdk.scopes.listScopes();

            if (res.listScopesResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## readScope

Read scope

### Example Usage

```java
package hello.world;

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.ReadScopeRequest;
import org.openapis.openapi.models.operations.ReadScopeResponse;
import org.openapis.openapi.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("eius") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ReadScopeRequest req = new ReadScopeRequest("maxime");            

            ReadScopeResponse res = sdk.scopes.readScope(req);

            if (res.readScopeResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## updateScope

Update scope

### Example Usage

```java
package hello.world;

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.UpdateScopeRequest;
import org.openapis.openapi.models.operations.UpdateScopeResponse;
import org.openapis.openapi.models.shared.Security;
import org.openapis.openapi.models.shared.UpdateScopeRequest;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("deleniti") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            UpdateScopeRequest req = new UpdateScopeRequest("facilis") {{
                updateScopeRequest = new UpdateScopeRequest("in") {{
                    metadata = new java.util.HashMap<String, Object>() {{
                        put("architecto", "repudiandae");
                    }};
                }};;
            }};            

            UpdateScopeResponse res = sdk.scopes.updateScope(req);

            if (res.updateScopeResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```
