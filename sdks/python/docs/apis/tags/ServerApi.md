<a name="__pageTop"></a>
# Formance.apis.tags.server_api.ServerApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**get_info**](#get_info) | **get** /api/ledger/_info | Show server information

# **get_info**
<a name="get_info"></a>
> ConfigInfo get_info()

Show server information

### Example

```python
import Formance
from Formance.apis.tags import server_api
from Formance.model.error_response import ErrorResponse
from Formance.model.config_info import ConfigInfo
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = Formance.Configuration(
    host = "http://localhost"
)

# Enter a context with an instance of the API client
with Formance.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = server_api.ServerApi(api_client)

    # example, this endpoint has no required or optional parameters
    try:
        # Show server information
        api_response = api_instance.get_info()
        pprint(api_response)
    except Formance.ApiException as e:
        print("Exception when calling ServerApi->get_info: %s\n" % e)
```
### Parameters
This endpoint does not need any parameter.

### Return Types, Responses

Code | Class | Description
------------- | ------------- | -------------
n/a | api_client.ApiResponseWithoutDeserialization | When skip_deserialization is True this response is returned
200 | [ApiResponseFor200](#get_info.ApiResponseFor200) | OK
default | [ApiResponseForDefault](#get_info.ApiResponseForDefault) | Error

#### get_info.ApiResponseFor200
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | typing.Union[SchemaFor200ResponseBodyApplicationJson, ] |  |
headers | Unset | headers were not defined |

# SchemaFor200ResponseBodyApplicationJson
Type | Description  | Notes
------------- | ------------- | -------------
[**ConfigInfo**](../../models/ConfigInfo.md) |  | 


#### get_info.ApiResponseForDefault
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | typing.Union[SchemaFor0ResponseBodyApplicationJson, ] |  |
headers | Unset | headers were not defined |

# SchemaFor0ResponseBodyApplicationJson
Type | Description  | Notes
------------- | ------------- | -------------
[**ErrorResponse**](../../models/ErrorResponse.md) |  | 


### Authorization

No authorization required

[[Back to top]](#__pageTop) [[Back to API list]](../../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../../README.md#documentation-for-models) [[Back to README]](../../../README.md)

