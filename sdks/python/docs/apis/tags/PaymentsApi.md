<a name="__pageTop"></a>
# Formance.apis.tags.payments_api.PaymentsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**connectors_stripe_transfer**](#connectors_stripe_transfer) | **post** /api/payments/connectors/stripe/transfer | Transfer funds between Stripe accounts
[**get_all_connectors**](#get_all_connectors) | **get** /api/payments/connectors | Get all installed connectors
[**get_all_connectors_configs**](#get_all_connectors_configs) | **get** /api/payments/connectors/configs | Get all available connectors configs
[**get_connector_task**](#get_connector_task) | **get** /api/payments/connectors/{connector}/tasks/{taskId} | Read a specific task of the connector
[**get_payment**](#get_payment) | **get** /api/payments/payments/{paymentId} | Returns a payment.
[**install_connector**](#install_connector) | **post** /api/payments/connectors/{connector} | Install connector
[**list_connector_tasks**](#list_connector_tasks) | **get** /api/payments/connectors/{connector}/tasks | List connector tasks
[**list_payments**](#list_payments) | **get** /api/payments/payments | Returns a list of payments.
[**paymentslist_accounts**](#paymentslist_accounts) | **get** /api/payments/accounts | Returns a list of accounts.
[**read_connector_config**](#read_connector_config) | **get** /api/payments/connectors/{connector}/config | Read connector config
[**reset_connector**](#reset_connector) | **post** /api/payments/connectors/{connector}/reset | Reset connector
[**uninstall_connector**](#uninstall_connector) | **delete** /api/payments/connectors/{connector} | Uninstall connector

# **connectors_stripe_transfer**
<a name="connectors_stripe_transfer"></a>
> connectors_stripe_transfer(stripe_transfer_request)

Transfer funds between Stripe accounts

Execute a transfer between two Stripe accounts

### Example

* OAuth Authentication (Authorization):
```python
import Formance
from Formance.apis.tags import payments_api
from Formance.model.stripe_transfer_request import StripeTransferRequest
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = Formance.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure OAuth2 access token for authorization: Authorization
configuration = Formance.Configuration(
    host = "http://localhost"
)
configuration.access_token = 'YOUR_ACCESS_TOKEN'
# Enter a context with an instance of the API client
with Formance.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = payments_api.PaymentsApi(api_client)

    # example passing only required values which don't have defaults set
    body = StripeTransferRequest(
        amount=100,
        asset="USD",
        destination="acct_1Gqj58KZcSIg2N2q",
        metadata=dict(),
    )
    try:
        # Transfer funds between Stripe accounts
        api_response = api_instance.connectors_stripe_transfer(
            body=body,
        )
    except Formance.ApiException as e:
        print("Exception when calling PaymentsApi->connectors_stripe_transfer: %s\n" % e)
```
### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
body | typing.Union[SchemaForRequestBodyApplicationJson] | required |
content_type | str | optional, default is 'application/json' | Selects the schema and serialization of the request body
stream | bool | default is False | if True then the response.content will be streamed and loaded from a file like object. When downloading a file, set this to True to force the code to deserialize the content to a FileSchema file
timeout | typing.Optional[typing.Union[int, typing.Tuple]] | default is None | the timeout used by the rest client
skip_deserialization | bool | default is False | when True, headers and body will be unset and an instance of api_client.ApiResponseWithoutDeserialization will be returned

### body

# SchemaForRequestBodyApplicationJson
Type | Description  | Notes
------------- | ------------- | -------------
[**StripeTransferRequest**](../../models/StripeTransferRequest.md) |  | 


### Return Types, Responses

Code | Class | Description
------------- | ------------- | -------------
n/a | api_client.ApiResponseWithoutDeserialization | When skip_deserialization is True this response is returned
200 | [ApiResponseFor200](#connectors_stripe_transfer.ApiResponseFor200) | Transfer has been executed

#### connectors_stripe_transfer.ApiResponseFor200
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | Unset | body was not defined |
headers | Unset | headers were not defined |

### Authorization

[Authorization](../../../README.md#Authorization)

[[Back to top]](#__pageTop) [[Back to API list]](../../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../../README.md#documentation-for-models) [[Back to README]](../../../README.md)

# **get_all_connectors**
<a name="get_all_connectors"></a>
> ListConnectorsResponse get_all_connectors()

Get all installed connectors

Get all installed connectors

### Example

* OAuth Authentication (Authorization):
```python
import Formance
from Formance.apis.tags import payments_api
from Formance.model.list_connectors_response import ListConnectorsResponse
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = Formance.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure OAuth2 access token for authorization: Authorization
configuration = Formance.Configuration(
    host = "http://localhost"
)
configuration.access_token = 'YOUR_ACCESS_TOKEN'
# Enter a context with an instance of the API client
with Formance.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = payments_api.PaymentsApi(api_client)

    # example, this endpoint has no required or optional parameters
    try:
        # Get all installed connectors
        api_response = api_instance.get_all_connectors()
        pprint(api_response)
    except Formance.ApiException as e:
        print("Exception when calling PaymentsApi->get_all_connectors: %s\n" % e)
```
### Parameters
This endpoint does not need any parameter.

### Return Types, Responses

Code | Class | Description
------------- | ------------- | -------------
n/a | api_client.ApiResponseWithoutDeserialization | When skip_deserialization is True this response is returned
200 | [ApiResponseFor200](#get_all_connectors.ApiResponseFor200) | List of installed connectors

#### get_all_connectors.ApiResponseFor200
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | typing.Union[SchemaFor200ResponseBodyApplicationJson, ] |  |
headers | Unset | headers were not defined |

# SchemaFor200ResponseBodyApplicationJson
Type | Description  | Notes
------------- | ------------- | -------------
[**ListConnectorsResponse**](../../models/ListConnectorsResponse.md) |  | 


### Authorization

[Authorization](../../../README.md#Authorization)

[[Back to top]](#__pageTop) [[Back to API list]](../../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../../README.md#documentation-for-models) [[Back to README]](../../../README.md)

# **get_all_connectors_configs**
<a name="get_all_connectors_configs"></a>
> ListConnectorsConfigsResponse get_all_connectors_configs()

Get all available connectors configs

Get all available connectors configs

### Example

* OAuth Authentication (Authorization):
```python
import Formance
from Formance.apis.tags import payments_api
from Formance.model.list_connectors_configs_response import ListConnectorsConfigsResponse
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = Formance.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure OAuth2 access token for authorization: Authorization
configuration = Formance.Configuration(
    host = "http://localhost"
)
configuration.access_token = 'YOUR_ACCESS_TOKEN'
# Enter a context with an instance of the API client
with Formance.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = payments_api.PaymentsApi(api_client)

    # example, this endpoint has no required or optional parameters
    try:
        # Get all available connectors configs
        api_response = api_instance.get_all_connectors_configs()
        pprint(api_response)
    except Formance.ApiException as e:
        print("Exception when calling PaymentsApi->get_all_connectors_configs: %s\n" % e)
```
### Parameters
This endpoint does not need any parameter.

### Return Types, Responses

Code | Class | Description
------------- | ------------- | -------------
n/a | api_client.ApiResponseWithoutDeserialization | When skip_deserialization is True this response is returned
200 | [ApiResponseFor200](#get_all_connectors_configs.ApiResponseFor200) | List of available connectors configs

#### get_all_connectors_configs.ApiResponseFor200
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | typing.Union[SchemaFor200ResponseBodyApplicationJson, ] |  |
headers | Unset | headers were not defined |

# SchemaFor200ResponseBodyApplicationJson
Type | Description  | Notes
------------- | ------------- | -------------
[**ListConnectorsConfigsResponse**](../../models/ListConnectorsConfigsResponse.md) |  | 


### Authorization

[Authorization](../../../README.md#Authorization)

[[Back to top]](#__pageTop) [[Back to API list]](../../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../../README.md#documentation-for-models) [[Back to README]](../../../README.md)

# **get_connector_task**
<a name="get_connector_task"></a>
> bool, date, datetime, dict, float, int, list, str, none_type get_connector_task(connectortask_id)

Read a specific task of the connector

Get a specific task associated to the connector

### Example

* OAuth Authentication (Authorization):
```python
import Formance
from Formance.apis.tags import payments_api
from Formance.model.task_descriptor_dummy_pay import TaskDescriptorDummyPay
from Formance.model.task_descriptor_wise import TaskDescriptorWise
from Formance.model.task_descriptor_modulr import TaskDescriptorModulr
from Formance.model.task_descriptor_stripe import TaskDescriptorStripe
from Formance.model.connectors import Connectors
from Formance.model.task_descriptor_banking_circle import TaskDescriptorBankingCircle
from Formance.model.task_descriptor_currency_cloud import TaskDescriptorCurrencyCloud
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = Formance.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure OAuth2 access token for authorization: Authorization
configuration = Formance.Configuration(
    host = "http://localhost"
)
configuration.access_token = 'YOUR_ACCESS_TOKEN'
# Enter a context with an instance of the API client
with Formance.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = payments_api.PaymentsApi(api_client)

    # example passing only required values which don't have defaults set
    path_params = {
        'connector': Connectors("STRIPE"),
        'taskId': "task1",
    }
    try:
        # Read a specific task of the connector
        api_response = api_instance.get_connector_task(
            path_params=path_params,
        )
        pprint(api_response)
    except Formance.ApiException as e:
        print("Exception when calling PaymentsApi->get_connector_task: %s\n" % e)
```
### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
path_params | RequestPathParams | |
accept_content_types | typing.Tuple[str] | default is ('application/json', ) | Tells the server the content type(s) that are accepted by the client
stream | bool | default is False | if True then the response.content will be streamed and loaded from a file like object. When downloading a file, set this to True to force the code to deserialize the content to a FileSchema file
timeout | typing.Optional[typing.Union[int, typing.Tuple]] | default is None | the timeout used by the rest client
skip_deserialization | bool | default is False | when True, headers and body will be unset and an instance of api_client.ApiResponseWithoutDeserialization will be returned

### path_params
#### RequestPathParams

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
connector | ConnectorSchema | | 
taskId | TaskIdSchema | | 

# ConnectorSchema
Type | Description  | Notes
------------- | ------------- | -------------
[**Connectors**](../../models/Connectors.md) |  | 


# TaskIdSchema

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
str,  | str,  |  | 

### Return Types, Responses

Code | Class | Description
------------- | ------------- | -------------
n/a | api_client.ApiResponseWithoutDeserialization | When skip_deserialization is True this response is returned
200 | [ApiResponseFor200](#get_connector_task.ApiResponseFor200) | The specified task

#### get_connector_task.ApiResponseFor200
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | typing.Union[SchemaFor200ResponseBodyApplicationJson, ] |  |
headers | Unset | headers were not defined |

# SchemaFor200ResponseBodyApplicationJson

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
dict, frozendict.frozendict, str, date, datetime, uuid.UUID, int, float, decimal.Decimal, bool, None, list, tuple, bytes, io.FileIO, io.BufferedReader,  | frozendict.frozendict, str, decimal.Decimal, BoolClass, NoneClass, tuple, bytes, FileIO |  | 

### Composed Schemas (allOf/anyOf/oneOf/not)
#### oneOf
Class Name | Input Type | Accessed Type | Description | Notes
------------- | ------------- | ------------- | ------------- | -------------
[TaskDescriptorStripe]({{complexTypePrefix}}TaskDescriptorStripe.md) | [**TaskDescriptorStripe**]({{complexTypePrefix}}TaskDescriptorStripe.md) | [**TaskDescriptorStripe**]({{complexTypePrefix}}TaskDescriptorStripe.md) |  | 
[TaskDescriptorWise]({{complexTypePrefix}}TaskDescriptorWise.md) | [**TaskDescriptorWise**]({{complexTypePrefix}}TaskDescriptorWise.md) | [**TaskDescriptorWise**]({{complexTypePrefix}}TaskDescriptorWise.md) |  | 
[TaskDescriptorCurrencyCloud]({{complexTypePrefix}}TaskDescriptorCurrencyCloud.md) | [**TaskDescriptorCurrencyCloud**]({{complexTypePrefix}}TaskDescriptorCurrencyCloud.md) | [**TaskDescriptorCurrencyCloud**]({{complexTypePrefix}}TaskDescriptorCurrencyCloud.md) |  | 
[TaskDescriptorDummyPay]({{complexTypePrefix}}TaskDescriptorDummyPay.md) | [**TaskDescriptorDummyPay**]({{complexTypePrefix}}TaskDescriptorDummyPay.md) | [**TaskDescriptorDummyPay**]({{complexTypePrefix}}TaskDescriptorDummyPay.md) |  | 
[TaskDescriptorModulr]({{complexTypePrefix}}TaskDescriptorModulr.md) | [**TaskDescriptorModulr**]({{complexTypePrefix}}TaskDescriptorModulr.md) | [**TaskDescriptorModulr**]({{complexTypePrefix}}TaskDescriptorModulr.md) |  | 
[TaskDescriptorBankingCircle]({{complexTypePrefix}}TaskDescriptorBankingCircle.md) | [**TaskDescriptorBankingCircle**]({{complexTypePrefix}}TaskDescriptorBankingCircle.md) | [**TaskDescriptorBankingCircle**]({{complexTypePrefix}}TaskDescriptorBankingCircle.md) |  | 

### Authorization

[Authorization](../../../README.md#Authorization)

[[Back to top]](#__pageTop) [[Back to API list]](../../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../../README.md#documentation-for-models) [[Back to README]](../../../README.md)

# **get_payment**
<a name="get_payment"></a>
> Payment get_payment(payment_id)

Returns a payment.

### Example

* OAuth Authentication (Authorization):
```python
import Formance
from Formance.apis.tags import payments_api
from Formance.model.payment import Payment
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = Formance.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure OAuth2 access token for authorization: Authorization
configuration = Formance.Configuration(
    host = "http://localhost"
)
configuration.access_token = 'YOUR_ACCESS_TOKEN'
# Enter a context with an instance of the API client
with Formance.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = payments_api.PaymentsApi(api_client)

    # example passing only required values which don't have defaults set
    path_params = {
        'paymentId': "XXX",
    }
    try:
        # Returns a payment.
        api_response = api_instance.get_payment(
            path_params=path_params,
        )
        pprint(api_response)
    except Formance.ApiException as e:
        print("Exception when calling PaymentsApi->get_payment: %s\n" % e)
```
### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
path_params | RequestPathParams | |
accept_content_types | typing.Tuple[str] | default is ('application/json', ) | Tells the server the content type(s) that are accepted by the client
stream | bool | default is False | if True then the response.content will be streamed and loaded from a file like object. When downloading a file, set this to True to force the code to deserialize the content to a FileSchema file
timeout | typing.Optional[typing.Union[int, typing.Tuple]] | default is None | the timeout used by the rest client
skip_deserialization | bool | default is False | when True, headers and body will be unset and an instance of api_client.ApiResponseWithoutDeserialization will be returned

### path_params
#### RequestPathParams

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
paymentId | PaymentIdSchema | | 

# PaymentIdSchema

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
str,  | str,  |  | 

### Return Types, Responses

Code | Class | Description
------------- | ------------- | -------------
n/a | api_client.ApiResponseWithoutDeserialization | When skip_deserialization is True this response is returned
200 | [ApiResponseFor200](#get_payment.ApiResponseFor200) | A payment

#### get_payment.ApiResponseFor200
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | typing.Union[SchemaFor200ResponseBodyApplicationJson, ] |  |
headers | Unset | headers were not defined |

# SchemaFor200ResponseBodyApplicationJson
Type | Description  | Notes
------------- | ------------- | -------------
[**Payment**](../../models/Payment.md) |  | 


### Authorization

[Authorization](../../../README.md#Authorization)

[[Back to top]](#__pageTop) [[Back to API list]](../../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../../README.md#documentation-for-models) [[Back to README]](../../../README.md)

# **install_connector**
<a name="install_connector"></a>
> install_connector(connectorconnector_config)

Install connector

Install connector

### Example

* OAuth Authentication (Authorization):
```python
import Formance
from Formance.apis.tags import payments_api
from Formance.model.connector_config import ConnectorConfig
from Formance.model.connectors import Connectors
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = Formance.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure OAuth2 access token for authorization: Authorization
configuration = Formance.Configuration(
    host = "http://localhost"
)
configuration.access_token = 'YOUR_ACCESS_TOKEN'
# Enter a context with an instance of the API client
with Formance.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = payments_api.PaymentsApi(api_client)

    # example passing only required values which don't have defaults set
    path_params = {
        'connector': Connectors("STRIPE"),
    }
    body = ConnectorConfig(None)
    try:
        # Install connector
        api_response = api_instance.install_connector(
            path_params=path_params,
            body=body,
        )
    except Formance.ApiException as e:
        print("Exception when calling PaymentsApi->install_connector: %s\n" % e)
```
### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
body | typing.Union[SchemaForRequestBodyApplicationJson] | required |
path_params | RequestPathParams | |
content_type | str | optional, default is 'application/json' | Selects the schema and serialization of the request body
stream | bool | default is False | if True then the response.content will be streamed and loaded from a file like object. When downloading a file, set this to True to force the code to deserialize the content to a FileSchema file
timeout | typing.Optional[typing.Union[int, typing.Tuple]] | default is None | the timeout used by the rest client
skip_deserialization | bool | default is False | when True, headers and body will be unset and an instance of api_client.ApiResponseWithoutDeserialization will be returned

### body

# SchemaForRequestBodyApplicationJson
Type | Description  | Notes
------------- | ------------- | -------------
[**ConnectorConfig**](../../models/ConnectorConfig.md) |  | 


### path_params
#### RequestPathParams

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
connector | ConnectorSchema | | 

# ConnectorSchema
Type | Description  | Notes
------------- | ------------- | -------------
[**Connectors**](../../models/Connectors.md) |  | 


### Return Types, Responses

Code | Class | Description
------------- | ------------- | -------------
n/a | api_client.ApiResponseWithoutDeserialization | When skip_deserialization is True this response is returned
204 | [ApiResponseFor204](#install_connector.ApiResponseFor204) | Connector has been installed

#### install_connector.ApiResponseFor204
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | Unset | body was not defined |
headers | Unset | headers were not defined |

### Authorization

[Authorization](../../../README.md#Authorization)

[[Back to top]](#__pageTop) [[Back to API list]](../../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../../README.md#documentation-for-models) [[Back to README]](../../../README.md)

# **list_connector_tasks**
<a name="list_connector_tasks"></a>
> [bool, date, datetime, dict, float, int, list, str, none_type] list_connector_tasks(connector)

List connector tasks

List all tasks associated with this connector.

### Example

* OAuth Authentication (Authorization):
```python
import Formance
from Formance.apis.tags import payments_api
from Formance.model.task_descriptor_dummy_pay import TaskDescriptorDummyPay
from Formance.model.task_descriptor_wise import TaskDescriptorWise
from Formance.model.task_descriptor_modulr import TaskDescriptorModulr
from Formance.model.task_descriptor_stripe import TaskDescriptorStripe
from Formance.model.connectors import Connectors
from Formance.model.task_descriptor_banking_circle import TaskDescriptorBankingCircle
from Formance.model.task_descriptor_currency_cloud import TaskDescriptorCurrencyCloud
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = Formance.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure OAuth2 access token for authorization: Authorization
configuration = Formance.Configuration(
    host = "http://localhost"
)
configuration.access_token = 'YOUR_ACCESS_TOKEN'
# Enter a context with an instance of the API client
with Formance.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = payments_api.PaymentsApi(api_client)

    # example passing only required values which don't have defaults set
    path_params = {
        'connector': Connectors("STRIPE"),
    }
    try:
        # List connector tasks
        api_response = api_instance.list_connector_tasks(
            path_params=path_params,
        )
        pprint(api_response)
    except Formance.ApiException as e:
        print("Exception when calling PaymentsApi->list_connector_tasks: %s\n" % e)
```
### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
path_params | RequestPathParams | |
accept_content_types | typing.Tuple[str] | default is ('application/json', ) | Tells the server the content type(s) that are accepted by the client
stream | bool | default is False | if True then the response.content will be streamed and loaded from a file like object. When downloading a file, set this to True to force the code to deserialize the content to a FileSchema file
timeout | typing.Optional[typing.Union[int, typing.Tuple]] | default is None | the timeout used by the rest client
skip_deserialization | bool | default is False | when True, headers and body will be unset and an instance of api_client.ApiResponseWithoutDeserialization will be returned

### path_params
#### RequestPathParams

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
connector | ConnectorSchema | | 

# ConnectorSchema
Type | Description  | Notes
------------- | ------------- | -------------
[**Connectors**](../../models/Connectors.md) |  | 


### Return Types, Responses

Code | Class | Description
------------- | ------------- | -------------
n/a | api_client.ApiResponseWithoutDeserialization | When skip_deserialization is True this response is returned
200 | [ApiResponseFor200](#list_connector_tasks.ApiResponseFor200) | Task list

#### list_connector_tasks.ApiResponseFor200
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | typing.Union[SchemaFor200ResponseBodyApplicationJson, ] |  |
headers | Unset | headers were not defined |

# SchemaFor200ResponseBodyApplicationJson

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
list, tuple,  | tuple,  |  | 

### Tuple Items
Class Name | Input Type | Accessed Type | Description | Notes
------------- | ------------- | ------------- | ------------- | -------------
[items](#items) | dict, frozendict.frozendict, str, date, datetime, uuid.UUID, int, float, decimal.Decimal, bool, None, list, tuple, bytes, io.FileIO, io.BufferedReader,  | frozendict.frozendict, str, decimal.Decimal, BoolClass, NoneClass, tuple, bytes, FileIO |  | 

# items

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
dict, frozendict.frozendict, str, date, datetime, uuid.UUID, int, float, decimal.Decimal, bool, None, list, tuple, bytes, io.FileIO, io.BufferedReader,  | frozendict.frozendict, str, decimal.Decimal, BoolClass, NoneClass, tuple, bytes, FileIO |  | 

### Composed Schemas (allOf/anyOf/oneOf/not)
#### oneOf
Class Name | Input Type | Accessed Type | Description | Notes
------------- | ------------- | ------------- | ------------- | -------------
[TaskDescriptorStripe]({{complexTypePrefix}}TaskDescriptorStripe.md) | [**TaskDescriptorStripe**]({{complexTypePrefix}}TaskDescriptorStripe.md) | [**TaskDescriptorStripe**]({{complexTypePrefix}}TaskDescriptorStripe.md) |  | 
[TaskDescriptorWise]({{complexTypePrefix}}TaskDescriptorWise.md) | [**TaskDescriptorWise**]({{complexTypePrefix}}TaskDescriptorWise.md) | [**TaskDescriptorWise**]({{complexTypePrefix}}TaskDescriptorWise.md) |  | 
[TaskDescriptorCurrencyCloud]({{complexTypePrefix}}TaskDescriptorCurrencyCloud.md) | [**TaskDescriptorCurrencyCloud**]({{complexTypePrefix}}TaskDescriptorCurrencyCloud.md) | [**TaskDescriptorCurrencyCloud**]({{complexTypePrefix}}TaskDescriptorCurrencyCloud.md) |  | 
[TaskDescriptorDummyPay]({{complexTypePrefix}}TaskDescriptorDummyPay.md) | [**TaskDescriptorDummyPay**]({{complexTypePrefix}}TaskDescriptorDummyPay.md) | [**TaskDescriptorDummyPay**]({{complexTypePrefix}}TaskDescriptorDummyPay.md) |  | 
[TaskDescriptorModulr]({{complexTypePrefix}}TaskDescriptorModulr.md) | [**TaskDescriptorModulr**]({{complexTypePrefix}}TaskDescriptorModulr.md) | [**TaskDescriptorModulr**]({{complexTypePrefix}}TaskDescriptorModulr.md) |  | 
[TaskDescriptorBankingCircle]({{complexTypePrefix}}TaskDescriptorBankingCircle.md) | [**TaskDescriptorBankingCircle**]({{complexTypePrefix}}TaskDescriptorBankingCircle.md) | [**TaskDescriptorBankingCircle**]({{complexTypePrefix}}TaskDescriptorBankingCircle.md) |  | 

### Authorization

[Authorization](../../../README.md#Authorization)

[[Back to top]](#__pageTop) [[Back to API list]](../../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../../README.md#documentation-for-models) [[Back to README]](../../../README.md)

# **list_payments**
<a name="list_payments"></a>
> ListPaymentsResponse list_payments()

Returns a list of payments.

### Example

* OAuth Authentication (Authorization):
```python
import Formance
from Formance.apis.tags import payments_api
from Formance.model.list_payments_response import ListPaymentsResponse
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = Formance.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure OAuth2 access token for authorization: Authorization
configuration = Formance.Configuration(
    host = "http://localhost"
)
configuration.access_token = 'YOUR_ACCESS_TOKEN'
# Enter a context with an instance of the API client
with Formance.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = payments_api.PaymentsApi(api_client)

    # example passing only optional values
    query_params = {
        'limit': 10,
        'skip': 100,
        'sort': [
        "status"
    ],
    }
    try:
        # Returns a list of payments.
        api_response = api_instance.list_payments(
            query_params=query_params,
        )
        pprint(api_response)
    except Formance.ApiException as e:
        print("Exception when calling PaymentsApi->list_payments: %s\n" % e)
```
### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
query_params | RequestQueryParams | |
accept_content_types | typing.Tuple[str] | default is ('application/json', ) | Tells the server the content type(s) that are accepted by the client
stream | bool | default is False | if True then the response.content will be streamed and loaded from a file like object. When downloading a file, set this to True to force the code to deserialize the content to a FileSchema file
timeout | typing.Optional[typing.Union[int, typing.Tuple]] | default is None | the timeout used by the rest client
skip_deserialization | bool | default is False | when True, headers and body will be unset and an instance of api_client.ApiResponseWithoutDeserialization will be returned

### query_params
#### RequestQueryParams

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
limit | LimitSchema | | optional
skip | SkipSchema | | optional
sort | SortSchema | | optional


# LimitSchema

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
decimal.Decimal, int,  | decimal.Decimal,  |  | 

# SkipSchema

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
decimal.Decimal, int,  | decimal.Decimal,  |  | 

# SortSchema

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
list, tuple,  | tuple,  |  | 

### Tuple Items
Class Name | Input Type | Accessed Type | Description | Notes
------------- | ------------- | ------------- | ------------- | -------------
items | str,  | str,  |  | 

### Return Types, Responses

Code | Class | Description
------------- | ------------- | -------------
n/a | api_client.ApiResponseWithoutDeserialization | When skip_deserialization is True this response is returned
200 | [ApiResponseFor200](#list_payments.ApiResponseFor200) | A JSON array of payments

#### list_payments.ApiResponseFor200
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | typing.Union[SchemaFor200ResponseBodyApplicationJson, ] |  |
headers | Unset | headers were not defined |

# SchemaFor200ResponseBodyApplicationJson
Type | Description  | Notes
------------- | ------------- | -------------
[**ListPaymentsResponse**](../../models/ListPaymentsResponse.md) |  | 


### Authorization

[Authorization](../../../README.md#Authorization)

[[Back to top]](#__pageTop) [[Back to API list]](../../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../../README.md#documentation-for-models) [[Back to README]](../../../README.md)

# **paymentslist_accounts**
<a name="paymentslist_accounts"></a>
> ListAccountsResponse paymentslist_accounts()

Returns a list of accounts.

### Example

* OAuth Authentication (Authorization):
```python
import Formance
from Formance.apis.tags import payments_api
from Formance.model.list_accounts_response import ListAccountsResponse
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = Formance.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure OAuth2 access token for authorization: Authorization
configuration = Formance.Configuration(
    host = "http://localhost"
)
configuration.access_token = 'YOUR_ACCESS_TOKEN'
# Enter a context with an instance of the API client
with Formance.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = payments_api.PaymentsApi(api_client)

    # example passing only optional values
    query_params = {
        'limit': 10,
        'skip': 100,
        'sort': [
        "status"
    ],
    }
    try:
        # Returns a list of accounts.
        api_response = api_instance.paymentslist_accounts(
            query_params=query_params,
        )
        pprint(api_response)
    except Formance.ApiException as e:
        print("Exception when calling PaymentsApi->paymentslist_accounts: %s\n" % e)
```
### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
query_params | RequestQueryParams | |
accept_content_types | typing.Tuple[str] | default is ('application/json', ) | Tells the server the content type(s) that are accepted by the client
stream | bool | default is False | if True then the response.content will be streamed and loaded from a file like object. When downloading a file, set this to True to force the code to deserialize the content to a FileSchema file
timeout | typing.Optional[typing.Union[int, typing.Tuple]] | default is None | the timeout used by the rest client
skip_deserialization | bool | default is False | when True, headers and body will be unset and an instance of api_client.ApiResponseWithoutDeserialization will be returned

### query_params
#### RequestQueryParams

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
limit | LimitSchema | | optional
skip | SkipSchema | | optional
sort | SortSchema | | optional


# LimitSchema

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
decimal.Decimal, int,  | decimal.Decimal,  |  | 

# SkipSchema

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
decimal.Decimal, int,  | decimal.Decimal,  |  | 

# SortSchema

## Model Type Info
Input Type | Accessed Type | Description | Notes
------------ | ------------- | ------------- | -------------
list, tuple,  | tuple,  |  | 

### Tuple Items
Class Name | Input Type | Accessed Type | Description | Notes
------------- | ------------- | ------------- | ------------- | -------------
items | str,  | str,  |  | 

### Return Types, Responses

Code | Class | Description
------------- | ------------- | -------------
n/a | api_client.ApiResponseWithoutDeserialization | When skip_deserialization is True this response is returned
200 | [ApiResponseFor200](#paymentslist_accounts.ApiResponseFor200) | A JSON array of accounts

#### paymentslist_accounts.ApiResponseFor200
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | typing.Union[SchemaFor200ResponseBodyApplicationJson, ] |  |
headers | Unset | headers were not defined |

# SchemaFor200ResponseBodyApplicationJson
Type | Description  | Notes
------------- | ------------- | -------------
[**ListAccountsResponse**](../../models/ListAccountsResponse.md) |  | 


### Authorization

[Authorization](../../../README.md#Authorization)

[[Back to top]](#__pageTop) [[Back to API list]](../../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../../README.md#documentation-for-models) [[Back to README]](../../../README.md)

# **read_connector_config**
<a name="read_connector_config"></a>
> ConnectorConfig read_connector_config(connector)

Read connector config

Read connector config

### Example

* OAuth Authentication (Authorization):
```python
import Formance
from Formance.apis.tags import payments_api
from Formance.model.connector_config import ConnectorConfig
from Formance.model.connectors import Connectors
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = Formance.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure OAuth2 access token for authorization: Authorization
configuration = Formance.Configuration(
    host = "http://localhost"
)
configuration.access_token = 'YOUR_ACCESS_TOKEN'
# Enter a context with an instance of the API client
with Formance.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = payments_api.PaymentsApi(api_client)

    # example passing only required values which don't have defaults set
    path_params = {
        'connector': Connectors("STRIPE"),
    }
    try:
        # Read connector config
        api_response = api_instance.read_connector_config(
            path_params=path_params,
        )
        pprint(api_response)
    except Formance.ApiException as e:
        print("Exception when calling PaymentsApi->read_connector_config: %s\n" % e)
```
### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
path_params | RequestPathParams | |
accept_content_types | typing.Tuple[str] | default is ('application/json', ) | Tells the server the content type(s) that are accepted by the client
stream | bool | default is False | if True then the response.content will be streamed and loaded from a file like object. When downloading a file, set this to True to force the code to deserialize the content to a FileSchema file
timeout | typing.Optional[typing.Union[int, typing.Tuple]] | default is None | the timeout used by the rest client
skip_deserialization | bool | default is False | when True, headers and body will be unset and an instance of api_client.ApiResponseWithoutDeserialization will be returned

### path_params
#### RequestPathParams

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
connector | ConnectorSchema | | 

# ConnectorSchema
Type | Description  | Notes
------------- | ------------- | -------------
[**Connectors**](../../models/Connectors.md) |  | 


### Return Types, Responses

Code | Class | Description
------------- | ------------- | -------------
n/a | api_client.ApiResponseWithoutDeserialization | When skip_deserialization is True this response is returned
200 | [ApiResponseFor200](#read_connector_config.ApiResponseFor200) | Connector config

#### read_connector_config.ApiResponseFor200
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | typing.Union[SchemaFor200ResponseBodyApplicationJson, ] |  |
headers | Unset | headers were not defined |

# SchemaFor200ResponseBodyApplicationJson
Type | Description  | Notes
------------- | ------------- | -------------
[**ConnectorConfig**](../../models/ConnectorConfig.md) |  | 


### Authorization

[Authorization](../../../README.md#Authorization)

[[Back to top]](#__pageTop) [[Back to API list]](../../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../../README.md#documentation-for-models) [[Back to README]](../../../README.md)

# **reset_connector**
<a name="reset_connector"></a>
> reset_connector(connector)

Reset connector

Reset connector. Will remove the connector and ALL PAYMENTS generated with it.

### Example

* OAuth Authentication (Authorization):
```python
import Formance
from Formance.apis.tags import payments_api
from Formance.model.connectors import Connectors
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = Formance.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure OAuth2 access token for authorization: Authorization
configuration = Formance.Configuration(
    host = "http://localhost"
)
configuration.access_token = 'YOUR_ACCESS_TOKEN'
# Enter a context with an instance of the API client
with Formance.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = payments_api.PaymentsApi(api_client)

    # example passing only required values which don't have defaults set
    path_params = {
        'connector': Connectors("STRIPE"),
    }
    try:
        # Reset connector
        api_response = api_instance.reset_connector(
            path_params=path_params,
        )
    except Formance.ApiException as e:
        print("Exception when calling PaymentsApi->reset_connector: %s\n" % e)
```
### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
path_params | RequestPathParams | |
stream | bool | default is False | if True then the response.content will be streamed and loaded from a file like object. When downloading a file, set this to True to force the code to deserialize the content to a FileSchema file
timeout | typing.Optional[typing.Union[int, typing.Tuple]] | default is None | the timeout used by the rest client
skip_deserialization | bool | default is False | when True, headers and body will be unset and an instance of api_client.ApiResponseWithoutDeserialization will be returned

### path_params
#### RequestPathParams

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
connector | ConnectorSchema | | 

# ConnectorSchema
Type | Description  | Notes
------------- | ------------- | -------------
[**Connectors**](../../models/Connectors.md) |  | 


### Return Types, Responses

Code | Class | Description
------------- | ------------- | -------------
n/a | api_client.ApiResponseWithoutDeserialization | When skip_deserialization is True this response is returned
204 | [ApiResponseFor204](#reset_connector.ApiResponseFor204) | Connector has been reset

#### reset_connector.ApiResponseFor204
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | Unset | body was not defined |
headers | Unset | headers were not defined |

### Authorization

[Authorization](../../../README.md#Authorization)

[[Back to top]](#__pageTop) [[Back to API list]](../../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../../README.md#documentation-for-models) [[Back to README]](../../../README.md)

# **uninstall_connector**
<a name="uninstall_connector"></a>
> uninstall_connector(connector)

Uninstall connector

Uninstall  connector

### Example

* OAuth Authentication (Authorization):
```python
import Formance
from Formance.apis.tags import payments_api
from Formance.model.connectors import Connectors
from pprint import pprint
# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = Formance.Configuration(
    host = "http://localhost"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure OAuth2 access token for authorization: Authorization
configuration = Formance.Configuration(
    host = "http://localhost"
)
configuration.access_token = 'YOUR_ACCESS_TOKEN'
# Enter a context with an instance of the API client
with Formance.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = payments_api.PaymentsApi(api_client)

    # example passing only required values which don't have defaults set
    path_params = {
        'connector': Connectors("STRIPE"),
    }
    try:
        # Uninstall connector
        api_response = api_instance.uninstall_connector(
            path_params=path_params,
        )
    except Formance.ApiException as e:
        print("Exception when calling PaymentsApi->uninstall_connector: %s\n" % e)
```
### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
path_params | RequestPathParams | |
stream | bool | default is False | if True then the response.content will be streamed and loaded from a file like object. When downloading a file, set this to True to force the code to deserialize the content to a FileSchema file
timeout | typing.Optional[typing.Union[int, typing.Tuple]] | default is None | the timeout used by the rest client
skip_deserialization | bool | default is False | when True, headers and body will be unset and an instance of api_client.ApiResponseWithoutDeserialization will be returned

### path_params
#### RequestPathParams

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
connector | ConnectorSchema | | 

# ConnectorSchema
Type | Description  | Notes
------------- | ------------- | -------------
[**Connectors**](../../models/Connectors.md) |  | 


### Return Types, Responses

Code | Class | Description
------------- | ------------- | -------------
n/a | api_client.ApiResponseWithoutDeserialization | When skip_deserialization is True this response is returned
204 | [ApiResponseFor204](#uninstall_connector.ApiResponseFor204) | Connector has been uninstalled

#### uninstall_connector.ApiResponseFor204
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
response | urllib3.HTTPResponse | Raw response |
body | Unset | body was not defined |
headers | Unset | headers were not defined |

### Authorization

[Authorization](../../../README.md#Authorization)

[[Back to top]](#__pageTop) [[Back to API list]](../../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../../README.md#documentation-for-models) [[Back to README]](../../../README.md)

