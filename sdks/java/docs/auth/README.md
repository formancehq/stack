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
                .setSecurity(new Security("deserunt") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            AddScopeToClientRequest req = new AddScopeToClientRequest("perferendis", "ipsam");            

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
                .setSecurity(new Security("repellendus") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            AddTransientScopeRequest req = new AddTransientScopeRequest("sapiente", "quo");            

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
                .setSecurity(new Security("odit") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            com.formance.formance_sdk.models.shared.CreateClientRequest req = new CreateClientRequest("at") {{
                description = "at";
                metadata = new java.util.HashMap<String, Object>() {{
                    put("molestiae", "quod");
                    put("quod", "esse");
                    put("totam", "porro");
                    put("dolorum", "dicta");
                }};
                postLogoutRedirectUris = new String[]{{
                    add("officia"),
                    add("occaecati"),
                    add("fugit"),
                }};
                public_ = false;
                redirectUris = new String[]{{
                    add("hic"),
                    add("optio"),
                    add("totam"),
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
                .setSecurity(new Security("beatae") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            com.formance.formance_sdk.models.shared.CreateScopeRequest req = new CreateScopeRequest("commodi") {{
                metadata = new java.util.HashMap<String, Object>() {{
                    put("modi", "qui");
                    put("impedit", "cum");
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
                .setSecurity(new Security("esse") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            CreateSecretRequest req = new CreateSecretRequest("ipsum") {{
                createSecretRequest = new CreateSecretRequest("excepturi") {{
                    metadata = new java.util.HashMap<String, Object>() {{
                        put("perferendis", "ad");
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
                .setSecurity(new Security("natus") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            DeleteClientRequest req = new DeleteClientRequest("sed");            

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
                .setSecurity(new Security("iste") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            DeleteScopeRequest req = new DeleteScopeRequest("dolor");            

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
                .setSecurity(new Security("natus") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            DeleteScopeFromClientRequest req = new DeleteScopeFromClientRequest("laboriosam", "hic");            

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
                .setSecurity(new Security("saepe") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            DeleteSecretRequest req = new DeleteSecretRequest("fuga", "in");            

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
                .setSecurity(new Security("corporis") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            DeleteTransientScopeRequest req = new DeleteTransientScopeRequest("iste", "iure");            

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
                .setSecurity(new Security("saepe") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
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
                .setSecurity(new Security("quidem") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
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
                .setSecurity(new Security("architecto") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
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
                .setSecurity(new Security("ipsa") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
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
                .setSecurity(new Security("reiciendis") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ReadClientRequest req = new ReadClientRequest("est");            

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
                .setSecurity(new Security("mollitia") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ReadScopeRequest req = new ReadScopeRequest("laborum");            

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
                .setSecurity(new Security("dolores") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ReadUserRequest req = new ReadUserRequest("dolorem");            

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
                .setSecurity(new Security("corporis") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            UpdateClientRequest req = new UpdateClientRequest("explicabo") {{
                updateClientRequest = new UpdateClientRequest("nobis") {{
                    description = "enim";
                    metadata = new java.util.HashMap<String, Object>() {{
                        put("nemo", "minima");
                        put("excepturi", "accusantium");
                        put("iure", "culpa");
                    }};
                    postLogoutRedirectUris = new String[]{{
                        add("sapiente"),
                        add("architecto"),
                        add("mollitia"),
                        add("dolorem"),
                    }};
                    public_ = false;
                    redirectUris = new String[]{{
                        add("consequuntur"),
                        add("repellat"),
                        add("mollitia"),
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
                .setSecurity(new Security("occaecati") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            UpdateScopeRequest req = new UpdateScopeRequest("numquam") {{
                updateScopeRequest = new UpdateScopeRequest("commodi") {{
                    metadata = new java.util.HashMap<String, Object>() {{
                        put("molestiae", "velit");
                        put("error", "quia");
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
