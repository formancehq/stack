# DummyPayConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FilePollingPeriod** | Pointer to **interface{}** | The frequency at which the connector will try to fetch new payment objects from the directory | [optional] [default to 10s]
**FileGenerationPeriod** | Pointer to **interface{}** | The frequency at which the connector will create new payment objects in the directory | [optional] [default to 10s]
**Directory** | **interface{}** |  | 

## Methods

### NewDummyPayConfig

`func NewDummyPayConfig(directory interface{}, ) *DummyPayConfig`

NewDummyPayConfig instantiates a new DummyPayConfig object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDummyPayConfigWithDefaults

`func NewDummyPayConfigWithDefaults() *DummyPayConfig`

NewDummyPayConfigWithDefaults instantiates a new DummyPayConfig object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFilePollingPeriod

`func (o *DummyPayConfig) GetFilePollingPeriod() interface{}`

GetFilePollingPeriod returns the FilePollingPeriod field if non-nil, zero value otherwise.

### GetFilePollingPeriodOk

`func (o *DummyPayConfig) GetFilePollingPeriodOk() (*interface{}, bool)`

GetFilePollingPeriodOk returns a tuple with the FilePollingPeriod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilePollingPeriod

`func (o *DummyPayConfig) SetFilePollingPeriod(v interface{})`

SetFilePollingPeriod sets FilePollingPeriod field to given value.

### HasFilePollingPeriod

`func (o *DummyPayConfig) HasFilePollingPeriod() bool`

HasFilePollingPeriod returns a boolean if a field has been set.

### SetFilePollingPeriodNil

`func (o *DummyPayConfig) SetFilePollingPeriodNil(b bool)`

 SetFilePollingPeriodNil sets the value for FilePollingPeriod to be an explicit nil

### UnsetFilePollingPeriod
`func (o *DummyPayConfig) UnsetFilePollingPeriod()`

UnsetFilePollingPeriod ensures that no value is present for FilePollingPeriod, not even an explicit nil
### GetFileGenerationPeriod

`func (o *DummyPayConfig) GetFileGenerationPeriod() interface{}`

GetFileGenerationPeriod returns the FileGenerationPeriod field if non-nil, zero value otherwise.

### GetFileGenerationPeriodOk

`func (o *DummyPayConfig) GetFileGenerationPeriodOk() (*interface{}, bool)`

GetFileGenerationPeriodOk returns a tuple with the FileGenerationPeriod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFileGenerationPeriod

`func (o *DummyPayConfig) SetFileGenerationPeriod(v interface{})`

SetFileGenerationPeriod sets FileGenerationPeriod field to given value.

### HasFileGenerationPeriod

`func (o *DummyPayConfig) HasFileGenerationPeriod() bool`

HasFileGenerationPeriod returns a boolean if a field has been set.

### SetFileGenerationPeriodNil

`func (o *DummyPayConfig) SetFileGenerationPeriodNil(b bool)`

 SetFileGenerationPeriodNil sets the value for FileGenerationPeriod to be an explicit nil

### UnsetFileGenerationPeriod
`func (o *DummyPayConfig) UnsetFileGenerationPeriod()`

UnsetFileGenerationPeriod ensures that no value is present for FileGenerationPeriod, not even an explicit nil
### GetDirectory

`func (o *DummyPayConfig) GetDirectory() interface{}`

GetDirectory returns the Directory field if non-nil, zero value otherwise.

### GetDirectoryOk

`func (o *DummyPayConfig) GetDirectoryOk() (*interface{}, bool)`

GetDirectoryOk returns a tuple with the Directory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDirectory

`func (o *DummyPayConfig) SetDirectory(v interface{})`

SetDirectory sets Directory field to given value.


### SetDirectoryNil

`func (o *DummyPayConfig) SetDirectoryNil(b bool)`

 SetDirectoryNil sets the value for Directory to be an explicit nil

### UnsetDirectory
`func (o *DummyPayConfig) UnsetDirectory()`

UnsetDirectory ensures that no value is present for Directory, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


