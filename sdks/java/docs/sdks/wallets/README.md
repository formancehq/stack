# wallets

### Available Operations

* [confirmHold](#confirmhold) - Confirm a hold
* [createBalance](#createbalance) - Create a balance
* [createWallet](#createwallet) - Create a new wallet
* [creditWallet](#creditwallet) - Credit a wallet
* [debitWallet](#debitwallet) - Debit a wallet
* [getBalance](#getbalance) - Get detailed balance
* [getHold](#gethold) - Get a hold
* [getHolds](#getholds) - Get all holds for a wallet
* [getTransactions](#gettransactions)
* [getWallet](#getwallet) - Get a wallet
* [getWalletSummary](#getwalletsummary) - Get wallet summary
* [listBalances](#listbalances) - List balances of a wallet
* [listWallets](#listwallets) - List all wallets
* [updateWallet](#updatewallet) - Update a wallet
* [voidHold](#voidhold) - Cancel a hold
* [walletsgetServerInfo](#walletsgetserverinfo) - Get server info

## confirmHold

Confirm a hold

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ConfirmHoldRequest;
import com.formance.formance_sdk.models.operations.ConfirmHoldResponse;
import com.formance.formance_sdk.models.shared.ConfirmHoldRequest;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("ea") {{
                    authorization = "";
                }})
                .build();

            ConfirmHoldRequest req = new ConfirmHoldRequest("excepturi") {{
                confirmHoldRequest = new ConfirmHoldRequest() {{
                    amount = 100L;
                    final_ = true;
                }};;
            }};            

            ConfirmHoldResponse res = sdk.wallets.confirmHold(req);

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
| `request`                                                                                                       | [com.formance.formance_sdk.models.operations.ConfirmHoldRequest](../../models/operations/ConfirmHoldRequest.md) | :heavy_check_mark:                                                                                              | The request object to use for the request.                                                                      |


### Response

**[com.formance.formance_sdk.models.operations.ConfirmHoldResponse](../../models/operations/ConfirmHoldResponse.md)**


## createBalance

Create a balance

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CreateBalanceRequest;
import com.formance.formance_sdk.models.operations.CreateBalanceResponse;
import com.formance.formance_sdk.models.shared.CreateBalanceRequest;
import com.formance.formance_sdk.models.shared.Security;
import java.time.OffsetDateTime;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("odit") {{
                    authorization = "";
                }})
                .build();

            CreateBalanceRequest req = new CreateBalanceRequest("ea") {{
                createBalanceRequest = new CreateBalanceRequest("accusantium") {{
                    expiresAt = OffsetDateTime.parse("2022-01-07T08:38:36.147Z");
                    priority = 697429L;
                }};;
            }};            

            CreateBalanceResponse res = sdk.wallets.createBalance(req);

            if (res.createBalanceResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                           | Type                                                                                                                | Required                                                                                                            | Description                                                                                                         |
| ------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                           | [com.formance.formance_sdk.models.operations.CreateBalanceRequest](../../models/operations/CreateBalanceRequest.md) | :heavy_check_mark:                                                                                                  | The request object to use for the request.                                                                          |


### Response

**[com.formance.formance_sdk.models.operations.CreateBalanceResponse](../../models/operations/CreateBalanceResponse.md)**


## createWallet

Create a new wallet

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CreateWalletResponse;
import com.formance.formance_sdk.models.shared.CreateWalletRequest;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("ipsam") {{
                    authorization = "";
                }})
                .build();

            com.formance.formance_sdk.models.shared.CreateWalletRequest req = new CreateWalletRequest(                new java.util.HashMap<String, String>() {{
                                put("autem", "nam");
                                put("eaque", "pariatur");
                            }}, "nemo");            

            CreateWalletResponse res = sdk.wallets.createWallet(req);

            if (res.createWalletResponse != null) {
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
| `request`                                                                                                 | [com.formance.formance_sdk.models.shared.CreateWalletRequest](../../models/shared/CreateWalletRequest.md) | :heavy_check_mark:                                                                                        | The request object to use for the request.                                                                |


### Response

**[com.formance.formance_sdk.models.operations.CreateWalletResponse](../../models/operations/CreateWalletResponse.md)**


## creditWallet

Credit a wallet

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CreditWalletRequest;
import com.formance.formance_sdk.models.operations.CreditWalletResponse;
import com.formance.formance_sdk.models.shared.CreditWalletRequest;
import com.formance.formance_sdk.models.shared.LedgerAccountSubject;
import com.formance.formance_sdk.models.shared.Monetary;
import com.formance.formance_sdk.models.shared.Security;
import com.formance.formance_sdk.models.shared.WalletSubject;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("voluptatibus") {{
                    authorization = "";
                }})
                .build();

            CreditWalletRequest req = new CreditWalletRequest("perferendis") {{
                creditWalletRequest = new CreditWalletRequest(                new Monetary(855804L, "amet");,                 new java.util.HashMap<String, String>() {{
                                    put("cumque", "corporis");
                                }},                 new Object[]{{
                                    add(new WalletSubject("totam", "dignissimos") {{
                                        balance = "nobis";
                                        identifier = "dolores";
                                        type = "quis";
                                    }}),
                                    add(new LedgerAccountSubject("eos", "perferendis") {{
                                        identifier = "quis";
                                        type = "nesciunt";
                                    }}),
                                    add(new LedgerAccountSubject("dolor", "vero") {{
                                        identifier = "minus";
                                        type = "quam";
                                    }}),
                                    add(new LedgerAccountSubject("omnis", "facilis") {{
                                        identifier = "hic";
                                        type = "recusandae";
                                    }}),
                                }}) {{
                    balance = "perspiciatis";
                    reference = "voluptatem";
                }};;
            }};            

            CreditWalletResponse res = sdk.wallets.creditWallet(req);

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
| `request`                                                                                                         | [com.formance.formance_sdk.models.operations.CreditWalletRequest](../../models/operations/CreditWalletRequest.md) | :heavy_check_mark:                                                                                                | The request object to use for the request.                                                                        |


### Response

**[com.formance.formance_sdk.models.operations.CreditWalletResponse](../../models/operations/CreditWalletResponse.md)**


## debitWallet

Debit a wallet

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.DebitWalletRequest;
import com.formance.formance_sdk.models.operations.DebitWalletResponse;
import com.formance.formance_sdk.models.shared.DebitWalletRequest;
import com.formance.formance_sdk.models.shared.LedgerAccountSubject;
import com.formance.formance_sdk.models.shared.Monetary;
import com.formance.formance_sdk.models.shared.Security;
import com.formance.formance_sdk.models.shared.WalletSubject;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("porro") {{
                    authorization = "";
                }})
                .build();

            DebitWalletRequest req = new DebitWalletRequest("consequuntur") {{
                debitWalletRequest = new DebitWalletRequest(                new Monetary(500026L, "error");,                 new java.util.HashMap<String, String>() {{
                                    put("occaecati", "rerum");
                                }}) {{
                    balances = new String[]{{
                        add("asperiores"),
                    }};
                    description = "earum";
                    destination = new LedgerAccountSubject("iste", "dolorum");;
                    pending = false;
                }};;
            }};            

            DebitWalletResponse res = sdk.wallets.debitWallet(req);

            if (res.debitWalletResponse != null) {
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
| `request`                                                                                                       | [com.formance.formance_sdk.models.operations.DebitWalletRequest](../../models/operations/DebitWalletRequest.md) | :heavy_check_mark:                                                                                              | The request object to use for the request.                                                                      |


### Response

**[com.formance.formance_sdk.models.operations.DebitWalletResponse](../../models/operations/DebitWalletResponse.md)**


## getBalance

Get detailed balance

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetBalanceRequest;
import com.formance.formance_sdk.models.operations.GetBalanceResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("deleniti") {{
                    authorization = "";
                }})
                .build();

            GetBalanceRequest req = new GetBalanceRequest("pariatur", "provident");            

            GetBalanceResponse res = sdk.wallets.getBalance(req);

            if (res.getBalanceResponse != null) {
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
| `request`                                                                                                     | [com.formance.formance_sdk.models.operations.GetBalanceRequest](../../models/operations/GetBalanceRequest.md) | :heavy_check_mark:                                                                                            | The request object to use for the request.                                                                    |


### Response

**[com.formance.formance_sdk.models.operations.GetBalanceResponse](../../models/operations/GetBalanceResponse.md)**


## getHold

Get a hold

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetHoldRequest;
import com.formance.formance_sdk.models.operations.GetHoldResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("nobis") {{
                    authorization = "";
                }})
                .build();

            GetHoldRequest req = new GetHoldRequest("libero");            

            GetHoldResponse res = sdk.wallets.getHold(req);

            if (res.getHoldResponse != null) {
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
| `request`                                                                                               | [com.formance.formance_sdk.models.operations.GetHoldRequest](../../models/operations/GetHoldRequest.md) | :heavy_check_mark:                                                                                      | The request object to use for the request.                                                              |


### Response

**[com.formance.formance_sdk.models.operations.GetHoldResponse](../../models/operations/GetHoldResponse.md)**


## getHolds

Get all holds for a wallet

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetHoldsRequest;
import com.formance.formance_sdk.models.operations.GetHoldsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("delectus") {{
                    authorization = "";
                }})
                .build();

            GetHoldsRequest req = new GetHoldsRequest() {{
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                metadata = new java.util.HashMap<String, String>() {{
                    put("quos", "aliquid");
                    put("dolorem", "dolorem");
                }};
                pageSize = 222443L;
                walletID = "qui";
            }};            

            GetHoldsResponse res = sdk.wallets.getHolds(req);

            if (res.getHoldsResponse != null) {
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
| `request`                                                                                                 | [com.formance.formance_sdk.models.operations.GetHoldsRequest](../../models/operations/GetHoldsRequest.md) | :heavy_check_mark:                                                                                        | The request object to use for the request.                                                                |


### Response

**[com.formance.formance_sdk.models.operations.GetHoldsResponse](../../models/operations/GetHoldsResponse.md)**


## getTransactions

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetTransactionsRequest;
import com.formance.formance_sdk.models.operations.GetTransactionsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("ipsum") {{
                    authorization = "";
                }})
                .build();

            GetTransactionsRequest req = new GetTransactionsRequest() {{
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                pageSize = 944373L;
                walletID = "excepturi";
            }};            

            GetTransactionsResponse res = sdk.wallets.getTransactions(req);

            if (res.getTransactionsResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

### Parameters

| Parameter                                                                                                               | Type                                                                                                                    | Required                                                                                                                | Description                                                                                                             |
| ----------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                               | [com.formance.formance_sdk.models.operations.GetTransactionsRequest](../../models/operations/GetTransactionsRequest.md) | :heavy_check_mark:                                                                                                      | The request object to use for the request.                                                                              |


### Response

**[com.formance.formance_sdk.models.operations.GetTransactionsResponse](../../models/operations/GetTransactionsResponse.md)**


## getWallet

Get a wallet

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetWalletRequest;
import com.formance.formance_sdk.models.operations.GetWalletResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("cum") {{
                    authorization = "";
                }})
                .build();

            GetWalletRequest req = new GetWalletRequest("voluptate");            

            GetWalletResponse res = sdk.wallets.getWallet(req);

            if (res.getWalletResponse != null) {
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
| `request`                                                                                                   | [com.formance.formance_sdk.models.operations.GetWalletRequest](../../models/operations/GetWalletRequest.md) | :heavy_check_mark:                                                                                          | The request object to use for the request.                                                                  |


### Response

**[com.formance.formance_sdk.models.operations.GetWalletResponse](../../models/operations/GetWalletResponse.md)**


## getWalletSummary

Get wallet summary

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetWalletSummaryRequest;
import com.formance.formance_sdk.models.operations.GetWalletSummaryResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("dignissimos") {{
                    authorization = "";
                }})
                .build();

            GetWalletSummaryRequest req = new GetWalletSummaryRequest("reiciendis");            

            GetWalletSummaryResponse res = sdk.wallets.getWalletSummary(req);

            if (res.getWalletSummaryResponse != null) {
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
| `request`                                                                                                                 | [com.formance.formance_sdk.models.operations.GetWalletSummaryRequest](../../models/operations/GetWalletSummaryRequest.md) | :heavy_check_mark:                                                                                                        | The request object to use for the request.                                                                                |


### Response

**[com.formance.formance_sdk.models.operations.GetWalletSummaryResponse](../../models/operations/GetWalletSummaryResponse.md)**


## listBalances

List balances of a wallet

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListBalancesRequest;
import com.formance.formance_sdk.models.operations.ListBalancesResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("amet") {{
                    authorization = "";
                }})
                .build();

            ListBalancesRequest req = new ListBalancesRequest("dolorum");            

            ListBalancesResponse res = sdk.wallets.listBalances(req);

            if (res.listBalancesResponse != null) {
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
| `request`                                                                                                         | [com.formance.formance_sdk.models.operations.ListBalancesRequest](../../models/operations/ListBalancesRequest.md) | :heavy_check_mark:                                                                                                | The request object to use for the request.                                                                        |


### Response

**[com.formance.formance_sdk.models.operations.ListBalancesResponse](../../models/operations/ListBalancesResponse.md)**


## listWallets

List all wallets

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListWalletsRequest;
import com.formance.formance_sdk.models.operations.ListWalletsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("numquam") {{
                    authorization = "";
                }})
                .build();

            ListWalletsRequest req = new ListWalletsRequest() {{
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                metadata = new java.util.HashMap<String, String>() {{
                    put("ipsa", "ipsa");
                }};
                name = "Viola Hahn";
                pageSize = 976405L;
            }};            

            ListWalletsResponse res = sdk.wallets.listWallets(req);

            if (res.listWalletsResponse != null) {
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
| `request`                                                                                                       | [com.formance.formance_sdk.models.operations.ListWalletsRequest](../../models/operations/ListWalletsRequest.md) | :heavy_check_mark:                                                                                              | The request object to use for the request.                                                                      |


### Response

**[com.formance.formance_sdk.models.operations.ListWalletsResponse](../../models/operations/ListWalletsResponse.md)**


## updateWallet

Update a wallet

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.UpdateWalletRequest;
import com.formance.formance_sdk.models.operations.UpdateWalletRequestBody;
import com.formance.formance_sdk.models.operations.UpdateWalletResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("voluptas") {{
                    authorization = "";
                }})
                .build();

            UpdateWalletRequest req = new UpdateWalletRequest("natus") {{
                requestBody = new UpdateWalletRequestBody(                new java.util.HashMap<String, String>() {{
                                    put("atque", "sit");
                                }});;
            }};            

            UpdateWalletResponse res = sdk.wallets.updateWallet(req);

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
| `request`                                                                                                         | [com.formance.formance_sdk.models.operations.UpdateWalletRequest](../../models/operations/UpdateWalletRequest.md) | :heavy_check_mark:                                                                                                | The request object to use for the request.                                                                        |


### Response

**[com.formance.formance_sdk.models.operations.UpdateWalletResponse](../../models/operations/UpdateWalletResponse.md)**


## voidHold

Cancel a hold

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.VoidHoldRequest;
import com.formance.formance_sdk.models.operations.VoidHoldResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("fugiat") {{
                    authorization = "";
                }})
                .build();

            VoidHoldRequest req = new VoidHoldRequest("ab");            

            VoidHoldResponse res = sdk.wallets.voidHold(req);

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

| Parameter                                                                                                 | Type                                                                                                      | Required                                                                                                  | Description                                                                                               |
| --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| `request`                                                                                                 | [com.formance.formance_sdk.models.operations.VoidHoldRequest](../../models/operations/VoidHoldRequest.md) | :heavy_check_mark:                                                                                        | The request object to use for the request.                                                                |


### Response

**[com.formance.formance_sdk.models.operations.VoidHoldResponse](../../models/operations/VoidHoldResponse.md)**


## walletsgetServerInfo

Get server info

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.WalletsgetServerInfoResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("soluta") {{
                    authorization = "";
                }})
                .build();

            WalletsgetServerInfoResponse res = sdk.wallets.walletsgetServerInfo();

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

**[com.formance.formance_sdk.models.operations.WalletsgetServerInfoResponse](../../models/operations/WalletsgetServerInfoResponse.md)**

