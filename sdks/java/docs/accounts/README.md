# accounts

### Available Operations

* [addMetadataToAccount](#addmetadatatoaccount) - Add metadata to an account
* [countAccounts](#countaccounts) - Count the accounts from a ledger
* [getAccount](#getaccount) - Get account by its address
* [listAccounts](#listaccounts) - List accounts from a ledger

## addMetadataToAccount

Add metadata to an account

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.AddMetadataToAccountRequest;
import com.formance.formance_sdk.models.operations.AddMetadataToAccountResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("distinctio") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            AddMetadataToAccountRequest req = new AddMetadataToAccountRequest(                new java.util.HashMap<String, String>() {{
                                put("unde", "nulla");
                                put("corrupti", "illum");
                                put("vel", "error");
                                put("deserunt", "suscipit");
                            }}, "users:001", "ledger001") {{
                idempotencyKey = "iure";
                async = true;
                dryRun = true;
            }};            

            AddMetadataToAccountResponse res = sdk.accounts.addMetadataToAccount(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## countAccounts

Count the accounts from a ledger

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CountAccountsRequest;
import com.formance.formance_sdk.models.operations.CountAccountsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("magnam") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            CountAccountsRequest req = new CountAccountsRequest("ledger001") {{
                address = "users:.+";
                metadata = new java.util.HashMap<String, Object>() {{
                    put("ipsa", "delectus");
                    put("tempora", "suscipit");
                    put("molestiae", "minus");
                    put("placeat", "voluptatum");
                }};
            }};            

            CountAccountsResponse res = sdk.accounts.countAccounts(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getAccount

Get account by its address

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetAccountRequest;
import com.formance.formance_sdk.models.operations.GetAccountResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("iusto") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetAccountRequest req = new GetAccountRequest("users:001", "ledger001");            

            GetAccountResponse res = sdk.accounts.getAccount(req);

            if (res.accountResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listAccounts

List accounts from a ledger, sorted by address in descending order.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListAccountsBalanceOperator;
import com.formance.formance_sdk.models.operations.ListAccountsRequest;
import com.formance.formance_sdk.models.operations.ListAccountsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("excepturi") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListAccountsRequest req = new ListAccountsRequest("ledger001") {{
                address = "users:.+";
                balance = 2400L;
                balanceOperator = ListAccountsBalanceOperator.GTE;
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                metadata = new java.util.HashMap<String, String>() {{
                    put("recusandae", "temporibus");
                    put("ab", "quis");
                }};
                pageSize = 87129L;
            }};            

            ListAccountsResponse res = sdk.accounts.listAccounts(req);

            if (res.accountsCursorResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```
