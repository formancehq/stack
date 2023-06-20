# payments

### Available Operations

* [connectorsStripeTransfer](#connectorsstripetransfer) - Transfer funds between Stripe accounts
* [connectorsTransfer](#connectorstransfer) - Transfer funds between Connector accounts
* [getConnectorTask](#getconnectortask) - Read a specific task of the connector
* [getPayment](#getpayment) - Get a payment
* [installConnector](#installconnector) - Install a connector
* [listAllConnectors](#listallconnectors) - List all installed connectors
* [listConfigsAvailableConnectors](#listconfigsavailableconnectors) - List the configs of each available connector
* [listConnectorTasks](#listconnectortasks) - List tasks from a connector
* [listConnectorsTransfers](#listconnectorstransfers) - List transfers and their statuses
* [listPayments](#listpayments) - List payments
* [paymentsgetServerInfo](#paymentsgetserverinfo) - Get server info
* [paymentslistAccounts](#paymentslistaccounts) - List accounts
* [readConnectorConfig](#readconnectorconfig) - Read the config of a connector
* [resetConnector](#resetconnector) - Reset a connector
* [uninstallConnector](#uninstallconnector) - Uninstall a connector
* [updateMetadata](#updatemetadata) - Update metadata

## connectorsStripeTransfer

Execute a transfer between two Stripe accounts.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ConnectorsStripeTransferResponse;
import com.formance.formance_sdk.models.shared.Security;
import com.formance.formance_sdk.models.shared.StripeTransferRequest;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("occaecati") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            com.formance.formance_sdk.models.shared.StripeTransferRequest req = new StripeTransferRequest() {{
                amount = 100L;
                asset = "USD";
                destination = "acct_1Gqj58KZcSIg2N2q";
                metadata = new java.util.HashMap<String, Object>() {{
                    put("accusamus", "delectus");
                    put("quidem", "provident");
                }};
            }};            

            ConnectorsStripeTransferResponse res = sdk.payments.connectorsStripeTransfer(req);

            if (res.stripeTransferResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

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
                .setSecurity(new Security("nam") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ConnectorsTransferRequest req = new ConnectorsTransferRequest(                new TransferRequest(100L, "USD", "acct_1Gqj58KZcSIg2N2q") {{
                                source = "acct_1Gqj58KZcSIg2N2q";
                            }};, Connector.MODULR);            

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
                .setSecurity(new Security("blanditiis") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetConnectorTaskRequest req = new GetConnectorTaskRequest(Connector.MODULR, "sapiente");            

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
                .setSecurity(new Security("amet") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetPaymentRequest req = new GetPaymentRequest("deserunt");            

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
import com.formance.formance_sdk.models.shared.ModulrConfig;
import com.formance.formance_sdk.models.shared.Security;
import com.formance.formance_sdk.models.shared.StripeConfig;
import com.formance.formance_sdk.models.shared.WiseConfig;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("nisi") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            InstallConnectorRequest req = new InstallConnectorRequest(                new WiseConfig("XXX") {{
                                apiKey = "XXX";
                            }}, Connector.MODULR);            

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
                .setSecurity(new Security("omnis") {{
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
                .setSecurity(new Security("molestiae") {{
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
                .setSecurity(new Security("perferendis") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListConnectorTasksRequest req = new ListConnectorTasksRequest(Connector.WISE) {{
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                pageSize = 301575L;
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
                .setSecurity(new Security("distinctio") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListConnectorsTransfersRequest req = new ListConnectorsTransfersRequest(Connector.MODULR);            

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
                .setSecurity(new Security("labore") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListPaymentsRequest req = new ListPaymentsRequest() {{
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                pageSize = 290077L;
                sort = new String[]{{
                    add("natus"),
                    add("nobis"),
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
                .setSecurity(new Security("eum") {{
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
                .setSecurity(new Security("vero") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            PaymentslistAccountsRequest req = new PaymentslistAccountsRequest() {{
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                pageSize = 135474L;
                sort = new String[]{{
                    add("magnam"),
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
                .setSecurity(new Security("et") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ReadConnectorConfigRequest req = new ReadConnectorConfigRequest(Connector.MODULR);            

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
                .setSecurity(new Security("ullam") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ResetConnectorRequest req = new ResetConnectorRequest(Connector.MODULR);            

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
                .setSecurity(new Security("quos") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            UninstallConnectorRequest req = new UninstallConnectorRequest(Connector.MODULR);            

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
                .setSecurity(new Security("accusantium") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            UpdateMetadataRequest req = new UpdateMetadataRequest(                new PaymentMetadata() {{
                                key = "mollitia";
                            }};, "reiciendis");            

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
