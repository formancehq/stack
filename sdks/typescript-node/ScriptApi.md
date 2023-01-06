# formance.ScriptApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**runScript**](ScriptApi.md#runScript) | **POST** /api/ledger/{ledger}/script | Execute a Numscript.


# **runScript**
> ScriptResult runScript(script)


### Example


```typescript
import { ScriptApi, createConfiguration } from '@formancehq/formance';
import * as fs from 'fs';

const configuration = createConfiguration();
const apiInstance = new ScriptApi(configuration);

let body:ScriptApiRunScriptRequest = {
  // string | Name of the ledger.
  ledger: "ledger001",
  // Script
  script: {
    reference: "order_1234",
    metadata: {
      "key": null,
    },
    plain: `vars {
account $user
}
send [COIN 10] (
	source = @world
	destination = $user
)
`,
    vars: {},
  },
  // boolean | Set the preview mode. Preview mode doesn't add the logs to the database or publish a message to the message broker. (optional)
  preview: true,
};

apiInstance.runScript(body).then((data:any) => {
  console.log('API called successfully. Returned data: ' + data);
}).catch((error:any) => console.error(error));
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **script** | **Script**|  |
 **ledger** | [**string**] | Name of the ledger. | defaults to undefined
 **preview** | [**boolean**] | Set the preview mode. Preview mode doesn&#39;t add the logs to the database or publish a message to the message broker. | (optional) defaults to undefined


### Return type

**ScriptResult**

### Authorization

[Authorization](README.md#Authorization)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**409** | Conflict |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

