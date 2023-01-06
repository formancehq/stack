# CurrencyCloudConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApiKey** | **interface{}** |  |
**LoginID** | **interface{}** | Username of the API Key holder |
**PollingPeriod** | Pointer to **interface{}** | The frequency at which the connector will fetch transactions | [optional]
**Endpoint** | Pointer to **interface{}** | The endpoint to use for the API. Defaults to https://devapi.currencycloud.com | [optional]

## Methods

### NewCurrencyCloudConfig

`func NewCurrencyCloudConfig(apiKey interface{}, loginID interface{}, ) *CurrencyCloudConfig`

NewCurrencyCloudConfig instantiates a new CurrencyCloudConfig object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCurrencyCloudConfigWithDefaults

`func NewCurrencyCloudConfigWithDefaults() *CurrencyCloudConfig`

NewCurrencyCloudConfigWithDefaults instantiates a new CurrencyCloudConfig object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetApiKey

`func (o *CurrencyCloudConfig) GetApiKey() interface{}`

GetApiKey returns the ApiKey field if non-nil, zero value otherwise.

### GetApiKeyOk

`func (o *CurrencyCloudConfig) GetApiKeyOk() (*interface{}, bool)`

GetApiKeyOk returns a tuple with the ApiKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiKey

`func (o *CurrencyCloudConfig) SetApiKey(v interface{})`

SetApiKey sets ApiKey field to given value.


### SetApiKeyNil

`func (o *CurrencyCloudConfig) SetApiKeyNil(b bool)`

 SetApiKeyNil sets the value for ApiKey to be an explicit nil

### UnsetApiKey
`func (o *CurrencyCloudConfig) UnsetApiKey()`

UnsetApiKey ensures that no value is present for ApiKey, not even an explicit nil
### GetLoginID

`func (o *CurrencyCloudConfig) GetLoginID() interface{}`

GetLoginID returns the LoginID field if non-nil, zero value otherwise.

### GetLoginIDOk

`func (o *CurrencyCloudConfig) GetLoginIDOk() (*interface{}, bool)`

GetLoginIDOk returns a tuple with the LoginID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLoginID

`func (o *CurrencyCloudConfig) SetLoginID(v interface{})`

SetLoginID sets LoginID field to given value.


### SetLoginIDNil

`func (o *CurrencyCloudConfig) SetLoginIDNil(b bool)`

 SetLoginIDNil sets the value for LoginID to be an explicit nil

### UnsetLoginID
`func (o *CurrencyCloudConfig) UnsetLoginID()`

UnsetLoginID ensures that no value is present for LoginID, not even an explicit nil
### GetPollingPeriod

`func (o *CurrencyCloudConfig) GetPollingPeriod() interface{}`

GetPollingPeriod returns the PollingPeriod field if non-nil, zero value otherwise.

### GetPollingPeriodOk

`func (o *CurrencyCloudConfig) GetPollingPeriodOk() (*interface{}, bool)`

GetPollingPeriodOk returns a tuple with the PollingPeriod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPollingPeriod

`func (o *CurrencyCloudConfig) SetPollingPeriod(v interface{})`

SetPollingPeriod sets PollingPeriod field to given value.

### HasPollingPeriod

`func (o *CurrencyCloudConfig) HasPollingPeriod() bool`

HasPollingPeriod returns a boolean if a field has been set.

### SetPollingPeriodNil

`func (o *CurrencyCloudConfig) SetPollingPeriodNil(b bool)`

 SetPollingPeriodNil sets the value for PollingPeriod to be an explicit nil

### UnsetPollingPeriod
`func (o *CurrencyCloudConfig) UnsetPollingPeriod()`

UnsetPollingPeriod ensures that no value is present for PollingPeriod, not even an explicit nil
### GetEndpoint

`func (o *CurrencyCloudConfig) GetEndpoint() interface{}`

GetEndpoint returns the Endpoint field if non-nil, zero value otherwise.

### GetEndpointOk

`func (o *CurrencyCloudConfig) GetEndpointOk() (*interface{}, bool)`

GetEndpointOk returns a tuple with the Endpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpoint

`func (o *CurrencyCloudConfig) SetEndpoint(v interface{})`

SetEndpoint sets Endpoint field to given value.

### HasEndpoint

`func (o *CurrencyCloudConfig) HasEndpoint() bool`

HasEndpoint returns a boolean if a field has been set.

### SetEndpointNil

`func (o *CurrencyCloudConfig) SetEndpointNil(b bool)`

 SetEndpointNil sets the value for Endpoint to be an explicit nil

### UnsetEndpoint
`func (o *CurrencyCloudConfig) UnsetEndpoint()`

UnsetEndpoint ensures that no value is present for Endpoint, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
