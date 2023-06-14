# orchestration

### Available Operations

* [cancelEvent](#cancelevent) - Cancel a running workflow
* [createWorkflow](#createworkflow) - Create workflow
* [getInstance](#getinstance) - Get a workflow instance by id
* [getInstanceHistory](#getinstancehistory) - Get a workflow instance history by id
* [getInstanceStageHistory](#getinstancestagehistory) - Get a workflow instance stage history
* [getWorkflow](#getworkflow) - Get a flow by id
* [listInstances](#listinstances) - List instances of a workflow
* [listWorkflows](#listworkflows) - List registered workflows
* [orchestrationgetServerInfo](#orchestrationgetserverinfo) - Get server info
* [runWorkflow](#runworkflow) - Run workflow
* [sendEvent](#sendevent) - Send an event to a running workflow

## cancelEvent

Cancel a running workflow

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CancelEventRequest;
import com.formance.formance_sdk.models.operations.CancelEventResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("incidunt") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            CancelEventRequest req = new CancelEventRequest("enim");            

            CancelEventResponse res = sdk.orchestration.cancelEvent(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## createWorkflow

Create a workflow

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CreateWorkflowResponse;
import com.formance.formance_sdk.models.shared.CreateWorkflowRequest;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("consequatur") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            com.formance.formance_sdk.models.shared.CreateWorkflowRequest req = new CreateWorkflowRequest(                new java.util.HashMap<String, Object>[]{{
                                add(new java.util.HashMap<String, Object>() {{
                                    put("explicabo", "deserunt");
                                    put("distinctio", "quibusdam");
                                    put("labore", "modi");
                                    put("qui", "aliquid");
                                }}),
                                add(new java.util.HashMap<String, Object>() {{
                                    put("quos", "perferendis");
                                    put("magni", "assumenda");
                                    put("ipsam", "alias");
                                }}),
                                add(new java.util.HashMap<String, Object>() {{
                                    put("dolorum", "excepturi");
                                }}),
                            }}) {{
                name = "Olivia Rice";
            }};            

            CreateWorkflowResponse res = sdk.orchestration.createWorkflow(req);

            if (res.createWorkflowResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getInstance

Get a workflow instance by id

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetInstanceRequest;
import com.formance.formance_sdk.models.operations.GetInstanceResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("eum") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetInstanceRequest req = new GetInstanceRequest("non");            

            GetInstanceResponse res = sdk.orchestration.getInstance(req);

            if (res.getWorkflowInstanceResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getInstanceHistory

Get a workflow instance history by id

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetInstanceHistoryRequest;
import com.formance.formance_sdk.models.operations.GetInstanceHistoryResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("eligendi") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetInstanceHistoryRequest req = new GetInstanceHistoryRequest("sint");            

            GetInstanceHistoryResponse res = sdk.orchestration.getInstanceHistory(req);

            if (res.getWorkflowInstanceHistoryResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getInstanceStageHistory

Get a workflow instance stage history

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetInstanceStageHistoryRequest;
import com.formance.formance_sdk.models.operations.GetInstanceStageHistoryResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("aliquid") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetInstanceStageHistoryRequest req = new GetInstanceStageHistoryRequest("provident", 896039L);            

            GetInstanceStageHistoryResponse res = sdk.orchestration.getInstanceStageHistory(req);

            if (res.getWorkflowInstanceHistoryStageResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getWorkflow

Get a flow by id

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetWorkflowRequest;
import com.formance.formance_sdk.models.operations.GetWorkflowResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("sint") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetWorkflowRequest req = new GetWorkflowRequest("officia");            

            GetWorkflowResponse res = sdk.orchestration.getWorkflow(req);

            if (res.getWorkflowResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listInstances

List instances of a workflow

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListInstancesRequest;
import com.formance.formance_sdk.models.operations.ListInstancesResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("dolor") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListInstancesRequest req = new ListInstancesRequest() {{
                running = false;
                workflowID = "debitis";
            }};            

            ListInstancesResponse res = sdk.orchestration.listInstances(req);

            if (res.listRunsResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listWorkflows

List registered workflows

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListWorkflowsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("a") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListWorkflowsResponse res = sdk.orchestration.listWorkflows();

            if (res.listWorkflowsResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## orchestrationgetServerInfo

Get server info

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.OrchestrationgetServerInfoResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("dolorum") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            OrchestrationgetServerInfoResponse res = sdk.orchestration.orchestrationgetServerInfo();

            if (res.serverInfo != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## runWorkflow

Run workflow

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.RunWorkflowRequest;
import com.formance.formance_sdk.models.operations.RunWorkflowResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("in") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            RunWorkflowRequest req = new RunWorkflowRequest("in") {{
                requestBody = new java.util.HashMap<String, String>() {{
                    put("maiores", "rerum");
                    put("dicta", "magnam");
                    put("cumque", "facere");
                    put("ea", "aliquid");
                }};
                wait = false;
            }};            

            RunWorkflowResponse res = sdk.orchestration.runWorkflow(req);

            if (res.runWorkflowResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## sendEvent

Send an event to a running workflow

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.SendEventRequest;
import com.formance.formance_sdk.models.operations.SendEventRequestBody;
import com.formance.formance_sdk.models.operations.SendEventResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("laborum") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            SendEventRequest req = new SendEventRequest("accusamus") {{
                requestBody = new SendEventRequestBody("non");;
            }};            

            SendEventResponse res = sdk.orchestration.sendEvent(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```
