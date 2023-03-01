# formance.OrchestrationApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**cancelEvent**](OrchestrationApi.md#cancelEvent) | **PUT** /api/orchestration/instances/{instanceID}/abort | Cancel a running workflow
[**createWorkflow**](OrchestrationApi.md#createWorkflow) | **POST** /api/orchestration/workflows | Create workflow
[**getInstance**](OrchestrationApi.md#getInstance) | **GET** /api/orchestration/instances/{instanceID} | Get a workflow instance by id
[**getInstanceHistory**](OrchestrationApi.md#getInstanceHistory) | **GET** /api/orchestration/instances/{instanceID}/history | Get a workflow instance history by id
[**getInstanceStageHistory**](OrchestrationApi.md#getInstanceStageHistory) | **GET** /api/orchestration/instances/{instanceID}/stages/{number}/history | Get a workflow instance stage history
[**getWorkflow**](OrchestrationApi.md#getWorkflow) | **GET** /api/orchestration/workflows/{flowId} | Get a flow by id
[**listInstances**](OrchestrationApi.md#listInstances) | **GET** /api/orchestration/instances | List instances of a workflow
[**listWorkflows**](OrchestrationApi.md#listWorkflows) | **GET** /api/orchestration/workflows | List registered workflows
[**orchestrationgetServerInfo**](OrchestrationApi.md#orchestrationgetServerInfo) | **GET** /api/orchestration/_info | Get server info
[**runWorkflow**](OrchestrationApi.md#runWorkflow) | **POST** /api/orchestration/workflows/{workflowID}/instances | Run workflow
[**sendEvent**](OrchestrationApi.md#sendEvent) | **POST** /api/orchestration/instances/{instanceID}/events | Send an event to a running workflow


# **cancelEvent**
> void cancelEvent()

Cancel a running workflow

### Example


```typescript
import { formance } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = formance.createConfiguration();
const apiInstance = new formance.OrchestrationApi(configuration);

let body:formance.OrchestrationApiCancelEventRequest = {
  // string | The instance id
  instanceID: "xxx",
};

apiInstance.cancelEvent(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **instanceID** | [**string**] | The instance id | defaults to undefined


### Return type

**void**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**204** | No content |  -  |
**0** | General error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **createWorkflow**
> CreateWorkflowResponse createWorkflow()

Create a workflow

### Example


```typescript
import { formance } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = formance.createConfiguration();
const apiInstance = new formance.OrchestrationApi(configuration);

let body:formance.OrchestrationApiCreateWorkflowRequest = {
  // WorkflowConfig (optional)
  body: {
    name: "name_example",
    stages: [
      {
        "key": null,
      },
    ],
  },
};

apiInstance.createWorkflow(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | **WorkflowConfig**|  |


### Return type

**CreateWorkflowResponse**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | Created workflow |  -  |
**0** | General error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **getInstance**
> GetWorkflowInstanceResponse getInstance()

Get a workflow instance by id

### Example


```typescript
import { formance } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = formance.createConfiguration();
const apiInstance = new formance.OrchestrationApi(configuration);

let body:formance.OrchestrationApiGetInstanceRequest = {
  // string | The instance id
  instanceID: "xxx",
};

apiInstance.getInstance(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **instanceID** | [**string**] | The instance id | defaults to undefined


### Return type

**GetWorkflowInstanceResponse**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | The workflow instance |  -  |
**0** | General error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **getInstanceHistory**
> GetWorkflowInstanceHistoryResponse getInstanceHistory()

Get a workflow instance history by id

### Example


```typescript
import { formance } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = formance.createConfiguration();
const apiInstance = new formance.OrchestrationApi(configuration);

let body:formance.OrchestrationApiGetInstanceHistoryRequest = {
  // string | The instance id
  instanceID: "xxx",
};

apiInstance.getInstanceHistory(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **instanceID** | [**string**] | The instance id | defaults to undefined


### Return type

**GetWorkflowInstanceHistoryResponse**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | The workflow instance history |  -  |
**0** | General error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **getInstanceStageHistory**
> GetWorkflowInstanceHistoryStageResponse getInstanceStageHistory()

Get a workflow instance stage history

### Example


```typescript
import { formance } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = formance.createConfiguration();
const apiInstance = new formance.OrchestrationApi(configuration);

let body:formance.OrchestrationApiGetInstanceStageHistoryRequest = {
  // string | The instance id
  instanceID: "xxx",
  // number | The stage number
  number: 0,
};

apiInstance.getInstanceStageHistory(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **instanceID** | [**string**] | The instance id | defaults to undefined
 **number** | [**number**] | The stage number | defaults to undefined


### Return type

**GetWorkflowInstanceHistoryStageResponse**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | The workflow instance stage history |  -  |
**0** | General error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **getWorkflow**
> GetWorkflowResponse getWorkflow()

Get a flow by id

### Example


```typescript
import { formance } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = formance.createConfiguration();
const apiInstance = new formance.OrchestrationApi(configuration);

let body:formance.OrchestrationApiGetWorkflowRequest = {
  // string | The flow id
  flowId: "xxx",
};

apiInstance.getWorkflow(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **flowId** | [**string**] | The flow id | defaults to undefined


### Return type

**GetWorkflowResponse**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | The workflow |  -  |
**0** | General error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **listInstances**
> ListRunsResponse listInstances()

List instances of a workflow

### Example


```typescript
import { formance } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = formance.createConfiguration();
const apiInstance = new formance.OrchestrationApi(configuration);

let body:formance.OrchestrationApiListInstancesRequest = {
  // string | A workflow id (optional)
  workflowID: "xxx",
  // boolean | Filter running instances (optional)
  running: true,
};

apiInstance.listInstances(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflowID** | [**string**] | A workflow id | (optional) defaults to undefined
 **running** | [**boolean**] | Filter running instances | (optional) defaults to undefined


### Return type

**ListRunsResponse**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | List of workflow instances |  -  |
**0** | General error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **listWorkflows**
> ListWorkflowsResponse listWorkflows()

List registered workflows

### Example


```typescript
import { formance } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = formance.createConfiguration();
const apiInstance = new formance.OrchestrationApi(configuration);

let body:any = {};

apiInstance.listWorkflows(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters
This endpoint does not need any parameter.


### Return type

**ListWorkflowsResponse**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | List of workflows |  -  |
**0** | General error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **orchestrationgetServerInfo**
> ServerInfo orchestrationgetServerInfo()


### Example


```typescript
import { formance } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = formance.createConfiguration();
const apiInstance = new formance.OrchestrationApi(configuration);

let body:any = {};

apiInstance.orchestrationgetServerInfo(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters
This endpoint does not need any parameter.


### Return type

**ServerInfo**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Server information |  -  |
**0** | General error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **runWorkflow**
> RunWorkflowResponse runWorkflow()

Run workflow

### Example


```typescript
import { formance } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = formance.createConfiguration();
const apiInstance = new formance.OrchestrationApi(configuration);

let body:formance.OrchestrationApiRunWorkflowRequest = {
  // string | The flow id
  workflowID: "xxx",
  // boolean | Wait end of the workflow before return (optional)
  wait: true,
  // { [key: string]: string; } (optional)
  requestBody: {
    "key": "key_example",
  },
};

apiInstance.runWorkflow(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestBody** | **{ [key: string]: string; }**|  |
 **workflowID** | [**string**] | The flow id | defaults to undefined
 **wait** | [**boolean**] | Wait end of the workflow before return | (optional) defaults to undefined


### Return type

**RunWorkflowResponse**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | The workflow instance |  -  |
**0** | General error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **sendEvent**
> void sendEvent()

Send an event to a running workflow

### Example


```typescript
import { formance } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = formance.createConfiguration();
const apiInstance = new formance.OrchestrationApi(configuration);

let body:formance.OrchestrationApiSendEventRequest = {
  // string | The instance id
  instanceID: "xxx",
  // SendEventRequest (optional)
  sendEventRequest: {
    name: "name_example",
  },
};

apiInstance.sendEvent(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **sendEventRequest** | **SendEventRequest**|  |
 **instanceID** | [**string**] | The instance id | defaults to undefined


### Return type

**void**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**204** | No content |  -  |
**0** | General error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)


