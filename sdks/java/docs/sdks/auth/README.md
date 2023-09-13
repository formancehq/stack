# auth

### Available Operations

* [addScopeToClient](#addscopetoclient) - Add scope to client
* [addTransientScope](#addtransientscope) - Add a transient scope to a scope
* [createClient](#createclient) - Create client
* [createScope](#createscope) - Create scope
* [createSecret](#createsecret) - Add a secret to a client
* [deleteClient](#deleteclient) - Delete client
* [deleteScope](#deletescope) - Delete scope
* [deleteScopeFromClient](#deletescopefromclient) - Delete scope from client
* [deleteSecret](#deletesecret) - Delete a secret from a client
* [deleteTransientScope](#deletetransientscope) - Delete a transient scope from a scope
* [getServerInfo](#getserverinfo) - Get server info
* [listClients](#listclients) - List clients
* [listScopes](#listscopes) - List scopes
* [listUsers](#listusers) - List users
* [readClient](#readclient) - Read client
* [readScope](#readscope) - Read scope
* [readUser](#readuser) - Read user
* [updateClient](#updateclient) - Update client
* [updateScope](#updatescope) - Update scope

## addScopeToClient

Add scope to client

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.AddScopeToClientRequest;
import com.formance.formance_sdk.models.operations.AddScopeToClientResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("delectus") {{
                    authorization = "";
                }})
                .build();

            AddScopeToClientRequest req = new AddScopeToClientRequest("tempora", "suscipit");            

            AddScopeToClientResponse res = sdk.auth.addScopeToClient(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                                 | Type                                                                                                                      | Required                                                                                                                  | Description                                                                                                               |
| ------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                                 | [com.formance.formance_sdk.models.operations.AddScopeToClientRequest](../../models/operations/AddScopeToClientRequest.md) | :heavy_check_mark:                                                                                                        | The request object to use for the request.                                                                                |


### Response

**[com.formance.formance_sdk.models.operations.AddScopeToClientResponse](../../models/operations/AddScopeToClientResponse.md)**


## addTransientScope

Add a transient scope to a scope

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.AddTransientScopeRequest;
import com.formance.formance_sdk.models.operations.AddTransientScopeResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("molestiae") {{
                    authorization = "";
                }})
                .build();

            AddTransientScopeRequest req = new AddTransientScopeRequest("minus", "placeat");            

            AddTransientScopeResponse res = sdk.auth.addTransientScope(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                                   | Type                                                                                                                        | Required                                                                                                                    | Description                                                                                                                 |
| --------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                                   | [com.formance.formance_sdk.models.operations.AddTransientScopeRequest](../../models/operations/AddTransientScopeRequest.md) | :heavy_check_mark:                                                                                                          | The request object to use for the request.                                                                                  |


### Response

**[com.formance.formance_sdk.models.operations.AddTransientScopeResponse](../../models/operations/AddTransientScopeResponse.md)**


## createClient

Create client

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CreateClientResponse;
import com.formance.formance_sdk.models.shared.CreateClientRequest;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("voluptatum") {{
                    authorization = "";
                }})
                .build();

            com.formance.formance_sdk.models.shared.CreateClientRequest req = new CreateClientRequest("iusto") {{
                description = "excepturi";
                metadata = new java.util.HashMap<String, Object>() {{
                    put("recusandae", "temporibus");
                    put("ab", "quis");
                }};
                postLogoutRedirectUris = new String[]{{
                    add("deserunt"),
                }};
                public_ = false;
                redirectUris = new String[]{{
                    add("ipsam"),
                }};
                trusted = false;
            }};            

            CreateClientResponse res = sdk.auth.createClient(req);

            if (res.createClientResponse != null) {
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
| `request`                                                                                                 | [com.formance.formance_sdk.models.shared.CreateClientRequest](../../models/shared/CreateClientRequest.md) | :heavy_check_mark:                                                                                        | The request object to use for the request.                                                                |


### Response

**[com.formance.formance_sdk.models.operations.CreateClientResponse](../../models/operations/CreateClientResponse.md)**


## createScope

Create scope

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CreateScopeResponse;
import com.formance.formance_sdk.models.shared.CreateScopeRequest;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("repellendus") {{
                    authorization = "";
                }})
                .build();

            com.formance.formance_sdk.models.shared.CreateScopeRequest req = new CreateScopeRequest("sapiente") {{
                metadata = new java.util.HashMap<String, Object>() {{
                    put("odit", "at");
                    put("at", "maiores");
                    put("molestiae", "quod");
                    put("quod", "esse");
                }};
            }};            

            CreateScopeResponse res = sdk.auth.createScope(req);

            if (res.createScopeResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                               | Type                                                                                                    | Required                                                                                                | Description                                                                                             |
| ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| `request`                                                                                               | [com.formance.formance_sdk.models.shared.CreateScopeRequest](../../models/shared/CreateScopeRequest.md) | :heavy_check_mark:                                                                                      | The request object to use for the request.                                                              |


### Response

**[com.formance.formance_sdk.models.operations.CreateScopeResponse](../../models/operations/CreateScopeResponse.md)**


## createSecret

Add a secret to a client

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CreateSecretRequest;
import com.formance.formance_sdk.models.operations.CreateSecretResponse;
import com.formance.formance_sdk.models.shared.CreateSecretRequest;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("totam") {{
                    authorization = "";
                }})
                .build();

            CreateSecretRequest req = new CreateSecretRequest("porro") {{
                createSecretRequest = new CreateSecretRequest("dolorum") {{
                    metadata = new java.util.HashMap<String, Object>() {{
                        put("nam", "officia");
                    }};
                }};;
            }};            

            CreateSecretResponse res = sdk.auth.createSecret(req);

            if (res.createSecretResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                         | Type                                                                                                              | Required                                                                                                          | Description                                                                                                       |
| ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                         | [com.formance.formance_sdk.models.operations.CreateSecretRequest](../../models/operations/CreateSecretRequest.md) | :heavy_check_mark:                                                                                                | The request object to use for the request.                                                                        |


### Response

**[com.formance.formance_sdk.models.operations.CreateSecretResponse](../../models/operations/CreateSecretResponse.md)**


## deleteClient

Delete client

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.DeleteClientRequest;
import com.formance.formance_sdk.models.operations.DeleteClientResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("occaecati") {{
                    authorization = "";
                }})
                .build();

            DeleteClientRequest req = new DeleteClientRequest("fugit");            

            DeleteClientResponse res = sdk.auth.deleteClient(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                         | Type                                                                                                              | Required                                                                                                          | Description                                                                                                       |
| ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                         | [com.formance.formance_sdk.models.operations.DeleteClientRequest](../../models/operations/DeleteClientRequest.md) | :heavy_check_mark:                                                                                                | The request object to use for the request.                                                                        |


### Response

**[com.formance.formance_sdk.models.operations.DeleteClientResponse](../../models/operations/DeleteClientResponse.md)**


## deleteScope

Delete scope

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.DeleteScopeRequest;
import com.formance.formance_sdk.models.operations.DeleteScopeResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("deleniti") {{
                    authorization = "";
                }})
                .build();

            DeleteScopeRequest req = new DeleteScopeRequest("hic");            

            DeleteScopeResponse res = sdk.auth.deleteScope(req);

            if (res.statusCode == 200) {
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
| `request`                                                                                                       | [com.formance.formance_sdk.models.operations.DeleteScopeRequest](../../models/operations/DeleteScopeRequest.md) | :heavy_check_mark:                                                                                              | The request object to use for the request.                                                                      |


### Response

**[com.formance.formance_sdk.models.operations.DeleteScopeResponse](../../models/operations/DeleteScopeResponse.md)**


## deleteScopeFromClient

Delete scope from client

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.DeleteScopeFromClientRequest;
import com.formance.formance_sdk.models.operations.DeleteScopeFromClientResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("optio") {{
                    authorization = "";
                }})
                .build();

            DeleteScopeFromClientRequest req = new DeleteScopeFromClientRequest("totam", "beatae");            

            DeleteScopeFromClientResponse res = sdk.auth.deleteScopeFromClient(req);

            if (res.statusCode == 200) {
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
| `request`                                                                                                                           | [com.formance.formance_sdk.models.operations.DeleteScopeFromClientRequest](../../models/operations/DeleteScopeFromClientRequest.md) | :heavy_check_mark:                                                                                                                  | The request object to use for the request.                                                                                          |


### Response

**[com.formance.formance_sdk.models.operations.DeleteScopeFromClientResponse](../../models/operations/DeleteScopeFromClientResponse.md)**


## deleteSecret

Delete a secret from a client

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.DeleteSecretRequest;
import com.formance.formance_sdk.models.operations.DeleteSecretResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("commodi") {{
                    authorization = "";
                }})
                .build();

            DeleteSecretRequest req = new DeleteSecretRequest("molestiae", "modi");            

            DeleteSecretResponse res = sdk.auth.deleteSecret(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                         | Type                                                                                                              | Required                                                                                                          | Description                                                                                                       |
| ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                         | [com.formance.formance_sdk.models.operations.DeleteSecretRequest](../../models/operations/DeleteSecretRequest.md) | :heavy_check_mark:                                                                                                | The request object to use for the request.                                                                        |


### Response

**[com.formance.formance_sdk.models.operations.DeleteSecretResponse](../../models/operations/DeleteSecretResponse.md)**


## deleteTransientScope

Delete a transient scope from a scope

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.DeleteTransientScopeRequest;
import com.formance.formance_sdk.models.operations.DeleteTransientScopeResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("qui") {{
                    authorization = "";
                }})
                .build();

            DeleteTransientScopeRequest req = new DeleteTransientScopeRequest("impedit", "cum");            

            DeleteTransientScopeResponse res = sdk.auth.deleteTransientScope(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                                         | Type                                                                                                                              | Required                                                                                                                          | Description                                                                                                                       |
| --------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                                         | [com.formance.formance_sdk.models.operations.DeleteTransientScopeRequest](../../models/operations/DeleteTransientScopeRequest.md) | :heavy_check_mark:                                                                                                                | The request object to use for the request.                                                                                        |


### Response

**[com.formance.formance_sdk.models.operations.DeleteTransientScopeResponse](../../models/operations/DeleteTransientScopeResponse.md)**


## getServerInfo

Get server info

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetServerInfoResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("esse") {{
                    authorization = "";
                }})
                .build();

            GetServerInfoResponse res = sdk.auth.getServerInfo();

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

**[com.formance.formance_sdk.models.operations.GetServerInfoResponse](../../models/operations/GetServerInfoResponse.md)**


## listClients

List clients

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListClientsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("ipsum") {{
                    authorization = "";
                }})
                .build();

            ListClientsResponse res = sdk.auth.listClients();

            if (res.listClientsResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```


### Response

**[com.formance.formance_sdk.models.operations.ListClientsResponse](../../models/operations/ListClientsResponse.md)**


## listScopes

List Scopes

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListScopesResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("excepturi") {{
                    authorization = "";
                }})
                .build();

            ListScopesResponse res = sdk.auth.listScopes();

            if (res.listScopesResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```


### Response

**[com.formance.formance_sdk.models.operations.ListScopesResponse](../../models/operations/ListScopesResponse.md)**


## listUsers

List users

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListUsersResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("aspernatur") {{
                    authorization = "";
                }})
                .build();

            ListUsersResponse res = sdk.auth.listUsers();

            if (res.listUsersResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```


### Response

**[com.formance.formance_sdk.models.operations.ListUsersResponse](../../models/operations/ListUsersResponse.md)**


## readClient

Read client

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ReadClientRequest;
import com.formance.formance_sdk.models.operations.ReadClientResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("perferendis") {{
                    authorization = "";
                }})
                .build();

            ReadClientRequest req = new ReadClientRequest("ad");            

            ReadClientResponse res = sdk.auth.readClient(req);

            if (res.readClientResponse != null) {
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
| `request`                                                                                                     | [com.formance.formance_sdk.models.operations.ReadClientRequest](../../models/operations/ReadClientRequest.md) | :heavy_check_mark:                                                                                            | The request object to use for the request.                                                                    |


### Response

**[com.formance.formance_sdk.models.operations.ReadClientResponse](../../models/operations/ReadClientResponse.md)**


## readScope

Read scope

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ReadScopeRequest;
import com.formance.formance_sdk.models.operations.ReadScopeResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("natus") {{
                    authorization = "";
                }})
                .build();

            ReadScopeRequest req = new ReadScopeRequest("sed");            

            ReadScopeResponse res = sdk.auth.readScope(req);

            if (res.readScopeResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                   | Type                                                                                                        | Required                                                                                                    | Description                                                                                                 |
| ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                   | [com.formance.formance_sdk.models.operations.ReadScopeRequest](../../models/operations/ReadScopeRequest.md) | :heavy_check_mark:                                                                                          | The request object to use for the request.                                                                  |


### Response

**[com.formance.formance_sdk.models.operations.ReadScopeResponse](../../models/operations/ReadScopeResponse.md)**


## readUser

Read user

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ReadUserRequest;
import com.formance.formance_sdk.models.operations.ReadUserResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("iste") {{
                    authorization = "";
                }})
                .build();

            ReadUserRequest req = new ReadUserRequest("dolor");            

            ReadUserResponse res = sdk.auth.readUser(req);

            if (res.readUserResponse != null) {
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
| `request`                                                                                                 | [com.formance.formance_sdk.models.operations.ReadUserRequest](../../models/operations/ReadUserRequest.md) | :heavy_check_mark:                                                                                        | The request object to use for the request.                                                                |


### Response

**[com.formance.formance_sdk.models.operations.ReadUserResponse](../../models/operations/ReadUserResponse.md)**


## updateClient

Update client

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.UpdateClientRequest;
import com.formance.formance_sdk.models.operations.UpdateClientResponse;
import com.formance.formance_sdk.models.shared.Security;
import com.formance.formance_sdk.models.shared.UpdateClientRequest;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("natus") {{
                    authorization = "";
                }})
                .build();

            UpdateClientRequest req = new UpdateClientRequest("laboriosam") {{
                updateClientRequest = new UpdateClientRequest("hic") {{
                    description = "saepe";
                    metadata = new java.util.HashMap<String, Object>() {{
                        put("in", "corporis");
                        put("iste", "iure");
                        put("saepe", "quidem");
                    }};
                    postLogoutRedirectUris = new String[]{{
                        add("ipsa"),
                    }};
                    public_ = false;
                    redirectUris = new String[]{{
                        add("est"),
                        add("mollitia"),
                        add("laborum"),
                        add("dolores"),
                    }};
                    trusted = false;
                }};;
            }};            

            UpdateClientResponse res = sdk.auth.updateClient(req);

            if (res.updateClientResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                         | Type                                                                                                              | Required                                                                                                          | Description                                                                                                       |
| ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                         | [com.formance.formance_sdk.models.operations.UpdateClientRequest](../../models/operations/UpdateClientRequest.md) | :heavy_check_mark:                                                                                                | The request object to use for the request.                                                                        |


### Response

**[com.formance.formance_sdk.models.operations.UpdateClientResponse](../../models/operations/UpdateClientResponse.md)**


## updateScope

Update scope

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.UpdateScopeRequest;
import com.formance.formance_sdk.models.operations.UpdateScopeResponse;
import com.formance.formance_sdk.models.shared.Security;
import com.formance.formance_sdk.models.shared.UpdateScopeRequest;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("dolorem") {{
                    authorization = "";
                }})
                .build();

            UpdateScopeRequest req = new UpdateScopeRequest("corporis") {{
                updateScopeRequest = new UpdateScopeRequest("explicabo") {{
                    metadata = new java.util.HashMap<String, Object>() {{
                        put("enim", "omnis");
                        put("nemo", "minima");
                        put("excepturi", "accusantium");
                        put("iure", "culpa");
                    }};
                }};;
            }};            

            UpdateScopeResponse res = sdk.auth.updateScope(req);

            if (res.updateScopeResponse != null) {
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
| `request`                                                                                                       | [com.formance.formance_sdk.models.operations.UpdateScopeRequest](../../models/operations/UpdateScopeRequest.md) | :heavy_check_mark:                                                                                              | The request object to use for the request.                                                                      |


### Response

**[com.formance.formance_sdk.models.operations.UpdateScopeResponse](../../models/operations/UpdateScopeResponse.md)**

