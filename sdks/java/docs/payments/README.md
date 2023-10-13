# payments

### Available Operations

* [connectorsTransfer](#connectorstransfer) - Transfer funds between Connector accounts
* [createBankAccount](#createbankaccount) - Create a BankAccount in Payments and on the PSP
* [createTransferInitiation](#createtransferinitiation) - Create a TransferInitiation
* [deleteTransferInitiation](#deletetransferinitiation) - Delete a transfer initiation
* [getAccountBalances](#getaccountbalances) - Get account balances
* [getBankAccount](#getbankaccount) - Get a bank account created by user on Formance
* [getConnectorTask](#getconnectortask) - Read a specific task of the connector
* [getPayment](#getpayment) - Get a payment
* [getTransferInitiation](#gettransferinitiation) - Get a transfer initiation
* [installConnector](#installconnector) - Install a connector
* [listAllConnectors](#listallconnectors) - List all installed connectors
* [listBankAccounts](#listbankaccounts) - List bank accounts created by user on Formance
* [listConfigsAvailableConnectors](#listconfigsavailableconnectors) - List the configs of each available connector
* [listConnectorTasks](#listconnectortasks) - List tasks from a connector
* [listConnectorsTransfers](#listconnectorstransfers) - List transfers and their statuses
* [listPayments](#listpayments) - List payments
* [listTransferInitiations](#listtransferinitiations) - List Transfer Initiations
* [paymentsgetAccount](#paymentsgetaccount) - Get an account
* [paymentsgetServerInfo](#paymentsgetserverinfo) - Get server info
* [paymentslistAccounts](#paymentslistaccounts) - List accounts
* [readConnectorConfig](#readconnectorconfig) - Read the config of a connector
* [resetConnector](#resetconnector) - Reset a connector
* [retryTransferInitiation](#retrytransferinitiation) - Retry a failed transfer initiation
* [udpateTransferInitiationStatus](#udpatetransferinitiationstatus) - Update the status of a transfer initiation
* [uninstallConnector](#uninstallconnector) - Uninstall a connector
* [updateMetadata](#updatemetadata) - Update metadata

## connectorsTransfer

Execute a transfer between two accounts.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ConnectorsTransferRequest;
import com.formance.formance_sdk.models.operations.ConnectorsTransferResponse;
import com.formance.formance_sdk.models.shared.Connector;
import com.formance.formance_sdk.models.shared.Security;
import com.formance.formance_sdk.models.shared.TransferRequest;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("sapiente") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ConnectorsTransferRequest req = new ConnectorsTransferRequest(                new TransferRequest(100L, "USD", "acct_1Gqj58KZcSIg2N2q") {{
                                source = "acct_1Gqj58KZcSIg2N2q";
                            }};, Connector.DUMMY_PAY);            

            ConnectorsTransferResponse res = sdk.payments.connectorsTransfer(req);

            if (res.transferResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## createBankAccount

Create a bank account in Payments and on the PSP.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CreateBankAccountResponse;
import com.formance.formance_sdk.models.shared.BankAccountRequest;
import com.formance.formance_sdk.models.shared.Connector;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("deserunt") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            com.formance.formance_sdk.models.shared.BankAccountRequest req = new BankAccountRequest("GB", "My account", Connector.MODULR) {{
                accountNumber = "vel";
                iban = "natus";
                swiftBicCode = "omnis";
            }};            

            CreateBankAccountResponse res = sdk.payments.createBankAccount(req);

            if (res.bankAccountResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## createTransferInitiation

Create a transfer initiation

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CreateTransferInitiationResponse;
import com.formance.formance_sdk.models.shared.Connector;
import com.formance.formance_sdk.models.shared.Security;
import com.formance.formance_sdk.models.shared.TransferInitiationRequest;
import com.formance.formance_sdk.models.shared.TransferInitiationRequestType;
import java.time.OffsetDateTime;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("molestiae") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            com.formance.formance_sdk.models.shared.TransferInitiationRequest req = new TransferInitiationRequest(19193L, "USD", OffsetDateTime.parse("2022-09-12T22:12:15.897Z"), "distinctio", "id", Connector.WISE, "XXX", "labore", TransferInitiationRequestType.TRANSFER, false);            

            CreateTransferInitiationResponse res = sdk.payments.createTransferInitiation(req);

            if (res.transferInitiationResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## deleteTransferInitiation

Delete a transfer initiation by its id.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.DeleteTransferInitiationRequest;
import com.formance.formance_sdk.models.operations.DeleteTransferInitiationResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("natus") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            DeleteTransferInitiationRequest req = new DeleteTransferInitiationRequest("nobis");            

            DeleteTransferInitiationResponse res = sdk.payments.deleteTransferInitiation(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getAccountBalances

Get account balances

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetAccountBalancesRequest;
import com.formance.formance_sdk.models.operations.GetAccountBalancesResponse;
import com.formance.formance_sdk.models.shared.Security;
import java.time.OffsetDateTime;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("eum") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetAccountBalancesRequest req = new GetAccountBalancesRequest("vero") {{
                asset = "aspernatur";
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                from = OffsetDateTime.parse("2022-09-14T03:02:47.808Z");
                limit = 92373L;
                pageSize = 569965L;
                sort = new String[]{{
                    add("provident"),
                    add("quos"),
                }};
                to = OffsetDateTime.parse("2022-12-07T10:53:17.121Z");
            }};            

            GetAccountBalancesResponse res = sdk.payments.getAccountBalances(req);

            if (res.balancesCursor != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getBankAccount

Get a bank account created by user on Formance

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetBankAccountRequest;
import com.formance.formance_sdk.models.operations.GetBankAccountResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("mollitia") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetBankAccountRequest req = new GetBankAccountRequest("reiciendis");            

            GetBankAccountResponse res = sdk.payments.getBankAccount(req);

            if (res.bankAccountResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getConnectorTask

Get a specific task associated to the connector.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetConnectorTaskRequest;
import com.formance.formance_sdk.models.operations.GetConnectorTaskResponse;
import com.formance.formance_sdk.models.shared.Connector;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("mollitia") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetConnectorTaskRequest req = new GetConnectorTaskRequest(Connector.WISE, "eum");            

            GetConnectorTaskResponse res = sdk.payments.getConnectorTask(req);

            if (res.taskResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getPayment

Get a payment

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetPaymentRequest;
import com.formance.formance_sdk.models.operations.GetPaymentResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("dolor") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetPaymentRequest req = new GetPaymentRequest("necessitatibus");            

            GetPaymentResponse res = sdk.payments.getPayment(req);

            if (res.paymentResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getTransferInitiation

Get a transfer initiation

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetTransferInitiationRequest;
import com.formance.formance_sdk.models.operations.GetTransferInitiationResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("odit") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetTransferInitiationRequest req = new GetTransferInitiationRequest("nemo");            

            GetTransferInitiationResponse res = sdk.payments.getTransferInitiation(req);

            if (res.transferInitiationResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## installConnector

Install a connector by its name and config.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.InstallConnectorRequest;
import com.formance.formance_sdk.models.operations.InstallConnectorResponse;
import com.formance.formance_sdk.models.shared.BankingCircleConfig;
import com.formance.formance_sdk.models.shared.Connector;
import com.formance.formance_sdk.models.shared.CurrencyCloudConfig;
import com.formance.formance_sdk.models.shared.DummyPayConfig;
import com.formance.formance_sdk.models.shared.MangoPayConfig;
import com.formance.formance_sdk.models.shared.ModulrConfig;
import com.formance.formance_sdk.models.shared.MoneycorpConfig;
import com.formance.formance_sdk.models.shared.Security;
import com.formance.formance_sdk.models.shared.StripeConfig;
import com.formance.formance_sdk.models.shared.WiseConfig;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("quasi") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            InstallConnectorRequest req = new InstallConnectorRequest(                new ModulrConfig("XXX", "XXX") {{
                                apiKey = "XXX";
                                apiSecret = "XXX";
                                endpoint = "XXX";
                                pollingPeriod = "60s";
                            }}, Connector.MONEYCORP);            

            InstallConnectorResponse res = sdk.payments.installConnector(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listAllConnectors

List all installed connectors.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListAllConnectorsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("debitis") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListAllConnectorsResponse res = sdk.payments.listAllConnectors();

            if (res.connectorsResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listBankAccounts

List all bank accounts created by user on Formance.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListBankAccountsRequest;
import com.formance.formance_sdk.models.operations.ListBankAccountsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("eius") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListBankAccountsRequest req = new ListBankAccountsRequest() {{
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                pageSize = 806194L;
                sort = new String[]{{
                    add("facilis"),
                    add("in"),
                    add("architecto"),
                }};
            }};            

            ListBankAccountsResponse res = sdk.payments.listBankAccounts(req);

            if (res.bankAccountsCursor != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listConfigsAvailableConnectors

List the configs of each available connector.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListConfigsAvailableConnectorsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("architecto") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListConfigsAvailableConnectorsResponse res = sdk.payments.listConfigsAvailableConnectors();

            if (res.connectorsConfigsResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listConnectorTasks

List all tasks associated with this connector.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListConnectorTasksRequest;
import com.formance.formance_sdk.models.operations.ListConnectorTasksResponse;
import com.formance.formance_sdk.models.shared.Connector;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("repudiandae") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListConnectorTasksRequest req = new ListConnectorTasksRequest(Connector.WISE) {{
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                pageSize = 714242L;
            }};            

            ListConnectorTasksResponse res = sdk.payments.listConnectorTasks(req);

            if (res.tasksCursor != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listConnectorsTransfers

List transfers

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListConnectorsTransfersRequest;
import com.formance.formance_sdk.models.operations.ListConnectorsTransfersResponse;
import com.formance.formance_sdk.models.shared.Connector;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("nihil") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListConnectorsTransfersRequest req = new ListConnectorsTransfersRequest(Connector.MONEYCORP);            

            ListConnectorsTransfersResponse res = sdk.payments.listConnectorsTransfers(req);

            if (res.transfersResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listPayments

List payments

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListPaymentsRequest;
import com.formance.formance_sdk.models.operations.ListPaymentsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("quibusdam") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListPaymentsRequest req = new ListPaymentsRequest() {{
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                pageSize = 149448L;
                sort = new String[]{{
                    add("pariatur"),
                    add("accusantium"),
                    add("consequuntur"),
                    add("praesentium"),
                }};
            }};            

            ListPaymentsResponse res = sdk.payments.listPayments(req);

            if (res.paymentsCursor != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listTransferInitiations

List Transfer Initiations

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListTransferInitiationsRequest;
import com.formance.formance_sdk.models.operations.ListTransferInitiationsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("natus") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListTransferInitiationsRequest req = new ListTransferInitiationsRequest() {{
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                pageSize = 166847L;
                query = "sunt";
                sort = new String[]{{
                    add("illum"),
                    add("pariatur"),
                    add("maxime"),
                    add("ea"),
                }};
            }};            

            ListTransferInitiationsResponse res = sdk.payments.listTransferInitiations(req);

            if (res.transferInitiationsCursor != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## paymentsgetAccount

Get an account

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.PaymentsgetAccountRequest;
import com.formance.formance_sdk.models.operations.PaymentsgetAccountResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("excepturi") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            PaymentsgetAccountRequest req = new PaymentsgetAccountRequest("odit");            

            PaymentsgetAccountResponse res = sdk.payments.paymentsgetAccount(req);

            if (res.paymentsAccountResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## paymentsgetServerInfo

Get server info

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.PaymentsgetServerInfoResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("ea") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            PaymentsgetServerInfoResponse res = sdk.payments.paymentsgetServerInfo();

            if (res.serverInfo != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## paymentslistAccounts

List accounts

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.PaymentslistAccountsRequest;
import com.formance.formance_sdk.models.operations.PaymentslistAccountsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("accusantium") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            PaymentslistAccountsRequest req = new PaymentslistAccountsRequest() {{
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                pageSize = 69167L;
                sort = new String[]{{
                    add("quidem"),
                    add("ipsam"),
                    add("voluptate"),
                    add("autem"),
                }};
            }};            

            PaymentslistAccountsResponse res = sdk.payments.paymentslistAccounts(req);

            if (res.accountsCursor != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## readConnectorConfig

Read connector config

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ReadConnectorConfigRequest;
import com.formance.formance_sdk.models.operations.ReadConnectorConfigResponse;
import com.formance.formance_sdk.models.shared.Connector;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("nam") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ReadConnectorConfigRequest req = new ReadConnectorConfigRequest(Connector.STRIPE);            

            ReadConnectorConfigResponse res = sdk.payments.readConnectorConfig(req);

            if (res.connectorConfigResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## resetConnector

Reset a connector by its name.
It will remove the connector and ALL PAYMENTS generated with it.


### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ResetConnectorRequest;
import com.formance.formance_sdk.models.operations.ResetConnectorResponse;
import com.formance.formance_sdk.models.shared.Connector;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("pariatur") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ResetConnectorRequest req = new ResetConnectorRequest(Connector.WISE);            

            ResetConnectorResponse res = sdk.payments.resetConnector(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## retryTransferInitiation

Retry a failed transfer initiation

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.RetryTransferInitiationRequest;
import com.formance.formance_sdk.models.operations.RetryTransferInitiationResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("voluptatibus") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            RetryTransferInitiationRequest req = new RetryTransferInitiationRequest("perferendis");            

            RetryTransferInitiationResponse res = sdk.payments.retryTransferInitiation(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## udpateTransferInitiationStatus

Update a transfer initiation status

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.UdpateTransferInitiationStatusRequest;
import com.formance.formance_sdk.models.operations.UdpateTransferInitiationStatusResponse;
import com.formance.formance_sdk.models.shared.Security;
import com.formance.formance_sdk.models.shared.UpdateTransferInitiationStatusRequest;
import com.formance.formance_sdk.models.shared.UpdateTransferInitiationStatusRequestStatus;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("fugiat") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            UdpateTransferInitiationStatusRequest req = new UdpateTransferInitiationStatusRequest(                new UpdateTransferInitiationStatusRequest(UpdateTransferInitiationStatusRequestStatus.PROCESSING);, "aut");            

            UdpateTransferInitiationStatusResponse res = sdk.payments.udpateTransferInitiationStatus(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## uninstallConnector

Uninstall a connector by its name.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.UninstallConnectorRequest;
import com.formance.formance_sdk.models.operations.UninstallConnectorResponse;
import com.formance.formance_sdk.models.shared.Connector;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("cumque") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            UninstallConnectorRequest req = new UninstallConnectorRequest(Connector.WISE);            

            UninstallConnectorResponse res = sdk.payments.uninstallConnector(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## updateMetadata

Update metadata

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.UpdateMetadataRequest;
import com.formance.formance_sdk.models.operations.UpdateMetadataResponse;
import com.formance.formance_sdk.models.shared.PaymentMetadata;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("hic") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            UpdateMetadataRequest req = new UpdateMetadataRequest(                new PaymentMetadata() {{
                                key = "libero";
                            }};, "nobis");            

            UpdateMetadataResponse res = sdk.payments.updateMetadata(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```
