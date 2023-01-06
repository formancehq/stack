# StripeConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PollingPeriod** | Pointer to **interface{}** | The frequency at which the connector will try to fetch new BalanceTransaction objects from Stripe api | [optional] [default to 120s]
**ApiKey** | **interface{}** |  |
**PageSize** | Pointer to **interface{}** | Number of BalanceTransaction to fetch at each polling interval.  | [optional] [default to 10]

## Methods

### NewStripeConfig

`func NewStripeConfig(apiKey interface{}, ) *StripeConfig`

NewStripeConfig instantiates a new StripeConfig object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStripeConfigWithDefaults

`func NewStripeConfigWithDefaults() *StripeConfig`

NewStripeConfigWithDefaults instantiates a new StripeConfig object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPollingPeriod

`func (o *StripeConfig) GetPollingPeriod() interface{}`

GetPollingPeriod returns the PollingPeriod field if non-nil, zero value otherwise.

### GetPollingPeriodOk

`func (o *StripeConfig) GetPollingPeriodOk() (*interface{}, bool)`

GetPollingPeriodOk returns a tuple with the PollingPeriod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPollingPeriod

`func (o *StripeConfig) SetPollingPeriod(v interface{})`

SetPollingPeriod sets PollingPeriod field to given value.

### HasPollingPeriod

`func (o *StripeConfig) HasPollingPeriod() bool`

HasPollingPeriod returns a boolean if a field has been set.

### SetPollingPeriodNil

`func (o *StripeConfig) SetPollingPeriodNil(b bool)`

 SetPollingPeriodNil sets the value for PollingPeriod to be an explicit nil

### UnsetPollingPeriod
`func (o *StripeConfig) UnsetPollingPeriod()`

UnsetPollingPeriod ensures that no value is present for PollingPeriod, not even an explicit nil
### GetApiKey

`func (o *StripeConfig) GetApiKey() interface{}`

GetApiKey returns the ApiKey field if non-nil, zero value otherwise.

### GetApiKeyOk

`func (o *StripeConfig) GetApiKeyOk() (*interface{}, bool)`

GetApiKeyOk returns a tuple with the ApiKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiKey

`func (o *StripeConfig) SetApiKey(v interface{})`

SetApiKey sets ApiKey field to given value.


### SetApiKeyNil

`func (o *StripeConfig) SetApiKeyNil(b bool)`

 SetApiKeyNil sets the value for ApiKey to be an explicit nil

### UnsetApiKey
`func (o *StripeConfig) UnsetApiKey()`

UnsetApiKey ensures that no value is present for ApiKey, not even an explicit nil
### GetPageSize

`func (o *StripeConfig) GetPageSize() interface{}`

GetPageSize returns the PageSize field if non-nil, zero value otherwise.

### GetPageSizeOk

`func (o *StripeConfig) GetPageSizeOk() (*interface{}, bool)`

GetPageSizeOk returns a tuple with the PageSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPageSize

`func (o *StripeConfig) SetPageSize(v interface{})`

SetPageSize sets PageSize field to given value.

### HasPageSize

`func (o *StripeConfig) HasPageSize() bool`

HasPageSize returns a boolean if a field has been set.

### SetPageSizeNil

`func (o *StripeConfig) SetPageSizeNil(b bool)`

 SetPageSizeNil sets the value for PageSize to be an explicit nil

### UnsetPageSize
`func (o *StripeConfig) UnsetPageSize()`

UnsetPageSize ensures that no value is present for PageSize, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
