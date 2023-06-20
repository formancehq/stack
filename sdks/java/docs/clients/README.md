# clients

### Available Operations

* [addScopeToClient](#addscopetoclient) - Add scope to client
* [createClient](#createclient) - Create client
* [createSecret](#createsecret) - Add a secret to a client
* [deleteClient](#deleteclient) - Delete client
* [deleteScopeFromClient](#deletescopefromclient) - Delete scope from client
* [deleteSecret](#deletesecret) - Delete a secret from a client
* [listClients](#listclients) - List clients
* [readClient](#readclient) - Read client
* [updateClient](#updateclient) - Update client

## addScopeToClient

Add scope to client

### Example Usage

```java
package hello.world;

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.AddScopeToClientRequest;
import org.openapis.openapi.models.operations.AddScopeToClientResponse;
import org.openapis.openapi.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("animi") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            AddScopeToClientRequest req = new AddScopeToClientRequest("enim", "odit");            

            AddScopeToClientResponse res = sdk.clients.addScopeToClient(req);

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

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.CreateClientResponse;
import org.openapis.openapi.models.shared.CreateClientRequest;
import org.openapis.openapi.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("quo") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            org.openapis.openapi.models.shared.CreateClientRequest req = new CreateClientRequest("sequi") {{
                description = "tenetur";
                metadata = new java.util.HashMap<String, Object>() {{
                    put("id", "possimus");
                    put("aut", "quasi");
                }};
                postLogoutRedirectUris = new String[]{{
                    add("temporibus"),
                    add("laborum"),
                    add("quasi"),
                }};
                public_ = false;
                redirectUris = new String[]{{
                    add("voluptatibus"),
                    add("vero"),
                    add("nihil"),
                    add("praesentium"),
                }};
                trusted = false;
            }};            

            CreateClientResponse res = sdk.clients.createClient(req);

            if (res.createClientResponse != null) {
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

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.CreateSecretRequest;
import org.openapis.openapi.models.operations.CreateSecretResponse;
import org.openapis.openapi.models.shared.CreateSecretRequest;
import org.openapis.openapi.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("voluptatibus") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            CreateSecretRequest req = new CreateSecretRequest("ipsa") {{
                createSecretRequest = new CreateSecretRequest("omnis") {{
                    metadata = new java.util.HashMap<String, Object>() {{
                        put("cum", "perferendis");
                        put("doloremque", "reprehenderit");
                    }};
                }};;
            }};            

            CreateSecretResponse res = sdk.clients.createSecret(req);

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

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.DeleteClientRequest;
import org.openapis.openapi.models.operations.DeleteClientResponse;
import org.openapis.openapi.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("ut") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            DeleteClientRequest req = new DeleteClientRequest("maiores");            

            DeleteClientResponse res = sdk.clients.deleteClient(req);

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

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.DeleteScopeFromClientRequest;
import org.openapis.openapi.models.operations.DeleteScopeFromClientResponse;
import org.openapis.openapi.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("dicta") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            DeleteScopeFromClientRequest req = new DeleteScopeFromClientRequest("corporis", "dolore");            

            DeleteScopeFromClientResponse res = sdk.clients.deleteScopeFromClient(req);

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

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.DeleteSecretRequest;
import org.openapis.openapi.models.operations.DeleteSecretResponse;
import org.openapis.openapi.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("iusto") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            DeleteSecretRequest req = new DeleteSecretRequest("dicta", "harum");            

            DeleteSecretResponse res = sdk.clients.deleteSecret(req);

            if (res.statusCode == 200) {
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

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.ListClientsResponse;
import org.openapis.openapi.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("enim") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListClientsResponse res = sdk.clients.listClients();

            if (res.listClientsResponse != null) {
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

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.ReadClientRequest;
import org.openapis.openapi.models.operations.ReadClientResponse;
import org.openapis.openapi.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("accusamus") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ReadClientRequest req = new ReadClientRequest("commodi");            

            ReadClientResponse res = sdk.clients.readClient(req);

            if (res.readClientResponse != null) {
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

import org.openapis.openapi.SDK;
import org.openapis.openapi.models.operations.UpdateClientRequest;
import org.openapis.openapi.models.operations.UpdateClientResponse;
import org.openapis.openapi.models.shared.Security;
import org.openapis.openapi.models.shared.UpdateClientRequest;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("repudiandae") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            UpdateClientRequest req = new UpdateClientRequest("quae") {{
                updateClientRequest = new UpdateClientRequest("ipsum") {{
                    description = "quidem";
                    metadata = new java.util.HashMap<String, Object>() {{
                        put("excepturi", "pariatur");
                        put("modi", "praesentium");
                        put("rem", "voluptates");
                    }};
                    postLogoutRedirectUris = new String[]{{
                        add("repudiandae"),
                    }};
                    public_ = false;
                    redirectUris = new String[]{{
                        add("veritatis"),
                        add("itaque"),
                        add("incidunt"),
                    }};
                    trusted = false;
                }};;
            }};            

            UpdateClientResponse res = sdk.clients.updateClient(req);

            if (res.updateClientResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```
