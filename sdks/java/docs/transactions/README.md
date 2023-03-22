# transactions

### Available Operations

* [addMetadataOnTransaction](#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [countTransactions](#counttransactions) - Count the transactions from a ledger
* [createTransaction](#createtransaction) - Create a new transaction to a ledger
* [getTransaction](#gettransaction) - Get transaction from a ledger by its ID
* [listTransactions](#listtransactions) - List transactions from a ledger
* [revertTransaction](#reverttransaction) - Revert a ledger transaction by its ID

## addMetadataOnTransaction

Set the metadata of a transaction by its ID

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.AddMetadataOnTransactionRequest;
import com.formance.formance_sdk.models.operations.AddMetadataOnTransactionResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("delectus") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            AddMetadataOnTransactionRequest req = new AddMetadataOnTransactionRequest("ledger001", 1234L) {{
                idempotencyKey = "quidem";
                requestBody = new java.util.HashMap<String, String>() {{
                    put("nam", "id");
                    put("blanditiis", "deleniti");
                    put("sapiente", "amet");
                }};
                async = true;
                dryRun = true;
            }};            

            AddMetadataOnTransactionResponse res = sdk.transactions.addMetadataOnTransaction(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## countTransactions

Count the transactions from a ledger

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CountTransactionsRequest;
import com.formance.formance_sdk.models.operations.CountTransactionsResponse;
import com.formance.formance_sdk.models.shared.Security;
import java.time.OffsetDateTime;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("deserunt") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            CountTransactionsRequest req = new CountTransactionsRequest("ledger001") {{
                account = "users:001";
                destination = "users:001";
                endTime = OffsetDateTime.parse("2022-07-30T07:01:46.758Z");
                metadata = new java.util.HashMap<String, String>() {{
                    put("omnis", "molestiae");
                    put("perferendis", "nihil");
                    put("magnam", "distinctio");
                }};
                reference = "ref:001";
                source = "users:001";
                startTime = OffsetDateTime.parse("2022-06-04T18:23:50.695Z");
            }};            

            CountTransactionsResponse res = sdk.transactions.countTransactions(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## createTransaction

Create a new transaction to a ledger

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CreateTransactionRequest;
import com.formance.formance_sdk.models.operations.CreateTransactionResponse;
import com.formance.formance_sdk.models.shared.PostTransaction;
import com.formance.formance_sdk.models.shared.PostTransactionScript;
import com.formance.formance_sdk.models.shared.Posting;
import com.formance.formance_sdk.models.shared.Security;
import java.time.OffsetDateTime;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("labore") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            CreateTransactionRequest req = new CreateTransactionRequest(                new PostTransaction(                new java.util.HashMap<String, String>() {{
                                                put("natus", "nobis");
                                                put("eum", "vero");
                                            }}) {{
                                postings = new com.formance.formance_sdk.models.shared.Posting[]{{
                                    add(new Posting(100L, "COIN", "users:002", "users:001") {{
                                        amount = 100L;
                                        asset = "COIN";
                                        destination = "users:002";
                                        source = "users:001";
                                    }}),
                                }};
                                reference = "ref:001";
                                script = new PostTransactionScript("vars {
                                account $user
                                }
                                send [COIN 10] (
                                	source = @world
                                	destination = $user
                                )
                                ") {{
                                    vars = new java.util.HashMap<String, Object>() {{
                                        put("magnam", "et");
                                    }};
                                }};;
                                timestamp = OffsetDateTime.parse("2022-04-17T13:06:08.135Z");
                            }};, "ledger001") {{
                idempotencyKey = "provident";
                async = true;
                dryRun = true;
            }};            

            CreateTransactionResponse res = sdk.transactions.createTransaction(req);

            if (res.createTransactionResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getTransaction

Get transaction from a ledger by its ID

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetTransactionRequest;
import com.formance.formance_sdk.models.operations.GetTransactionResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("quos") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetTransactionRequest req = new GetTransactionRequest("ledger001", 1234L);            

            GetTransactionResponse res = sdk.transactions.getTransaction(req);

            if (res.getTransactionResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listTransactions

List transactions from a ledger, sorted by txid in descending order.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListTransactionsRequest;
import com.formance.formance_sdk.models.operations.ListTransactionsResponse;
import com.formance.formance_sdk.models.shared.Security;
import java.time.OffsetDateTime;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("sint") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListTransactionsRequest req = new ListTransactionsRequest("ledger001") {{
                account = "users:001";
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                destination = "users:001";
                endTime = OffsetDateTime.parse("2022-05-07T13:57:38.738Z");
                metadata = new java.util.HashMap<String, String>() {{
                    put("mollitia", "ad");
                    put("eum", "dolor");
                    put("necessitatibus", "odit");
                    put("nemo", "quasi");
                }};
                pageSize = 435865L;
                reference = "ref:001";
                source = "users:001";
                startTime = OffsetDateTime.parse("2020-04-29T08:15:14.819Z");
            }};            

            ListTransactionsResponse res = sdk.transactions.listTransactions(req);

            if (res.transactionsCursorResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## revertTransaction

Revert a ledger transaction by its ID

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.RevertTransactionRequest;
import com.formance.formance_sdk.models.operations.RevertTransactionResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("eius") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            RevertTransactionRequest req = new RevertTransactionRequest("ledger001", 1234L);            

            RevertTransactionResponse res = sdk.transactions.revertTransaction(req);

            if (res.revertTransactionResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```
