# OrchestrationApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**cancelEvent**](OrchestrationApi.md#cancelEvent) | **PUT** api/orchestration/instances/{instanceID}/abort | Cancel a running workflow |
| [**createWorkflow**](OrchestrationApi.md#createWorkflow) | **POST** api/orchestration/workflows | Create workflow |
| [**getInstance**](OrchestrationApi.md#getInstance) | **GET** api/orchestration/instances/{instanceID} | Get a workflow instance by id |
| [**getInstanceHistory**](OrchestrationApi.md#getInstanceHistory) | **GET** api/orchestration/instances/{instanceID}/history | Get a workflow instance history by id |
| [**getInstanceStageHistory**](OrchestrationApi.md#getInstanceStageHistory) | **GET** api/orchestration/instances/{instanceID}/stages/{number}/history | Get a workflow instance stage history |
| [**getWorkflow**](OrchestrationApi.md#getWorkflow) | **GET** api/orchestration/workflows/{flowId} | Get a flow by id |
| [**listInstances**](OrchestrationApi.md#listInstances) | **GET** api/orchestration/instances | List instances of a workflow |
| [**listWorkflows**](OrchestrationApi.md#listWorkflows) | **GET** api/orchestration/workflows | List registered workflows |
| [**orchestrationgetServerInfo**](OrchestrationApi.md#orchestrationgetServerInfo) | **GET** api/orchestration/_info | Get server info |
| [**runWorkflow**](OrchestrationApi.md#runWorkflow) | **POST** api/orchestration/workflows/{workflowID}/instances | Run workflow |
| [**sendEvent**](OrchestrationApi.md#sendEvent) | **POST** api/orchestration/instances/{instanceID}/events | Send an event to a running workflow |



## cancelEvent

> cancelEvent(instanceID)

Cancel a running workflow

Cancel a running workflow

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.OrchestrationApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        OrchestrationApi apiInstance = new OrchestrationApi(defaultClient);
        String instanceID = "xxx"; // String | The instance id
        try {
            apiInstance.cancelEvent(instanceID);
        } catch (ApiException e) {
            System.err.println("Exception when calling OrchestrationApi#cancelEvent");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **instanceID** | **String**| The instance id | |

### Return type

null (empty response body)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | No content |  -  |
| **0** | General error |  -  |


## createWorkflow

> CreateWorkflowResponse createWorkflow(body)

Create workflow

Create a workflow

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.OrchestrationApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        OrchestrationApi apiInstance = new OrchestrationApi(defaultClient);
        WorkflowConfig body = new WorkflowConfig(); // WorkflowConfig | 
        try {
            CreateWorkflowResponse result = apiInstance.createWorkflow(body);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling OrchestrationApi#createWorkflow");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **body** | **WorkflowConfig**|  | [optional] |

### Return type

[**CreateWorkflowResponse**](CreateWorkflowResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created workflow |  -  |
| **0** | General error |  -  |


## getInstance

> GetWorkflowInstanceResponse getInstance(instanceID)

Get a workflow instance by id

Get a workflow instance by id

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.OrchestrationApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        OrchestrationApi apiInstance = new OrchestrationApi(defaultClient);
        String instanceID = "xxx"; // String | The instance id
        try {
            GetWorkflowInstanceResponse result = apiInstance.getInstance(instanceID);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling OrchestrationApi#getInstance");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **instanceID** | **String**| The instance id | |

### Return type

[**GetWorkflowInstanceResponse**](GetWorkflowInstanceResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | The workflow instance |  -  |
| **0** | General error |  -  |


## getInstanceHistory

> GetWorkflowInstanceHistoryResponse getInstanceHistory(instanceID)

Get a workflow instance history by id

Get a workflow instance history by id

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.OrchestrationApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        OrchestrationApi apiInstance = new OrchestrationApi(defaultClient);
        String instanceID = "xxx"; // String | The instance id
        try {
            GetWorkflowInstanceHistoryResponse result = apiInstance.getInstanceHistory(instanceID);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling OrchestrationApi#getInstanceHistory");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **instanceID** | **String**| The instance id | |

### Return type

[**GetWorkflowInstanceHistoryResponse**](GetWorkflowInstanceHistoryResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | The workflow instance history |  -  |
| **0** | General error |  -  |


## getInstanceStageHistory

> GetWorkflowInstanceHistoryStageResponse getInstanceStageHistory(instanceID, number)

Get a workflow instance stage history

Get a workflow instance stage history

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.OrchestrationApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        OrchestrationApi apiInstance = new OrchestrationApi(defaultClient);
        String instanceID = "xxx"; // String | The instance id
        Integer number = 0; // Integer | The stage number
        try {
            GetWorkflowInstanceHistoryStageResponse result = apiInstance.getInstanceStageHistory(instanceID, number);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling OrchestrationApi#getInstanceStageHistory");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **instanceID** | **String**| The instance id | |
| **number** | **Integer**| The stage number | |

### Return type

[**GetWorkflowInstanceHistoryStageResponse**](GetWorkflowInstanceHistoryStageResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | The workflow instance stage history |  -  |
| **0** | General error |  -  |


## getWorkflow

> GetWorkflowResponse getWorkflow(flowId)

Get a flow by id

Get a flow by id

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.OrchestrationApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        OrchestrationApi apiInstance = new OrchestrationApi(defaultClient);
        String flowId = "xxx"; // String | The flow id
        try {
            GetWorkflowResponse result = apiInstance.getWorkflow(flowId);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling OrchestrationApi#getWorkflow");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **flowId** | **String**| The flow id | |

### Return type

[**GetWorkflowResponse**](GetWorkflowResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | The workflow |  -  |
| **0** | General error |  -  |


## listInstances

> ListRunsResponse listInstances(workflowID, running)

List instances of a workflow

List instances of a workflow

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.OrchestrationApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        OrchestrationApi apiInstance = new OrchestrationApi(defaultClient);
        String workflowID = "xxx"; // String | A workflow id
        Boolean running = xxx; // Boolean | Filter running instances
        try {
            ListRunsResponse result = apiInstance.listInstances(workflowID, running);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling OrchestrationApi#listInstances");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **workflowID** | **String**| A workflow id | [optional] |
| **running** | **Boolean**| Filter running instances | [optional] |

### Return type

[**ListRunsResponse**](ListRunsResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of workflow instances |  -  |
| **0** | General error |  -  |


## listWorkflows

> ListWorkflowsResponse listWorkflows()

List registered workflows

List registered workflows

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.OrchestrationApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        OrchestrationApi apiInstance = new OrchestrationApi(defaultClient);
        try {
            ListWorkflowsResponse result = apiInstance.listWorkflows();
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling OrchestrationApi#listWorkflows");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**ListWorkflowsResponse**](ListWorkflowsResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of workflows |  -  |
| **0** | General error |  -  |


## orchestrationgetServerInfo

> ServerInfo orchestrationgetServerInfo()

Get server info

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.OrchestrationApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        OrchestrationApi apiInstance = new OrchestrationApi(defaultClient);
        try {
            ServerInfo result = apiInstance.orchestrationgetServerInfo();
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling OrchestrationApi#orchestrationgetServerInfo");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**ServerInfo**](ServerInfo.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Server information |  -  |
| **0** | General error |  -  |


## runWorkflow

> RunWorkflowResponse runWorkflow(workflowID, wait, requestBody)

Run workflow

Run workflow

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.OrchestrationApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        OrchestrationApi apiInstance = new OrchestrationApi(defaultClient);
        String workflowID = "xxx"; // String | The flow id
        Boolean wait = true; // Boolean | Wait end of the workflow before return
        Map<String, String> requestBody = new HashMap(); // Map<String, String> | 
        try {
            RunWorkflowResponse result = apiInstance.runWorkflow(workflowID, wait, requestBody);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling OrchestrationApi#runWorkflow");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **workflowID** | **String**| The flow id | |
| **wait** | **Boolean**| Wait end of the workflow before return | [optional] |
| **requestBody** | [**Map&lt;String, String&gt;**](String.md)|  | [optional] |

### Return type

[**RunWorkflowResponse**](RunWorkflowResponse.md)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | The workflow instance |  -  |
| **0** | General error |  -  |


## sendEvent

> sendEvent(instanceID, sendEventRequest)

Send an event to a running workflow

Send an event to a running workflow

### Example

```java
// Import classes:
import com.formance.formance.ApiClient;
import com.formance.formance.ApiException;
import com.formance.formance.Configuration;
import com.formance.formance.auth.*;
import com.formance.formance.models.*;
import com.formance.formance.api.OrchestrationApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");
        
        // Configure OAuth2 access token for authorization: Authorization
        OAuth Authorization = (OAuth) defaultClient.getAuthentication("Authorization");
        Authorization.setAccessToken("YOUR ACCESS TOKEN");

        OrchestrationApi apiInstance = new OrchestrationApi(defaultClient);
        String instanceID = "xxx"; // String | The instance id
        SendEventRequest sendEventRequest = new SendEventRequest(); // SendEventRequest | 
        try {
            apiInstance.sendEvent(instanceID, sendEventRequest);
        } catch (ApiException e) {
            System.err.println("Exception when calling OrchestrationApi#sendEvent");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **instanceID** | **String**| The instance id | |
| **sendEventRequest** | [**SendEventRequest**](SendEventRequest.md)|  | [optional] |

### Return type

null (empty response body)

### Authorization

[Authorization](../README.md#Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | No content |  -  |
| **0** | General error |  -  |

