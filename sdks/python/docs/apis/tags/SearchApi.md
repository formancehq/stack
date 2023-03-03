<a name="__pageTop"></a>
# FormanceHQ.apis.tags.search_api.SearchApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**search**](#search) | **post** /api/search/ | Search

# **search**
<a name="search"></a>
> Response search(query)

Search

ElasticSearch query engine

### Example

* OAuth Authentication (Authorization):
```python
import FormanceHQ
from FormanceHQ.apis.tags import search_api
from FormanceHQ.model.response import Response
from FormanceHQ.model.query import Query
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = FormanceHQ.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure OAuth2 access token for authorization: Authorization
configuration = FormanceHQ.Configuration(
    host = "http://localhost",
    access_token = 'YOUR_ACCESS_TOKEN'
)
# Enter a context with an instance of the API client
with FormanceHQ.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = search_api.SearchApi(api_client)

    # example passing only required values which don't have defaults set
    body = Query(
        ledgers=[
            "quickstart"
        ],
        after=[
            "users:002"
        ],
        page_size=0,
        terms=[
            "destination=central_bank1"
        ],
        sort="txid:asc",
        policy="OR",
        target="target_example",
        cursor="YXVsdCBhbmQgYSBtYXhpbXVtIG1heF9yZXN1bHRzLol=",
        raw=dict(),
    )
    try:
        # Search
        api_response = api_instance.search(
            body=body,
        )
        pprint(api_response)
    except FormanceHQ.ApiException as e:
        print("Exception when calling SearchApi->search: %s\n" % e)
```
### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
body | typing.Union[SchemaForRequestBodyApplicationJson] | required |
content_type | str | optional, default is 'application/json' | Selects the schema and serialization of the request body
accept_content_types | typing.Tuple[str] | default is ('application/json', ) | Tells the server the content type(s) that are accepted by the client
stream | bool | default is False | if True then the response.content will be streamed and loaded from a file like object. When downloading a file, set this to True to force the code to deserialize the content to a FileSchema file
timeout | typing.Optional[typing.Union[int, typing.Tuple]] | default is None | the timeout used by the rest client
skip_deserialization | bool | default is False | when True, headers and body will be unset and an instance of api_client.ApiResponseWithoutDeserialization will be returned

### body

# SchemaForRequestBodyApplicationJson
Type | Description  | Notes
------------- | ------------- | -------------
[**Query**](../../models/Query.md) |  | 


### Return Types, Responses

Code | Class | Description
------------- | ------------- | -------------
n/a | api_client.ApiResponseWithoutDeserialization | When skip_deserialization is True this response is returned
200 | [ApiResponseFor200](#search.ApiResponseFor200) | Success
default | [ApiResponseForDefault](#search.ApiResponseForDefault) | Error

#### search.ApiResponseFor200
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | typing.Union[SchemaFor200ResponseBodyApplicationJson, ] |  |
headers | Unset | headers were not defined |

# SchemaFor200ResponseBodyApplicationJson
Type | Description  | Notes
------------- | ------------- | -------------
[**Response**](../../models/Response.md) |  | 


#### search.ApiResponseForDefault
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | typing.Union[] |  |
headers | Unset | headers were not defined |

### Authorization

[Authorization](../../../README.md#Authorization)

[[Back to top]](#__pageTop) [[Back to API list]](../../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../../README.md#documentation-for-models) [[Back to README]](../../../README.md)

