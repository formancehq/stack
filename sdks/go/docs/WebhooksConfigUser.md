# WebhooksConfigUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Endpoint** | **string** |  |
**Secret** | Pointer to **string** |  | [optional]
**EventTypes** | **[]string** |  |

## Methods

### NewWebhooksConfigUser

`func NewWebhooksConfigUser(endpoint string, eventTypes []string, ) *WebhooksConfigUser`

NewWebhooksConfigUser instantiates a new WebhooksConfigUser object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebhooksConfigUserWithDefaults

`func NewWebhooksConfigUserWithDefaults() *WebhooksConfigUser`

NewWebhooksConfigUserWithDefaults instantiates a new WebhooksConfigUser object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEndpoint

`func (o *WebhooksConfigUser) GetEndpoint() string`

GetEndpoint returns the Endpoint field if non-nil, zero value otherwise.

### GetEndpointOk

`func (o *WebhooksConfigUser) GetEndpointOk() (*string, bool)`

GetEndpointOk returns a tuple with the Endpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpoint

`func (o *WebhooksConfigUser) SetEndpoint(v string)`

SetEndpoint sets Endpoint field to given value.


### GetSecret

`func (o *WebhooksConfigUser) GetSecret() string`

GetSecret returns the Secret field if non-nil, zero value otherwise.

### GetSecretOk

`func (o *WebhooksConfigUser) GetSecretOk() (*string, bool)`

GetSecretOk returns a tuple with the Secret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecret

`func (o *WebhooksConfigUser) SetSecret(v string)`

SetSecret sets Secret field to given value.

### HasSecret

`func (o *WebhooksConfigUser) HasSecret() bool`

HasSecret returns a boolean if a field has been set.

### GetEventTypes

`func (o *WebhooksConfigUser) GetEventTypes() []string`

GetEventTypes returns the EventTypes field if non-nil, zero value otherwise.

### GetEventTypesOk

`func (o *WebhooksConfigUser) GetEventTypesOk() (*[]string, bool)`

GetEventTypesOk returns a tuple with the EventTypes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEventTypes

`func (o *WebhooksConfigUser) SetEventTypes(v []string)`

SetEventTypes sets EventTypes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
