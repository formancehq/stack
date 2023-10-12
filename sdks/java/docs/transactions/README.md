# transactions

### Available Operations

* [createTransactions](#createtransactions) - Create a new batch of transactions to a ledger
* [addMetadataOnTransaction](#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [countTransactions](#counttransactions) - Count the transactions from a ledger
* [createTransaction](#createtransaction) - Create a new transaction to a ledger
* [getTransaction](#gettransaction) - Get transaction from a ledger by its ID
* [listTransactions](#listtransactions) - List transactions from a ledger
* [revertTransaction](#reverttransaction) - Revert a ledger transaction by its ID

## createTransactions

Create a new batch of transactions to a ledger

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CreateTransactionsRequest;
import com.formance.formance_sdk.models.operations.CreateTransactionsResponse;
import com.formance.formance_sdk.models.shared.Posting;
import com.formance.formance_sdk.models.shared.Security;
import com.formance.formance_sdk.models.shared.TransactionData;
import com.formance.formance_sdk.models.shared.Transactions;
import java.time.OffsetDateTime;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("rerum") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            CreateTransactionsRequest req = new CreateTransactionsRequest(                new Transactions(                new com.formance.formance_sdk.models.shared.TransactionData[]{{
                                                add(new TransactionData(                new com.formance.formance_sdk.models.shared.Posting[]{{
                                                                    add(new Posting(100L, "COIN", "users:002", "users:001") {{
                                                                        amount = 100L;
                                                                        asset = "COIN";
                                                                        destination = "users:002";
                                                                        source = "users:001";
                                                                    }}),
                                                                    add(new Posting(100L, "COIN", "users:002", "users:001") {{
                                                                        amount = 100L;
                                                                        asset = "COIN";
                                                                        destination = "users:002";
                                                                        source = "users:001";
                                                                    }}),
                                                                    add(new Posting(100L, "COIN", "users:002", "users:001") {{
                                                                        amount = 100L;
                                                                        asset = "COIN";
                                                                        destination = "users:002";
                                                                        source = "users:001";
                                                                    }}),
                                                                }}) {{
                                                    metadata = new java.util.HashMap<String, Object>() {{
                                                        put("earum", "modi");
                                                        put("iste", "dolorum");
                                                        put("deleniti", "pariatur");
                                                        put("provident", "nobis");
                                                    }};
                                                    postings = new com.formance.formance_sdk.models.shared.Posting[]{{
                                                        add(new Posting(100L, "COIN", "users:002", "users:001") {{
                                                            amount = 100L;
                                                            asset = "COIN";
                                                            destination = "users:002";
                                                            source = "users:001";
                                                        }}),
                                                        add(new Posting(100L, "COIN", "users:002", "users:001") {{
                                                            amount = 100L;
                                                            asset = "COIN";
                                                            destination = "users:002";
                                                            source = "users:001";
                                                        }}),
                                                        add(new Posting(100L, "COIN", "users:002", "users:001") {{
                                                            amount = 100L;
                                                            asset = "COIN";
                                                            destination = "users:002";
                                                            source = "users:001";
                                                        }}),
                                                    }};
                                                    reference = "ref:001";
                                                    timestamp = OffsetDateTime.parse("2022-01-24T10:05:07.174Z");
                                                }}),
                                            }});, "ledger001");            

            CreateTransactionsResponse res = sdk.transactions.createTransactions(req);

            if (res.transactionsResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

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
                .setSecurity(new Security("aliquid") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            AddMetadataOnTransactionRequest req = new AddMetadataOnTransactionRequest("ledger001", 1234L) {{
                requestBody = new java.util.HashMap<String, Object>() {{
                    put("dolorem", "dolor");
                }};
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
                .setSecurity(new Security("qui") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            CountTransactionsRequest req = new CountTransactionsRequest("ledger001") {{
                account = "users:001";
                destination = "users:001";
                endTime = OffsetDateTime.parse("2022-01-21T07:17:52.299Z");
                metadata = new java.util.HashMap<String, Object>() {{
                    put("cum", "voluptate");
                    put("dignissimos", "reiciendis");
                    put("amet", "dolorum");
                }};
                reference = "ref:001";
                source = "users:001";
                startTime = OffsetDateTime.parse("2022-11-30T20:48:38.083Z");
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
                .setSecurity(new Security("ipsa") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            CreateTransactionRequest req = new CreateTransactionRequest(                new PostTransaction() {{
                                metadata = new java.util.HashMap<String, Object>() {{
                                    put("iure", "odio");
                                }};
                                postings = new com.formance.formance_sdk.models.shared.Posting[]{{
                                    add(new Posting(100L, "COIN", "users:002", "users:001") {{
                                        amount = 100L;
                                        asset = "COIN";
                                        destination = "users:002";
                                        source = "users:001";
                                    }}),
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
                                        put("quidem", "voluptatibus");
                                        put("voluptas", "natus");
                                        put("eos", "atque");
                                        put("sit", "fugiat");
                                    }};
                                }};;
                                timestamp = OffsetDateTime.parse("2022-04-04T12:00:33.616Z");
                            }};, "ledger001") {{
                idempotencyKey = "dolorum";
                preview = true;
            }};            

            CreateTransactionResponse res = sdk.transactions.createTransaction(req);

            if (res.transactionsResponse != null) {
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
                .setSecurity(new Security("iusto") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetTransactionRequest req = new GetTransactionRequest("ledger001", 1234L);            

            GetTransactionResponse res = sdk.transactions.getTransaction(req);

            if (res.transactionResponse != null) {
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
                .setSecurity(new Security("voluptate") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListTransactionsRequest req = new ListTransactionsRequest("ledger001") {{
                account = "users:001";
                after = "dolorum";
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                destination = "users:001";
                endTime = OffsetDateTime.parse("2021-10-14T20:34:13.522Z");
                metadata = new java.util.HashMap<String, Object>() {{
                    put("distinctio", "asperiores");
                    put("nihil", "ipsum");
                    put("voluptate", "id");
                    put("saepe", "eius");
                }};
                pageSize = 137220L;
                reference = "ref:001";
                source = "users:001";
                startTime = OffsetDateTime.parse("2022-10-09T08:02:18.659Z");
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
                .setSecurity(new Security("optio") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            RevertTransactionRequest req = new RevertTransactionRequest("ledger001", 1234L);            

            RevertTransactionResponse res = sdk.transactions.revertTransaction(req);

            if (res.transactionResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```
